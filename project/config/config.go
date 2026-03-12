package config

import (
	"time"
)

const (
	NumFloors    = 4
	NumElevators = 3

	NumElevatorButtons = 3
	NumOrderTypes      = 2 + NumElevators
	PeersPortNumber    = 58735
	BcastPortNumber    = 58730
	Buffer             = 1024

	DisconnectTime    = 1 * time.Second
	DoorOpenDuration  = 3 * time.Second
	ElevatorMotorTime = 4 * time.Second
	HeartbeatTime     = 15 * time.Millisecond
)
