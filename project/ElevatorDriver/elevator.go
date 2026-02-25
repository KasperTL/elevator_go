package ElevatorDriver

import (
	"fmt"
	"project/config"
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
	requests    Orders // NumButtons is BT_HallDown, BT_HallUp and BT_Cab
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
			if e.requests[f][b] {
				fmt.Printf("  #   ")
			} else {
				fmt.Printf("  -   ")
			}
		}
		fmt.Printf("| %d\n", f)
	}
	fmt.Printf("  +--------------------+\n")
}

func InitializeElevator() Elevator {
	return Elevator{
		floor:     -1,
		direction: ED_Down,
		behaviour: EB_Idle,
		requests:  [config.NumFloors][config.NumButtons]bool{},
	}
}
