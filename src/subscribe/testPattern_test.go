package subscribe

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"math/rand"
)

type Subscription struct {
	mu   sync.Mutex
	stop chan<- bool
}

func (s *Subscription) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stop != nil {
		close(s.stop)
	}
	s.stop = nil
}

type processer struct {
	stopChan <-chan bool
	callback func(s string)
}
type sbManager struct {
	mu      sync.RWMutex
	clients []*processer
}

func (sb *sbManager) Subscribe(foo func(s string)) *Subscription {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	stop := make(chan bool)

	pc := &processer{
		stopChan: stop,
		callback: foo,
	}

	sb.clients = append(sb.clients, pc)

	go func() {
		<-pc.stopChan
		sb.mu.Lock()
		for i, v := range sb.clients {
			if v == pc {

				fmt.Println("stopChan", i)
				sb.clients = append(sb.clients[:i], sb.clients[i+1:]...)

				break
			}
		} //for
		sb.mu.Unlock()
	}()

	return &Subscription{
		stop: stop,
	}
}

func (sb *sbManager) Notification(msg string) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	if len(sb.clients) == 0 {
		fmt.Println("empty")
		return
	}

	for _, v := range sb.clients {
		v.callback(msg)
	}
}

func TestSubscribePattern(t *testing.T) {
	sb := sbManager{
		clients: []*processer{},
	}

	re1 := sb.Subscribe(func(s string) {
		fmt.Println("foo1 :", s)
	})
	re2 := sb.Subscribe(func(s string) {
		fmt.Println("foo2 :", s)
	})
	re3 := sb.Subscribe(func(s string) {
		fmt.Println("foo3 :", s)
	})

	go func() {
		time.Sleep(time.Millisecond * 1510)
		_ = re1
		//re1.Stop()
	}()

	go func() {
		time.Sleep(time.Millisecond * 1510)
		_ = re2
		//re2.Stop()
	}()

	go func() {
		time.Sleep(time.Millisecond * 2510)
		_ = re3
		//re3.Stop()
	}()

	wg := sync.WaitGroup{}

	notiFunc := func(id int, wg *sync.WaitGroup) {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			//time.Sleep(time.Millisecond * 500)
			fmt.Println("[", id, "] notification", i)
			sb.Notification(fmt.Sprintf("[%v]%v", id, i))
		} //for
	}

	loopCount := 5
	wg.Add(loopCount)
	for i := 0; i < loopCount; i++ {
		go notiFunc(i+1, &wg)
	}
	wg.Wait()

	re1.Stop()
	re2.Stop()
	re3.Stop()

	fmt.Println("end!!")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	_ = r

	for i := 0; i < 10; i++ {
		fmt.Printf("%v ", r.Intn(100))
	}

}
