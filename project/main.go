package main

import (
	"flag"
	"fmt"
	"project/ElevatorDriver"
	"project/Network/bcast"
	"project/Network/peers"
	"project/OrderDispatcher"
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
	localElevator := ElevatorDriver.InitializeElevator()

	fmt.Println("Elevator initialized with id", id, "on port", Port)
	fmt.Println("System has", config.NumFloors, "floors and", config.NumElevators, "elevators")

	peersRx := make(chan peers.PeerUpdate, config.Buffer)
	peersTx := make(chan bool, config.Buffer)
	go peers.Transmitter(config.PeersPortNumber, strconv.Itoa(id), peersTx)
	go peers.Receiver(config.PeersPortNumber, peersRx)

	networkTx := make(chan WorldView.WorldView, config.Buffer)
	networkRx := make(chan WorldView.WorldView, config.Buffer)
	go bcast.Transmitter(config.BcastPortNumber, networkTx)
	go bcast.Receiver(config.BcastPortNumber, networkRx)

	requestOrderC := make(chan elevio.ButtonEvent, config.Buffer)
	go elevio.PollButtons(requestOrderC)

	localElevatorStateC := make(chan ElevatorDriver.Elevator, config.Buffer)
	servedOrderC := make(chan elevio.ButtonEvent, config.Buffer)
	newOrderC := make(chan ElevatorDriver.Orders, config.Buffer)
	assignerInputC := make(chan WorldView.WorldView, config.Buffer)

	go WorldView.WorldViewManager(
		networkRx,
		networkTx,
		localElevatorStateC,
		servedOrderC,
		assignerInputC,
		peersRx,
		requestOrderC,
		id)

	go ElevatorDriver.ElevatorController(
		newOrderC,
		localElevatorStateC,
		servedOrderC,
		localElevator)

	go OrderDispatcher.AssignLocalOrdersFromWorldView(
		assignerInputC,
		newOrderC,
		id)

	select {}
}
