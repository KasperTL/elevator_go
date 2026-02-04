
package ElevatorDriver
import (
	"fmt"
)

// This is just for development purposes
// In the real system, these constants and types would be provided by the config file
const N_FLOORS = 4
const N_BUTTONS = 3

type ElevatorDirection int 
const (
	ED_Up = 1
	ED_Stop = 0
	ED_Down = -1
)

type ElevatorBehaviour int 
const (
	EB_Moving = 0
	EB_Idle = 1
	EB_DoorOpen = 2
)



type Elevator struct {
	floor int
	direction ElevatorDirection
	requests [N_FLOORS][N_BUTTONS] bool
	behaviour ElevatorBehaviour
}



func DirnToString(ed ElevatorDirection) string {
	switch ed {
	case ED_Up: return "Up"
	case ED_Down: return "Down"
	case ED_Stop: return "Stop"
	default: return "Unknown"
	}
} 



func StateToString(b ElevatorBehaviour) string {
    switch b {
    case EB_Idle:     return "Idle"
    case EB_Moving:   return "Moving"
    case EB_DoorOpen: return "DoorOpen"
    default:       return "Unknown"
    }
}



func ElevatorPrint(e Elevator) {
    fmt.Printf("  +--------------------+\n")
    fmt.Printf("  |floor = %-2d          |\n", e.floor)
    fmt.Printf("  |Direction  = %-12s|\n", DirnToString(e.direction))
    fmt.Printf("  |Behaviour = %-12s|\n", StateToString(e.behaviour))
    fmt.Printf("  +--------------------+\n")
    fmt.Printf("  |  Up  | Down |  Cab |\n")

    for f := N_FLOORS - 1; f >= 0; f-- {
        fmt.Printf("  |")
        for b := 0; b < N_BUTTONS; b++ {
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
		floor: -1,
		direction: ED_Stop,
		behaviour: EB_Idle,
		requests: [N_FLOORS][N_BUTTONS] bool{},
	}
}

