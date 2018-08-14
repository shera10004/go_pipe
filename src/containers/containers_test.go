package containers

import (
	"bytes"
	"container/heap"
	"container/list"
	"container/ring"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"vendorpkg/abc"
)

func Test_Ip(t *testing.T) {
	conn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(conn.LocalAddr().String())

	smap := map[int]*string{}
	rValue := "aaa"
	smap[1] = &rValue
	smap[2] = &rValue

	rValue = "bbb"

	for k, v := range smap {
		fmt.Printf("k:%v , v:%v\n", k, *v)
	}
}

/////////////// 임베딩 test ////////////////////////

func Test_Dummy0(t *testing.T) {
	sl := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	pf := func(ss []int) {
		for i := 1; i < len(ss); i++ {
			fmt.Printf("%#v\n", ss)
		}

	}

	for i := range sl {
		pf(sl[i : i+1])
	}

}

func Test_Dummy(t *testing.T) {

	const MAXCOUNT = 3

	resultEntries := []string{
		"a",
		"b",
		"c",
		"d",
	}
	pendingQueries := 0
	reply := make(chan string, MAXCOUNT)

	for {
		for i := 0; i < len(resultEntries) && pendingQueries < MAXCOUNT; i++ {
			n := resultEntries[i]
			pendingQueries++
			go func(v string) {
				//fmt.Println("value :", v)
				time.Sleep(time.Millisecond * 500)
				reply <- v
			}(n)
		} //for

		if pendingQueries == 0 {
			break
		}

		re := <-reply
		if len(resultEntries) < 10 {
			resultEntries = append(resultEntries, re)
			fmt.Printf("%#v\n", resultEntries)
		} else {
			resultEntries = nil
		}

		pendingQueries--

	} //for

}

/////////////// 채널 select close ////////////////////////

func Test_SelectClose(t *testing.T) {

	cDone := make(chan bool)
	cDone2 := make(chan bool)
	isOut := false

	go func() {
		time.Sleep(time.Second * 60)
		isOut = true
	}()

	// sleep and done channel!
	go func() {
		time.Sleep(time.Second * 3)
		fmt.Println("close chanell <-")
		cDone <- true

		time.Sleep(time.Second * 6)
		fmt.Println("close chanell")
		close(cDone)
	}()
	getResult := func() string {

		select {
		case <-cDone2:
			return "cDone2"

		case _, ok := <-cDone:
			if ok == false {
				return "cDone close"
			}
			return "cDone true"
		default:
			return "default"
		}
	}

	for {
		re := getResult()

		fmt.Println("re :", re)
		if isOut == true {
			break
		}

		time.Sleep(time.Millisecond * 300)
	}

}

/////////////// 채널 select 2 ////////////////////////

func Test_SelectInn(t *testing.T) {
	sc := SelectConn{
		addpending: make(chan *string),
		closing:    make(chan error),
	}

	go func() {
		input := ""
		fmt.Scanf("%v", input)
	}()

	go func() {
		fmt.Println("wait pending")

		pValue := <-sc.addpending
		fmt.Println("pending value:", *pValue)
		sc.callback.print()
		sc.callback.closing <- nil

	}()

	time.Sleep(time.Second * 2)

	go func() {
		fmt.Println("wait 5 time")
		time.Sleep(time.Second * 1)
		fmt.Println("closing chan<-")
		sc.closing <- errClosed
	}()
	//*/

	time.Sleep(time.Second * 2)

	cEnd := sc.pending(10, func() {
		fmt.Println("callback", sc.id)
	})

	end := <-cEnd

	fmt.Println("main end :", end)
}

type SelectConn struct {
	id         int
	callback   *sCallBack
	addpending chan *string
	closing    chan error
}
type sCallBack struct {
	closing chan error
	print   func()
}

var errClosed = errors.New("app closed")

func (sc *SelectConn) pending(id int, callback func()) <-chan error {
	ch := make(chan error, 1)

	sc.id = id
	sc.callback = &sCallBack{
		closing: ch,
		print:   callback,
	}

	pd := "pending"

	fmt.Println("in pending")
	select {
	case <-sc.closing:
		fmt.Println("end err")
		ch <- errClosed

	case sc.addpending <- &pd:

	}
	return ch
}

/////////////////////////////////////////////////////
func Test_SelectChannel(t *testing.T) {

	cc := make(chan string)
	dd := make(chan string)
	go func(a chan string, b chan string) {
		for {
			select {
			case v, ok := <-a:
				if ok {
					fmt.Println("Acc :", v)
				}

			case v := <-b:
				fmt.Println("dd :", v)
			}

			fmt.Println("default")
			time.Sleep(time.Millisecond * 100)
		}

	}(cc, dd)

	go func(x chan string) {
		v := <-x
		fmt.Println("Bcc :", v)
	}(cc)

	//cc <- "aaa"
	timer := time.NewTimer(time.Second * 2)

	<-timer.C

	cc <- "cccc"

	isStop := timer.Stop()
	fmt.Println("isStop:", isStop)

	dd <- "ddd"

	//close(cc)

	time.Sleep(time.Second * 3)

	fmt.Println("isStop:", isStop)
	fmt.Println("end")

}

// MAP Test
func inputMap(m map[string]int, key string, v int) {
	m[key] = v
}
func Test_Map(t *testing.T) {
	cMap := map[string]int{}

	inputMap(cMap, "a", 1)
	inputMap(cMap, "b", 2)
	inputMap(cMap, "c", 3)

	fmt.Printf("%#v \n", cMap)
}

func Test_Switch(t *testing.T) {
	a := 0
	b := "b"
	c := true

	switch {
	case a == 0:
		fmt.Println("a ok")
	case b == "b":
		fmt.Println("b is ok")
	case c == false:
		fmt.Println("c is ok")
	default:
		fmt.Println("default")

	}
}

func Test_Closer(t *testing.T) {
	funcA := func(s string) func(int) string {
		re := fmt.Sprintf("func:%v", s)
		return func(v int) string {
			return fmt.Sprintf("%s -> %v", re, v)
		}
	}

	subA := funcA("A")
	fmt.Println("subA:", subA(18))

	fmt.Println("subA:", funcA("B")(20))

	funcB := func(i int, foo func(int)) func(int) {
		teni := i + 10
		return func(j int) {
			foo(teni + j)
		}
	}

	resultB := funcB(10, func(v int) {
		fmt.Println("resultB :", v)
	})
	resultB = funcB(10, resultB)
	resultB = funcB(10, resultB)
	resultB = funcB(10, resultB)

	resultB(10)
}

// vendor 폴더 패키지 임포트 테스트
func Test_VendorPkg(t *testing.T) {
	fmt.Println(abc.GetString())
}

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

// 슬라이스 부분 삽입
func Test_SliceAppend(t *testing.T) {
	sa := []string{
		"a", "b", "c",
	}
	sb := []string{
		"d", "e", "f",
	}

	//sa = append(sa[:0], sb...)
	sa = append(sa[:1], sb...)

	fmt.Printf("%#v\n", sa)
}

// 슬라이스 삭제
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

func TestChanFor(t *testing.T) {

	//wg := sync.WaitGroup{}
	chs := make(chan int) //[]chan string

	go func() {
		defer func() {
			close(chs)
			fmt.Println("end func1")
		}()

		for i := 0; i < 10; i++ {
			chs <- i //fmt.Sprintf("%v", i)
			time.Sleep(time.Millisecond * 100)
		} //for

	}()

EXIT:
	for {
		select {
		case v, ok := <-chs:
			fmt.Println("chain value:", v, ok)
			if ok == false {
				break EXIT
			}
		} //select
	} //for

	//wg.Wait()

	fmt.Println("END")
}

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

///////////// 슬라이스 복사 //////////////////
func TestSliceCopy(t *testing.T) {
	sl := []string{
		"a", //0
		"b", //1
		"c", //2
		"d", //3
	}

	copy(sl[1:], sl[:1]) // 0,1번을 1번인덱스로 복사
	//sl[0] = "x"

	fmt.Printf("%#v\n", sl)

	a := 1
	b := 0
	v := a ^ b
	fmt.Println("value :", v)
}

func TestSliceRef(t *testing.T) {
	sl := []byte{1, 2, 3, 4}

	sl2 := []byte{}
	sl2 = sl

	sl2[1] = 10 //주소 값은 다르지만 같은 데이타를 가리킨다

	fmt.Printf("&sl %p \n", &sl)
	fmt.Printf("&sl2 %p \n", &sl2)
	fmt.Println("----")
	fmt.Printf("sl value : %v \n", sl)

	//슬라이스 복사는 이렇게...
	//sl2 := make([]byte, len(sl))
	//copy(sl2, sl)

	//슬라이스 값 모두 비우기... nil을 대입하면 모두 클리어 된다.
	//sl = nil

}

//////////////////////////////////////////////////////////////////////////

func TestInterfaceEmbeding(t *testing.T) {
	testA := &printA{}
	testB := printB{}

	prints := []iPrint{testA, &testB}
	for i, v := range prints {
		val := strconv.Itoa(i)
		v.Print(val)
		v.Input(val + "- xxx")
		v.Show()
	}

	fmt.Println("-----------------")

	sets := []iSet{testA, testB}
	testA.ISet("a", func(s string) {
		testA.name = s
	})
	testB.ISet("b", func(s string) {
		testB.id = s
	})
	for _, v := range sets {
		v.Show()
	}

}

type (
	iPrint interface {
		Input(s string)
		Print(s string)
		Show()
	}
	iSet interface {
		ISet(s string, f func(s string))
		Show()
	}

	printA struct {
		name string
	}

	printB struct {
		id string
	}
)

func (this_printA printA) XX() {
	x := this_printA.name
	_ = x
}
func (c *printA) Input(s string) {
	c.name = s
}
func (c printA) Input2(s string) {
	c.name = s
}

func (c printA) Print(s string) {
	fmt.Println("printA :", s)
}

func (c printA) ISet(s string, f func(string)) {
	f(s)
}
func (c printA) Show() {
	fmt.Println("PrintA :", c.name)
}

func (c *printB) Input(s string) {
	c.id = s
}

func (c printB) Print(s string) {
	fmt.Println("PrintB :", s)
}

func (c printB) ISet(s string, f func(string)) {
	f(s)
}
func (c printB) Show() {
	fmt.Println("PrintB :", c.id)
}

//////////////////////////////////////////////////////////
func TestStringToByte(t *testing.T) {
	//str := "eth.accounts\n"
	str := "테스터"
	src := []byte(str)

	fmt.Printf("%v\n", string([]byte{0x65, 0x74, 0x68, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0xa}))
	fmt.Printf("%v\n", string([]byte{0x4a, 0x53, 0x4f, 0x4e}))
	fmt.Printf("%v\n", string([]byte{0x28, 0x7b, 0x22, 0x6a, 0x73, 0x6f, 0x6e, 0x72, 0x70, 0x63, 0x22, 0x3a, 0x22, 0x32, 0x2e, 0x30, 0x22, 0x7d, 0x29}))
	fmt.Printf("----------------- \n")
	//fmt.Printf("%v , %v\n", src, len(bytes.Runes(src)))
	fmt.Printf("%v , %v\n", string(src), len([]rune(str)))

	jstrs := []string{"ab", "cd", "ef"}
	fmt.Printf("%v \n", strings.Join(jstrs, "-"))

	buffer := bytes.Buffer{}
	for _, v := range jstrs {
		buffer.WriteString(v)
	}
	fmt.Println("buffer:", buffer.String())

	var sourcemapSource interface{}
	sourcemapSource = nil

	if sourcemapSource == nil {

		lines := bytes.Split(src, []byte("\n"))
		lastLine := lines[len(lines)-1]

		if bytes.HasPrefix(lastLine, []byte("//# sourceMappingURL=data:application/json")) {
			bits := bytes.SplitN(lastLine, []byte(","), 2)

			if len(bits) == 2 {
				if d, err := base64.StdEncoding.DecodeString(string(bits[1])); err == nil {
					sourcemapSource = d
				}
			}

		}
	}

	fmt.Printf("%v \n", sourcemapSource)

}

//////////////////////////////////////////////////////////

type ActionFunc func(int) int

func selectType(t interface{}) string {

	if _, ok := t.(func(int) int); ok == true {
		return "func(int) int"
	} else if _, ok := t.(ActionFunc); ok == true {
		return "ActionFunc"
	} else {
		return "nothing"
	}
}

func TestFuncType(t *testing.T) {
	var actionFunc ActionFunc
	var funcFunc func(int) int
	var number int32
	sl := []interface{}{actionFunc, funcFunc, number}

	for i, v := range sl {
		fmt.Println("id:", i, ", value:", selectType(v))
	}
}
