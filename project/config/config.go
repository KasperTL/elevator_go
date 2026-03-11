package config

import (
	"time"
)

const (
	NumFloors       = 4
	NumElevators    = 3

	// Number of distinct order buttons
	NumButtons      = 2 + NumElevators 
	PeersPortNumber = 58735
	BcastPortNumber = 58750
	Buffer          = 1024

	DisconnectTime    = 1 * time.Second
	DoorOpenDuration  = 3 * time.Second
	ElevatorMotorTime = 4 * time.Second
	HeartbeatTime     = 1500 * time.Millisecond
)
