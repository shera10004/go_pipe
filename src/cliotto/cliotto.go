package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/robertkrimen/otto"
)

var (
	Otto = otto.New()

	exitOtto = make(chan struct{})
)

func main() {
	fmt.Println("otto\n")
	loop()
}

func dispatch(entry string) string {
	if len(entry) == 0 {
		return entry
	}

	value, err := Otto.Run(entry)

	fmt.Println("-----------------------------------")
	if err != nil {
		fmt.Printf("%v", value)
		fmt.Println("error :", err.Error())

	} else {
		if value.IsUndefined() {
			return ""
		}
		if value.IsString() {
			if value.String() == "exit" {
				exitOtto <- struct{}{}
				return ""
			}
		}

		fmt.Printf("dispatch> result type-%v\n", value.TypeString())
		//fmt.Printf("dispatch> result value : %+#v\n", value.String())
		fmt.Printf("dispatch> result value : %v\n", value)
	}
	return ""
}

func loop() {

	Otto.Set("exit", "exit")

	Otto.Set("log", func(call otto.FunctionCall) otto.Value {
		fmt.Printf("log> %v\n", call.ArgumentList)

		v, _ := Otto.ToValue(call.ArgumentList)
		return v

	})

	Otto.Set("sum", func(call otto.FunctionCall) otto.Value {
		sum := int64(0)
		for _, v := range call.ArgumentList {
			iv, _ := v.ToInteger()
			sum += iv
		}

		v, _ := Otto.ToValue(sum)
		return v
	})

	Otto.Set("pow2", func(call otto.FunctionCall) otto.Value {
		v, _ := call.Argument(0).ToInteger()
		vv := v * v
		result, _ := Otto.ToValue(vv)
		return result
	})

	_, _ = Otto.Run(`
		abc = 2 + 2;
		console.log("The value of abc is " + abc);	//4
	`)

	Otto.Set("member", "This value is member.")

	Otto.Run(`
			console.log("exit :" + exit)
			console.log("abc :" + abc)
			console.log("member length:" + member.length)
			log(123)
			`)

	{
		value, _ := Otto.Run("member.length")
		v, _ := value.ToInteger()
		fmt.Println("v :", v)
	}

	if value, err := Otto.Get("abc"); err == nil {
		if value_int, err := value.ToInteger(); err == nil {
			fmt.Println("value_int:", value_int, err)
		}
	}

	//////////////////////////////////////////////

	go func() {
		defer fmt.Println("Exit Otto !!")
	EXIT:
		for {
			select {
			case <-exitOtto:
				goto EXIT
			default:
				fmt.Print("> ")
				in := bufio.NewReader(os.Stdin)
				entered, err := in.ReadString('\n')
				if err != nil {
					fmt.Println(err)
					break
				}
				entry := strings.TrimLeft(entered[:len(entered)-1], "\t ") // without tabs,spaces and newline
				_ = dispatch(entry)

				in.Reset(in)
			} //select
		} //for
	}()

	<-exitOtto
	fmt.Println("Program OUT")
	/*
		for {
			time.Sleep(time.Second)
		}
		//*/
}
