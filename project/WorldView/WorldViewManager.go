package WorldView

import (
	"fmt"
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
	myWorldView.AliveList[myNodeID] = true

	orderRequest := make(chan elevio.ButtonEvent, config.Buffer)
	go elevio.PollButtons(orderRequest)

	var peers peers.PeerUpdate

	heartbeat := time.NewTicker(config.HeartbeatTime)

	for {
		select {
		case <-heartbeat.C:
			networkTx <- myWorldView

		case peers = <-peersC:
			fmt.Println(peers)
			worldViewConfirmed <- myWorldView

		case peerWorldView := <-networkRx:
			if peerWorldView.SenderID == myNodeID {
				continue
			}
			// print functions
			fmt.Printf("=== Received from Elevator %d ===\n", peerWorldView.SenderID)
			PrintHallOrders(peerWorldView)

			myWorldView = updatePeerStatusInMyWorldView(myWorldView, peerWorldView)
			myWorldView.HallOrders = updateHallOrders(myWorldView.HallOrders, myNodeID, peers.Peers)
			setHallOrderLights(myWorldView)
			worldViewConfirmed <- myWorldView

		case myElevatorState := <-newLocalElevatorState:

			myWorldView.ElevatorStates[myNodeID].Elevator = myElevatorState
			worldViewConfirmed <- myWorldView

		case newOrder := <-orderRequest:

			if newOrder.Button == elevio.BT_Cab {
				myWorldView.ElevatorStates[myNodeID].CabOrders[newOrder.Floor] = true
				elevio.SetButtonLamp(elevio.BT_Cab, newOrder.Floor, true)
				worldViewConfirmed <- myWorldView
				continue
			}

			switch myWorldView.HallOrders[myNodeID][newOrder.Floor][newOrder.Button] {
			case OrderIdle:

				var peersOrderView []OrderState
				for _, peerIDstr := range peers.Peers {
					peerID, err := strconv.Atoi(peerIDstr)
					if err != nil {
						continue
					}
					if peerID != myNodeID {
						peersOrderView = append(peersOrderView, myWorldView.HallOrders[peerID][newOrder.Floor][newOrder.Button])
					}
				}

				if allPeersUpToDateOrAhead(peersOrderView, OrderIdle, OrderPending) {
					myWorldView.HallOrders[myNodeID][newOrder.Floor][newOrder.Button] = OrderPending
				} else {
					continue
				}

			case OrderPending:
				continue
			case OrderConfirmed:
				continue
			}

		case comleteOrder := <-orderComplete:

			if comleteOrder.Button == elevio.BT_Cab {
				myWorldView.ElevatorStates[myNodeID].CabOrders[comleteOrder.Floor] = false
				elevio.SetButtonLamp(comleteOrder.Button, comleteOrder.Floor, false)
				worldViewConfirmed <- myWorldView
				continue
			}

			switch myWorldView.HallOrders[myNodeID][comleteOrder.Floor][comleteOrder.Button] {
			case OrderConfirmed:

				var peersOrderView []OrderState
				for _, peerIDstr := range peers.Peers {
					peerID, err := strconv.Atoi(peerIDstr)
					if err != nil {
						continue
					}
					if peerID != myNodeID {
						peersOrderView = append(peersOrderView, myWorldView.HallOrders[peerID][comleteOrder.Floor][comleteOrder.Button])
					}
				}

				if allPeersUpToDateOrAhead(peersOrderView, OrderConfirmed, OrderIdle) {
					myWorldView.HallOrders[myNodeID][comleteOrder.Floor][comleteOrder.Button] = OrderIdle
					setHallOrderLights(myWorldView)

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
