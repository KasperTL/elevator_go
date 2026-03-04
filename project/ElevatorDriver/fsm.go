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
) {

	newFloorC := make(chan int)

	openDoorC := make(chan bool)
	doorObstructedc := make(chan bool)
	doorClosingc := make(chan bool)

	elevatorMotorTimer := time.NewTimer(config.ElevatorMotorTime)
	elevatorMotorTimer.Stop()

	go door_fsm(openDoorC, doorObstructedc, doorClosingc)

	go elevio.PollFloorSensor(newFloorC)

	var orders Orders

	elevator := Elevator{
		floor:       -1,
		direction:   ED_Down,
		behaviour:   EB_Moving,
		obstruction: false,
	}

	elevio.SetMotorDirection(elevio.MD_Down)

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
