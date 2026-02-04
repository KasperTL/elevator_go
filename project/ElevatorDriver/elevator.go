
package ElevatorDriver

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



