package containers

import (
	"container/heap"
	"container/list"
	"container/ring"
	"context"
	"fmt"
	"runtime"
	"sync"
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

type ttt struct {
	A int
	B int
	C string
}

func (x *ttt) foo() {
	x.A = 10
	x.C = "abc"
}

func Test_Scanf(t *testing.T) {

	c := make(chan int)
	//go func(cc chan<- int) {
	go func(cc <-chan int) {
		i := 0
		for {
			time.Sleep(1 * time.Second)
			fmt.Println("i:", i, ", c:", <-cc)
			i++
		}
		//cc <- 10
	}(c)

	x := &ttt{} //new(ttt)
	x.foo()

	Ttt(*x)
	fmt.Println(x)

	c <- 10

	c <- 20
	time.Sleep(1 * time.Second)
	c <- 30
}

func Ttt(param ttt) {
	param.A = 20

}

func Test_SliceRemove(t *testing.T) {
	sl := make([]string, 0)
	sl = append(sl, "a")
	sl = append(sl, "b")
	sl = append(sl, "c")

	removeindex := 1
	if len(sl) > removeindex {
		sl = append(sl[:removeindex], sl[removeindex+1:]...)
	}

	fmt.Println(len(sl))
	for i, v := range sl {
		fmt.Printf("%d  value:%v\n", i, v)
	}

}

func TestUint32(t *testing.T) {

	v := uint32(1) //65535)
	d := uint32(0x01000193)
	c := uint32(213)

	r := v*d ^ c

	fmt.Println(d)
	fmt.Println(r)

}

//////////////////////////////////////////////////////////////////////////

func longFuncWithCtx(ctx context.Context) (string, error) {
	select {
	default:
		return "success", nil
	case <-ctx.Done():
		re, ok := ctx.Value("key").(string)
		if ok == true {
			return re, ctx.Err()
		}
		return "false", ctx.Err()
	}
}

func TestContext(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	const (
		c_key  = "key"
		c_key2 = "key2"
	)

	ctx, cancel := context.WithTimeout(context.WithValue(context.Background(), c_key, "Done"), time.Second*3)
	//cancel()

	type ctxData struct {
		value string
	}

	ctx2, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Second*5))
	ctx2 = context.WithValue(ctx2, c_key, &ctxData{value: "우하하하"})
	ctx2 = context.WithValue(ctx2, c_key2, []string{"xxxx"})

	go func(c context.Context) {
		fmt.Println("wait ctx2 -----")
		loop := 0
		for {
			select {
			default:
				if v := c.Value(c_key); v != nil {
					if re, ok := v.(*ctxData); ok {
						re.value = fmt.Sprintf("우하하하%v", loop)
						if v := c.Value(c_key2); v != nil {
							if re2, ok := v.([]string); ok {
								re2[0] = fmt.Sprintf("sss%v", loop)
							}
						}

						loop++
						time.Sleep(time.Second)

						continue
					}
				}
				fmt.Println("Context Value ERROR")
				goto FUNC_OUT

			case <-(c).Done():
				fmt.Println("Done ctx2 -----")
				goto FUNC_OUT
			}
		}
	FUNC_OUT:
	}(ctx2)

	go func(c context.Context) {
		<-(c).Done()
		if v := c.Value(c_key); v != nil {
			if re, ok := v.(*ctxData); ok {
				fmt.Printf(" ctx2 key: %v \n", re)
			}
		}
		if v := c.Value(c_key2); v != nil {
			if re, ok := v.([]string); ok {
				fmt.Printf(" ctx2 key2: %v \n", re[0])
			}
		}

	}(ctx2)

	_ = ctx
	_ = cancel

	jobCount := 10
	var wg sync.WaitGroup
	for i := 0; i < jobCount; i++ {

		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := longFuncWithCtx(ctx)

			if err != nil {
				fmt.Println(i, ":", result, "/", err)
			} else {
				fmt.Println(i, ":", result, "/", err)
			}

		}()

		time.Sleep(time.Second)
	}

	wg.Wait()

}

//////////////////////////////////////////////////////////////////////////

func TestChanGetSet(t *testing.T) {

	runtime.GOMAXPROCS(runtime.NumCPU())

	about := make(chan struct{})
	cv := make(chan string)
	counter := 0
	wg := sync.WaitGroup{}

	wg.Add(1)
	//go func(a <-chan struct{}, c chan<- string) {
	go func() {
		defer wg.Done()
		//EXITA:
		for {
			select {
			case <-about:
				fmt.Println("about")
				close(cv)
				//break EXITA
				goto EXITB

			default:
				str := fmt.Sprintf("str:%v", counter)
				cv <- str
				time.Sleep(100 * time.Millisecond)
				counter++
			}

		}

	EXITB:
		fmt.Println(":: EXIT!!")
	}()
	//}(about, cv)

	//fmt.Println(about)
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			if view, ok := <-cv; ok {
				fmt.Println(view)
			} else {
				break
			}
		}
		//*/

		fmt.Println("func println out!")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if counter == 5 {
				about <- struct{}{}
				break
			}
		}
		fmt.Println("func count out!")
	}()

	wg.Wait()

	fmt.Println("func exit!")

}

//////////////////////////////////////////////////////////////////////////
type stref struct {
	name string
	id   int
	sl   []int
}

func (sr *stref) IncSlice(v int) {
	sr.sl = append(sr.sl, v)
}

func strefCall(s *stref, name string, id int) {
	s.name = name
	s.id = id
}
func appendSliceRef(sl *[]int, v int) {
	*sl = append(*sl, v)
}
func appendSliceVal(sl []int, v int) {
	sl = append(sl, v)
}

func appendStrefVal(sr *stref, v int) {
	sr.sl = append(sr.sl, v)
}
func addSliceVal(sl []int, v int) {
	sl = append(sl, v)
	for i, val := range sl {
		sl[i] = val + v
	}

	_ = sl

}
func TestRef(t *testing.T) {

	sr := stref{}
	sl := []int{}
	strefCall(&sr, "a", 1)
	appendSliceRef(&sl, 1)
	appendSliceRef(&sl, 2)
	appendSliceRef(&sl, 3)
	appendSliceVal(sl, 4)

	addSliceVal(sl, 10)

	appendStrefVal(&sr, 1)

	sr.IncSlice(2)

	fmt.Printf("stref -> %#v\n", sr)
	fmt.Printf("sl -> %#v\n", sl)

}
