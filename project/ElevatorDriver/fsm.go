package ElevatorDriver

import (
	"fmt"
	"project/config"
	"project/elevio"
	"time"
)

func Elevator_fsm(
	newOrder <-chan Orders,
	updatedElevatorState chan<- Elevator,
	deliveredOrder chan<- elevio.ButtonEvent,
	elevator Elevator,
) {

	newFloorC := make(chan int, config.Buffer)
	openDoorC := make(chan bool, config.Buffer)
	doorObstructedc := make(chan bool, config.Buffer)
	doorClosingc := make(chan bool, config.Buffer)

	elevatorMotorTimer := time.NewTimer(config.ElevatorMotorTime)
	elevatorMotorTimer.Stop()

	go door_fsm(openDoorC, doorObstructedc, doorClosingc)

	go elevio.PollFloorSensor(newFloorC)

	var orders Orders

	ElevatorPrint(elevator)
	for {
		select {
		case floor := <-newFloorC:
			handleFloorArrival(&elevator, orders, openDoorC, deliveredOrder, elevatorMotorTimer, floor)
			elevator.MotorStop = false
			updatedElevatorState <- elevator

		case <-doorClosingc:
			handleDoorClosing(&elevator, orders, openDoorC, deliveredOrder, elevatorMotorTimer)
			updatedElevatorState <- elevator

		case orders = <-newOrder:
			handleNewOrder(&elevator, orders, openDoorC, deliveredOrder, elevatorMotorTimer, updatedElevatorState)
			

		case <-elevatorMotorTimer.C:
			fmt.Println("Gets motorstop")
			elevator.MotorStop = true
			updatedElevatorState <- elevator

		case obstrucion := <-doorObstructedc:
			if obstrucion != elevator.Obstruction {
				elevator.Obstruction = obstrucion
				updatedElevatorState <- elevator
			}

		}

	}
}
