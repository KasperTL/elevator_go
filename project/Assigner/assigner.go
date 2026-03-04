package assigner

import (
	"project/config"
)

type AssignerFormatedState struct {
	Behaviour   string                 `json:"behaviour"`
	Floor       int                    `json:"floor"`
	Direction   string                 `json:"direction"`
	CabRequests [config.NumFloors]bool `json:"cabRequests"`
}

type AssignerInput struct {
	HallRequests [config.NumFloors][2]bool 				`json:"hallRequests"`
	States       map[string]AssignerFormatedState       `json:"states"`
}

func 