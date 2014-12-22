// GawainBroker project main.go
package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"time"
)

const (
	reqPort = "5555"
	monPort = "5556"
)

func Monitor(sock *zmq.Socket) error {
	fmt.Println(sock.Recv(0))
	return nil
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

	url = fmt.Sprintf("tcp://*:%s", monPort)
	if err := mon.Bind(url); err != nil {
		panic(err)
	}

	monitor := zmq.NewReactor()
	monitor.AddSocket(mon, zmq.POLLIN, func(e zmq.State) error { return Monitor(mon) })
	monitor.Run(200 * time.Millisecond)

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
