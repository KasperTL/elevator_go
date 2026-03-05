package main

import (
	"flag"
	"fmt"
	"project/Assigner"
	"project/ElevatorDriver"
	//"project/Network/bcast"
	"project/Network/peers"
	"project/WorldView"
	"project/config"
	"project/elevio"
	"strconv"
)

var Port int
var id int

func main() {

	port := flag.Int("port", 15657, "<--Default value, override with command line argument -port=xxxxx")
	ElevatorId := flag.Int("id", 0, "<--Default value, override with command line argument -id=xxxxx")

	flag.Parse()

	Port = *port
	id = *ElevatorId

	elevio.Init("localhost:"+strconv.Itoa(Port), config.NumFloors)


	fmt.Println("Elevator initialized with id", id, "on port", Port)
	fmt.Println("System has", config.NumFloors, "floors and", config.NumElevators, "elevators")

	peersRx := make(chan peers.PeerUpdate, config.Buffer)
	// peersTx := make(chan bool, config.Buffer)

	//starts peers moudle which will detect which leevators are alive on network
	// go peers.Transmitter(config.PeersPortNumber, strconv.Itoa(id), peersTx)
	// go peers.Receiver(config.PeersPortNumber, peersRx)

	//connect networktx and rx to actual UDP network for node talking

	networkTx := make(chan WorldView.WorldView, config.Buffer)
	networkRx := make(chan WorldView.WorldView, config.Buffer)

	// go bcast.Transmitter(Port, networkTx)
	// go bcast.Receiver(Port, networkRx)

	newLocalElevatorState := make(chan ElevatorDriver.Elevator, config.Buffer)
	orderComplete := make(chan elevio.ButtonEvent, config.Buffer)

	worldViewOut := make(chan WorldView.WorldView, config.Buffer)
	newOrder := make(chan ElevatorDriver.Orders, config.Buffer)


	go func() {
		for wv := range worldViewOut {
			orders := assigner.CalculateOptimalOrders(wv, id)
			newOrder <- orders
		}
	}()

	go ElevatorDriver.Elevator_fsm(
		newOrder,
		newLocalElevatorState,
		orderComplete)

	go WorldView.WorldViewManager(
		networkRx,
		networkTx,
		newLocalElevatorState,
		orderComplete,
		worldViewOut,
		peersRx,
		id)

	select {}

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
