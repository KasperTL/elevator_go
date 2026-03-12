package WorldView

import (
	"project/ElevatorDriver"
	"project/Network/peers"
	"project/config"
	"project/elevio"
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

	buttonEventCh := make(chan elevio.ButtonEvent, config.Buffer)
	go elevio.PollButtons(buttonEventCh)

	var peers peers.PeerUpdate

	heartbeat := time.NewTicker(config.HeartbeatTime)

	var online bool
	var stable bool

	for {
		select {
		case <-heartbeat.C:
			networkTx <- myWorldView
			worldViewConfirmed <- myWorldView

		case peers = <-peersC:
	
			myWorldView.setAliveElevators(peers)
			online = amIOnline(peers)
			if peers.New != "" {
				myWorldView.setStatusToReconnect()
			}

		case peerWorldView := <-networkRx:

			myWorldView = updatePeerStatusInMyWorldView(myWorldView, peerWorldView)
			stable = myWorldView.areAllPeerWorlViewsNominal()

			if stable {

				myWorldView.Orders = updateOrders(myWorldView.Orders, myNodeID, myWorldView.AliveList)
				setOrderLights(myWorldView, myNodeID)

			} else {

				myWorldView.Orders = syncOnRejon(myWorldView.Orders, myNodeID, peerWorldView.SenderID)
				if isHallOrdersConsistent(myWorldView.Orders, myWorldView.AliveList) {
					myWorldView.Status[myNodeID] = nominal
				}
			}

		case myElevatorState := <-newLocalElevatorState:
			myWorldView.ElevatorStates[myNodeID] = myElevatorState

		case buttonEvent := <-buttonEventCh:

			orderType := orderTypeFromButton(buttonEvent, myNodeID)
			orderFloor := buttonEvent.Floor

			if online {
				myWorldView.tryPromoteIdleOrderToPending(orderFloor, orderType)

			} else {
				myWorldView.Orders[myNodeID][orderFloor][orderType] = OrderConfirmed
				setOrderLights(myWorldView, myNodeID)

			}

		case completeOrder := <-orderComplete:
			
			orderType := orderTypeFromButton(completeOrder, myNodeID)
			orderFloor := completeOrder.Floor

			if online {

				myWorldView.tryMarkCofirmedOrderCompleted(orderFloor, orderType)
				setOrderLights(myWorldView, myNodeID)
			} else {
				myWorldView.Orders[myNodeID][orderFloor][orderType] = OrderIdle
				setOrderLights(myWorldView, myNodeID)
			}
		}
	}
}
