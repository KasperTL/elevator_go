package WorldView

import (
	"project/config"
	"project/ElevatorDriver"
)


type ConsensusMode int 
const (
	Standalone = 0
	Networked  = 1
)


type OrderState int
const (
	OrderIdle      = 0
	OrderPending   = 1
	OrderConfirmed = 2
	OrderComplete  = 3
)

type WorldView struct {
	NodeID           int
	AliveList        [config.NumElevators]bool
	ElevatorStates   [config.NumElevators]ElevatorDriver.Elevator
	Orders           [config.NumElevators][config.NumFloors][config.NumOrderTypes]OrderState
	CabOrderRecovery [config.NumElevators][config.NumFloors]OrderState
}

func InitWorldView(nodeID int) WorldView {
	view := WorldView{NodeID: nodeID}
	view.AliveList[nodeID] = true
	return view
}
