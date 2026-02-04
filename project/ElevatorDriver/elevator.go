
package ElevatorDriver

import "fmt"

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
	last_floor int
	direction ElevatorDirection
	requests [N_FLOORS][N_BUTTONS] int
	behaviour ElevatorBehaviour
}


func StateToString(b ElevatorBehaviour) string {
    switch b {
    case Idle:     return "Idle"
    case Moving:   return "Moving"
    case DoorOpen: return "DoorOpen"
    default:       return "Unknown"
    }
}

func ElevatorPrint(e Elevator) {
    fmt.Printf("  +--------------------+\n")
    fmt.Printf("  |floor = %-2d          |\n", e.lastFloor)
    fmt.Printf("  |Direction  = %-12s|\n", DirnToString(e.Direction))
    fmt.Printf("  |Behaviour = %-12s|\n", StateToString(e.Behaviour))
    fmt.Printf("  +--------------------+\n")
    fmt.Printf("  |  Up  | Down |  Cab |\n")

    for f := NumFloors - 1; f >= 0; f-- {
        fmt.Printf("  |")
        for b := 0; b < NumButtons; b++ {
            if e.Requests[f][b] {
                fmt.Printf("  #   ")
            } else {
                fmt.Printf("  -   ")
            }
        }
        fmt.Printf("| %d\n", f)
    }
    fmt.Printf("  +--------------------+\n")
}