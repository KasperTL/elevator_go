package WorldView

import (
	"project/ElevatorDriver"
	"project/config"
)

type OrderState int 
const (
	OrderIdle       = 0 
	OrderPending    = 1
	OrderConfirmed  = 2
)

type WorldView struct {
    SenderID int 
    AliveList [config.NumElevators] bool
    ElevatorStates [config.NumElevators]ElevatorDriver.Elevator
    Orders [config.NumElevators][config.NumFloors][config.NumButtons]OrderState
}

func InitWorldView(nodeID int) WorldView {
    view := WorldView{SenderID: nodeID}
    return view 
}

func consesus(peers []OrderState, stateA OrderState, stateB OrderState) bool {
    for _, p := range peers {
        if p != stateA || p != stateB {
            return false 
        }
    }
    return true
}

func anyIs(peers []OrderState, state OrderState) bool {
    for _, p := range peers {
        if p == state {
            return true 
        }
    }
    return false
}

func updateHalOrders(
    orders [config.NumElevators][config.NumFloors][config.NumButtons] OrderState,
    NodeID int ,
) [config.NumElevators][config.NumFloors][config.NumButtons] OrderState {

    for floor := 0; floor < config.NumFloors; floor++ {
        for button := 0; button < config.NumButtons; button++ {
            currentOrderState := orders[NodeID][floor][button]
            newOrderState := currentOrderState

            var peers []OrderState

            for elevator := 0; elevator < config.NumElevators; elevator++ {
                if elevator != NodeID {
                    peers = append(peers, orders[elevator][floor][button])
                }
            }

            switch currentOrderState {
            case OrderIdle:
                if consesus(peers, OrderIdle, OrderPending) && anyIs(peers, OrderPending) {
                    newOrderState = OrderPending
                }
            case OrderPending:
                if consesus(peers, OrderPending, OrderConfirmed) && anyIs(peers, OrderConfirmed) {
                    newOrderState = OrderConfirmed
                }
            case OrderConfirmed:
                if consesus(peers, OrderConfirmed, OrderIdle) && anyIs(peers, OrderIdle) {
                    newOrderState = OrderIdle
                }
            }
            orders[NodeID][floor][button] = newOrderState
        }
    }
    return orders
}

 
