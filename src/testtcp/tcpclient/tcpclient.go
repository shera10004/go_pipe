package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const conn_tcp = "tcp"
const NetLocalHost = "127.0.0.1"
const NetPort = 8000

const (
	NET_CLOSE = iota
	NET_TRYCONNECT
	NET_SUCCESS
)

type GSocket struct {
	session     net.Conn
	connState   int
	rBufferSize int
	readBuffer  []byte
	userMessage func([]byte, int)
	mu          sync.Mutex
}

func (s *GSocket) Init() {
	s.connState = NET_CLOSE
	s.rBufferSize = 4096
	s.readBuffer = make([]byte, s.rBufferSize)
}
func (s *GSocket) Close() {
	s.mu.Lock()
	if s.connState == NET_SUCCESS {
		s.session.Close()
	}
	s.connState = NET_CLOSE
	s.mu.Unlock()
}
func (s *GSocket) MessageCallback(f func(msg []byte, size int)) {
	s.userMessage = f
}
func (s GSocket) Request(packet []byte) {
	_, err := s.session.Write(packet)
	if err != nil {
		fmt.Println("Request error :", err)
		s.Close()
		return
	}
}
func (s GSocket) ConnectionState() int {
	return s.connState
}
func (s *GSocket) Connect_to_server(ip string, port int) {
	s.connState = NET_TRYCONNECT

	addr := fmt.Sprintf("%s:%d", ip, port)
	fmt.Println("connect to", addr)

	var err error
	s.session, err = net.Dial(conn_tcp, addr)
	if err != nil {
		fmt.Println("Dial_err", err)
		s.connState = NET_CLOSE
	}

	fmt.Println("connect ok!")
	s.connState = NET_SUCCESS
	go func(gs *GSocket) {
		for {
			size, err := gs.session.Read(gs.readBuffer)
			if err != nil {
				fmt.Println("Read error :", err)
				gs.Close()
				return
			}
			//fmt.Println("recv size :", size)
			s.userMessage(s.readBuffer, size)

		}
	}(s)

	input_msg := "hello_"
	s.Request([]byte(input_msg))

	return
}

func run() {
	socket := GSocket{}
	socket.Init()
	counter := 0
	defer socket.Close()

	for {

		cState := socket.ConnectionState()

		if cState == NET_CLOSE {
			fmt.Println("socket is closed")
			socket.Connect_to_server(NetLocalHost, NetPort)
			socket.MessageCallback(func(msg []byte, size int) {
				recvmsg := string(msg[:size])
				_ = recvmsg
				//fmt.Println("recvmsg:", recvmsg)

				rmsg := fmt.Sprintf("good_%v", counter)
				counter++
				socket.Request([]byte(rmsg))
			})

			timesleep(0.5)
			continue
		} else if cState == NET_TRYCONNECT {
			fmt.Println("socket is NET_TRYCONNECT")
			timesleep(0.05)
			continue
		}

		timesleep(0.05)
	}
}

func main() {

	for i := 0; i < 2000; i++ {
		go run()
	}

	for {
		time.Sleep(100 * time.Second)
	}
}

func timesleep(t float32) {
	time.Sleep(time.Duration(t * float32(time.Second)))
}
