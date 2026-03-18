# Elevator Project Description

# Introduction

This is the completed code for the Elevator project in TTK4145, Real Time Programming, NTNU Spring 2026

The elevator project, according to the [specification](https://github.com/TTK4145/Project) provided, is to create a fault-tolerant distributed system of multiple elevators cooperating to provide a seemless user experience, even with packet loss, power outages, crashes and loss of network connectivity.

# Design

The code is written in Go language to communicate between the network and modules. The network is peer-to-peer and uses UDP for node connection. The system consists of modules Assigner, ElevatorDriver, Network, WorldView.

## [OrderDispatcher](https://github.com/KasperTL/elevator_go/tree/main/project/OrderDispatcher)

This module assigns orders using the [cost function](https://github.com/TTK4145/Project-resources/tree/master/cost_fns#alternative-2-reassigning-all-requests) (reassigning all requests) provided. 

### dispatcher.go
The main assignment logic. On every update from WorldView it:
- Builds an input for the HRA executable containing:
  - All confirmed hall orders
  - The current state of every alive elevator (floor, direction, behaviour, cab orders)

## [ElevatorDriver](https://github.com/KasperTL/elevator_go/tree/main/project/ElevatorDriver)

This module handles controls the pyhsical elevators by taking orders as input and drives the elevators to serve them. It also send back state updates for complete orders to the rest of the system

## [Network](https://github.com/KasperTL/elevator_go/tree/main/project/Network)

This module handles the communication between the elevators via worldview. This was taken from [project resources](https://github.com/TTK4145/Network-go) provided

## [WorldView](https://github.com/KasperTL/elevator_go/tree/main/project/WorldView)

This module is responsible for maintaining a consistent shared representation of the system state across the elevators.

### types.go
Defines the core data structures:
- `ConsensusMode`: whether the elevator is running `Standalone` 
  (without peers) or `Networked` (connected to peers). Determines 
  how orders are handled. 
- `OrderState`: the four states an order can be in: 
  `Idle → Pending → Confirmed → Complete`
- `WorldView`: the main data structure containing:
  - `NodeID`: which elevator this worldview belongs to
  - `AliveList`: which elevators are currently online
  - `ElevatorStates`: floor, direction and behaviour of every elevator
  - `Orders`: order states for every elevator, floor and button type
  - `CabOrderRecovery`: stores cab orders of a disconnected elevator 
    so they can be restored when it reconnects

### worldview.go
Logic for updating and reading the worldview. Handles merging peer 
states, advancing order states through the consensus protocol, and 
collecting peer order views.

### manager.go
The main event loop for worldview. 
- Button presses from the local elevator
- Peer worldviews from the network
- Served orders from the FSM
- Peer connection/disconnection updates
- Heartbeat for broadcasting to the network

### lights.go
Sets button light states based on current order states.

# How to run




