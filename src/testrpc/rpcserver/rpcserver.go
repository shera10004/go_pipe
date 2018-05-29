package main

import (
	"fmt"
	"net"
	"net/rpc"
)

type Calc int // RPC Server에 등록하기 위해 임의의 타입으로 정의.

//매개변수 구조체
//구조체의 접근 제한자는 반드시 public(대문자)로 명시 해야 한다.
type Args struct {
	A, B int
}

//리턴값 구조체
//구조체의 접근 제한자는 반드시 public(대문자)로 명시 해야 한다.
type Reply struct {
	C int
}

// rpc서버에 함수를 등록하려면 함수만으로는 안 되고, 구조체나 일반 자료형과 같은 타입에
// 메서드 형태로 구성되어 있어야 합니다. 여기서는 Calc 타입을 int 형으로 정의했는데, 다른
// 자료형이나 빈 구조체로 정의해도 됩니다.
func (c *Calc) Sum(args Args, reply *Reply) error {
	fmt.Println("process Sum function", args.A, args.B)
	reply.C = args.A + args.B
	return nil
}

func main() {
	rpc.Register(new(Calc)) //Calc 타입의 인스턴스를 생성하여 RPC 서버에 등록

	ln, err := net.Listen("tcp", ":6000") //TCP 프로토콜에 6000번 포트로 연결 받음.
	if err != nil {
		fmt.Println("Listen err :", err)
		return
	}
	defer ln.Close()

	fmt.Println("RPC Server run...")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept err :", err)
			continue
		}
		defer conn.Close()

		fmt.Println("Accept client :", conn)
		go rpc.ServeConn(conn) //RPC를 처리하는 함수를 고루틴을 실행
	}
}
