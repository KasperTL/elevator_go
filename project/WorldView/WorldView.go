package WorldView

import (
	"project/ElevatorDriver"
	"project/Network/peers"
	"project/config"
	"project/elevio"
	"strconv"
)

type OrderState int

const (
	OrderIdle      = 0
	OrderPending   = 1
	OrderConfirmed = 2
	OrderComplete  = 3
)

type WorldView struct {
	NodeID           int
	AliveList        [config.NumElevators]bool
	ElevatorStates   [config.NumElevators]ElevatorDriver.Elevator
	Orders           [config.NumElevators][config.NumFloors][config.NumOrderTypes]OrderState
	CabOrderRecovery [config.NumElevators][config.NumFloors]OrderState
}

func InitWorldView(nodeID int) WorldView {
	view := WorldView{NodeID: nodeID}
	view.AliveList[nodeID] = true
	return view
}

func (wv *WorldView) setAliveElevators(peers peers.PeerUpdate) {
	wv.AliveList = [config.NumElevators]bool{}
	wv.AliveList[wv.NodeID] = true
	for _, peerIDstr := range peers.Peers {
		peerID, err := strconv.Atoi(peerIDstr)
		if err != nil {
			continue
		}
		wv.AliveList[peerID] = true
	}
}

func (wv *WorldView) stashLostNodesCabOrders(peers peers.PeerUpdate) {
	for _, lostIDstr := range peers.Lost {
		lostID, err := strconv.Atoi(lostIDstr)
		if err != nil {
			continue
		}
		for floor := range config.NumFloors {
			if wv.NodeID == lostID {
				continue
			}
			wv.CabOrderRecovery[lostID][floor] = wv.Orders[lostID][floor][lostID+2]
			wv.Orders[lostID][floor][lostID+2] = OrderIdle
		}
	}
}

func amIOnline(peers peers.PeerUpdate) bool {
	var online bool
	if len(peers.Peers) > 1 {
		online = true
	} else {
		online = false
	}
	return online
}

func getHighestOrderState(peers []OrderState) OrderState {
	var highestOrderState OrderState
	for _, p := range peers {
		if p > highestOrderState {
			highestOrderState = p
		}
	}
	return highestOrderState
}

func updateOrders(
	orders [config.NumElevators][config.NumFloors][config.NumOrderTypes]OrderState,
	NodeID int,
	alivePeers [config.NumElevators]bool,
) [config.NumElevators][config.NumFloors][config.NumOrderTypes]OrderState {

	for floor := 0; floor < config.NumFloors; floor++ {
		for button := 0; button < config.NumOrderTypes; button++ {
			peersOrderView := alivePeersOrderView(orders, alivePeers, floor, button)

			currentOrderState := orders[NodeID][floor][button]

			switch currentOrderState {
			case OrderIdle:
				if getHighestOrderState(peersOrderView) != OrderComplete {
					orders[NodeID][floor][button] = getHighestOrderState(peersOrderView)
				}
			case OrderPending:
				orders[NodeID][floor][button] = getHighestOrderState(peersOrderView)
				if allPeersUpToDateOrAhead(peersOrderView, OrderPending, OrderConfirmed) {
					orders[NodeID][floor][button] = OrderConfirmed
				}
			case OrderConfirmed:
				orders[NodeID][floor][button] = getHighestOrderState(peersOrderView)

			case OrderComplete:
				if allPeersUpToDateOrAhead(peersOrderView, OrderComplete, OrderIdle) {
					orders[NodeID][floor][button] = OrderIdle
				}
			}
		}
	}
	return orders
}

func allPeersUpToDateOrAhead(peers []OrderState, myState OrderState, aheadState OrderState) bool {
	for _, p := range peers {
		if p != myState && p != aheadState {
			return false
		}
	}
	return true
}

func updatePeerStatusInMyWorldView(myWorldView WorldView, peerWorldView WorldView) WorldView {
	myWorldView.ElevatorStates[peerWorldView.NodeID] = peerWorldView.ElevatorStates[peerWorldView.NodeID]
	myWorldView.AliveList[peerWorldView.NodeID] = peerWorldView.AliveList[peerWorldView.NodeID]
	myWorldView.Orders[peerWorldView.NodeID] = peerWorldView.Orders[peerWorldView.NodeID]
	return myWorldView
}

func CabOrdersAsBool(Orders [config.NumFloors][config.NumOrderTypes]OrderState, nodeID int) [config.NumFloors]bool {
	var cabOrders [config.NumFloors]bool
	for floor := 0; floor < config.NumFloors; floor++ {
		cabOrders[floor] = (Orders[floor][2+nodeID] == OrderConfirmed) || (Orders[floor][2+nodeID] == OrderComplete)
	}
	return cabOrders
}

func HallOrdersAsBool(Orders [config.NumFloors][config.NumOrderTypes]OrderState) [config.NumFloors][2]bool {
	var hallOrders [config.NumFloors][2]bool
	for floor := 0; floor < config.NumFloors; floor++ {
		hallOrders[floor][0] = (Orders[floor][elevio.BT_HallUp] == OrderConfirmed) || (Orders[floor][elevio.BT_HallUp] == OrderComplete)
		hallOrders[floor][1] = (Orders[floor][elevio.BT_HallDown] == OrderConfirmed) || (Orders[floor][elevio.BT_HallDown] == OrderComplete)
	}
	return hallOrders
}

func setOrderLights(myWorldView WorldView, myNodeID int) {
	for floor := 0; floor < config.NumFloors; floor++ {
		for button := 0; button < config.NumElevatorButtons; button++ {
			var buttonValue int
			if button == elevio.BT_Cab {
				buttonValue = 2 + myNodeID
			} else {
				buttonValue = button
			}
			orderState := myWorldView.Orders[myNodeID][floor][buttonValue]
			buttonType := elevio.ButtonType(button)
			switch orderState {
			case OrderConfirmed:
				elevio.SetButtonLamp(buttonType, floor, true)
			case OrderIdle:
				elevio.SetButtonLamp(buttonType, floor, false)
			case OrderPending:
				continue
			case OrderComplete:
				continue
			}
		}
	}
}

func orderTypeFromEvent(order elevio.ButtonEvent, nodeID int) int {
	var orderType int
	if order.Button == elevio.BT_Cab {
		orderType = 2 + nodeID
	} else {
		orderType = int(order.Button)
	}

	return orderType
}

func alivePeersOrderView(orders [config.NumElevators][config.NumFloors][config.NumOrderTypes]OrderState, alivePeers [config.NumElevators]bool, floor int, button int) []OrderState {
	var alivePeersOrderView []OrderState
	for peerID, alive := range alivePeers {
		if alive {
			alivePeersOrderView = append(alivePeersOrderView, orders[peerID][floor][button])
		}
	}
	return alivePeersOrderView
}

func (wv *WorldView) tryPromoteIdleOrderToPending(orderFloor int, orderType int) {
	if wv.Orders[wv.NodeID][orderFloor][orderType] == OrderIdle {
		var peersOrderView = alivePeersOrderView(wv.Orders, wv.AliveList, orderFloor, orderType)
		if allPeersUpToDateOrAhead(peersOrderView, OrderIdle, OrderPending) {
			wv.Orders[wv.NodeID][orderFloor][orderType] = OrderPending
			return
		}
	}
}

func (wv *WorldView) tryMarkConfirmedOrderCompleted(orderFloor int, orderType int) {
	if wv.Orders[wv.NodeID][orderFloor][orderType] == OrderConfirmed {
		var peersOrderView = alivePeersOrderView(wv.Orders, wv.AliveList, orderFloor, orderType)
		if allPeersUpToDateOrAhead(peersOrderView, OrderConfirmed, OrderIdle) {
			wv.Orders[wv.NodeID][orderFloor][orderType] = OrderComplete
			return
		}
	}
}
