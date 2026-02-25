package ElevatorDriver

import (
	"project/config"
	"project/elevio"
)


type Orders [config.NumFloors][config.NumButtons] bool


func (o Orders) orderInSameDirection(dir ElevatorDirection) bool {

	switch dir {
	case ED_Up:
		for floor := e.floor + 1; floor < config.NumFloors; floor++ {
			for button := 0; button < config.NumButtons; button++ {
				if o[floor][button] {
					return true
				}
			}
		}
		return false
	case ED_Down:
		for floor := 0; floor < e.floor; floor++ {
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





func orderDone(e Elevator, floor int, direction direction, orderDoneC chan <- elevio.ButtonEvent) {
	if e.requests[floor][BT_Cab] {
		e.requests[floor][BT_Cab] = false
		orderDoneC <- elevio.ButtonEvent{Floor: floor, Button: BT_Cab}
	}
	if e.requests[floor][direction] {
		e.requests[floor][direction] = false
		orderDoneC <- elevio.ButtonEvent{Floor: floor, Button: direction}
	}
}


//This should be written as an switch case, not seperate functions
func request_above(e Elevator) int {
	for floor := e.floor +1; floor < config.NumFloors; floor++ {
		for button := 0; button < config.NumButtons; button++ {
			if (e.requests[floor][button]) {
				return 1
			}
		}
	}
	return 0
}

func request_below(e Elevator) int {
	for floor := 0; floor < e.floor; floor++ {
		for button := 0; button < config.NumButtons; button++ {
			if (e.requests[floor][button]) {
				return 1
			}
		}
	}
	return 0
}

func requestAtCurrentFloor(orders Orders, e Elevator) bool {
	return orders[e.floor][e.direction] || orders[e.floor][elevio.BT_Cab]
}

func requests_chooseDirection(e Elevator) DirnBehaviourPair {
	switch e.direction {
	case ED_Up:
		if request_above(e) {
			return ED_Up, EB_Moving
		} else if requestAtCurrentFloor(e.orders, e) {
			return ED_Down, EB_DoorOpen
		} else if request_below(e) {
			return ED_Down, EB_Moving
		} else {
			return EB_Idle
		}
	case ED_Down:
		if request_below(e) {
			return ED_Down, EB_Moving
		} else if requestAtCurrentFloor(e.orders, e) {
			return ED_Up, EB_DoorOpen
		} else if request_above(e) {
			return ED_Up, EB_Moving
		} else {
			return EB_Idle
		}

	default:
		return EB_Idle
	}
}

func requests_shouldStop(e Elevator) int {
	switch e.direction {
	case ED_Down:
		return e.requests[e.floor][BT_HallDown] ||
		e.requests[e.floor][BT_Cab] ||
		!request_below(e) 
	case ED_Up:
		return e.requests[e.floor][BT_HallUp] ||
		e.requests[e.floor][BT_Cab] ||
		!request_above(e)
	default:
		return 1
	}
}

func requests_shouldClearImmediately(e Elevator, btn_floor int, btn_type Button) bool {
	return (e.floor == btn_floor) &&
		(
			(e.direction == ED_Up && btn_type == BT_HallUp) ||
			(e.direction == ED_Down && btn_type == BT_HallDown) ||
			(e.direction == ED_Stop) ||
			(btn_type == B_Cab)
		)
}

func requests_clearAtCurrentFloor(e Elevator) Elevator {
	e.requests[e.floor][B_Cab] = 0
	switch e.dirn {
	case D_Up:
		if (!request_above(e)) && (!e.requests[e.floor][B_HallUp]) {
			e.requests[e.floor][B_HallDown] = 0
		}
		e.requests[e.floor][B_HallUp] = 0
	case D_Down:
		if(!request_below(e)) && (!e.requests[e.floor][B_HallDown]) {
			e.requests[e.floor][B_HallUp] = 0
		}
		e.requests[e.floor][B_HallDown] = 0
	case D_Stop:
	default:
		e.requests[e.floor][B_HallUp] = 0
		e.requests[e.floor][B_HallDown] = 0
	}
	return e
}

