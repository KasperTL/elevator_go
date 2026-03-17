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
	Floor       int
	Direction   ElevatorDirection
	Behaviour   ElevatorBehaviour
	Obstruction bool
	MotorStop   bool
}

func oppositeDirection(dir ElevatorDirection) ElevatorDirection {
	switch dir {
	case ED_Up:
		return ED_Down
	case ED_Down:
		return ED_Up
	default:
		panic("Invaid direction")
	}
}

func DirnToString(ed ElevatorDirection) string {
	switch ed {
	case ED_Up:
		return "up"
	case ED_Down:
		return "down"
	default:
		panic("Invalid direction")
	}

}

func BehaviourToString(b ElevatorBehaviour) string {
	switch b {
	case EB_Idle:
		return "idle"
	case EB_Moving:
		return "moving"
	case EB_DoorOpen:
		return "doorOpen"
	default:
		return "unknown"
	}
}

func ElevatorPrint(e Elevator) {
	fmt.Printf("  +--------------------+\n")
	fmt.Printf("  |floor = %-2d          |\n", e.GetFloor())
	fmt.Printf("  |Direction  = %-12s|\n", DirnToString(e.GetDirection()))
	fmt.Printf("  |Behaviour = %-12s|\n", BehaviourToString(e.GetBehaviour()))
	fmt.Printf("  +--------------------+\n")
	fmt.Printf("  |  Up  | Down |  Cab |\n")

	for f := config.NumFloors - 1; f >= 0; f-- {
		fmt.Printf("  |")
		for b := 0; b < config.NumElevatorButtons; b++ {

		}
		fmt.Printf("| %d\n", f)
	}
	fmt.Printf("  +--------------------+\n")
}

func InitializeElevator() Elevator {
	elevio.SetDoorOpenLamp(false)
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
		Floor:     elevatorFloor,
		Direction: ED_Down,
		Behaviour: EB_Idle,
	}
}

func (e Elevator) GetFloor() int                   { return e.Floor }
func (e Elevator) GetDirection() ElevatorDirection { return e.Direction }
func (e Elevator) GetBehaviour() ElevatorBehaviour { return e.Behaviour }
func (e Elevator) GetMotorStop() bool              { return e.MotorStop }
func (e Elevator) GetObstructed() bool             { return e.Obstruction }

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
