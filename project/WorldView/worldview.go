package WorldView

import (
	"project/Network/peers"
	"project/config"
	"project/elevio"
	"strconv"
)

func getAliveElevatorsFromPeers(peers peers.PeerUpdate) [config.NumElevators]bool {
	aliveList := [config.NumElevators]bool{}
	for _, peerIDstr := range peers.Peers {
		peerID, err := strconv.Atoi(peerIDstr)
		if err != nil {
			continue
		}
		aliveList[peerID] = true
	}
	return aliveList
}

func markLostIDOrdersIdle(orders [config.NumElevators][config.NumFloors][config.NumOrderTypes]OrderState, lostID int, myNodeID int) [config.NumElevators][config.NumFloors][config.NumOrderTypes]OrderState {
	for elevator := range config.NumElevators {
		for floor := range config.NumFloors {
			orders[elevator][floor][2+lostID] = OrderIdle
		}
	}
	return orders
}

func getCabOrdersFromNodeID(orders [config.NumElevators][config.NumFloors][config.NumOrderTypes]OrderState, nodeID int) [config.NumFloors]OrderState {
	var cabOrders [config.NumFloors]OrderState
	for floor := range config.NumFloors {
		cabOrders[floor] = orders[nodeID][floor][2+nodeID]
	}
	return cabOrders
}

func tryPromoteIdleOrderToPending(wv WorldView, orderFloor int, orderType int) OrderState {
	if wv.Orders[wv.NodeID][orderFloor][orderType] == OrderIdle {
		var peersOrderView = collectPeerOrderStates(wv.Orders, wv.AliveList, orderFloor, orderType)
		if peersReadyToAdvance(peersOrderView, OrderIdle, OrderPending) {
			return OrderPending
		}
	}
	return wv.Orders[wv.NodeID][orderFloor][orderType]
}

func tryMarkConfirmedOrderCompleted(wv WorldView, orderFloor int, orderType int) OrderState {
	if wv.Orders[wv.NodeID][orderFloor][orderType] == OrderConfirmed {
		var peersOrderView = collectPeerOrderStates(wv.Orders, wv.AliveList, orderFloor, orderType)
		if peersReadyToAdvance(peersOrderView, OrderConfirmed, OrderComplete) {
			return OrderComplete
		}
	}
	return wv.Orders[wv.NodeID][orderFloor][orderType]
}

func updatedOrders(wv WorldView) [config.NumElevators][config.NumFloors][config.NumOrderTypes]OrderState {
	updateOrders := wv.Orders
	for floor := 0; floor < config.NumFloors; floor++ {
		for button := 0; button < config.NumOrderTypes; button++ {
			peersOrderView := collectPeerOrderStates(wv.Orders, wv.AliveList, floor, button)

			currentOrderState := wv.Orders[wv.NodeID][floor][button]

			switch currentOrderState {
			case OrderIdle:
				if mostAdvancedOrderState(peersOrderView) != OrderComplete {
					updateOrders[wv.NodeID][floor][button] = mostAdvancedOrderState(peersOrderView)
				}
			case OrderPending:
				updateOrders[wv.NodeID][floor][button] = mostAdvancedOrderState(peersOrderView)
				if peersReadyToAdvance(peersOrderView, OrderPending, OrderConfirmed) {
					updateOrders[wv.NodeID][floor][button] = OrderConfirmed
				}
			case OrderConfirmed:
				updateOrders[wv.NodeID][floor][button] = mostAdvancedOrderState(peersOrderView)

			case OrderComplete:
				if peersReadyToAdvance(peersOrderView, OrderComplete, OrderIdle) {
					updateOrders[wv.NodeID][floor][button] = OrderIdle
				}
			}
		}
	}
	return updateOrders
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
