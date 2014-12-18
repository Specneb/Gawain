// GawainBroker project main.go
package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
)

const (
	reqPort = "5555"
)

func main() {
	sock, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		panic(err)
	}
	defer sock.Close()

	url := fmt.Sprintf("tcp://*:%s", reqPort)
	if err := sock.Bind(url); err != nil {
		panic(err)
	}

	for {
		s, _ := sock.Recv(0)
		if s != "" {
			fmt.Println(s)
			sock.Send("PONG", 0)
		}
	}
}
