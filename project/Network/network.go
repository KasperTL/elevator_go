package Network

import (
	"encoding/json"
	"net"
	"time"
	"project/config"
	"project/WorldView"
)

const bufsize = config.Buffer


func WordlView_broadcast(tx <- chan WorldView.WorldView, port int) error {
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{IP: net.IPv4bcast, Port: port})

	if err != nil {
		return err
	}

	defer conn.Close()

	HeartbeatTimer := time.NewTicker(config.HeartbeatTime) 
    defer HeartbeatTimer.Stop()

    var last WorldView.WorldView
    hasState := false

    for {
        select {

        case view := <-tx:
            last = view
            hasState = true

        case <-HeartbeatTimer.C:
            if hasState {
                data, err := json.Marshal(last)
    		if err != nil {
        		continue
    		}
    		conn.Write(data)
            }
        }
    }
}




func Reciever(rx chan <- WorldView.WorldView, port int) error {

	addr := net.UDPAddr{IP: net.IPv4zero, Port: port}

	conn, err := net.ListenUDP("udp4",&addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	buf := make([]byte,config.Buffer)

	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		var view WorldView.WorldView 
		err = json.Unmarshal(buf[:n], &view)	
		if err != nil {
			continue
		}
		rx <- view 
	}
}

