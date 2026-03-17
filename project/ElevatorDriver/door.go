package ElevatorDriver

import (
	"project/config"
	"project/elevio"
	"time"
)

type DoorState int

const (
	DS_Open       = 0
	DS_Closed     = 1
	DS_Obstructed = 2
)

func door_fsm(
	openDoorC <-chan bool,
	doorObstructedC chan<- bool,
	doorClosingC chan<- bool,
) {
	obstructionC := make(chan bool, config.Buffer)
	myDoorState := DS_Closed
	doorIsObstructed := false
	doorOpenTimer := time.NewTimer(config.DoorOpenDuration)
	doorOpenTimer.Stop()

	go elevio.PollObstructionSwitch(obstructionC)

	for {
		select {
		//TODO: Could remove if statments and just send doorObstructedC <- doorIsObstructed
		case doorIsObstructed = <-obstructionC:
			if doorIsObstructed {
				myDoorState = DS_Obstructed
				doorObstructedC <- true
			} else {
				doorObstructedC <- false
			}

		case <-openDoorC:
			switch myDoorState {
			case DS_Closed:
				doorOpenTimer.Reset(config.DoorOpenDuration)
				myDoorState = DS_Open
				elevio.SetDoorOpenLamp(true)
			case DS_Open:
				doorOpenTimer.Reset(config.DoorOpenDuration)
				elevio.SetDoorOpenLamp(true)
			case DS_Obstructed:
				doorOpenTimer.Reset(config.DoorOpenDuration)
				myDoorState = DS_Obstructed
				elevio.SetDoorOpenLamp(true)
			}

		case <-doorOpenTimer.C:
			switch myDoorState {
			case DS_Open:
				myDoorState = DS_Closed
				doorClosingC <- true
				elevio.SetDoorOpenLamp(false)

			case DS_Obstructed:
				if !doorIsObstructed {
					myDoorState = DS_Closed
					doorClosingC <- true
					elevio.SetDoorOpenLamp(false)
				} else {
					doorOpenTimer.Reset(config.DoorOpenDuration)
				}
			}
		}
	}
}
