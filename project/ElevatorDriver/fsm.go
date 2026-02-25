package ElevatorDriver

import (
	"project/config"
	"project/elevio"
	"time"
)

func elevator_fsm(
	newOrder <-chan Orders,
	updatedElevatorState chan<- Elevator,
	deliveredOrder chan<- elevio.ButtonEvent,
) {

	newFloorC := make(chan int)
	elevatorState := InitializeElevator()

	openDoorC := make(chan bool)
	doorObstructedc := make(chan bool)
	doorClosingc := make(chan bool)
	go door_fsm(openDoorC, doorObstructedc, doorClosingc)

	elevio.PollFloorSensor(newFloorC)

	var orders Orders

	elevatorMotorTimer := time.NewTimer(config.ElevatorMotorTime)
	elevatorMotorTimer.Stop()

	for {
		select {
		case floor := <-newFloorC:
			elevatorMotorTimer.Reset(config.ElevatorMotorTime)
			elevio.SetFloorIndicator(floor)
			elevatorState.floor = floor
			switch elevatorState.behaviour {
			case EB_Moving:
				switch {
				case orders[floor][elevatorState.direction] || orders[floor][elevio.BT_Cab]:
					elevatorState.behaviour = EB_DoorOpen
					elevio.SetMotorDirection(elevio.MD_Stop)
					openDoorC <- true
					orderDone(elevatorState, floor, deliveredOrder)

				case orders.orderInSameDirection(elevatorState):

				case orders[floor][oppositeDirection(elevatorState.direction)]:
					elevatorState.behaviour = EB_DoorOpen
					elevio.SetMotorDirection(elevio.MD_Stop)
					openDoorC <- true
					elevatorState.direction = oppositeDirection(elevatorState.direction)
					orderDone(elevatorState, floor, deliveredOrder)

				case orders.orderInOppositeDirection(elevatorState):
					elevatorState.direction = oppositeDirection(elevatorState.direction)
					elevio.SetMotorDirection(elevatorState.direction.toMD())

				default:
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevatorState.behaviour = EB_Idle
				}
			}
			updatedElevatorState <- elevatorState

		case obstrucion := <-doorObstructedc:
			if obstrucion != elevatorState.obstruction {
				elevatorState.obstruction = obstrucion
				updatedElevatorState <- elevatorState
			}

		case <-doorClosingc:
			switch {
			case orders.orderInSameDirection(elevatorState):
				elevatorState.behaviour = EB_Moving
				elevio.SetMotorDirection(elevatorState.direction.toMD())
				elevatorMotorTimer.Reset(config.ElevatorMotorTime)

			case orders[elevatorState.floor][oppositeDirection(elevatorState.direction)]:
				openDoorC <- true
				elevatorState.direction = oppositeDirection(elevatorState.direction)
				orderDone(elevatorState, elevatorState.floor, deliveredOrder)

				//Need to change the names of functions, cant call orderInSameDirection(oppositeDirection)
			case orders.orderInOppositeDirection(elevatorState):
				elevatorState.direction = oppositeDirection(elevatorState.direction)
				elevatorState.behaviour = EB_Moving
				elevatorMotorTimer.Reset(config.ElevatorMotorTime)
				elevio.SetMotorDirection(elevatorState.direction.toMD())

			default:
				elevatorState.behaviour = EB_Idle
			}
			updatedElevatorState <- elevatorState

		case orders = <-newOrder:
			switch elevatorState.behaviour {
			case EB_Idle:
				switch {
				case orders[elevatorState.floor][elevatorState.direction] || orders[elevatorState.floor][elevio.BT_Cab]:
					openDoorC <- true
					orderDone(elevatorState, elevatorState.floor, deliveredOrder)
					elevatorState.behaviour = EB_DoorOpen
					//Wonder if we should combine the two first cases, as they have the same logic, and the second only need too change direction

				case orders[elevatorState.floor][oppositeDirection(elevatorState.direction)] || orders[elevatorState.floor][elevio.BT_Cab]:
					openDoorC <- true
					elevatorState.direction = oppositeDirection(elevatorState.direction)
					orderDone(elevatorState, elevatorState.floor, deliveredOrder)
					elevatorState.behaviour = EB_DoorOpen

				case orders.orderInSameDirection(elevatorState):
					elevatorState.behaviour = EB_Moving
					elevio.SetMotorDirection(elevatorState.direction.toMD())

					//Need to change the names of functions, cant call orderInSameDirection(oppositeDirection)
				case orders.orderInOppositeDirection(elevatorState):
					elevatorState.direction = oppositeDirection(elevatorState.direction)
					elevatorState.behaviour = EB_Moving
					elevio.SetMotorDirection(elevatorState.direction.toMD())

				default:
					elevatorState.behaviour = EB_Idle
				}
				updatedElevatorState <- elevatorState

			case EB_DoorOpen:
				switch {
				case orders[elevatorState.floor][elevatorState.direction] || orders[elevatorState.floor][elevio.BT_Cab]:
					openDoorC <- true
					orderDone(elevatorState, elevatorState.floor, deliveredOrder)

				default:

				}
			default:
				panic("New order not handled")

			}
		}
	}
}
