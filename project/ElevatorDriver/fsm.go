
package ElevatorDriver

import (
	"project/config"
	"project/elevio"
	"time"
)





func elevator_fsm(
	newOrder	           <- chan Orders,
	updatedElevatorState   chan <- Elevator,
	deliverOrder           chan <- elevio.ButtonEvent,
) {

	
	orderDoneC := make(chan elevio.ButtonEvent)
	newFloorC := make(chan int)
	elevatorState := InitializeElevator()
	
	doorOpen := make(chan bool)
	doorObstructed := make(chan bool)
	doorClosing := make(chan bool)
	go door_fsm(doorOpen, doorObstructed, doorClosing)


	elevio.PollFloorSensor(newFloorC)
	

	var orders Orders

	

	for {
		select {
		case floor := <- newFloorC:
			switch {
			case EB_Moving:
				switch {
				case orders[floor][elevatorState.direction]:
					elevatorState.behaviour = EB_DoorOpen
					elevio.SetMotorDirection(elevio.MD_Stop)
					doorOpen <- true
					orderDone(elevatorState, floor, elevatorState.direction, orderDoneC)
				
				case orders[floor][elevio.BT_Cab]: // && orders.orderInSameDirection(elevator.direction):
					elevatorState.behaviour = EB_DoorOpen
					elevio.SetMotorDirection(elevio.MD_Stop)
					doorOpen <- true
					orderDone(elevatorState, floor, elevio.BT_Cab, orderDoneC)
				}
			}

		case obstrucion := <- doorObstructed:
			if obstrucion != elevatorState.obstruction {
				elevatorState.obstruction = obstrucion
				updatedElevatorState <- elevatorState
			}

		case orders = <- newOrder:
			switch elevatorState.behaviour {
				case EB_Idle:
					switch {
						case orders[elevatorState.floor][elevatorState.direction] || orders[elevatorState.floor][elevio.BT_Cab]: 
							elevatorState.behaviour = EB_DoorOpen
							orderDone(elevatorState, elevatorState.floor, elevatorState.direction, orderDoneC)
						
						case orders[elevatorState.floor][oppositeDirection(elevatorState.direction)] || orders[elevatorState.floor][elevio.BT_Cab]:
							elevatorState.direction = oppositeDirection(elevatorState.direction)
							elevatorState.behaviour = EB_DoorOpen
							orderDone(elevatorState, elevatorState.floor, oppositeDirection(elevatorState.direction), orderDoneC)

						case orders.orderInSameDirection(elevatorState.direction):
							elevatorState.behaviour = EB_Moving

						case orders.orderInSameDirection(oppositeDirection(elevatorState.direction)):
							elevatorState.direction = oppositeDirection(elevatorState.direction)
							elevatorState.behaviour = EB_Moving

						default:
							elevatorState.behaviour = EB_Idle
					}
			case EB_Moving:
				switch {

				}
				
			}
		}
	}
}
	











































//This will be moved to Orders module later
func setAllLights(e Elevator) {
	for floor := 0; floor < config.NumFloors; floor++ {
		for button := 0; button < config.NumButtons; button++ {
			elevio.SetButtonLamp(button, floor, e.requests[floor][button])
		}
	}
}



func fsm_onInitBetweenFloors(e Elevator) {
	elevio.SetMotorDirection(elevio.MD_Down)
	e.direction = ED_Down
	e.behaviour = EB_Moving
} 


//We will change the input values. floor and button type will be an Order struct later
func fsm_onNewOrderRequest(e *Elevator, btn_floor int, btn_type int) {
	switch e.behaviour {
	case EB_Idle:
		e.requests[btn_floor][btn_type] = true
		
		newDirection, newBehaviour := requests_chooseDirection(e)
		
		e.direction = newDirection
		e.behaviour = newBehaviour

		switch e.behaviour {
		case EB_DoorOpen:
			elevio.SetDoorOpenLamp(1)
			TimerStart(3)
			elevio.SetDoorOpenLamp(0)
			e = requests_clearAtCurrentFloor(e)
		case EB_Moving:
			elevio.SetMotorDirection(e.direction)
		}
	case EB_Moving:
		e.requests[btn_floor][btn_type] = true
	case EB_DoorOpen:
		swi
		e.requests[btn_floor][btn_type] = true
	}	
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