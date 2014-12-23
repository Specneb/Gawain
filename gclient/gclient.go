// GawainClient project main.go
package main

import (
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
var brokers = make(map[*zmq.Socket]bool)

func getBrokerList() {
	var err error
	brokerList, err = net.LookupHost(brokerName)
	if err != nil {
		fmt.Println("ERROR: Brokers List unavailable")
		os.Exit(-1)
	}
}

func monitor(s *zmq.Socket) {
	for {
		time.Sleep(keepalive * time.Millisecond)
		brokers[s] = false
		ret, err := sendMsg(s, "PING")
		if err != nil {
			fmt.Println("Recv Error: PING")
			continue
		}
		brokers[s] = true
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
	for _, i := range brokerList {
		sock, err := zmq.NewSocket(zmq.REQ)
		if err != nil {
			fmt.Println(err)
			continue
		}
		url := fmt.Sprintf("tcp://%s:%s", i, brokerMon)
		if errs := sock.Connect(url); errs != nil {
			fmt.Println("Connect Error " + url)
			sock.Close()
			continue
		}
		brokers[sock] = true
		defer sock.Close()
		go monitor(sock)
	}
	fmt.Println(brokers)

	for {
		for sock, stat := range brokers {
			if stat {
				sendMsg(sock, "LOGLOG LOGLOG")
				break
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}
