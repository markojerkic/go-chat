package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

var registeredClients = map[int]string{}

type Peer struct {
	ID   string "json:id"
	Port int    "json:port"
}

func main() {
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		port, err := strconv.Atoi(r.FormValue("address"))
		if err != nil {
			log.Println("Invalid port", err)
			http.Error(w, "Invalid port", http.StatusBadRequest)
			return
		}
		registeredClients[port] = r.FormValue("id")
		log.Println("Registered", r.FormValue("id"), "at port", port, registeredClients)
	})

	http.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var peers []Peer
		for k, v := range registeredClients {
			peers = append(peers, Peer{v, k})
		}

		clients, err := json.Marshal(peers)

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
