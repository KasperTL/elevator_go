package WorldView

import (
	"project/ElevatorDriver"
	"project/config"
)

type ElevatorState struct {
	Elevator  ElevatorDriver.Elevator
	CabOrders [config.NumFloors]bool
}

type OrderState int

const (
	OrderIdle      = 0
	OrderPending   = 1
	OrderConfirmed = 2
)

type WorldView struct {
	SenderID       int
	AliveList      [config.NumElevators]bool
	ElevatorStates [config.NumElevators]ElevatorState
	HallOrders     [config.NumElevators][config.NumFloors][2]OrderState
}

func InitWorldView(nodeID int) WorldView {
	view := WorldView{SenderID: nodeID}
	return view
}

func allPeersUpToDateOrAhead(peers []OrderState, stateA OrderState, stateB OrderState) bool {
	for _, p := range peers {
		if p != stateA || p != stateB {
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
	localOrders [config.NumElevators][config.NumFloors][2]OrderState,
	alivePeers []int,
) [config.NumElevators][config.NumFloors][2]OrderState {

	if len(alivePeers) == 0 {
		return localOrders
	}

	for elevator := 0; elevator < config.NumElevators; elevator++ {
		for floor := 0; floor < config.NumFloors; floor++ {
			for button := 0; button < 2; button++ {

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
func updateHallOrders(
	orders [config.NumElevators][config.NumFloors][2]OrderState,
	NodeID int,
	alivePeers []int,
) [config.NumElevators][config.NumFloors][2]OrderState {

	for floor := 0; floor < config.NumFloors; floor++ {
		for button := 0; button < 2; button++ {
			currentOrderState := orders[NodeID][floor][button]
			newOrderState := currentOrderState

			var peersOrderView []OrderState

			for peerID := range alivePeers {
				if peerID != NodeID {
					peersOrderView = append(peersOrderView, orders[peerID][floor][button])
				}
			}

			switch currentOrderState {
			case OrderIdle:
				if allPeersUpToDateOrAhead(peersOrderView, OrderIdle, OrderPending) && anyPeerAhead(peersOrderView, OrderPending) {
					newOrderState = OrderPending
				}
			case OrderPending:
				if allPeersUpToDateOrAhead(peersOrderView, OrderPending, OrderConfirmed) && anyPeerAhead(peersOrderView, OrderConfirmed) {
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
	myWorldView.HallOrders[peerWorldView.SenderID] = peerWorldView.HallOrders[peerWorldView.SenderID]

	return myWorldView
}

func FromOrderStateToBool(Orders [config.NumFloors][2]OrderState) [config.NumFloors][2]bool {
	var boolOrders [config.NumFloors][2]bool

	for floor := 0; floor < config.NumFloors; floor++ {
		for direction := 0; direction < 2; direction++ {
			boolOrders[floor][direction] = Orders[floor][direction] == OrderConfirmed
		}
	}
	return boolOrders
}
