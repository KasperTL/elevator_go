package WorldView

import (
	"project/ElevatorDriver"
	"project/elevio"
)

func WorldViewManager(
	networkRx 			<-chan WorldView,
	networkTx			chan<- WorldView,
	newElevatorState	<-chan ElevatorDriver.Elevator,
	orderRequest		<-chan elevio.ButtonEvent,
	orderComplete 		chan<- elevio.ButtonEvent,
	orderConfirmed 		chan<- elevio.ButtonEvent,
	nodeID 				int,
) {

	localWorldView := InitWorldView(nodeID)

	for {
		select {
			case remoteView := <- networkRx:
			
				localWorldView.ElevatorStates[remoteView.SenderID] 	= remoteView.ElevatorStates[remoteView.SenderID]
				localWorldView.AliveList[remoteView.SenderID] 		= remoteView.AliveList[remoteView.SenderID]
				localWorldView.Orders[remoteView.SenderID] 			= remoteView.Orders[remoteView.SenderID]

				localWorldView.Orders								= updateHalOrders(localWorldView.Orders, nodeID)

			case newState := <- newElevatorState:
				localWorldView.ElevatorStates[nodeID] = newState

			case newOrder := <- orderRequest:

				switch localWorldView.Orders[nodeID][newOrder.Floor][newOrder.Button] {
				case OrderIdle:
					localWorldView.Orders[nodeID][newOrder.Floor][newOrder.Button] = OrderPending
				case OrderPending:
					continue
				case OrderConfirmed:
					//need to discuss this 
					continue
				}
			
		}
	}
}