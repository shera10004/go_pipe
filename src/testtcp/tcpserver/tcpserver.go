package main

import (
	"bufio"
	"fmt"
	"net"
	"runtime"
	"time"
)

var clientCounter int64

func requestHandler(c net.Conn) {
	data := make([]byte, 4096)
	bf := bufio.NewReader(c)

	for {

		read_bytes, err := bf.Read(data)
		if err != nil {
			fmt.Println("read_err : ", err)
			break
		}

		/*
			read_bytes, read_err := c.Read(data)
			if read_err != nil {
				fmt.Println("read_err : ", read_err)
				break
			}
			//*/

		//fmt.Println(fmt.Sprintf("%d", &c), string(data[:read_bytes]))
		_, write_err := c.Write(data[:read_bytes])
		if write_err != nil {
			fmt.Println("write_err : ", write_err)
			break
		}
	}

	closeConn(c)
}

func closeConn(c net.Conn) {
	c.Close()

	clientCounter--
	//atomic.AddInt64(&clientCounter, -1)
}

func acceptClient(id int, ln net.Listener) {
	for {
		fmt.Println("wait client...")

		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept_err : ", err)
			continue
		}
		defer conn.Close()

		//atomic.AddInt64(&clientCounter, 1)
		clientCounter++

		fmt.Println("id:", id, "client Accepted :", conn)
		go requestHandler(conn)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//tcp 프로토콜에 8000포트로 연결을 받는다.
	//"198.168.0.1:8000" 처럼 IP 주소와 함께 설정하면 특정 NIC에서만 연결을 받는다.
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("listen_err : ", err)
		return
	}
	defer ln.Close()

	clientCounter = 0

	for i := 0; i < runtime.NumCPU(); i++ {
		go acceptClient(i, ln)
	}

	for {
		time.Sleep(1000 * time.Millisecond)
		//acounter := atomic.LoadInt64(&clientCounter)

		acounter := clientCounter
		fmt.Println("client count:", acounter)
	}
	/*
		for {
			fmt.Println("wait client...")

			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("Accept_err : ", err)
				continue
			}
			defer conn.Close()

			fmt.Println("client Accepted :", conn)
			go requestHandler(conn)
		}

		//*/
}
