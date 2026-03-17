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
	myNodeID int,
) {

	myWorldView := InitWorldView(myNodeID)

	buttonEventCh := make(chan elevio.ButtonEvent, config.Buffer)
	go elevio.PollButtons(buttonEventCh)

	heartbeat := time.NewTicker(config.HeartbeatTime)

	initial := true
	online := false

	for {
		select {
		case <-heartbeat.C:
			networkTx <- myWorldView
			assignerInputC <- myWorldView

		case peers := <-peersC:
			myWorldView.setAliveElevators(peers)
			online = amIOnline(peers)
			myWorldView.stashLostNodesCabOrders(peers)

		case peerWorldView := <-networkRx:
			if initial {
				for floor := range config.NumFloors {
					myWorldView.Orders[myNodeID][floor][myNodeID+2] = peerWorldView.CabOrderRecovery[myNodeID][floor]
				}
				initial = false
			}
			myWorldView = updatePeerStatusInMyWorldView(myWorldView, peerWorldView)
			myWorldView.Orders = updateOrders(myWorldView.Orders, myNodeID, myWorldView.AliveList)
			setOrderLights(myWorldView, myNodeID)

		case myNewElevatorState := <- localElevatorStateC:
			myWorldView.ElevatorStates[myNodeID] = myNewElevatorState

		case buttonEvent := <-buttonEventCh:
			orderType := orderTypeFromEvent(buttonEvent, myNodeID)
			orderFloor := buttonEvent.Floor

			if online {
				myWorldView.tryPromoteIdleOrderToPending(orderFloor, orderType)

			} else {
				myWorldView.Orders[myNodeID][orderFloor][orderType] = OrderConfirmed
				setOrderLights(myWorldView, myNodeID)
			}

		case servedOrderEvent := <-servedOrderC:
			orderType := orderTypeFromEvent(servedOrderEvent, myNodeID)
			orderFloor := servedOrderEvent.Floor

			if online {
				myWorldView.tryMarkConfirmedOrderCompleted(orderFloor, orderType)
				setOrderLights(myWorldView, myNodeID)
			} else {
				myWorldView.Orders[myNodeID][orderFloor][orderType] = OrderIdle
				setOrderLights(myWorldView, myNodeID)
			}
		}
	}
}
