package main

import (
	"fmt"
	"net"
)

func requestHandler(c net.Conn) {
	data := make([]byte, 4096)

	for {
		read_bytes, read_err := c.Read(data)
		if read_err != nil {
			fmt.Println("read_err : ", read_err)
			return
		}

		fmt.Println(string(data[:read_bytes]))
		_, write_err := c.Write(data[:read_bytes])
		if write_err != nil {
			fmt.Println("write_err : ", write_err)
			return
		}
	}
}

func main() {

	//tcp 프로토콜에 8000포트로 연결을 받는다.
	//"198.168.0.1:8000" 처럼 IP 주소와 함께 설정하면 특정 NIC에서만 연결을 받는다.
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("listen_err : ", err)
		return
	}
	defer ln.Close()

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
}
