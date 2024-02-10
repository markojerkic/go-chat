package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func readMessages(conn *websocket.Conn) {
	defer conn.Close()

	for {
		_, p, err := conn.ReadMessage()

		if err != nil {
			log.Println("Error reading message", err)
		}

		log.Println("Message received: ", string(p))

	}

}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	go readMessages(conn)
}

var registeredClients = map[string]bool{}

func main() {
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		adr := r.FormValue("address")
		log.Println("Registering", adr)

		registeredClients[adr] = true
	})

	http.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		keys := make([]string, 0, len(registeredClients))
		for k := range registeredClients {
			keys = append(keys, k)
		}

		clients, err := json.Marshal(keys)

		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		io.WriteString(w, string(clients))
	})

	server := &http.Server{
		Addr:              "localhost:8080",
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Println("Server started")
	err := server.ListenAndServe()

	defer server.Close()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
