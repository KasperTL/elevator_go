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
)

type WorldViewStatus int

const (
	reconnectState = 0
	nominal        = 1
)

type WorldView struct {
	SenderID       int
	AliveList      [config.NumElevators]bool
	Status         [config.NumElevators]WorldViewStatus
	ElevatorStates [config.NumElevators]ElevatorDriver.Elevator
	Orders         [config.NumElevators][config.NumFloors][config.NumOrderTypes]OrderState
}

func InitWorldView(nodeID int) WorldView {
	view := WorldView{SenderID: nodeID}
	view.AliveList[nodeID] = true
	return view
}

func (wv *WorldView) setAliveElevators(peers peers.PeerUpdate) {
	wv.AliveList = [config.NumElevators]bool{}
	wv.AliveList[wv.SenderID] = true
	for _, peerIDstr := range peers.Peers {
		peerID, err := strconv.Atoi(peerIDstr)
		if err != nil {
			continue
		}
		wv.AliveList[peerID] = true
	}
}

func (wv *WorldView) setStatusToReconnect() {
	for i := range wv.Status {
		wv.Status[i] = reconnectState
	}
}

func amIOnline(peers peers.PeerUpdate) bool {
	var online bool
	if len(peers.Peers) > 0 {
		online = true
	} else {
		online = false
	}
	return online
}

func (wv *WorldView) areAllPeerWorlViewsNominal() bool {
	var stable bool
	stable = true
	for i := range config.NumElevators {
		if wv.AliveList[i] {
			if wv.Status[i] == reconnectState {
				stable = false
			}
		}
	}
	return stable
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
	default:
		return false
	}
}

func syncOnRejon(
	localOrders [config.NumElevators][config.NumFloors][config.NumOrderTypes]OrderState,
	myNodeID int,
	peerID int,
) [config.NumElevators][config.NumFloors][config.NumOrderTypes]OrderState {

	for floor := 0; floor < config.NumFloors; floor++ {
		for button := 0; button < config.NumOrderTypes; button++ {
			if isPeerAhead(localOrders[peerID][floor][button], localOrders[myNodeID][floor][button]) {
				localOrders[myNodeID][floor][button] = localOrders[peerID][floor][button]
			}
		}
	}
	return localOrders
}

func isHallOrdersConsistent(
	orders [config.NumElevators][config.NumFloors][config.NumOrderTypes]OrderState,
	alivePeers [config.NumElevators]bool,
) bool {

	for floor := 0; floor < config.NumFloors; floor++ {
		for button := 0; button < config.NumOrderTypes; button++ {

			peersOrderView := alivePeersOrderView(orders,alivePeers, floor, button)

			if len(peersOrderView) == 0 {
				continue
			}
			for _, state := range peersOrderView[1:] {
				if state != peersOrderView[0] {
					return false
				}
			}
		}
	}
	return true
}

func updateOrders(
	orders [config.NumElevators][config.NumFloors][config.NumOrderTypes]OrderState,
	NodeID int,
	alivePeers [config.NumElevators]bool,
) [config.NumElevators][config.NumFloors][config.NumOrderTypes]OrderState {

	for floor := 0; floor < config.NumFloors; floor++ {
		for button := 0; button < config.NumOrderTypes; button++ {
			currentOrderState := orders[NodeID][floor][button]
			newOrderState := currentOrderState

			peersOrderView := alivePeersOrderView(orders, alivePeers, floor, button)

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
	myWorldView.Status[peerWorldView.SenderID] = peerWorldView.Status[peerWorldView.SenderID]
	return myWorldView
}

func CabOrdersAsBool(Orders [config.NumFloors][config.NumOrderTypes]OrderState, nodeID int) [config.NumFloors]bool {
	var cabOrders [config.NumFloors]bool
	for floor := 0; floor < config.NumFloors; floor++ {
		cabOrders[floor] = Orders[floor][2+nodeID] == OrderConfirmed
	}
	return cabOrders
}

func HallOrdersAsBool(Orders [config.NumFloors][config.NumOrderTypes]OrderState) [config.NumFloors][2]bool {
	var hallOrders [config.NumFloors][2]bool
	for floor := 0; floor < config.NumFloors; floor++ {
		hallOrders[floor][0] = Orders[floor][elevio.BT_HallUp] == OrderConfirmed
		hallOrders[floor][1] = Orders[floor][elevio.BT_HallDown] == OrderConfirmed
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
			}
		}
	}
}

func orderTypeFromButton(order elevio.ButtonEvent, nodeID int) int {
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
	switch wv.Orders[wv.SenderID][orderFloor][orderType] {
		case OrderIdle:
			var peersOrderView = alivePeersOrderView(wv.Orders, wv.AliveList, orderFloor, orderType)
			if allPeersUpToDateOrAhead(peersOrderView, OrderIdle, OrderPending) {
				wv.Orders[wv.SenderID][orderFloor][orderType] = OrderPending
				return 
			} 
		case OrderPending:
			return
		case OrderConfirmed:
			return 
		}
}

func (wv *WorldView) tryMarkCofirmedOrderCompleted(orderFloor int, orderType int) {
	switch wv.Orders[wv.SenderID][orderFloor][orderType] {
		case OrderConfirmed:
			var peersOrderView = alivePeersOrderView(wv.Orders, wv.AliveList, orderFloor, orderType)

			if allPeersUpToDateOrAhead(peersOrderView, OrderConfirmed, OrderIdle) {
				wv.Orders[wv.SenderID][orderFloor][orderType] = OrderIdle
				return 
			} 
		case OrderPending:
			return 
		case OrderIdle:
			return 
		}
}