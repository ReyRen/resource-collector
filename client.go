package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type Client struct {
	conn *websocket.Conn
	addr string
	rm   *recvMsg
	sm   *sendMsg
}

func (c *Client) handler() {
	defer func() {
		err := c.conn.Close()
		Trace.Printf("%s disconnected!", c.addr)
		if err != nil {
			Error.Printf("readPump conn close err: %s\n", err)
		}
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage() // This is a block func, once ws closed, this would be get err
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				Error.Printf("readMessage error: %s\n", err)
			}
			Error.Printf("readMessage error: %s\n", err)
			//flush websites and close website would caused ReadMessage err and trigger defer func
			return
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		Trace.Printf("received messages: %s\n", message)
		jsonHandler(message, c.rm)

		getGpuRsInfo(c)
		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			Error.Printf("handle log nextWriter error:%s\n", err)
			return
		}
		sdmsg, _ := json.Marshal(c.sm)
		_, err = w.Write(sdmsg)
		if err != nil {
			Error.Printf("write err: %s\n", err)
		}
		if err := w.Close(); err != nil {
			Error.Printf("websocket closed err: %s\n", err)
			return
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// mute : websocket: the client is not using the websocket protocol: 'upgrade' token not found in 'Connection' header
		return
	}

	var rmtmp recvMsg
	var smtmp sendMsg

	client := &Client{
		conn: conn,
		addr: conn.RemoteAddr().String(),
		rm:   &rmtmp,
		sm:   &smtmp,
	}

	//addr
	Trace.Printf("%s connected!\n", client.addr)
	go client.handler()
}
