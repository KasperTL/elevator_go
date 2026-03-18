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

			myWorldView.AliveList = getAliveElevatorsFromPeers(peers)
			myWorldView.AliveList[myNodeID] = true
			for _, lostIDstr := range peers.Lost {
				lostID, err := strconv.Atoi(lostIDstr)
				if err != nil {
					continue
				}
				if lostID == myNodeID {
					continue
				}
				myWorldView.CabOrderRecovery[lostID] = getCabOrdersFromNodeID(myWorldView.Orders, lostID)
				myWorldView.Orders = markLostIDOrdersIdle(myWorldView.Orders, lostID, myNodeID)

				println("Order of lostID:", myWorldView.Orders[myNodeID][0][2+lostID])
			}
			mode = deriveConsensusMode(peers)

		case peerWorldView := <-networkRx:
			if mode == Standalone {
				continue
			}
			if !cabOrdersRecovered {
				for floor := range config.NumFloors {
					myWorldView.Orders[myNodeID][floor][myNodeID+2] = peerWorldView.CabOrderRecovery[myNodeID][floor]
				}
				cabOrdersRecovered = true
			}
			myWorldView.ElevatorStates[peerWorldView.NodeID] = peerWorldView.ElevatorStates[peerWorldView.NodeID]
			myWorldView.AliveList[peerWorldView.NodeID] = peerWorldView.AliveList[peerWorldView.NodeID]
			myWorldView.Orders[peerWorldView.NodeID] = peerWorldView.Orders[peerWorldView.NodeID]

			myWorldView.Orders = updatedOrders(myWorldView)
			setOrderLights(myWorldView, myNodeID)

		case myNewElevatorState := <-localElevatorStateC:
			myWorldView.ElevatorStates[myNodeID] = myNewElevatorState

		case requestEvent := <-requestOrderCh:
			orderType := orderTypeFromEvent(requestEvent, myNodeID)
			orderFloor := requestEvent.Floor

			switch mode {
			case Networked:
				myWorldView.Orders[myNodeID][orderFloor][orderType] = tryPromoteIdleOrderToPending(myWorldView, orderFloor, orderType)

			case Standalone:
				myWorldView.Orders[myNodeID][orderFloor][orderType] = OrderConfirmed
				setOrderLights(myWorldView, myNodeID)
			}

		case servedOrderEvent := <-servedOrderC:
			orderType := orderTypeFromEvent(servedOrderEvent, myNodeID)
			orderFloor := servedOrderEvent.Floor

			switch mode {
			case Networked:
				myWorldView.Orders[myNodeID][orderFloor][orderType] = tryMarkConfirmedOrderCompleted(myWorldView, orderFloor, orderType)

			case Standalone:
				myWorldView.Orders[myNodeID][orderFloor][orderType] = OrderIdle
				setOrderLights(myWorldView, myNodeID)
			}
		}
	}
}
