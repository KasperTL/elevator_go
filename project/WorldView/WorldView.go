package WorldView

import (
	"fmt"
	"project/ElevatorDriver"
	"project/config"
	"project/elevio"
	"strconv"
)

type OrderState int

const (
	OrderIdle      = 0
	OrderPending   = 1
	OrderConfirmed = 2
)

type WorldView struct {
	SenderID       int
	AliveList      [config.NumElevators]bool
	ElevatorStates [config.NumElevators]ElevatorDriver.Elevator
	Orders         [config.NumElevators][config.NumFloors][2 + config.NumElevators]OrderState
}

func InitWorldView(nodeID int) WorldView {
	view := WorldView{SenderID: nodeID}
	view.AliveList[nodeID] = true
	return view
}

func PrintHallOrders(wv WorldView) {
	fmt.Println("\n=== Hall Orders ===")
	for floor := 0; floor < config.NumFloors; floor++ {
		fmt.Printf("Floor %d: ", floor)
		for elevator := 0; elevator < config.NumElevators; elevator++ {
			upState := orderStateToString(wv.Orders[elevator][floor][elevio.BT_HallUp])
			downState := orderStateToString(wv.Orders[elevator][floor][elevio.BT_HallDown])
			fmt.Printf("E%d[↑%s ↓%s] ", elevator, upState, downState)
		}
		fmt.Println()
	}
	fmt.Println()
}

func orderStateToString(state OrderState) string {
	switch state {
	case OrderIdle:
		return "Idle"
	case OrderPending:
		return "Pend"
	case OrderConfirmed:
		return "Conf"
	default:
		return "?"
	}
}

func allPeersUpToDateOrAhead(peers []OrderState, stateA OrderState, stateB OrderState) bool {
	for _, p := range peers {
		if p != stateA && p != stateB {
			return false
		}
	}
	return true
}

func anyPeerAhead(peers []OrderState, state OrderState) bool {
	for _, p := range peers {
		if p == state {
			return true
		}
	}
	return false
}

func isPeerAhead(peerOrderState OrderState, myOrderState OrderState) bool {
	if peerOrderState == myOrderState {
		return false
	}
	switch myOrderState {
	case OrderIdle:
		return peerOrderState == OrderPending || peerOrderState == OrderConfirmed

	case OrderPending:
		return peerOrderState == OrderConfirmed

	case OrderConfirmed:
		return peerOrderState == OrderIdle

	default:
		return false
	}
}

func syncOnRejon(
	localOrders [config.NumElevators][config.NumFloors][config.NumButtons]OrderState,
	alivePeers []int,
) [config.NumElevators][config.NumFloors][config.NumButtons]OrderState {

	if len(alivePeers) == 0 {
		return localOrders
	}

	for elevator := 0; elevator < config.NumElevators; elevator++ {
		for floor := 0; floor < config.NumFloors; floor++ {
			for button := 0; button < config.NumButtons; button++ {

				for peerID := range alivePeers {
					peerState := localOrders[peerID][floor][button]

					if isPeerAhead(peerState, localOrders[elevator][floor][button]) {
						localOrders[elevator][floor][button] = peerState
					}
				}

			}
		}
	}
	return localOrders
}
func updateOrders(
	orders [config.NumElevators][config.NumFloors][2 + config.NumElevators]OrderState,
	NodeID int,
	alivePeers []string,
) [config.NumElevators][config.NumFloors][2 + config.NumElevators]OrderState {

	for floor := 0; floor < config.NumFloors; floor++ {
		for button := 0; button < 2+config.NumElevators; button++ {
			currentOrderState := orders[NodeID][floor][button]
			newOrderState := currentOrderState

			var peersOrderView []OrderState

			for _, alivePeerStr := range alivePeers {
				alivePeerID, err := strconv.Atoi(alivePeerStr)
				if err != nil {
					continue
				}
				if alivePeerID != NodeID {
					peersOrderView = append(peersOrderView, orders[alivePeerID][floor][button])
				}
			}

			switch currentOrderState {
			case OrderIdle:
				if allPeersUpToDateOrAhead(peersOrderView, OrderIdle, OrderPending) && anyPeerAhead(peersOrderView, OrderPending) {
					newOrderState = OrderPending
				}
			case OrderPending:
				if allPeersUpToDateOrAhead(peersOrderView, OrderPending, OrderConfirmed) {
					newOrderState = OrderConfirmed
				}
			case OrderConfirmed:
				if allPeersUpToDateOrAhead(peersOrderView, OrderConfirmed, OrderIdle) && anyPeerAhead(peersOrderView, OrderIdle) {
					newOrderState = OrderIdle
				}
			}
			orders[NodeID][floor][button] = newOrderState
		}
	}
	return orders
}

func updatePeerStatusInMyWorldView(myWorldView WorldView, peerWorldView WorldView) WorldView {
	myWorldView.ElevatorStates[peerWorldView.SenderID] = peerWorldView.ElevatorStates[peerWorldView.SenderID]
	myWorldView.AliveList[peerWorldView.SenderID] = peerWorldView.AliveList[peerWorldView.SenderID]
	myWorldView.Orders[peerWorldView.SenderID] = peerWorldView.Orders[peerWorldView.SenderID]
	return myWorldView
}

func CabOrdersAsBool(Orders [config.NumFloors][2 + config.NumElevators]OrderState, nodeID int) [config.NumFloors]bool {
	var cabOrders [config.NumFloors]bool
	for floor := 0; floor < config.NumFloors; floor++ {
		cabOrders[floor] = Orders[floor][2+nodeID] == OrderConfirmed
	}
	return cabOrders
}

func HallOrdersAsBool(Orders [config.NumFloors][2 + config.NumElevators]OrderState) [config.NumFloors][2]bool {
	var hallOrders [config.NumFloors][2]bool
	for floor := 0; floor < config.NumFloors; floor++ {
		hallOrders[floor][0] = Orders[floor][elevio.BT_HallUp] == OrderConfirmed
		hallOrders[floor][1] = Orders[floor][elevio.BT_HallDown] == OrderConfirmed
	}
	return hallOrders
}

func setOrderLights(myWorldView WorldView) {
	for floor := 0; floor < config.NumFloors; floor++ {
		for button := 0; button < config.NumButtons; button++ {
			orderState := myWorldView.Orders[myWorldView.SenderID][floor][button]
			buttonType := elevio.ButtonType(button)
			switch orderState {
			case OrderConfirmed:
				elevio.SetButtonLamp(buttonType, floor, true)
			case OrderIdle:
				elevio.SetButtonLamp(buttonType, floor, false)
			case OrderPending:
				continue
			}
		}
	}
}
