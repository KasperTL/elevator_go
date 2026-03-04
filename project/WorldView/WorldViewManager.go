package WorldView

import (
	"project/ElevatorDriver"
	"project/config"
	"project/elevio"
)

func WorldViewManager(
	networkRx 			    <-chan WorldView,
	networkTx			    chan<- WorldView,
	newLocalElevatorState	<-chan ElevatorDriver.Elevator,
	orderComplete 		    chan<- elevio.ButtonEvent,
	orderConfirmed 		    chan<- elevio.ButtonEvent,
	alivePeersInput  	    <-chan []int,
	myNodeID 			    int,
) {

	myWorldView := InitWorldView(myNodeID)

	orderRequest := make(chan elevio.ButtonEvent, config.Buffer)
	
	alivePeers := []int{}

	go elevio.PollButtons(orderRequest)


	for {
		select {
			case alivePeers = <- alivePeersInput:

				myWorldView.HallOrders = syncOnRejon(myWorldView.HallOrders, alivePeers)

			case peerWorldView := <- networkRx:

				myWorldView        = updatePeerStatusInMyWorldView(myWorldView, peerWorldView)
				myWorldView.HallOrders = updateHallOrders(myWorldView.HallOrders, myNodeID, alivePeers)

			case myElevatorState := <- newLocalElevatorState:

				myWorldView.ElevatorStates[myNodeID].Elevator = myElevatorState

			case newOrder := <- orderRequest:

				switch myWorldView.HallOrders[myNodeID][newOrder.Floor][newOrder.Button] {
				case OrderIdle:
					var peersOrderView []OrderState

            		for peerID := range(alivePeers) {
                		if peerID != myNodeID {
                    		peersOrderView = append(peersOrderView, myWorldView.HallOrders[peerID][newOrder.Floor][newOrder.Button])
               			}
            		} 
					// there may be some problems regarding the cab orders here 
					if allPeersUpToDateOrAhead(peersOrderView, OrderIdle, OrderPending){
						myWorldView.HallOrders[myNodeID][newOrder.Floor][newOrder.Button] = OrderPending
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