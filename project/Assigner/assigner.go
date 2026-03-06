package Assigner

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"project/ElevatorDriver"
	"project/WorldView"
	"project/config"
	"runtime"
	"strconv"
)

type HRAElevState struct {
	Behaviour   string                 `json:"behaviour"`
	Floor       int                    `json:"floor"`
	Direction   string                 `json:"direction"`
	CabRequests [config.NumFloors]bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [config.NumFloors][2]bool `json:"hallRequests"`
	States       map[string]HRAElevState   `json:"states"`
}

func CalculateOptimalOrders(wv WorldView.WorldView, nodeID int) ElevatorDriver.Orders {

	hraExecutable := ""
	switch runtime.GOOS {
	case "linux":
		hraExecutable = "hall_request_assigner"
	case "windows":
		hraExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	stateMap := make(map[string]HRAElevState)
	for i, e := range wv.ElevatorStates {
		if !wv.AliveList[i] {
			continue
		}
		if e.Elevator.GetObstructed() || e.Elevator.GetMotorStop() {
			continue
		} else {
			stateMap[strconv.Itoa(i)] = HRAElevState{
				Behaviour:   ElevatorDriver.StateToString(e.Elevator.GetBehaviour()),
				Floor:       e.Elevator.GetFloor(),
				Direction:   ElevatorDriver.DirnToString(e.Elevator.GetDirection()),
				CabRequests: e.CabOrders,
			}
		}
	}

	orders := WorldView.FromOrderStateToBool(wv.HallOrders[nodeID])
	input := HRAInput{orders, stateMap}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
		panic("json.Marshal error")
	}

	ret, err := exec.Command("Assigner/Excecutables/"+hraExecutable, "-i", "--includeCab", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
		panic("exec.Command error")
	}

	output := new(map[string]ElevatorDriver.Orders)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		panic("json.Unmarshal error")
	}
	// fmt.Printf("output: \n")
    // for k, v := range *output {
    //     fmt.Printf("%6v :  %+v\n", k, v)
    // }

	return (*output)[strconv.Itoa(nodeID)]
}

func Assigner(
	incomingC <-chan WorldView.WorldView,
	assignedOrdersC chan<- ElevatorDriver.Orders,
	myNodeId int,
) {
	for wv := range incomingC {
		assignedOrders := CalculateOptimalOrders(wv, myNodeId)
		assignedOrdersC <- assignedOrders
	}
}
