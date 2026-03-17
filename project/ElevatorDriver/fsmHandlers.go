package ElevatorDriver

import (
	"fmt"
	"project/config"
	"project/elevio"
	"time"
)

func handleNewOrder(
	elevator *Elevator,
	orders Orders,
	openDoorC chan<- bool,
	deliveredOrder chan<- elevio.ButtonEvent,
	elevatorMotorTimer *time.Timer,
) {

	switch elevator.Behaviour {
	case EB_Idle:

		if orderAtFloorInDirection(elevator.Floor, elevator.Direction, orders) {
			stopElevator(elevatorMotorTimer)
			openDoor(elevator, openDoorC)
			clearOrderAtFloor(elevator.Floor, elevator.Direction, orders, deliveredOrder)
			return
		}
		if orderAtFloorOppositeDirection(elevator.Floor, elevator.Direction, orders) {
			reverseDirection(elevator)
			stopElevator(elevatorMotorTimer)
			openDoor(elevator, openDoorC)
			clearOrderAtFloor(elevator.Floor, elevator.Direction, orders, deliveredOrder)
			return
		}
		if orderInDirection(elevator.Floor, elevator.Direction, orders) {
			startMoving(elevator, elevatorMotorTimer)
			return
		}
		if orderInOppositeDirection(elevator.Floor, elevator.Direction, orders) {
			reverseDirection(elevator)
			startMoving(elevator, elevatorMotorTimer)
			return
		}

		enterIdle(elevator, elevatorMotorTimer)

	case EB_DoorOpen:

		if orderAtFloorInDirection(elevator.Floor, elevator.Direction, orders) {
			stopElevator(elevatorMotorTimer)
			openDoor(elevator, openDoorC)
			clearOrderAtFloor(elevator.Floor, elevator.Direction, orders, deliveredOrder)
			return
		}

	case EB_Moving:
		return

	default:
		fmt.Println("New order not handled")
	}

}

func handleFloorArrival(
	elevator *Elevator,
	orders Orders,
	openDoorC chan<- bool,
	deliveredOrder chan<- elevio.ButtonEvent,
	elevatorMotorTimer *time.Timer,
	floor int,
) {
	elevatorMotorTimer.Reset(config.ElevatorMotorTime)
	elevio.SetFloorIndicator(floor)
	elevator.Floor = floor

	if elevator.Obstruction {
		stopElevator(elevatorMotorTimer)
		openDoor(elevator, openDoorC)
		return
	}

	if orderAtFloorInDirection(elevator.Floor, elevator.Direction, orders) {
		if !orderInDirection(elevator.Floor, elevator.Direction, orders) && orderAtFloorOppositeDirection(elevator.Floor, elevator.Direction, orders) && !orderAtFloorInDirection(elevator.Floor, elevator.Direction, orders) {
			stopElevator(elevatorMotorTimer)
			openDoor(elevator, openDoorC)
			reverseDirection(elevator)
			clearOrderAtFloor(elevator.Floor, elevator.Direction, orders, deliveredOrder)
			//TODO: Should we return here?
			return
		}
		stopElevator(elevatorMotorTimer)
		openDoor(elevator, openDoorC)
		clearOrderAtFloor(elevator.Floor, elevator.Direction, orders, deliveredOrder)
		return
	}

	if orderInDirection(elevator.Floor, elevator.Direction, orders) {
		return
	}

	if orderAtFloorOppositeDirection(elevator.Floor, elevator.Direction, orders) {
		stopElevator(elevatorMotorTimer)
		openDoor(elevator, openDoorC)
		reverseDirection(elevator)
		clearOrderAtFloor(elevator.Floor, elevator.Direction, orders, deliveredOrder)
		return
	}

	if orderInOppositeDirection(elevator.Floor, elevator.Direction, orders) {
		reverseDirection(elevator)
		startMoving(elevator, elevatorMotorTimer)
		return
	}

	enterIdle(elevator, elevatorMotorTimer)
}

func handleDoorClosing(
	elevator *Elevator,
	orders Orders,
	openDoorC chan<- bool,
	deliveredOrder chan<- elevio.ButtonEvent,
	elevatorMotorTimer *time.Timer,
) {

	if orderInDirection(elevator.Floor, elevator.Direction, orders) {
		startMoving(elevator, elevatorMotorTimer)
		return
	}

	if orderAtFloorOppositeDirection(elevator.Floor, elevator.Direction, orders) {
		reverseDirection(elevator)
		stopElevator(elevatorMotorTimer)
		openDoor(elevator, openDoorC)
		clearOrderAtFloor(elevator.Floor, elevator.Direction, orders, deliveredOrder)
		return
	}

	if orderInOppositeDirection(elevator.Floor, elevator.Direction, orders) {
		reverseDirection(elevator)
		startMoving(elevator, elevatorMotorTimer)
		return
	}

}
