package containers

import (
	"container/heap"
	"container/list"
	"container/ring"
	"fmt"
	"testing"
	"time"
)

// 이중 연결리스트 테스트

func Test_linkedlist(t *testing.T) {
	elements := []*list.Element{}
	l := list.New() //연결 리스트 생성
	v1 := "message"

	elements = append(elements, l.PushBack(v1))
	elements = append(elements, l.PushBack(20))
	elements = append(elements, l.PushBack(1))
	elV1 := l.PushBack(100)
	elements = append(elements, elV1)
	elements = append(elements, l.PushBack(10))

	l.Remove(elV1)

	fmt.Println("Front", l.Front().Value)
	fmt.Println("Back", l.Back().Value)

	for i := l.Front(); i != nil; i = i.Next() { //연결 리스트의 맨 앞부터 끝까지 순회
		//fmt.Println(i.Value)
		switch i.Value.(type) {
		case string:
			fmt.Println(i.Value.(string))
		default:
			fmt.Println(i.Value)
		}
	}

	fmt.Println("-------")
	for i, v := range elements {
		fmt.Printf("[%d]element : %+v\n", i, v)

	}

}

// 힙 테스트

type MinHeap []int //힙을 int 슬라이스로 정의

func (h MinHeap) Len() int {
	return len(h) //슬라이스의 길이를 구함
}
func (h MinHeap) Less(i, j int) bool {
	r := h[i] < h[j] //대소관계 판단
	fmt.Printf("Less %d < %d %t \n", h[i], h[j], r)
	return r
}
func (h MinHeap) Swap(i, j int) {
	fmt.Printf("Swap %d %d\n", h[i], h[j])
	h[i], h[j] = h[j], h[i] //값의 위치 바꿈
}
func (h *MinHeap) Push(x interface{}) {
	fmt.Println("Push", x)
	*h = append(*h, x.(int)) //맨 마지막에 값 추가
}
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]     //슬라이스의 맨 마지막 값을 가져옴
	*h = old[0 : n-1] //맨 마지막 값을 제외한 슬라이스를 다시 저장
	return x
}

func Test_Heap(t *testing.T) {
	data := new(MinHeap)

	heap.Init(data)
	heap.Push(data, 5)
	heap.Push(data, 2)
	heap.Push(data, 7)
	heap.Push(data, 3)

	fmt.Println(data, "최솟값 :", (*data)[0])

}

//링 테스트
func Test_Ring(t *testing.T) {
	data := []string{"Maria", "John", "Andrew", "James"}

	r := ring.New(len(data))
	for i := 0; i < r.Len(); i++ {
		r.Value = data[i] //링 노드의 개수만큼 반복해서 값 넣기
		r = r.Next()      //다음노드 이동
	}

	r.Do(func(x interface{}) { //링의 모든 노드 순휘
		fmt.Println(x)
	})

	fmt.Println("Move forward:")
	r = r.Move(1) //링을 시계 방향으로 1노드 만큼 회전

	fmt.Println("Curr :", r.Value)
	fmt.Println("Prev :", r.Prev().Value)
	fmt.Println("Next :", r.Next().Value)

}

func f(x chan chan int) int {
	c := make(chan int) //int 형 채널을 생성
	fmt.Println("f() ... 1")
	time.Sleep(10 * time.Microsecond)
	fmt.Println("f() ... 2")
	x <- c //생성된 채널을 채널 x에 보냄
	fmt.Println("f() ... 3")
	//time.Sleep(1000 * time.Microsecond)
	re := <-c
	fmt.Println("f() ... 4")
	return re //채널 c에 들어온 값을 꺼내서 리턴
}
func Test_Channel(t *testing.T) {
	x := make(chan chan int) //int 형 채널을 보내는 채널
	fmt.Println("make chan - x")

	go func() {
		fmt.Println("go func()... in")
		c := <-x //x에서 채널을 꺼냄.
		fmt.Println("go func()... 2")
		c <- 11 //채널 c에 값을 보냄.
		fmt.Println("go func()... 3")
	}()

	fmt.Println("result :", f(x)) // 고루틴에서 보낸값 출력
}

func Test_Scanf(t *testing.T) {

	fmt.Println("------ input")

	inStr := ""
	fmt.Scanf("%s", inStr)

	fmt.Println("------ result :", inStr)

}
