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
) {

	switch elevator.Behaviour {

	//TODO: Changed the if-statments with orderAtFloorInDirection
	//Earlier it was a check for both hall and cab order, this is combined to orderAtFloorInDirection

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
		//TODO: Added at cleanup, verify that this is correct
		break

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

	//TODO: This must be wrong. Shouldnt be possible to arrive at a floor while not moving?
	//elevator.Behaviour != EB_Moving {
	// TODO : maybo wrong?
	//	enterIdle(elevator, elevatorMotorTimer)
	//	return
	//}

	if orderAtFloorInDirection(elevator.Floor, elevator.Direction, orders) {
		if !orderInDirection(elevator.Floor, elevator.Direction, orders) && orderAtFloorOppositeDirection(elevator.Floor, elevator.Direction, orders) { //TODO, verifiser at dette fungerer!&& !orderAtFloorInDirection(elevator.Floor, elevator.Direction, orders) {
			stopElevator(elevatorMotorTimer)
			openDoor(elevator, openDoorC)
			clearOrderAtFloor(elevator.Floor, elevator.Direction, orders, deliveredOrder)
			//TODO: Changed the order of the functions, now it clears order before reversing direction, verify that this is correct
			reverseDirection(elevator)

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

	if elevator.Behaviour != EB_DoorOpen {
		panic("Door closed while not open")
	}

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

	enterIdle(elevator, elevatorMotorTimer)
}
