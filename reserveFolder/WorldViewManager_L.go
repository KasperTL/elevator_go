// package WorldView

// import (
// 	"project/ElevatorDriver"
// 	"project/elevio"
// )

// func WorldViewManager(
// 	networkRx <-chan WorldView,
// 	networkTx chan<- WorldView,
// 	newLocalElevatorState <-chan ElevatorDriver.Elevator,
// 	orderRequest <-chan elevio.ButtonEvent,
// 	orderComplete <-chan elevio.ButtonEvent,
// 	confirmedOrdersOut chan<- WorldView,
// 	//orderConfirmed chan<- elevio.ButtonEvent,
// 	alivePeersInput <-chan []int,
// 	myNodeID int,
// ) {

// 	myWorldView := InitWorldView(myNodeID)
// 	//alivePeers := make([]int, 0)

// 	// peerViews stores the last received worldview from each peer.
// 	// These are acknowledgent and can act as consensus list
// 	peerViews := make(map[int]WorldView)
// 	alivePeers := []int{}

// 	// broadcast is a helper to send our current worldview out on the network
// 	broadcast := func() {
// 		myWorldView.Epoch++
// 		networkTx <- myWorldView
// 	}

// 	// confirmAndNotify checks for newly confirmable orders and notifies the assigner
// 	confirmAndNotify := func() {
// 		checkConfirmation(&myWorldView, peerViews, alivePeers)
// 		confirmedOrdersOut <- myWorldView
// 	}

// 	//taken out since it is in peers.go
// 	//heartbeat := time.NewTicker(config.HeartbeatTime)
// 	//defer heartbeat.Stop()

// 	for {
// 		select {

// 		//incoming worldview from peer,merge using cyclic counter
// 		case peerWorldView := <-networkRx:
// 			peerViews[peerWorldView.SenderID] = peerWorldView
// 			MergeWorldView(&myWorldView, peerWorldView)
// 			confirmAndNotify()
// 			broadcast()

// 		//Alive peer list updated
// 		case peers := <-alivePeersInput:
// 			alivePeers = peers
// 			confirmAndNotify()
// 			broadcast()

// 		case myElevatorState := <-newLocalElevatorState:
// 			myWorldView.ElevatorStates[myNodeID] = ElevatorState{
// 				Floor:     myElevatorState.GetFloor(),
// 				Direction: myElevatorState.GetDirection(),
// 				Behaviour: myElevatorState.GetBehaviour(),
// 			}
// 			broadcast()

// 		case newOrder := <-orderRequest:
// 			OnButtonPress(&myWorldView, newOrder)
// 			confirmAndNotify()
// 			broadcast()

// 		case completed := <-orderComplete:
// 			OnOrderComplete(&myWorldView, completed)
// 			broadcast()
// 		}
// 	}
// }
