// GawainClient project main.go
package main

import (
	nc "code.google.com/p/goncurses"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"net"
	"os"
	"time"
)

const (
	brokerName = "riak.mail.tiscali.sys"
	brokerPort = "5555"
	brokerMon  = "5556"
	keepalive  = 1000
)

var brokerList []string

type broker struct {
	mon    *zmq.Socket
	sock   *zmq.Socket
	id     int
	ip     string
	status bool
}

var brokers = make(map[string]*broker)

func getBrokerList() {
	var err error
	brokerList, err = net.LookupHost(brokerName)
	if err != nil {
		fmt.Println("ERROR: Brokers List unavailable")
		os.Exit(-1)
	}
}

func monitor(s *broker) {
	for {
		time.Sleep(keepalive * time.Millisecond)
		s.status = false
		ret, err := sendMsg(s.mon, "PING")
		if err != nil {
			fmt.Println("Recv Error: PING")
			continue
		}
		s.status = true
		fmt.Println("OK: ", ret)
	}
}

func sendMsg(s *zmq.Socket, msg string) (retval string, err error) {
	(*s).Send(msg, 0)
	retval, err = (*s).Recv(0)
	if err != nil {
		fmt.Println("SendMsg Error: ", msg, err)
		return retval, err
	}

	return retval, nil
}

func main() {
	getBrokerList()
	for ind, i := range brokerList {
		sock, err := zmq.NewSocket(zmq.REQ)
		if err != nil {
			fmt.Println(err)
			continue
		}
		url := fmt.Sprintf("tcp://%s:%s", i, brokerPort)
		if errs := sock.Connect(url); errs != nil {
			fmt.Println("Connect Error " + url)
			sock.Close()
			continue
		}
		defer sock.Close()

		mon, err := zmq.NewSocket(zmq.REQ)
		if err != nil {
			fmt.Println(err)
			continue
		}
		url = fmt.Sprintf("tcp://%s:%s", ind, i, brokerMon)
		if errs := mon.Connect(url); errs != nil {
			fmt.Println("Connect Error " + url)
			mon.Close()
			continue
		}
		defer mon.Close()

		brokers[i] = &broker{mon, sock, ind, i, true}

		go monitor(brokers[i])
	}
	fmt.Println(brokers)

	for {
		for ip, b := range brokers {
			if b.status {
				sendMsg(b.sock, "LOGLOG "+ip)
				break
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}
