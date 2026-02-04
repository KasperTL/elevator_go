package elev_algo
import "time"

var timerEndTime time.Time
var timerActive int

func get_wall_time() float64 {
	return float64(time.Now().UnixNano()) / 1e9
}

func timer_start(duration float64) {
	timerEndTime = get_wall_time() + duration
	timerActive = 1
}

func timer_stop() {
	timerActive = 0	
}

func timer_expired() int {
	return (timerActive && get_wall_time() > timerEndTime)
}


