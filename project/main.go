package main

import (
	"flag"
	"fmt"
	"os"
	"project/Network/bcast"
	"project/Network/localip"
	"project/Network/peers"
	"project/WorldView"
	"project/config"
	"project/elevio"
	"project/ElevatorDriver"
	"time"
)

var Port int
var id int

func main() {

	port := flag.Int("port", 15657, "UDP port to use for communication" "<--Default value, override with command line argument -port=xxxxx")
	ElevatorId := flag.Int("id", 0, "Unique ID for this elevator" "<--Default value, override with command line argument -id=xxxxx")

	flag.Parse()

	Port = *port
	id = *ElevatorId

	elevio.Init("localhost:"+strconv.Itoa(Port), config.NumFloors)
	
	fmt.Println("Elevator initialized with id", id, "on port", Port)
	fmt.Println("System has", config.NumFloors, "floors and", config.NumElevators, "elevators")
	
	peersRx 		:= make(chan peers.PeerUpdate, config.Buffer)
	peersTx 		:= make(chan bool, config.Buffer)
	
	networkTx 		:= make(chan WorldView.WorldView, config.Buffer)
	networkRx 		:= make(chan WorldView.WorldView, config.Buffer)
	newLocalElevatorState := make(chan ElevatorDriver.Elevator, config.Buffer)
	orderComplete 	:= make(chan elevio.ButtonEvent, config.Buffer)
	orderConfirmed 	:= make(chan elevio.ButtonEvent, config.Buffer)
	alivePeersInput := make(chan []int, config.Buffer)

	newOrder			 := make(chan ElevatorDriver.Orders, config.Buffer)
	

	go ElevatorDriver.Elevator_fsm(newOrder, newLocalElevatorState, orderComplete)

	go worldViewManager(
		networkRx, 
		networkTx, 
		newLocalElevatorState, 
		orderComplete, 
		orderConfirmed, 
		alivePeersInput, 
		myNodeID)


}




/*
import (
	"project/ElevatorDriver"
	"project/config"
	"project/elevio"
)

func main() {

	newOrder             := make(chan ElevatorDriver.Orders, config.Buffer)
	updatedElevatorState := make(chan ElevatorDriver.Elevator, config.Buffer)
	deliveredOrder       := make(chan elevio.ButtonEvent, config.Buffer)
	pollButtonC          := make(chan elevio.ButtonEvent, config.Buffer)

	elevio.Init("localhost:15657", 4)

	go elevio.PollButtons(pollButtonC)

	go handleOrder(newOrder, pollButtonC, deliveredOrder)

	go ElevatorDriver.Elevator_fsm(newOrder, updatedElevatorState, deliveredOrder)

	select {} // Kjøp continue forever
}

func handleOrder(newOrder chan<- ElevatorDriver.Orders, pollButtonC <-chan elevio.ButtonEvent, deliveredOrder <-chan elevio.ButtonEvent) {
	var orders ElevatorDriver.Orders

	for {
		select {
		case buttonEvent := <-pollButtonC:
			orders[buttonEvent.Floor][buttonEvent.Button] = true
			newOrder <- orders

		case delivered := <-deliveredOrder:
			orders[delivered.Floor][delivered.Button] = false
			newOrder <- orders
		}
	}
}
*/
