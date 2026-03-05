package WorldView

import (
	"fmt"
	"project/ElevatorDriver"
	"project/config"
	"project/elevio"
)

func WorldViewManager(
	networkRx <-chan WorldView,
	networkTx chan<- WorldView,
	newLocalElevatorState <-chan ElevatorDriver.Elevator,
	orderComplete  <-chan elevio.ButtonEvent,
	orderConfirmed chan<- elevio.ButtonEvent,
	worldViewOut chan<- WorldView,
	alivePeersInput <-chan []int,
	myNodeID int,
) {

	myWorldView := InitWorldView(myNodeID)
	myWorldView.AliveList[myNodeID] = true

	orderRequest := make(chan elevio.ButtonEvent, config.Buffer)

	alivePeers := []int{}

	go elevio.PollButtons(orderRequest)

	for {
		select {
		case alivePeers = <-alivePeersInput:
			fmt.Println("alivePeers", alivePeers)
			//Reset all to False
			for i := range myWorldView.AliveList {
				myWorldView.AliveList[i] = false
			}

			//mark alive as true
			for _, peerID := range alivePeers {
				myWorldView.AliveList[peerID] = true
			}

			myWorldView.AliveList[myNodeID] = true

			myWorldView.HallOrders = syncOnRejon(myWorldView.HallOrders, alivePeers)
			networkTx <- myWorldView

		case peerWorldView := <-networkRx:
			myWorldView = updatePeerStatusInMyWorldView(myWorldView, peerWorldView)
			myWorldView.HallOrders = updateHallOrders(myWorldView.HallOrders, myNodeID, alivePeers)
			worldViewOut <- myWorldView
			networkTx <- myWorldView

		case myElevatorState := <-newLocalElevatorState:

			myWorldView.ElevatorStates[myNodeID].Elevator = myElevatorState
			networkTx <- myWorldView

		case newOrder := <-orderRequest:
			if newOrder.Button == elevio.BT_Cab {
				myWorldView.ElevatorStates[myNodeID].CabOrders[newOrder.Floor] = true
				elevio.SetButtonLamp(newOrder.Button, newOrder.Floor, true)
				worldViewOut <- myWorldView
				continue
			}
			
			switch myWorldView.HallOrders[myNodeID][newOrder.Floor][newOrder.Button] {
			case OrderIdle:
				
				var peersOrderView []OrderState
				for _, peerID := range alivePeers {
					if peerID != myNodeID {
						peersOrderView = append(peersOrderView, myWorldView.HallOrders[peerID][newOrder.Floor][newOrder.Button])
					}
				}
				// there may be some problems regarding the cab orders here
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
		case order := <- orderComplete:
			fmt.Println("Order complete:", order)
			if order.Button == elevio.BT_Cab {
				myWorldView.ElevatorStates[myNodeID].CabOrders[order.Floor] = false
				elevio.SetButtonLamp(order.Button, order.Floor, false)
				worldViewOut <- myWorldView
				continue
			}
			// NOT CORRECT, should be updated to reflect the new order states
			myWorldView.HallOrders[myNodeID][order.Floor][order.Button] = OrderIdle
			orderConfirmed <- order
			worldViewOut <- myWorldView

		}
		networkTx <- myWorldView
	}
}
