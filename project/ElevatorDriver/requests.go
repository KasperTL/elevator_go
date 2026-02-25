package ElevatorDriver

import (
	"fmt"
	"project/config"
	"project/elevio"
)

type Orders [config.NumFloors][config.NumButtons]bool

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

func orderDone(direction ElevatorDirection, floor int, orders Orders, orderDoneC chan<- elevio.ButtonEvent) {
	fmt.Println("orderDone called - Floor:", floor, "Direction:", direction)
	fmt.Println("  orders[floor]:", orders[floor])
	if orders[floor][elevio.BT_Cab] {
		fmt.Println("  Sending BT_Cab to deliveredOrder")
		orderDoneC <- elevio.ButtonEvent{Floor: floor, Button: elevio.BT_Cab}
	}
	if orders[floor][direction.toBT()] { // ← Fikset: direction.toBT() i stedet for direction
		fmt.Println("  Sending direction button to deliveredOrder")
		orderDoneC <- elevio.ButtonEvent{Floor: floor, Button: direction.toBT()}
	}
}


func printOrders(orders Orders) {
	fmt.Printf("  +--------------------+\n")
	for f := config.NumFloors - 1; f >= 0; f-- {
		fmt.Printf("  |")
		for b := 0; b < config.NumButtons; b++ {
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