package ElevatorDriver

import (
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

	switch elevator.behaviour {

	case EB_Idle:

		if orderAtFloorInDirection(elevator.floor, elevator.direction, orders) || cabOrderAtFloor(elevator.floor, orders) {
			stopAndOpenDoor(elevator, openDoorC)
			clearOrderAtFloor(elevator.floor, elevator.direction, orders, deliveredOrder)
			return
		}
		if orderAtFloorOppositeDirection(elevator.floor, elevator.direction, orders) {
			reverseDirection(elevator)
			stopAndOpenDoor(elevator, openDoorC)
			clearOrderAtFloor(elevator.floor, elevator.direction, orders, deliveredOrder)
			return
		}
		if orderInDirection(elevator.floor, elevator.direction, orders) {
			startMoving(elevator, elevatorMotorTimer)
			return
		}
		if orderInOppositeDirection(elevator.floor, elevator.direction, orders) {
			reverseDirection(elevator)
			startMoving(elevator, elevatorMotorTimer)
			return
		}

		enterIdle(elevator)

	case EB_DoorOpen:

		if orderAtFloorInDirection(elevator.floor, elevator.direction, orders) || cabOrderAtFloor(elevator.floor, orders) {
			stopAndOpenDoor(elevator, openDoorC)
			clearOrderAtFloor(elevator.floor, elevator.direction, orders, deliveredOrder)
			return
		}

	case EB_Moving:
		//Handled at floor arrival

	default:
		panic("New order not handled")
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
	elevio.SetFloorIndicator(floor)
	elevator.floor = floor

	if elevator.behaviour != EB_Moving {
		//panic("Floor arrival in non-moving state")
		return
	}

	if orderAtFloorInDirection(elevator.floor, elevator.direction, orders) || cabOrderAtFloor(elevator.floor, orders) {
		// TODO
		// This is toooooooo long, need to fix
		if !orderInDirection(elevator.floor, elevator.direction, orders) && orderAtFloorOppositeDirection(elevator.floor, elevator.direction, orders) && !orderAtFloorInDirection(elevator.floor, elevator.direction, orders) {
			stopAndOpenDoor(elevator, openDoorC)
			reverseDirection(elevator)
			clearOrderAtFloor(elevator.floor, elevator.direction, orders, deliveredOrder)
		}
		stopAndOpenDoor(elevator, openDoorC)
		clearOrderAtFloor(elevator.floor, elevator.direction, orders, deliveredOrder)
		return
	}

	if orderInDirection(elevator.floor, elevator.direction, orders) {
		return
	}

	if orderAtFloorOppositeDirection(elevator.floor, elevator.direction, orders) {
		stopAndOpenDoor(elevator, openDoorC)
		reverseDirection(elevator)
		clearOrderAtFloor(elevator.floor, elevator.direction, orders, deliveredOrder)
		return
	}

	if orderInOppositeDirection(elevator.floor, elevator.direction, orders) {
		reverseDirection(elevator)
		startMoving(elevator, elevatorMotorTimer)
		return
	}

	enterIdle(elevator)
}

func handleDoorClosing(
	elevator *Elevator,
	orders Orders,
	openDoorC chan<- bool,
	deliveredOrder chan<- elevio.ButtonEvent,
	elevatorMotorTimer *time.Timer,
) {

	if elevator.behaviour != EB_DoorOpen {
		panic("Door closed while not open")
	}

	if orderInDirection(elevator.floor, elevator.direction, orders) {
		startMoving(elevator, elevatorMotorTimer)
		return
	}

	if orderAtFloorOppositeDirection(elevator.floor, elevator.direction, orders) {
		reverseDirection(elevator)
		stopAndOpenDoor(elevator, openDoorC)
		clearOrderAtFloor(elevator.floor, elevator.direction, orders, deliveredOrder)
		return
	}

	if orderInOppositeDirection(elevator.floor, elevator.direction, orders) {
		reverseDirection(elevator)
		startMoving(elevator, elevatorMotorTimer)
		return
	}

	enterIdle(elevator)
}
