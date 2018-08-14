package testnatpmp

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/huin/goupnp"
	"github.com/huin/goupnp/dcps/internetgateway1"
	"github.com/huin/goupnp/dcps/internetgateway2"
)

const soapRequestTimeout = 3 * time.Second

func internalAddress(host string) net.IP {
	devaddr, err := net.ResolveUDPAddr("udp4", host)
	if err != nil {
		return nil
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil
		}
		for _, addr := range addrs {
			switch x := addr.(type) {
			case *net.IPNet:
				if x.Contains(devaddr.IP) {
					return x.IP
				}
			}
		}
	}
	return nil
}

type WanIPConnecter interface {
	GetExternalIPAddress() (NewExternalIPAddress string, err error)
	DeletePortMapping(NewRemoteHost string, NewExternalPort uint16, NewProtocol string) (err error)
	AddPortMapping(NewRemoteHost string, NewExternalPort uint16, NewProtocol string, NewInternalPort uint16, NewInternalClient string, NewEnabled bool, NewPortMappingDescription string, NewLeaseDuration uint32) (err error)

	GetServiceClient() *goupnp.ServiceClient
}

func processMapper(protocol string, client WanIPConnecter) {
	ipString, _ := client.GetExternalIPAddress()

	ip := net.ParseIP(ipString)

	if ip != nil {
		fmt.Println("externalIP :", ip)

		remoteHost := ""
		protocol = strings.ToUpper(protocol)
		port := uint16(30303)
		leaveDur := uint32(1200)
		fmt.Println("gateway :", client.GetServiceClient().RootDevice.URLBase.Host)
		inAddr := internalAddress(client.GetServiceClient().RootDevice.URLBase.Host)
		fmt.Println("localIP :", inAddr)

		if inAddr != nil {
			dErr := client.DeletePortMapping(remoteHost, port, protocol)
			if dErr == nil {
				fmt.Println("delete port ok", port)
				mErr := client.AddPortMapping(remoteHost, port, protocol, port, inAddr.String(), true, "testportMap", leaveDur)
				if mErr == nil {
					fmt.Println("new port map ok", port)
				} else {
					fmt.Println("mErr :", mErr)
				}

			} else {
				fmt.Println("dErr :", dErr)
			}
		}
		//client.AddPortMapping(ip , 30303 , ip ,30303 , )
	}
}

func Test_ExternalIP(t *testing.T) {

	{
		waitTick := 2
		connectors := make(chan WanIPConnecter, waitTick)

		go func() {
			clients, errors, err := internetgateway1.NewWANIPConnection1Clients()
			_, _ = errors, err
			if connectors == nil {
				return
			}
			if len(clients) == 0 {
				connectors <- nil
			}

			for _, client := range clients {
				connectors <- client
			} //for

		}()

		go func() {
			clients, errors, err := internetgateway2.NewWANIPConnection1Clients()
			_, _ = errors, err
			if connectors == nil {
				return
			}
			if len(clients) == 0 {
				connectors <- nil
			}

			for _, client := range clients {
				connectors <- client
			} //for
		}()

	OUT:
		for conn := range connectors {
			if conn == nil {
				waitTick--
				if waitTick == 0 {
					break
				} else {
					continue
				}
			}
			processMapper("udp", conn)
			break OUT
		} //for

	}
	return
	//*/

	{
		clients, errors, err := internetgateway1.NewWANIPConnection1Clients()
		if len(errors) > 0 {
			fmt.Println("errors")
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, client := range clients {

			ipString, _ := client.GetExternalIPAddress()

			ip := net.ParseIP(ipString)

			if ip != nil {
				fmt.Println("externalIP :", ip)

				remoteHost := ""
				protocol := "UDP"
				port := uint16(30303)
				leaveDur := uint32(1200)

				fmt.Println("gateway :", client.RootDevice.URLBase.Host)
				inAddr := internalAddress(client.RootDevice.URLBase.Host)
				fmt.Println("localIP :", inAddr)

				if inAddr != nil {
					dErr := client.DeletePortMapping(remoteHost, port, protocol)
					if dErr == nil {
						fmt.Println("delete port ok", port)
						mErr := client.AddPortMapping(remoteHost, port, protocol, port, inAddr.String(), true, "testportMap", leaveDur)
						if mErr == nil {
							fmt.Println("new port map ok", port)
						} else {
							fmt.Println("mErr :", mErr)
						}

					} else {
						fmt.Println("dErr :", dErr)
					}
				}
				//client.AddPortMapping(ip , 30303 , ip ,30303 , )
			}
		} //for
	}

	//gatewayIP := net.IP{}
	/*
		gatewayIP, err := gateway.DiscoverGateway()
		if err != nil {
			fmt.Println("err : ", err)
			return
		}
		//*/
	/*
		client := natpmp.NewClient(gatewayIP)
		response, err := client.GetExternalAddress()
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Println("External IP address:", response.ExternalIPAddress)

		//*/
}
