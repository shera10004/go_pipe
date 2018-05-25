package main

import (
	"fmt"
	"runtime"
	"time"
)

func timesleep(t float32) {
	time.Sleep(time.Duration(t * float32(time.Second)))
}

func main() {
	//runtime.GOMAXPROCS(1)
	runtime.GOMAXPROCS(runtime.NumCPU())

	done := make(chan bool, 2)
	count := 4

	go func() {
		for i := 0; i < count; i++ {
			done <- true
			fmt.Println("Go func :", i)
			timesleep(0.00001)
		}
	}()

	for i := 0; i < count; i++ {
		<-done
		fmt.Println("main func :", i)
	}

	fmt.Println(" ---- exit ---- ")
}
