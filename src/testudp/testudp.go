package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(0)
	}
}

func main() {
	typeSlice := []string{
		"client",
		"server",
		"p2p",
	}

	rFuncMap := map[string]func(){}
	rFuncMap[typeSlice[0]] = runClient
	rFuncMap[typeSlice[1]] = runServer
	rFuncMap[typeSlice[2]] = p2p

	kind := flag.String("t", "",
		`
	[option]
		--t=client
		or
		--t=server
	`)

	flag.Parse()
	/*
		if flag.NFlag() == 0 {
			flag.Usage()
			return
		}
		//*/

	runType := *kind

	isok := false
	for _, v := range typeSlice {
		if v == runType {
			isok = true
			break
		}
	}

	if isok {
		rFuncMap[runType]()
	} else {
		flag.Usage()
	}

}

const (
	LOCAL_IP   = "localhost"
	BROAD_IP   = "255.255.255.255" //"localhost"
	PORT       = "2345"
	UDP        = "udp"
	BUFFERSIZE = 1024
)

type udpMsg struct {
	buf  []byte
	addr *net.UDPAddr
}
type Session struct {
	conn map[string]time.Time
}

func newSession() *Session {
	session := &Session{
		conn: make(map[string]time.Time),
	}
	return session
}

func (s *Session) Refresh(ip string) {
	//t := time.NewTimer(time.Second * 5)
	s.conn[ip] = time.Now()

}

func (s *Session) Len(limitTime int) int {

	count := len(s.conn)
	if count > 0 {
		rsl := []string{}
		for ip, t := range s.conn {
			elipsed := time.Now().Sub(t)
			if elipsed > time.Second*time.Duration(limitTime) {
				fmt.Println("elipsed : ", elipsed)
				rsl = append(rsl, ip)
			}
		} //for
		if len(rsl) > 0 {
			for _, ip := range rsl {
				delete(s.conn, ip)
				fmt.Println("remove : ", ip)
			}
		}

		return len(s.conn)
	}
	return count
}

func receiveBroadPeer(localIP string, chnInn chan<- *net.UDPAddr) {
	broadAddr, err := net.ResolveUDPAddr(UDP, ":"+PORT)
	CheckError(err)

	broadConn, err := net.ListenUDP(UDP, broadAddr)
	//realAddr := broadConn.LocalAddr().(*net.UDPAddr)

	CheckError(err)
	defer broadConn.Close()

	for {
		buf := make([]byte, BUFFERSIZE)
		//fmt.Println("wait listen client...")
		_, addr, err := broadConn.ReadFromUDP(buf)

		if err != nil {
			fmt.Println("Error:", err)
		} else {

			if addr.String() == localIP {
				fmt.Println("Search Peer")

			} else {
				fmt.Println("broad addr :", addr)
				chnInn <- addr
			}
		}
	} //for
}

func receiveChoicePeer(serverConn *net.UDPConn, cMsg chan<- udpMsg) {
	for {
		buf := make([]byte, BUFFERSIZE)
		n, addr, err := serverConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("error:", err)
		}
		cMsg <- udpMsg{buf[:n], addr}
	}
}

func p2p() {
	fmt.Println("----- p2p -----")

	broadAddr, err := net.ResolveUDPAddr(UDP, BROAD_IP+":"+PORT)
	CheckError(err)

	broadConn, err := net.DialUDP(UDP, nil, broadAddr)
	CheckError(err)

	localIP := broadConn.LocalAddr().String()
	fmt.Println("local Addr :", localIP)

	defer broadConn.Close()

	chnPacket := make(chan udpMsg)
	chnInn := make(chan *net.UDPAddr)
	session := newSession()

	serverAddr, err := net.ResolveUDPAddr(UDP, localIP)
	CheckError(err)

	serverConn, err := net.ListenUDP(UDP, serverAddr)
	CheckError(err)
	defer serverConn.Close()

	go receiveChoicePeer(serverConn, chnPacket)
	go receiveBroadPeer(localIP, chnInn)

	icount := 0
	for {
		select {
		default:
			if session.Len(2) == 0 {
				//fmt.Println("broad write")
				_, err := broadConn.Write([]byte("nn"))
				if err != nil {
					fmt.Println("error:", err)
				}
			}
			time.Sleep(time.Millisecond * 500)

		case addr := <-chnInn:
			session.Refresh(addr.String())

			_, err := serverConn.WriteToUDP([]byte("ACK"), addr)
			if err != nil {
				fmt.Println("error:", err)
			}

		case msgpack := <-chnPacket:
			fmt.Println("get:", string(msgpack.buf), " , addr:", msgpack.addr)

			session.Refresh(msgpack.addr.String())

			_, err := serverConn.WriteToUDP([]byte(fmt.Sprintf("%v", icount)), msgpack.addr)
			if err != nil {
				fmt.Println("error:", err)
			}

			icount++
		} //select
	} //for

}

func runClient() {
	fmt.Println("----- runClient -----")

	RemoteAddr, err := net.ResolveUDPAddr(UDP, LOCAL_IP+":"+PORT)
	CheckError(err)

	//LocalAddr, err := net.ResolveUDPAddr(UDP, ":0")
	//CheckError(err)

	Conn, err := net.DialUDP(UDP, nil, RemoteAddr)
	fmt.Println("local Addr :", Conn.LocalAddr().String())
	CheckError(err)

	defer Conn.Close()

	go func() {

		for {

			fmt.Println("wait read...")
			buf := make([]byte, BUFFERSIZE)
			n, addr, err := Conn.ReadFromUDP(buf)

			if err != nil {
				fmt.Println("error:", err)
			}
			fmt.Println("Received ", string(buf[:n]), " from ", addr)
		}
	}()

	i := 0
	for {

		msg := strconv.Itoa(i)
		i++
		buf := []byte(msg)
		_, err := Conn.Write(buf)
		//n, err := Conn.WriteToUDP(buf, RemoteAddr)

		if err != nil {
			fmt.Println(msg, err)
		} else {
			fmt.Println("write :", msg, ", localAddr:", Conn.LocalAddr(), " --> remoteAddr :", Conn.RemoteAddr())
		}
		time.Sleep(time.Millisecond * 500)

	}

}

func runServer() {
	fmt.Println("----- runServer -----")

	ServerAddr, err := net.ResolveUDPAddr(UDP, ":"+PORT)
	CheckError(err)

	ServerConn, err := net.ListenUDP(UDP, ServerAddr)
	CheckError(err)
	defer ServerConn.Close()

	for {
		fmt.Println("wait listen client...")

		buf := make([]byte, BUFFERSIZE)
		n, addr, err := ServerConn.ReadFromUDP(buf)

		if err != nil {
			fmt.Println("Error:", err)
		}

		//time.Sleep(time.Millisecond * 500)

		n, err = ServerConn.WriteToUDP([]byte("ack"), addr)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Received :", string(buf[:n]), " from:", addr, " --> write Ack!")
		}
	}

}
