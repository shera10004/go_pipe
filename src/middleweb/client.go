package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn //웹소켓 컨넥션
	send chan *Message   //메시지 전송용 채널

	roomId string //현재 접속한 채팅방 아이디
	user   *User  //현재 접속한 사용자 정보
}

func (c *Client) Close() {
	for i, client := range clients {
		if client == c {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}

	//send 채널 닫음
	close(c.send)

	//웹소켓 커넥션 종료
	c.conn.Close()
	log.Printf("close connection addr:%s", c.conn.RemoteAddr())
}
func (c *Client) read() (*Message, error) {
	var msg *Message

	if err := c.conn.ReadJSON(&msg); err != nil {
		return nil, err
	}

	msg.CreateAt = time.Now()
	msg.User = c.user

	log.Println("read from websocket:", msg)

	return msg, nil
}
func (c *Client) write(m *Message) error {
	log.Println("write to websocket:", m)

	return c.conn.WriteJSON(m)
}
func (c *Client) readLoop() {
	//메시지 수신 대기
	for {
		m, err := c.read()
		if err != nil {
			log.Println("read message error:", err)
			break
		}

		//메시지가 수신되면 수신된 메시지를 mongoDB에 생성하고 모든 clients에 전달
		m.create()
		boradcast(m)
	}
	c.Close()
}
func (c *Client) writeLoop() {
	//클라이언트의 send 채널 메시지 수신 대기
	for msg := range c.send {
		//클라이언트의 채팅방 아이디와 전달된 메시지의 채팅방 아이디가 일치하면 웹소켓에 메시지 전달
		if c.roomId == msg.RoomId.Hex() {
			c.write(msg)
		}
	}
}

func boradcast(m *Message) {
	//모든 클라이언트의 send채널에 메시지 전달
	for _, client := range clients {
		client.send <- m
	}
}

//현재 접속 중인 전체 클라이언트 리스트
var clients []*Client

const messageBufferSize = 256

func newClient(conn *websocket.Conn, roomId string, u *User) {

	c := &Client{
		conn:   conn,
		send:   make(chan *Message, messageBufferSize),
		roomId: roomId,
		user:   u,
	}

	//clients 목록에 새로 생성한 클라이언트 추가
	clients = append(clients, c)

	//메시지 수신/전송 대기
	go c.readLoop()
	go c.writeLoop()
}
