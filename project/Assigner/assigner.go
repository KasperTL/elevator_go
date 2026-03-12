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

	hallOrders := WorldView.HallOrdersAsBool(wv.Orders[nodeID])
	cabOrders := WorldView.CabOrdersAsBool(wv.Orders[nodeID], nodeID)

	stateMap := make(map[string]HRAElevState)
	for i, elevator := range wv.ElevatorStates {
		if !wv.AliveList[i] {
			continue
		}
		var movementDirection string
		switch elevator.Behaviour {
		case ElevatorDriver.EB_Idle:
			movementDirection = "stop"
		case ElevatorDriver.EB_DoorOpen:
			movementDirection = "stop"
		default:
			movementDirection = ElevatorDriver.DirnToString(elevator.Direction)
		}
		if elevator.GetObstructed() || elevator.GetMotorStop() {
			continue
		} else {
			stateMap[strconv.Itoa(i)] = HRAElevState{
				Behaviour:   ElevatorDriver.StateToString(elevator.GetBehaviour()),
				Floor:       elevator.GetFloor(),
				Direction:   movementDirection,
				CabRequests: cabOrders,
			}
		}

	}

	if len(stateMap) == 0 {
		return ElevatorDriver.Orders{}
	}

	input := HRAInput{hallOrders, stateMap}
	PrintHRAInput(input)
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
	PrintHRAOutput(*output)

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

func PrintHRAInput(input HRAInput) {
	prettyJSON, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}
	fmt.Println("=== HRA Input ===")
	fmt.Println(string(prettyJSON))
}

func PrintHRAOutput(output map[string]ElevatorDriver.Orders) {
    prettyJSON, err := json.MarshalIndent(output, "", "  ")
    if err != nil {
        fmt.Println("Error formatting JSON:", err)
        return
    }
    fmt.Println("=== HRA Output ===")
    fmt.Println(string(prettyJSON))
}