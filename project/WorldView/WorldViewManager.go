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

	var online bool
	var stable bool

	for {
		select {
		case <-heartbeat.C:
			networkTx <- myWorldView

		case peers = <-peersC:
	
			myWorldView.setAliveElevators(peers)
			online = amIOnline(peers)
			if peers.New != "" {
				myWorldView.setStatusToReconnect()
			}
			worldViewConfirmed <- myWorldView

		case peerWorldView := <-networkRx:

			myWorldView = updatePeerStatusInMyWorldView(myWorldView, peerWorldView)
			stable = myWorldView.areAllPeerWorlViewsNominal()

			if stable {

				myWorldView.Orders = updateOrders(myWorldView.Orders, myNodeID, peers.Peers)
				setOrderLights(myWorldView, myNodeID)
				worldViewConfirmed <- myWorldView

			} else {

				myWorldView.Orders = syncOnRejon(myWorldView.Orders, myNodeID, peerWorldView.SenderID)
				if isHallOrdersConsistent(myWorldView.Orders, myWorldView.AliveList) {
					myWorldView.Status[myNodeID] = nominal
				}
			}

		case myElevatorState := <-newLocalElevatorState:
			myWorldView.ElevatorStates[myNodeID] = myElevatorState
			worldViewConfirmed <- myWorldView

		case newOrder := <-orderRequest:

			var buttonValue int
			if newOrder.Button == elevio.BT_Cab {
				buttonValue = 2 + myNodeID
			} else {
				buttonValue = int(newOrder.Button)
			}

			switch online {
			case true:
				switch myWorldView.Orders[myNodeID][newOrder.Floor][buttonValue] {
				case OrderIdle:

					var peersOrderView []OrderState
					for _, peerIDstr := range peers.Peers {
						peerID, err := strconv.Atoi(peerIDstr)
						if err != nil {
							continue
						}
						if peerID != myNodeID {
							peersOrderView = append(peersOrderView, myWorldView.Orders[peerID][newOrder.Floor][buttonValue])
						}
					}

					if allPeersUpToDateOrAhead(peersOrderView, OrderIdle, OrderPending) {
						myWorldView.Orders[myNodeID][newOrder.Floor][buttonValue] = OrderPending
					} else {
						continue
					}

				case OrderPending:
					continue
				case OrderConfirmed:
					continue
				}

			case false:
				myWorldView.Orders[myNodeID][newOrder.Floor][buttonValue] = OrderConfirmed
				setOrderLights(myWorldView, myNodeID)
				worldViewConfirmed <- myWorldView

			}

		case completeOrder := <-orderComplete:
			var buttonValue int
			if completeOrder.Button == elevio.BT_Cab {
				buttonValue = 2 + myNodeID
			} else {
				buttonValue = int(completeOrder.Button)

			}

			switch online {
			case true:

				switch myWorldView.Orders[myNodeID][completeOrder.Floor][buttonValue] {
				case OrderConfirmed:

					var peersOrderView []OrderState
					for _, peerIDstr := range peers.Peers {
						peerID, err := strconv.Atoi(peerIDstr)
						if err != nil {
							continue
						}
						if peerID != myNodeID {
							peersOrderView = append(peersOrderView, myWorldView.Orders[peerID][completeOrder.Floor][buttonValue])
						}
					}

					if allPeersUpToDateOrAhead(peersOrderView, OrderConfirmed, OrderIdle) {
						myWorldView.Orders[myNodeID][completeOrder.Floor][buttonValue] = OrderIdle
						setOrderLights(myWorldView, myNodeID)
						worldViewConfirmed <- myWorldView

					} else {
						continue
					}

				case OrderPending:
					continue
				case OrderIdle:
					continue
				}
			case false:
				myWorldView.Orders[myNodeID][completeOrder.Floor][buttonValue] = OrderIdle
				setOrderLights(myWorldView, myNodeID)
				worldViewConfirmed <- myWorldView
			}
		}
	}
}
