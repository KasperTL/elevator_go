package ElevatorDriver

import (
	"project/config"
	"project/elevio"
	"time"
)

func startMoving(elevator *Elevator, timer *time.Timer) {
	elevator.Behaviour = EB_Moving
	elevio.SetMotorDirection(elevator.Direction.toMD()) //TODO:Migh want to use a channel for this instead
	timer.Reset(config.ElevatorMotorTime)
}

func stopElevator(timer *time.Timer) {
	elevio.SetMotorDirection(elevio.MD_Stop)
	timer.Stop()
}

func openDoor(elevator *Elevator, openDoorC chan<- bool) {
	elevator.Behaviour = EB_DoorOpen
	openDoorC <- true
}

//TODO: Verify that openDoor and stopElevator works.
//func stopAndOpenDoor(elevator *Elevator, openDoorC chan<- bool, timer *time.Timer) {
//	elevator.Behaviour = EB_DoorOpen
//	elevio.SetMotorDirection(elevio.MD_Stop)
//	openDoorC <- true
//	timer.Stop()
//}

func enterIdle(elevator *Elevator, timer *time.Timer) {
	elevator.Behaviour = EB_Idle
	elevio.SetMotorDirection(elevio.MD_Stop)
	timer.Stop()
}

func reverseDirection(elevator *Elevator) {
	elevator.Direction = oppositeDirection(elevator.Direction)
}
