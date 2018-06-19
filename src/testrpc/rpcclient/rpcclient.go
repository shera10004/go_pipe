package main

import (
	"fmt"
	"net/rpc"
)

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

func main() {

	client, err := rpc.Dial("tcp", "127.0.0.1:6000")
	if err != nil {
		fmt.Println("Dial err :", err)
		return
	}
	defer client.Close()

	//동기호출
	args := Args{1, 2}
	reply := new(Reply)
	err = client.Call("Calc.Sum", args, reply)
	if err != nil {
		fmt.Println("Call err:", err)
		return
	}
	fmt.Println("Call result:", reply.C)

	//비동기 호출
	args.A = 4
	args.B = 9

	//마지막 매개변수는 함수 실행이 끝났는지 확인하기 위한 채널입니다.
	//여기에 nil을 넣으면 채널이 새로 할당되어 리턴됩니다.
	sumCall := client.Go("Calc.Sum", args, reply, nil)

	<-sumCall.Done //함수가 끝날 때까지 대기

	fmt.Println("Call result2:", reply.C)

}
