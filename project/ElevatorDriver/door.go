
package ElevatorDriver
import (
	"fmt"
    "project/config"
)

type Door struct {
	doorIsOpen bool
	doorIsObstructed bool
}

func door_fsm(
	doorIsOpen <- chan bool,
	doorIsObstructed <- chan bool,
) {
	door := Door{doorIsOpen: false, doorIsObstructed: false}

	for {
		select {
			case 
		}
	}