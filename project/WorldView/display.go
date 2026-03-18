package WorldView

import (
	"project/config"
	"project/elevio"
)

func CabOrdersAsBool(Orders [config.NumFloors][config.NumOrderTypes]OrderState, nodeID int) [config.NumFloors]bool {
	var cabOrders [config.NumFloors]bool
	for floor := 0; floor < config.NumFloors; floor++ {
		cabOrders[floor] = (Orders[floor][2+nodeID] == OrderConfirmed) || (Orders[floor][2+nodeID] == OrderComplete)
	}
	return cabOrders
}

func HallOrdersAsBool(Orders [config.NumFloors][config.NumOrderTypes]OrderState) [config.NumFloors][2]bool {
	var hallOrders [config.NumFloors][2]bool
	for floor := 0; floor < config.NumFloors; floor++ {
		hallOrders[floor][0] = (Orders[floor][elevio.BT_HallUp] == OrderConfirmed) || (Orders[floor][elevio.BT_HallUp] == OrderComplete)
		hallOrders[floor][1] = (Orders[floor][elevio.BT_HallDown] == OrderConfirmed) || (Orders[floor][elevio.BT_HallDown] == OrderComplete)
	}
	return hallOrders
}

func setOrderLights(myWorldView WorldView, myNodeID int) {
	for floor := 0; floor < config.NumFloors; floor++ {
		for button := 0; button < config.NumElevatorButtons; button++ {
			var buttonValue int
			if button == elevio.BT_Cab {
				buttonValue = 2 + myNodeID
			} else {
				buttonValue = button
			}
			orderState := myWorldView.Orders[myNodeID][floor][buttonValue]
			buttonType := elevio.ButtonType(button)
			switch orderState {
			case OrderConfirmed:
				elevio.SetButtonLamp(buttonType, floor, true)
			case OrderIdle:
				elevio.SetButtonLamp(buttonType, floor, false)
			case OrderPending:
				continue
			case OrderComplete:
				continue
			}
		}
	}
}