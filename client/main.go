package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func getAvailableClients() ([]string, error) {
	response, err := http.Get("http://localhost:8080/clients")

	if err != nil {
		return nil, err
	}

	clients, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var availableClients []string

	err = json.Unmarshal(clients, &availableClients)
	if err != nil {
		return nil, err
	}

	return availableClients, nil
}

func main() {

	log.Println("Starting client")

	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Connection request received", r.RemoteAddr)
	})

	server := &http.Server{
		Addr: ":0",
	}

	defer server.Close()

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	requestBody := fmt.Sprintf("address=%d", listener.Addr().(*net.TCPAddr).Port)

	response, err := http.Post("http://localhost:8080/register", "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(requestBody)))

	log.Println("Request sent", response.Status)

	if err != nil {
		log.Println("Error sending request", err)
		panic(err)
	}

	log.Println("Using port:", listener.Addr().(*net.TCPAddr).Port)

	clients, err := getAvailableClients()

	if err != nil {
		log.Println("Error sending request", err)
		panic(err)
	}

	log.Println("Available clients", clients)

	panic(http.Serve(listener, nil))

}
