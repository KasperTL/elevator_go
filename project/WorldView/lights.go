package WorldView

import (
	"project/config"
	"project/elevio"
)

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