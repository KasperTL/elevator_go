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

	orderRequest := make(chan elevio.ButtonEvent, config.Buffer)
	go elevio.PollButtons(orderRequest)

	var peers peers.PeerUpdate

	heartbeat := time.NewTicker(config.HeartbeatTime)

	for {
		select {
		case <-heartbeat.C:
			networkTx <- myWorldView

		case peers = <-peersC:
			fmt.Println("HEI")
			myWorldView.AliveList = [config.NumElevators]bool{}
			for _, peerIDstr := range peers.Peers {
				peerID, err := strconv.Atoi(peerIDstr)
				if err != nil {
					continue
				}
				myWorldView.AliveList[peerID] = true
			}

			// rejoin logic
			if peers.New != "" {
				for i := range myWorldView.StableWorldview {
					myWorldView.StableWorldview[i] = false
				}
			}

			fmt.Println("PeersC: ", peers)
			// need to reset my view of the lost elevators.

			worldViewConfirmed <- myWorldView

		case peerWorldView := <-networkRx:

			myWorldView = updatePeerStatusInMyWorldView(myWorldView, peerWorldView)

			stable := true
			for i := range config.NumElevators {
				if myWorldView.AliveList[i] {
					if !myWorldView.StableWorldview[i] {
						stable = false
					}

				}
			}
			fmt.Println("stable: ", stable)
			if stable {

				myWorldView.Orders = updateOrders(myWorldView.Orders, myNodeID, peers.Peers)
				setOrderLights(myWorldView)
				worldViewConfirmed <- myWorldView

			} else {

				myWorldView.Orders = syncOnRejon(myWorldView.Orders, myNodeID, peerWorldView.SenderID)

				if isWorldViewStable(myWorldView.Orders, myWorldView.AliveList) {
					myWorldView.StableWorldview[myNodeID] = true
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

		case completeOrder := <-orderComplete:
			var buttonValue int
			if completeOrder.Button == elevio.BT_Cab {
				buttonValue = 2 + myNodeID
			} else {
				buttonValue = int(completeOrder.Button)

			}

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

		}
	}
}
