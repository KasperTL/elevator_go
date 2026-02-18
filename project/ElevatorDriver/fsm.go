
package ElevatorDriver

import (
	"project/config"
	"project/elevio"
	"time"
)





func elevator_fsm(
	newOrder	     <- chan Orders,
	newElevator	     chan <- Elevator,
	deliverOrder     chan <- elevio.ButtonEvent,
) {

	
	orderDoneC := make(chan elevio.ButtonEvent)
	newFloorC := make(chan int)
	elevator := InitializeElevator()
	newObstructionC := make(chan bool)


	elevio.PollFloorSensor(newFloorC)
	elevio.PollObstructionSwitch(newObstructionC)

	var orders Orders

	

	for {
		select {
		case floor := <- newFloorC:
			switch {
			case EB_Moving:
				switch {
				case orders[floor][elevator.direction]:
					elevator.behaviour = EB_DoorOpen
					elevio.SetMotorDirection(elevio.MD_Stop)
					orderDone(elevator, floor, elevator.direction, orderDoneC)
				
				case orders[floor][elevio.BT_Cab]: // && orders.orderInSameDirection(elevator.direction):
					elevator.behaviour = EB_DoorOpen
					elevio.SetMotorDirection(elevio.MD_Stop)
					orderDone(elevator, floor, elevator.direction, orderDoneC)
				}
			}

		case Obstruction := <- newObstructionC:
			switch {
				case EB_DoorOpen:
		
		

		case orders = <- newOrder:
			switch newElevator.behaviour {
			case EB_Idle:
				switch {
					case orders[elevator.floor][elevator.direction] || orders[elevator.floor][BT_Cab]: 
						elevator.behaviour = EB_DoorOpen
						orderDone(elevator, elevator.floor, elevator.direction, orderDoneC)
					
					case orders[elevator.floor][oppositeDirection(elevator.direction)] || orders[elevator.floor][BT_Cab]:
						elevator.direction = EB_DoorOpen
						orderDone(elevator, elevator.floor, oppositeDirection(elevator.direction), orderDoneC)

					case orders.orderInSameDirection(elevator.direction):
						elevator.behaviour = EB_Moving

					case orders.orderInSameDirection(oppositeDirection(elevator.direction)):
						elevator.direction = oppositeDirection(elevator.direction)
						elevator.behaviour = EB_Moving

					default:
						elevator.behaviour = EB_Idle
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