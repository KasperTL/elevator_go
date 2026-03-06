package ElevatorDriver

import (
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

		case <-doorClosingc:
			handleDoorClosing(&elevator, orders, openDoorC, deliveredOrder, elevatorMotorTimer)

		case orders = <-newOrder:
			handleNewOrder(&elevator, orders, openDoorC, deliveredOrder, elevatorMotorTimer)

		case <-elevatorMotorTimer.C:

		case obstrucion := <-doorObstructedc:
			if obstrucion != elevator.obstruction {
				elevator.obstruction = obstrucion
				updatedElevatorState <- elevator
			}

		}
		updatedElevatorState <- elevator
	}
}
