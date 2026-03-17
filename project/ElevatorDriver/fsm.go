package ElevatorDriver

import (
	"fmt"
	"project/config"
	"project/elevio"
	"time"
)

func Elevator_fsm(
	newOrderC <-chan Orders,
	updatedElevatorStateC chan<- Elevator,
	deliveredOrderC chan<- elevio.ButtonEvent,
	elevator Elevator,
) {

	newFloorC := make(chan int, config.Buffer)
	openDoorC := make(chan bool, config.Buffer)
	doorObstructedc := make(chan bool, config.Buffer)
	doorClosingc := make(chan bool, config.Buffer)

	elevatorMotorTimer := time.NewTimer(config.ElevatorMotorTime)
	elevatorMotorTimer.Stop()

	var orders Orders

	go door_fsm(openDoorC, doorObstructedc, doorClosingc)
	go elevio.PollFloorSensor(newFloorC)

	ElevatorPrint(elevator)
	for {
		select {
		case floor := <-newFloorC:
			handleFloorArrival(&elevator, orders, openDoorC, deliveredOrderC, elevatorMotorTimer, floor)
			elevator.MotorStop = false
			updatedElevatorStateC <- elevator

		case <-doorClosingc:
			handleDoorClosing(&elevator, orders, openDoorC, deliveredOrderC, elevatorMotorTimer)
			updatedElevatorStateC <- elevator
		case orders = <-newOrderC:
			handleNewOrder(&elevator, orders, openDoorC, deliveredOrderC, elevatorMotorTimer)
			updatedElevatorStateC <- elevator

		case <-elevatorMotorTimer.C:
			fmt.Println("Gets motorstop")
			elevator.MotorStop = true
			updatedElevatorStateC <- elevator

		case obstrucion := <-doorObstructedc:

			if elevator.Behaviour == EB_Idle {
				openDoorC <- true
			}
			elevator.Obstruction = obstrucion
			updatedElevatorStateC <- elevator

		}
	}
}
