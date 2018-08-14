package main

import (
	"fmt"

	"github.com/jackpal/gateway"
	natpmp "github.com/jackpal/go-nat-pmp"
)

func main() {
	gatewayIP, err := gateway.DiscoverGateway()
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	client := natpmp.NewClient(gatewayIP)
	response, err := client.GetExternalAddress()
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Printf("External IP address: %#v\n", response.ExternalIPAddress)
}
