// GawainBroker project main.go
package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
)

const (
	reqPort = "5555"
	monPort = "5556"
)

func Monitor(state zmq.State) {
	fmt.Println("State: ", state)
}

func main() {
	req, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		panic(err)
	}
	defer req.Close()

	url := fmt.Sprintf("tcp://*:%s", reqPort)
	if err := req.Bind(url); err != nil {
		panic(err)
	}

	mon, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		panic(err)
	}
	defer mon.Close()

	url := fmt.Sprintf("tcp://*:%s", monPort)
	if err := mon.Bind(url); err != nil {
		panic(err)
	}

	monitor = zmq.Reactor()
	monitor.AddSocket(mon, zmq.POLLIN, Monitor)
	monitor.Run(time.Second)

	/*
		for {
			s, _ := mon.Recv(0)
			if s != "" {
				fmt.Println(s)
				sock.Send("PONG", 0)
			}
		}
	*/
}
