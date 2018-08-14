package externalip

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"syscall"
	"testing"

	"github.com/jackpal/gateway"
	natpmp "github.com/jackpal/go-nat-pmp"
)

func Test_IP(t *testing.T) {

	type result struct {
		id  int
		ip  string
		err error
	}

	counter := 2
	ips := make(chan result, counter)
	cok := make(chan int)

	for i := 0; i < counter; i++ {
		go func(id int) {
			ip, err := ExternalIP()
			gr := result{id: id, ip: ip, err: err}
			ips <- gr
		}(i)
	}

	go func() {
		okCounter := 0
		for {
			okCounter = okCounter + <-cok

			if okCounter >= counter {
				close(ips)
			}
		}
	}()

	for r := range ips {
		if r.err != nil {
			fmt.Println(r.id, ">", r.err)
		} else {
			fmt.Println(r.id, "> myip:", r.ip)
		}
		cok <- 1
	}

	//goupnp.NewServiceClients
}
func DiscoverGateway() net.IP {
	routeCmd := exec.Command("route", "-4", "print", "0.0.0.0")
	parseWindowsRoutePrint := func(output []byte) net.IP {
		var (
			networktarget string
			networkmask   string
			gateway       string
			privateip     string
		)
		_ = networktarget
		_ = networkmask
		_ = gateway
		_ = privateip

		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "0.0.0.0") {
				datas := strings.Split(line, " ")
				index := 0
				for _, data := range datas {
					if data != "" {
						switch index {
						case 0:
							networktarget = data
						case 1:
							networkmask = data
						case 2:
							gateway = data
						case 3:
							privateip = data
						} //switch
						index++
					} //if
				} //for
			} //if
		} //for
		return net.ParseIP(gateway)
	} //end func
	routeCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := routeCmd.CombinedOutput()
	if err != nil {
		return nil
	}

	return parseWindowsRoutePrint(output)
}
func Test_IP2(t *testing.T) {
	gatewayIP, err := gateway.DiscoverGateway()
	if err != nil {
		fmt.Println("err:", err)

		gatewayIP = DiscoverGateway()
		if gatewayIP == nil {
			fmt.Println("gateway ip is nil")
			return
		}
		return
	}

	fmt.Println("gatewayIP:", gatewayIP)

	client := natpmp.NewClient(gatewayIP)

	result, err := client.AddPortMapping("udp", 30303, 30303, 0)
	if err != nil {
		return
	}
	fmt.Println("InternalPort", result.InternalPort)
	fmt.Println("MappedExternalPort", result.MappedExternalPort)
	fmt.Println("SecondsSinceStartOfEpoc", result.SecondsSinceStartOfEpoc)

	response, err := client.GetExternalAddress()

	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Printf("External IP address: %#v\n", response.ExternalIPAddress)
}
