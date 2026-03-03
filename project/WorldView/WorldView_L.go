package WorldView

import (
	"project/ElevatorDriver"
	"project/config"
	"project/elevio"
)

//Added in import from elevio

type OrderState int

const (
	OrderIdle      = 0
	OrderPending   = 1
	OrderConfirmed = 2
)

// Added in struct for Order state and cyclic counter for order
type OrderInfo struct {
	State OrderState
	Epoch uint64 //cyclic counter for each order
	//Helps to tell which order is newer, higher = newer
}

// ElevatorState -information about the elevator
type ElevatorState struct {
	Floor     int
	Direction ElevatorDriver.ElevatorDirection
	Behaviour ElevatorDriver.ElevatorBehaviour
}

// what gets broadcasted over UDP
type WorldView struct {
	SenderID int
	//Epoch    uint64 //- can be useful but use cyclic counter only for remove all epoch from worldview later

	//Hall orders [floor][0=up,1=down]
	//global
	HallOrders [config.NumFloors][2]OrderInfo

	//Elevator state needed by assigner
	ElevatorStates [config.NumElevators]ElevatorState

	//cab calls
	//local for each eleveator but others need to know if busy or not
	MyCabCalls [config.NumFloors]bool
}

func InitWorldView(nodeID int) WorldView {
	view := WorldView{
		SenderID: nodeID,
		//Epoch:    0,
	}
	for f := 0; f < config.NumFloors; f++ {
		view.HallOrders[f][0] = OrderInfo{State: OrderIdle, Epoch: 0}
		view.HallOrders[f][1] = OrderInfo{State: OrderIdle, Epoch: 0}
		view.MyCabCalls[f] = false
	}
	// Initialize elevator states to sane defaults
	for i := 0; i < config.NumElevators; i++ {
		view.ElevatorStates[i] = ElevatorState{
			Floor:     -1,
			Direction: ElevatorDriver.ED_Down,
			Behaviour: ElevatorDriver.EB_Idle,
		}
	}
	return view
}

// button press
// hall -->Order pending + light ON once confirmed
// cab --> Stored in MyCabCalls (Elevator driver handles cab calls?)
// dir for direction
func OnButtonPress(view *WorldView, btn elevio.ButtonEvent) {
	if btn.Button == elevio.BT_Cab {
		view.MyCabCalls[btn.Floor] = true
		//does cab call turn on light?
		return
	}

	dir := buttonToDir(btn.Button)
	if view.HallOrders[btn.Floor][dir].State == OrderIdle {
		//view.Epoch++
		view.HallOrders[btn.Floor][dir] = OrderInfo{
			State: OrderPending,
			Epoch: view.HallOrders[btn.Floor][dir].Epoch + 1,
		}
		//elevio.SetButtonLamp(btn.Button, btn.Floor, true)
		//dont turn on hall light here, wait for orderConfirmed
	}
}

// order complete
// hall --> OrderIdle = light OFF
// Cab -->clear from MyCabCalls
func OnOrderComplete(view *WorldView, btn elevio.ButtonEvent) {
	if btn.Button == elevio.BT_Cab {
		view.MyCabCalls[btn.Floor] = false
		//need cab light?
		return
	}

	dir := buttonToDir(btn.Button)
	//view.Epoch++
	view.HallOrders[btn.Floor][dir] = OrderInfo{
		State: OrderIdle,
		Epoch: view.HallOrders[btn.Floor][dir].Epoch + 1,
	}
	elevio.SetButtonLamp(btn.Button, btn.Floor, false)
}

//Function called when receiving WorldView from another elevator
//higher epoch/cyclic counter =newer = use it

func MergeWorldView(mine *WorldView, peer WorldView) {
	if peer.SenderID == mine.SenderID {
		return
	}

	for f := 0; f < config.NumFloors; f++ {
		for d := 0; d < 2; d++ {
			before := mine.HallOrders[f][d].State
			merged := mergeOrder(mine.HallOrders[f][d], peer.HallOrders[f][d])
			mine.HallOrders[f][d] = merged

			// If the merged state changed to Idle, turn light off
			if merged.State == OrderIdle && before != OrderIdle {
				btn := dirToButton(d, f)
				elevio.SetButtonLamp(btn.Button, f, false)
			}
		}
		// Merge peer elevator state
		mine.ElevatorStates[peer.SenderID] = peer.ElevatorStates[peer.SenderID]
	}
	//checkConfirmation(mine, peer)
}

// Higher cyclic counter/epoch/newer order =higher state wins
func mergeOrder(mine OrderInfo, peer OrderInfo) OrderInfo {
	if peer.Epoch > mine.Epoch {
		return peer
	}
	if mine.Epoch > peer.Epoch {
		return mine
	}
	if peer.State > mine.State {
		return peer
	}
	return mine
}

// pending to confirmed when peer also knows about order
// OrderPending --> OrderConfimred when all Alive peers
// Will help WorldView with consensus
func checkConfirmation(mine *WorldView, peerViews map[int]WorldView, alivePeers []int) {
	for f := 0; f < config.NumFloors; f++ {
		for d := 0; d < 2; d++ {
			if mine.HallOrders[f][d].State != OrderPending {
				continue
			}

			allSeen := true
			for _, peerID := range alivePeers {
				if peerID == mine.SenderID {
					continue
				}
				pv, exists := peerViews[peerID]
				if !exists || pv.HallOrders[f][d].State == OrderIdle {
					allSeen = false
					break
				}
			}

			if allSeen {
				//mine.Epoch++
				mine.HallOrders[f][d] = OrderInfo{
					State: OrderConfirmed,
					Epoch: mine.HallOrders[f][d].Epoch + 1,
				}
				// Turn on hall light only now, order is safely distributed
				btn := dirToButton(d, f)
				elevio.SetButtonLamp(btn.Button, f, true)
			}
		}
	}
}

// For assigner to use, only confirm orders are ready to be assigned
func GetConfirmedOrders(view WorldView) [config.NumFloors][2]bool {
	var confirmed [config.NumFloors][2]bool
	for f := 0; f < config.NumFloors; f++ {
		for d := 0; d < 2; d++ {
			confirmed[f][d] = (view.HallOrders[f][d].State == OrderConfirmed)
		}
	}
	return confirmed
}

func buttonToDir(btn elevio.ButtonType) int {
	if btn == elevio.BT_HallUp {
		return 0
	}
	return 1
}

// Helps with repetitive code to help indicate if it is a up/down button
func dirToButton(dir int, floor int) elevio.ButtonEvent {
	if dir == 0 {
		return elevio.ButtonEvent{Floor: floor, Button: elevio.BT_HallUp}
	}
	return elevio.ButtonEvent{Floor: floor, Button: elevio.BT_HallDown}
}
