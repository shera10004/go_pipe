package main

import (
	"fmt"
	"net"
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
}

var socket GSocket

func (s *GSocket) Init() {
	s.connState = NET_CLOSE
	s.rBufferSize = 4096
	s.readBuffer = make([]byte, s.rBufferSize)
}
func (s *GSocket) Close() {
	if s.connState == NET_SUCCESS {
		s.session.Close()
	}
	s.connState = NET_CLOSE
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
	socket.connState = NET_TRYCONNECT

	addr := fmt.Sprintf("%s:%d", ip, port)
	fmt.Println("connect to", addr)

	var err error
	socket.session, err = net.Dial(conn_tcp, addr)
	if err != nil {
		fmt.Println("Dial_err", err)
		socket.connState = NET_CLOSE
	}

	socket.connState = NET_SUCCESS
	go func(s *GSocket) {
		for {
			size, err := s.session.Read(s.readBuffer)
			if err != nil {
				s.Close()
				return
			}
			fmt.Println("recv size :", size)
			s.userMessage(s.readBuffer, size)
		}
	}(s)

	return
}

func main() {

	socket = GSocket{}
	socket.Init()

	socket.Connect_to_server(NetLocalHost, NetPort)
	socket.MessageCallback(func(msg []byte, size int) {
		recvmsg := string(msg[:size])
		fmt.Println("recvmsg:", recvmsg)
	})

	defer socket.Close()

	w_count := 0
	for {

		cState := socket.ConnectionState()

		if cState == NET_CLOSE {
			fmt.Println("socket is closed")
			socket.Connect_to_server(NetLocalHost, NetPort)
			timesleep(0.5)
		} else if cState == NET_TRYCONNECT {
			continue
		}

		/*
			fmt.Println("command message :")
			var input_msg string
			_, sErr := fmt.Scanln(&input_msg)
			if sErr != nil {
				fmt.Println("Scanln error :", sErr)
				continue
			}
			fmt.Println("input_msg :", input_msg)
			//*/

		input_msg := "hello_"
		wmsg := fmt.Sprintf("%s _ %d", input_msg, w_count)
		w_count++

		socket.Request([]byte(wmsg))
		timesleep(0.05)

	}
}

func timesleep(t float32) {
	time.Sleep(time.Duration(t * float32(time.Second)))
}
