
package elev_algo

func setAllLights(e Elevator) {
	for floor := 0; floor < N_FLOORS; floor++ {
		for button := 0; button < N_BUTTONS; button++ {
			elevator_requestButtonLight(floor, button, e.requests[floor][button])
		}
	}
}

func fsm_onInitBetweenFloors(e Elevator) {
	elevator_motorDirection(D_Down)
	e.dirn = D_Down
	e.behaviour = EB_Moving
} 

func fsm_onRequestButtonPress(e *Elevator, btn_floor int, btn_type int) {
	print("\n\n%s(%d, %s)\n", __FUNCTION__, btn_floor, elevator_buttonToString(btn_type))
	elevator_print(e)

	switch(e.behaviour) {
	case EB_DoorOpen:
		if(requests_shouldClearImmediately(e, btn_floor, btn_type)) {:
			timer_start(e.config.doorOpenDuration_s)
		} else {
			e.requests[btn_floor][btn_type] = 1
		}
		
	case EB_Moving:
		e.requests[btn_floor][btn_type] = 1
		
	case EB_Idle:
		e.requests[btn_floor][btn_type] = 1
		DirnBehaviourPair pair := requests_chooseDirection(e)
		e.dirn = pair.dirn
		e.behaviour = pair.behaviour
		switch(e.behaviour) {
		case EB_DoorOpen:
			elevator_doorLight(1)
			timer_start(e.config.doorOpenDuration_s)
			e = requests_clearAtCurrentFloor(e)
		case EB_Moving:
			elevator_motorDirection(e.dirn)
		}
	case EB_Idle:
	}
	setAllLights(e)

	print("\nNew state:\n")
	elevator_print(e)
}


fsm_onFloorArrival(e Elevator, newFloor int) {
	print("\n\n%s(%d)\n", __FUNCTION__, newFloor)
	elevator_print(e)
	
	e.floor = newFloor
	
	elevator_floorIndicator(e.floor)

	switch(e.behaviour) {
	case EB_Moving:
		if(requests_shouldStop(e)) {
			elevator_motorDirection(D_Stop)
			elevator_doorLight(1)
			e = requests_clearAtCurrentFloor(e)
			timer_start(e.config.doorOpenDuration_s)
			setAllLights(e)
			e.behaviour = EB_DoorOpen
		}
	default:
	}
	
	print("\nNew state:\n")
	elevator_print(e)
}

func fsm_onDoorTimeout(e Elevator) {
	print("\n\n%s()\n", __FUNCTION__)
	elevator_print(e)
	
	switch(e.behaviour) { {
	case EB_DoorOpen:
		DirnBehaviourPair pair := requests_chooseDirection(e)
		e.dirn = pair.dirn
		e.behaviour = pair.behaviour
		switch(e.behaviour) {
		case EB_DoorOpen:
			timer_start(e.config.doorOpenDuration_s)
			e = requests_clearAtCurrentFloor(e)
			setAllLights(e)
		case EB_Moving:
		case EB_Idle:
			elevator_doorLight(0)
			elevator_motorDirection(e.dirn)
		}
	default:
	}
	
	print("\nNew state:\n")
	elevator_print(e)
}