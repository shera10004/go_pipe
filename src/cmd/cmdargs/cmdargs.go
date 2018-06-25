package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(os.Args)

	for i, arg := range os.Args {
		fmt.Println("[", i, "]", arg)
	}
}
