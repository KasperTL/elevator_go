package elev_algo

func request_above(e Elevator) int {
	for (floor := e.floor +1; floor < N_FLOORS; floor++) {
		for (button := 0; button < N_BUTTONS; button++) {
			if (e.requests[floor][button]) {
				return 1
			}
		}
	}
	return 0
}

func request_below(e Elevator) int {
	for (floor := 0; floor < e.floor; floor++) {
		for (button := 0; button < N_BUTTONS; button++) {
			if (e.requests[floor][button]) {
				return 1
			}
		}
	}
	return 0
}

func request_here(e Elevator) int {
	for (button := 0; button < N_BUTTONS; button++) {
		if (e.requests[e.floor][button]) {
			return 1
		}
	}
	return 0
}

func requests_chooseDirection(e Elevator) DirnBehaviourPair {
	switch e.dirn {
	case D_Up:
		if request_above(e) {
			return DirnBehaviourPair(D_Up, EB_Moving)
		} else if request_here(e) {
			return DirnBehaviourPair(D_Down, EB_DoorOpen)
		} else if request_below(e) {
			return DirnBehaviourPair(D_Down, EB_Moving)
		} else {
			return DirnBehaviourPair(D_Stop, EB_Idle)
		}
	case D_Down:
		if request_below(e) {
			return DirnBehaviourPair(D_Down, EB_Moving)
		} else if request_here(e) {
			return DirnBehaviourPair(D_Up, EB_DoorOpen)
		} else if request_above(e) {
			return DirnBehaviourPair(D_Up, EB_Moving)
		} else {
			return DirnBehaviourPair(D_Stop, EB_Idle)
		}
	case D_Stop:
		if request_here(e) {
			return DirnBehaviourPair(D_Stop, EB_DoorOpen)
		} else if request_above(e) {
			return DirnBehaviourPair(D_Up, EB_Moving)
		} else if request_below(e) {
			return DirnBehaviourPair(D_Down, EB_Moving)
		} else {
			return DirnBehaviourPair(D_Stop, EB_Idle)
		}
	default:
		return DirnBehaviourPair(D_Stop, EB_Idle)
	}
}

func requests_shouldStop(e Elevator) int {
	switch e.dirn {
	case D_Down:
		return e.requests[e.floor][B_Hallown] ||
		e.requests[e.floor][B_Cab] ||
		!request_below(e) 
	case D_Up:
		return e.requests[e.floor][B_Hallup] ||
		e.requests[e.floor][B_Cab] ||
		!request_above(e)
	case D_Stop:
		return 0
	default:
		return 1
	}
}

func requests_shouldClearImmediately(e Elevator, btn_floor int, btn_type Button) int {
	return
		(e.florr == btn_floor) &&
		(
			(e.dirn == D_Up && btn_type == B_HallUp) ||
			(e.dirn == D_Down && btn_type == B_HallDown) ||
			(e.dirn == D_Stop) ||
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

