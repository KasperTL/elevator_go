package ElevatorDriver

import (
	"project/config"
	"project/elevio"
)


type Orders [config.NumFloors][config.NumButtons] bool


func (o Orders) orderInSameDirection(elevator Elevator) bool {

	switch elevator.direction {
	case ED_Up:
		for floor := elevator.floor + 1; floor < config.NumFloors; floor++ {
			for button := 0; button < config.NumButtons; button++ {
				if o[floor][button] {
					return true
				}
			}
		}
		return false
	case ED_Down:
		for floor := 0; floor < elevator.floor; floor++ {
			for button := 0; button < config.NumButtons; button++ {
				if o[floor][button] {
					return true
				}
			}
		}
		return false
	default:
		panic("Invalid elevator direction")
	}
}



func (o Orders) orderInOppositeDirection(elevator Elevator) bool {
	switch elevator.direction {
	case ED_Down:
		for floor := elevator.floor + 1; floor < config.NumFloors; floor++ {
			for button := 0; button < config.NumButtons; button++ {
				if o[floor][button] {
					return true
				}
			}
		}
		return false
	case ED_Up:
		for floor := 0; floor < elevator.floor; floor++ {
			for button := 0; button < config.NumButtons; button++ {
				if o[floor][button] {
					return true
				}
			}
		}
		return false
	default:
		panic("Invalid elevator direction")
	}
}




func orderDone(elevator Elevator, floor int, orderDoneC chan <- elevio.ButtonEvent) {
	if elevator.requests[floor][elevio.BT_Cab] {
		elevator.requests[floor][elevio.BT_Cab] = false
		orderDoneC <- elevio.ButtonEvent{Floor: floor, Button: elevio.BT_Cab}
	}
	if elevator.requests[floor][elevator.direction] {
		elevator.requests[floor][elevator.direction] = false
		orderDoneC <- elevio.ButtonEvent{Floor: floor, Button: e.direction}
	}
}
