package ElevatorDriver

import (	
	"project/config"
	"project/elevio"
	"time"
)


func startMoving(elevator *Elevator, timer *time.Timer) {
	elevator.Behaviour = EB_Moving
	elevio.SetMotorDirection(elevator.Direction.toMD()) //Migh want to use a channel for this instead
	timer.Reset(config.ElevatorMotorTime)
}


func stopAndOpenDoor(elevator *Elevator, openDoorC chan<- bool) {
	elevator.Behaviour = EB_DoorOpen
	elevio.SetMotorDirection(elevio.MD_Stop)
	openDoorC <- true
}

func enterIdle(elevator *Elevator) {
	elevator.Behaviour = EB_Idle
	elevio.SetMotorDirection(elevio.MD_Stop)
}

func reverseDirection(elevator *Elevator) {
	elevator.Direction = oppositeDirection(elevator.Direction)
}

