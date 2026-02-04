package ElevatorDriver

import (
	"time"
)

var timerEndTime time.Time
var timerActive bool

// TimerStart starts a timer for the given duration in seconds.
func TimerStart(duration float64) {
	timerEndTime = time.Now().Add(time.Duration(duration * float64(time.Second)))
	timerActive = true
}

func TimerStop() {
	timerActive = false
}

func TimerTimedOut() bool {
	return timerActive && time.Now().After(timerEndTime)
}

