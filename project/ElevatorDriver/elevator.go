package ElevatorDriver

import (
	"fmt"
	"project/config"
	"project/elevio"
)

type ElevatorDirection int

const (
	ED_Up   = 1
	ED_Down = 2
)

type ElevatorBehaviour int

const (
	EB_Idle     = 0
	EB_DoorOpen = 1
	EB_Moving   = 2
)

type Elevator struct {
	floor       int
	direction   ElevatorDirection
	behaviour   ElevatorBehaviour
	obstruction bool
	motorStop   bool // Must be implemented
}

func oppositeDirection(dir ElevatorDirection) ElevatorDirection {
	switch dir {
	case ED_Up:
		return ED_Down
	case ED_Down:
		return ED_Up
	default:
		panic("Invalid direction")
	}
}

// Look ower erikIsChamp, they have a clean solution in elevatorFsm.go func ToString(). Same for StateToString()
func DirnToString(ed ElevatorDirection) string {
	switch ed {
	case ED_Up:
		return "Up"
	case ED_Down:
		return "Down"
	default:
		return "Unknown"
	}
}

func StateToString(b ElevatorBehaviour) string {
	switch b {
	case EB_Idle:
		return "Idle"
	case EB_Moving:
		return "Moving"
	case EB_DoorOpen:
		return "DoorOpen"
	default:
		return "Unknown"
	}
}

func ElevatorPrint(e Elevator) {
	fmt.Printf("  +--------------------+\n")
	fmt.Printf("  |floor = %-2d          |\n", e.floor)
	fmt.Printf("  |Direction  = %-12s|\n", DirnToString(e.direction))
	fmt.Printf("  |Behaviour = %-12s|\n", StateToString(e.behaviour))
	fmt.Printf("  +--------------------+\n")
	fmt.Printf("  |  Up  | Down |  Cab |\n")

	for f := config.NumFloors - 1; f >= 0; f-- {
		fmt.Printf("  |")
		for b := 0; b < config.NumButtons; b++ {

		}
		fmt.Printf("| %d\n", f)
	}
	fmt.Printf("  +--------------------+\n")
}

func InitializeElevator() Elevator {
	var elevatorFloor int
	if elevio.GetFloor() != -1 {
		elevatorFloor = elevio.GetFloor()

	} else {
		elevio.SetMotorDirection(elevio.MD_Down)

		for {
			if elevio.GetFloor() != -1 {
				elevatorFloor = elevio.GetFloor()
				elevio.SetMotorDirection(elevio.MD_Stop)
				break
			}
		}
	}
	return Elevator{
		floor:     elevatorFloor,
		direction: ED_Down,
		behaviour: EB_Idle,
	}
}

func (ed ElevatorDirection) toMD() elevio.MotorDirection {
	switch ed {
	case ED_Up:
		return elevio.MD_Up
	case ED_Down:
		return elevio.MD_Down
	default:
		panic("toMD called with invalid ElevatorDirection")
	}
}

func (ed ElevatorDirection) toBT() elevio.ButtonType {
	switch ed {
	case ED_Up:
		return elevio.BT_HallUp
	case ED_Down:
		return elevio.BT_HallDown
	default:
		panic("toBT called with invalid ElevatorDirection")
	}
}
