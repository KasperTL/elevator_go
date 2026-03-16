package ElevatorDriver

import (
	"fmt"
	"project/config"
	"project/elevio"
)

type Orders [config.NumFloors][config.NumElevatorButtons]bool

func orderInDirection(thisFloor int, direction ElevatorDirection, orders Orders) bool {

	switch direction {
	case ED_Up:
		for floor := thisFloor + 1; floor < config.NumFloors; floor++ {
			for button := 0; button < config.NumElevatorButtons; button++ {
				if orders[floor][button] {
					return true
				}
			}
		}
		return false
	case ED_Down:
		for floor := 0; floor < thisFloor; floor++ {
			for button := 0; button < config.NumElevatorButtons; button++ {
				if orders[floor][button] {
					return true
				}
			}
		}
		return false
	default:
		panic("Invalid elevator direction")
	}
}

func orderInOppositeDirection(thisFloor int, direction ElevatorDirection, orders Orders) bool {
	switch direction {
	case ED_Down:
		for floor := thisFloor + 1; floor < config.NumFloors; floor++ {
			for button := 0; button < config.NumElevatorButtons; button++ {
				if orders[floor][button] {
					return true
				}
			}
		}
		return false
	case ED_Up:
		for floor := 0; floor < thisFloor; floor++ {
			for button := 0; button < config.NumElevatorButtons; button++ {
				if orders[floor][button] {
					return true
				}
			}
		}
		return false
	default:
		panic("Invalid elevator direction")
	}
}

func clearOrderAtFloor(floor int, direction ElevatorDirection, orders Orders, clearOrderAtFloorC chan<- elevio.ButtonEvent) {
	if orders[floor][elevio.BT_Cab] {
		clearOrderAtFloorC <- elevio.ButtonEvent{Floor: floor, Button: elevio.BT_Cab}
	}
	if orders[floor][direction.toBT()] { //TODO: ← Fikset: direction.toBT() i stedet for direction
		clearOrderAtFloorC <- elevio.ButtonEvent{Floor: floor, Button: direction.toBT()}
	}
}

func cabOrderAtFloor(currentFloor int, orders Orders) bool {
	return orders[currentFloor][elevio.BT_Cab]
}

func orderAtFloorInDirection(currentFloor int, direction ElevatorDirection, orders Orders) bool {
	return orders[currentFloor][direction.toBT()] || orders[currentFloor][elevio.BT_Cab]
}

func orderAtFloorOppositeDirection(currentFloor int, direction ElevatorDirection, orders Orders) bool {
	return orders[currentFloor][oppositeDirection(direction).toBT()] || orders[currentFloor][elevio.BT_Cab] //TODO: Added cab orders.
}

func PrintOrders(orders Orders) {
	fmt.Printf("  +--------------------+\n")
	for f := config.NumFloors - 1; f >= 0; f-- {
		fmt.Printf("  |")
		for b := 0; b < config.NumElevatorButtons; b++ {
			if orders[f][b] {
				fmt.Printf("  #   ")
			} else {
				fmt.Printf("  -   ")
			}
		}
		fmt.Printf("| %d\n", f)
	}
	fmt.Printf("  +--------------------+\n")
}
