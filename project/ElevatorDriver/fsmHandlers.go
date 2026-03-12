package ElevatorDriver

import (
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
	updatedElevatorState chan<- Elevator,
) {

	switch elevator.Behaviour {

	case EB_Idle:

		if orderAtFloorInDirection(elevator.Floor, elevator.Direction, orders) || cabOrderAtFloor(elevator.Floor, orders) {
			stopAndOpenDoor(elevator, openDoorC, elevatorMotorTimer)
			clearOrderAtFloor(elevator.Floor, elevator.Direction, orders, deliveredOrder)
			updatedElevatorState <- *elevator
			return
		}
		if orderAtFloorOppositeDirection(elevator.Floor, elevator.Direction, orders) {
			reverseDirection(elevator)
			stopAndOpenDoor(elevator, openDoorC, elevatorMotorTimer)
			clearOrderAtFloor(elevator.Floor, elevator.Direction, orders, deliveredOrder)
			updatedElevatorState <- *elevator
			return
		}
		if orderInDirection(elevator.Floor, elevator.Direction, orders) {
			startMoving(elevator, elevatorMotorTimer)
			updatedElevatorState <- *elevator
			return
		}
		if orderInOppositeDirection(elevator.Floor, elevator.Direction, orders) {
			reverseDirection(elevator)
			startMoving(elevator, elevatorMotorTimer)
			updatedElevatorState <- *elevator
			return
		}

		enterIdle(elevator, elevatorMotorTimer)

	case EB_DoorOpen:

		if orderAtFloorInDirection(elevator.Floor, elevator.Direction, orders) || cabOrderAtFloor(elevator.Floor, orders) {
			stopAndOpenDoor(elevator, openDoorC, elevatorMotorTimer)
			clearOrderAtFloor(elevator.Floor, elevator.Direction, orders, deliveredOrder)
			updatedElevatorState <- *elevator
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
	elevatorMotorTimer.Reset(config.ElevatorMotorTime)
	elevio.SetFloorIndicator(floor)
	elevator.Floor = floor

	if elevator.Behaviour != EB_Moving {
		enterIdle(elevator, elevatorMotorTimer)
		return
	}

	if orderAtFloorInDirection(elevator.Floor, elevator.Direction, orders) || cabOrderAtFloor(elevator.Floor, orders) {
		// TODO
		// This is toooooooo long, need to fix
		if !orderInDirection(elevator.Floor, elevator.Direction, orders) && orderAtFloorOppositeDirection(elevator.Floor, elevator.Direction, orders) && !orderAtFloorInDirection(elevator.Floor, elevator.Direction, orders) {
			stopAndOpenDoor(elevator, openDoorC, elevatorMotorTimer)
			reverseDirection(elevator)
			clearOrderAtFloor(elevator.Floor, elevator.Direction, orders, deliveredOrder)
		}
		stopAndOpenDoor(elevator, openDoorC, elevatorMotorTimer)
		clearOrderAtFloor(elevator.Floor, elevator.Direction, orders, deliveredOrder)
		return
	}

	if orderInDirection(elevator.Floor, elevator.Direction, orders) {
		return
	}

	if orderAtFloorOppositeDirection(elevator.Floor, elevator.Direction, orders) {
		stopAndOpenDoor(elevator, openDoorC, elevatorMotorTimer)
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

	if elevator.Behaviour != EB_DoorOpen {
		panic("Door closed while not open")
	}

	if orderInDirection(elevator.Floor, elevator.Direction, orders) {
		startMoving(elevator, elevatorMotorTimer)
		return
	}

	if orderAtFloorOppositeDirection(elevator.Floor, elevator.Direction, orders) {
		reverseDirection(elevator)
		stopAndOpenDoor(elevator, openDoorC, elevatorMotorTimer)
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
