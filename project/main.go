package main

import (
	"flag"
	"fmt"
	"project/Assigner"
	"project/ElevatorDriver"
	"project/Network/bcast"
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
	localElevator := ElevatorDriver.InitializeElevator()

	peersRx := make(chan peers.PeerUpdate, config.Buffer)
	peersTx := make(chan bool, config.Buffer)

	go peers.Transmitter(config.PeersPortNumber, strconv.Itoa(id), peersTx)
	go peers.Receiver(config.PeersPortNumber, peersRx)

	networkTx := make(chan WorldView.WorldView, config.Buffer)
	networkRx := make(chan WorldView.WorldView, config.Buffer)

	go bcast.Transmitter(config.BcastPortNumber, networkTx)
	go bcast.Receiver(config.BcastPortNumber, networkRx)

	newLocalElevatorState := make(chan ElevatorDriver.Elevator, config.Buffer)
	orderComplete := make(chan elevio.ButtonEvent, config.Buffer)
	newOrder := make(chan ElevatorDriver.Orders, config.Buffer)

	worldViewOut := make(chan WorldView.WorldView, config.Buffer)

	worldViewToAssigner := make(chan WorldView.WorldView, config.Buffer)

	go WorldView.WorldViewManager(
		networkRx,
		networkTx,
		newLocalElevatorState,
		orderComplete,
		worldViewOut,
		peersRx,
		id)

	go ElevatorDriver.Elevator_fsm(
		newOrder,
		newLocalElevatorState,
		orderComplete,
		localElevator)

	go Assigner.Assigner(
		worldViewToAssigner,
		newOrder,
		id)

	go func() {
		for w := range worldViewOut {
			fmt.Println("---- WorldView Update ----")
			fmt.Print("\033[H\033[2J")
			//PrintWorldView(w) // if you made the print function earlier
			worldViewToAssigner <- w
		}
	}()
	select {}

}

// func PrintWorldView(w WorldView.WorldView) {
// 	fmt.Println("------ WORLD VIEW ------")
// 	fmt.Printf("SenderID: %d\n", w.SenderID)

// 	fmt.Println("\nAlive List:")
// 	for i, alive := range w.AliveList {
// 		fmt.Printf("Elevator %d: %t\n", i, alive)
// 	}

// 	fmt.Println("\nStatus:")
// 	for i, s := range w.Status {
// 		fmt.Printf("Elevator %d: %d\n", i, s)
// 	}

// 	fmt.Println("\nElevator States:")
// 	for i, e := range w.ElevatorStates {
// 		fmt.Printf("Elevator %d -> Floor:%d Direction:%d Behaviour:%d Obstruction:%t MotorStop:%t\n",
// 			i, e.Floor, e.Direction, e.Behaviour, e.Obstruction, e.MotorStop)
// 	}

// 	fmt.Println("\nOrders:")
// 	for elev := 0; elev < config.NumElevators; elev++ {
// 		fmt.Printf("Elevator %d:\n", elev)
// 		for floor := 0; floor < config.NumFloors; floor++ {
// 			for btn := 0; btn < config.NumButtons; btn++ {
// 				order := w.Orders[elev][floor][btn]
// 				if order != WorldView.OrderIdle {
// 					fmt.Printf("  Floor %d Button %d -> %d\n", floor, btn, order)
// 				}
// 			}
// 		}
// 	}

// 	fmt.Println("------------------------")
// }
