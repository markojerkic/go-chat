package main

import (
	"container/list"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func readMessages(conn *websocket.Conn)  {
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
