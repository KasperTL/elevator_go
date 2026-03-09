package WorldView

import (
	"project/ElevatorDriver"
	"project/Network/peers"
	"project/config"
	"project/elevio"
	"strconv"
	"time"
)

func WorldViewManager(
	networkRx <-chan WorldView,
	networkTx chan<- WorldView,
	newLocalElevatorState <-chan ElevatorDriver.Elevator,
	orderComplete <-chan elevio.ButtonEvent,
	worldViewConfirmed chan<- WorldView,
	peersC <-chan peers.PeerUpdate,
	myNodeID int,
) {

	myWorldView := InitWorldView(myNodeID)

	orderRequest := make(chan elevio.ButtonEvent, config.Buffer)
	go elevio.PollButtons(orderRequest)

	var peers peers.PeerUpdate

	heartbeat := time.NewTicker(config.HeartbeatTime)

	for {
		select {
		case <-heartbeat.C:
			networkTx <- myWorldView

		case peers = <-peersC:
			for _, peerIDstr := range peers.Peers {
				peerID, err := strconv.Atoi(peerIDstr)
				if err != nil {
					continue
				}
				myWorldView.AliveList[peerID] = true
			}

			worldViewConfirmed <- myWorldView

		case peerWorldView := <-networkRx:
			if peerWorldView.SenderID == myNodeID {
				continue
			}

			myWorldView = updatePeerStatusInMyWorldView(myWorldView, peerWorldView) // this does not work, or get wrong in the assigner
			myWorldView.Orders = updateOrders(myWorldView.Orders, myNodeID, peers.Peers)
			setOrderLights(myWorldView)
			worldViewConfirmed <- myWorldView

		case myElevatorState := <-newLocalElevatorState:
			myWorldView.ElevatorStates[myNodeID] = myElevatorState
			worldViewConfirmed <- myWorldView

		case newOrder := <-orderRequest:

			switch myWorldView.Orders[myNodeID][newOrder.Floor][newOrder.Button] {
			case OrderIdle:

				var peersOrderView []OrderState
				for _, peerIDstr := range peers.Peers {
					peerID, err := strconv.Atoi(peerIDstr)
					if err != nil {
						continue
					}
					if peerID != myNodeID {
						peersOrderView = append(peersOrderView, myWorldView.Orders[peerID][newOrder.Floor][newOrder.Button])
					}
				}

				if allPeersUpToDateOrAhead(peersOrderView, OrderIdle, OrderPending) {
					myWorldView.Orders[myNodeID][newOrder.Floor][newOrder.Button] = OrderPending
				} else {
					continue
				}

			case OrderPending:
				continue
			case OrderConfirmed:
				continue
			}

		case comleteOrder := <-orderComplete:

			switch myWorldView.Orders[myNodeID][comleteOrder.Floor][comleteOrder.Button] {
			case OrderConfirmed:

				var peersOrderView []OrderState
				for _, peerIDstr := range peers.Peers {
					peerID, err := strconv.Atoi(peerIDstr)
					if err != nil {
						continue
					}
					if peerID != myNodeID {
						peersOrderView = append(peersOrderView, myWorldView.Orders[peerID][comleteOrder.Floor][comleteOrder.Button])
					}
				}

				if allPeersUpToDateOrAhead(peersOrderView, OrderConfirmed, OrderIdle) {
					myWorldView.Orders[myNodeID][comleteOrder.Floor][comleteOrder.Button] = OrderIdle
					setOrderLights(myWorldView)
					worldViewConfirmed <- myWorldView

				} else {
					continue
				}

			case OrderPending:
				continue
			case OrderIdle:
				continue
			}

		}
	}
}
