package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)


type Connection struct {
	id               string
	conn             *websocket.Conn
	incomingMessages chan string
}

func (c *Connection) readMessages() error {

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message", err)
			return err
		}

		fmt.Println("-----")
		fmt.Printf("Message from %s: %s", c.id, message)
		fmt.Println("-----")
	}
}

func (c *Connection) sendMessage(message string) {
	c.incomingMessages <- message
}

func (c *Connection) watchIncomingMessages() error {
	for {
		message := <-c.incomingMessages

		if message == "" {
			continue
		}

		err := c.conn.WriteMessage(websocket.TextMessage, []byte(message))

		c.incomingMessages <- ""

		if err != nil {
			log.Println("Error writing message")
			return err
		}

		return nil
	}
}

func (c *Connection) openConnection(port int) error {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://localhost:%d/connect", port), nil)

	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

func (c *Connection) acceptConnection(w http.ResponseWriter, r *http.Request) error {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Error upgrading connection", err)
		return err
	}

	peerId := r.FormValue("id")
	c.id = peerId
	c.conn = conn

	return nil
}

func (c *Connection) closeConnection() {
	c.conn.Close()
}

