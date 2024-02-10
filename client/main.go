package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
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

func promptForClient(availableClients []string, myPort int) string {
	var selectedClient string

	for {

		for index, client := range availableClients {
			if client == fmt.Sprintf("%d", myPort) {
				continue
			}

			fmt.Println(index+1, ": ", client)
		}

		fmt.Print("Select client: ")
		fmt.Scanln(&selectedClient)

		selectedClientIndex := int(selectedClient[0]-'0') - 1

		if selectedClientIndex < 0 || selectedClientIndex >= len(availableClients) {
			fmt.Println("Invalid index")
			continue
		}

		return availableClients[selectedClientIndex]
	}
}

func connectToClient(port string) {

	fmt.Println("Connecting to client", fmt.Sprintf("http://localhost:%s/connect", port));
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://localhost:%s/connect", port), nil)

	if err != nil {
		log.Println("Error connecting to client", err)
		return
	}

	go func() {
		fmt.Println("Enter message: ")
		for {
			var message string
			fmt.Scanln(&message)

			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println("Error writing message", err)
				return
			}
		}
	}()

}

func main() {

	log.Println("Starting client")

	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Connection request received", r.RemoteAddr)

		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Println("Error upgrading connection", err)
			return
		}

		go func() {
			defer conn.Close()
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					log.Println("Error reading message", err)
					return
				}
				log.Printf("Received message: %s", message)

			}
		}()

	})

	server := &http.Server{
		Addr: ":0",
	}

	defer server.Close()

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	myPort := listener.Addr().(*net.TCPAddr).Port
	requestBody := fmt.Sprintf("address=%d", myPort)

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

	go func() {
		selectedClient := promptForClient(clients, myPort)

		connectToClient(selectedClient)

		log.Println("Selected client:", selectedClient)
	}()

	panic(http.Serve(listener, nil))

}
