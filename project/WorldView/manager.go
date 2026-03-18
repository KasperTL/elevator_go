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
	localElevatorStateC <-chan ElevatorDriver.Elevator,
	servedOrderC <-chan elevio.ButtonEvent,
	assignerInputC chan<- WorldView,
	peersC <-chan peers.PeerUpdate,
	requestOrderCh <-chan elevio.ButtonEvent,
	myNodeID int,
) {

	myWorldView := InitWorldView(myNodeID)

	heartbeat := time.NewTicker(config.HeartbeatTime)

	cabOrdersRecovered := false
	var mode ConsensusMode

	for {
		select {
		case <-heartbeat.C:
			networkTx <- myWorldView
			assignerInputC <- myWorldView

		case peers := <-peersC:
			myWorldView.setAliveElevators(peers)
			mode = deriveConsensusMode(peers)
			myWorldView.stashLostNodesCabOrders(peers)

		case peerWorldView := <-networkRx:
			if !cabOrdersRecovered {
				for floor := range config.NumFloors {
					myWorldView.Orders[myNodeID][floor][myNodeID+2] = peerWorldView.CabOrderRecovery[myNodeID][floor]
				}
				cabOrdersRecovered = true
			}
			myWorldView.updatePeerStatusInMyWorldView(peerWorldView)
			myWorldView.updateOrders()
			setOrderLights(myWorldView, myNodeID)

		case myNewElevatorState := <- localElevatorStateC:
			myWorldView.ElevatorStates[myNodeID] = myNewElevatorState

		case requestEvent := <-requestOrderCh:
			orderType := orderTypeFromEvent(requestEvent, myNodeID)
			orderFloor := requestEvent.Floor

			switch mode {
			case Networked:
				myWorldView.tryPromoteIdleOrderToPending(orderFloor, orderType)

			case Standalone:
				myWorldView.Orders[myNodeID][orderFloor][orderType] = OrderConfirmed
				setOrderLights(myWorldView, myNodeID)
			}

		case servedOrderEvent := <-servedOrderC:
			orderType := orderTypeFromEvent(servedOrderEvent, myNodeID)
			orderFloor := servedOrderEvent.Floor

			switch mode {
			case Networked:
				myWorldView.tryMarkConfirmedOrderCompleted(orderFloor, orderType)

			case Standalone:
				myWorldView.Orders[myNodeID][orderFloor][orderType] = OrderIdle
				setOrderLights(myWorldView, myNodeID)
			}
		}
	}
}
