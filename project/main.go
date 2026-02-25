package main

import (
	"fmt"
	"project/ElevatorDriver"
	"project/config"
	"project/elevio"
)

func main() {
	//fmt.Println("Starting elevator control system...")

	//fmt.Println("DEBUG: Creating channels...")
	newOrder := make(chan ElevatorDriver.Orders, config.Buffer)
	updatedElevatorState := make(chan ElevatorDriver.Elevator, config.Buffer)
	deliveredOrder := make(chan elevio.ButtonEvent, config.Buffer)
	pollButtonC := make(chan elevio.ButtonEvent, config.Buffer)

	elevio.Init("localhost:15657", 4)

	//fmt.Println("DEBUG: Starting PollButtons...")
	go elevio.PollButtons(pollButtonC)

	//fmt.Println("DEBUG: Starting handleOrder...")
	go handleOrder(newOrder, pollButtonC, deliveredOrder)

	//fmt.Println("DEBUG: Starting Elevator_fsm...")
	go ElevatorDriver.Elevator_fsm(newOrder, updatedElevatorState, deliveredOrder)

	//fmt.Println("DEBUG: Main function done - waiting...")
	select {} // Kjøp continue forever
}

func handleOrder(newOrder chan<- ElevatorDriver.Orders, pollButtonC <-chan elevio.ButtonEvent, deliveredOrder <-chan elevio.ButtonEvent) {
	var orders ElevatorDriver.Orders

	for {
		select {
		case buttonEvent := <-pollButtonC:
			fmt.Println("Button pressed - Floor:", buttonEvent.Floor, "Button:", buttonEvent.Button)
			orders[buttonEvent.Floor][buttonEvent.Button] = true
			newOrder <- orders
			fmt.Println("Button event")

		case delivered := <-deliveredOrder:
			fmt.Println("Order delivered - Floor:", delivered.Floor, "Button:", delivered.Button)
			orders[delivered.Floor][delivered.Button] = false
			newOrder <- orders
			fmt.Println("Delivered order")
		}
	}
}
