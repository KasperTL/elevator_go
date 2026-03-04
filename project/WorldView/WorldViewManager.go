package WorldView

import (
	"project/ElevatorDriver"
	"project/config"
	"project/elevio"
)

func WorldViewManager(
	networkRx 			<-chan WorldView,
	networkTx			chan<- WorldView,
	newLocalElevatorState	<-chan ElevatorDriver.Elevator,
	orderRequest		<-chan elevio.ButtonEvent,
	orderComplete 		chan<- elevio.ButtonEvent,
	orderConfirmed 		chan<- elevio.ButtonEvent,
	alivePeersInput  	<-chan []int,
	myNodeID 			int,
) {

	myWorldView := InitWorldView(myNodeID)

	for {
		select {
			case alivePeers := <- alivePeersInput:

				myWorldView.Orders = syncOnRejon(myWorldView.Orders, alivePeers)

			case peerWorldView := <- networkRx:

				myWorldView        = updatePeerStatusInMyWorldView(myWorldView, peerWorldView)
				myWorldView.Orders = updateHalOrders(myWorldView.Orders, myNodeID, alivePeers)

			case myElevatorState := <- newLocalElevatorState:

				myWorldView.ElevatorStates[myNodeID] = myElevatorState

			case newOrder := <- orderRequest:

				switch myWorldView.Orders[myNodeID][newOrder.Floor][newOrder.Button] {
				case OrderIdle:
					var peersOrderView []OrderState

            		for peerID := range(alivePeers) {
                		if peerID != myNodeID {
                    		peersOrderView = append(peersOrderView, myWorldView.Orders[peerID][newOrder.Floor][newOrder.Button])
               			}
            		} 
					// there may be some problems regarding the cab orders here 
					if allPeersUpToDateOrAhead(peersOrderView, OrderIdle, OrderPending){
						myWorldView.Orders[myNodeID][newOrder.Floor][newOrder.Button] = OrderPending
					} else {
						continue
					}

				case OrderPending:
					continue
				case OrderConfirmed: 
					continue
				}
			
		}
	}
}