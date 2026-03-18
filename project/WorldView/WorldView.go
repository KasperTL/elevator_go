package WorldView

import (
	"project/Network/peers"
	"project/config"
	"project/elevio"
	"strconv"
)

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

func (wv *WorldView) updatePeerStatusInMyWorldView(peerWorldView WorldView) {
	wv.ElevatorStates[peerWorldView.NodeID] = peerWorldView.ElevatorStates[peerWorldView.NodeID]
	wv.AliveList[peerWorldView.NodeID] = peerWorldView.AliveList[peerWorldView.NodeID]
	wv.Orders[peerWorldView.NodeID] = peerWorldView.Orders[peerWorldView.NodeID]
}

func (wv *WorldView) tryPromoteIdleOrderToPending(orderFloor int, orderType int) {
	if wv.Orders[wv.NodeID][orderFloor][orderType] == OrderIdle {
		var peersOrderView = collectPeerOrderStates(wv.Orders, wv.AliveList, orderFloor, orderType)
		if peersReadyToAdvance(peersOrderView, OrderIdle, OrderPending) {
			wv.Orders[wv.NodeID][orderFloor][orderType] = OrderPending
			return
		}
	}
}

func (wv *WorldView) tryMarkConfirmedOrderCompleted(orderFloor int, orderType int) {
	if wv.Orders[wv.NodeID][orderFloor][orderType] == OrderConfirmed {
		var peersOrderView = collectPeerOrderStates(wv.Orders, wv.AliveList, orderFloor, orderType)
		if peersReadyToAdvance(peersOrderView, OrderConfirmed, OrderComplete) {
			wv.Orders[wv.NodeID][orderFloor][orderType] = OrderComplete
			return
		}
	}
}

func (wv *WorldView) updateOrders() {

	for floor := 0; floor < config.NumFloors; floor++ {
		for button := 0; button < config.NumOrderTypes; button++ {
			peersOrderView := collectPeerOrderStates(wv.Orders, wv.AliveList, floor, button)

			currentOrderState := wv.Orders[wv.NodeID][floor][button]

			switch currentOrderState {
			case OrderIdle:
				if mostAdvancedOrderState(peersOrderView) != OrderComplete {
					wv.Orders[wv.NodeID][floor][button] = mostAdvancedOrderState(peersOrderView)
				}
			case OrderPending:
				wv.Orders[wv.NodeID][floor][button] = mostAdvancedOrderState(peersOrderView)
				if peersReadyToAdvance(peersOrderView, OrderPending, OrderConfirmed) {
					wv.Orders[wv.NodeID][floor][button] = OrderConfirmed
				}
			case OrderConfirmed:
				wv.Orders[wv.NodeID][floor][button] = mostAdvancedOrderState(peersOrderView)

			case OrderComplete:
				if peersReadyToAdvance(peersOrderView, OrderComplete, OrderIdle) {
					wv.Orders[wv.NodeID][floor][button] = OrderIdle
				}
			}
		}
	}
}

func peersReadyToAdvance(peers []OrderState, myState OrderState, aheadState OrderState) bool {
	for _, p := range peers {
		if p != myState && p != aheadState {
			return false
		}
	}
	return true
}

func mostAdvancedOrderState(peers []OrderState) OrderState {
	var highestOrderState OrderState
	for _, p := range peers {
		if p > highestOrderState {
			highestOrderState = p
		}
	}
	return highestOrderState
}

func deriveConsensusMode(peers peers.PeerUpdate) ConsensusMode {
	if len(peers.Peers) > 1 {
		return Networked
	} else {
		return Standalone
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

func collectPeerOrderStates(orders [config.NumElevators][config.NumFloors][config.NumOrderTypes]OrderState, alivePeers [config.NumElevators]bool, floor int, button int) []OrderState {
	var alivePeersOrderView []OrderState
	for peerID, alive := range alivePeers {
		if alive {
			alivePeersOrderView = append(alivePeersOrderView, orders[peerID][floor][button])
		}
	}
	return alivePeersOrderView
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
