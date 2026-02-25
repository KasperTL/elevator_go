package main

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
