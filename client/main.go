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

func promptForReceiverClient(availableClients []string, myPort int) (string, bool) {
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

		if selectedClient == "refresh" {
			return "", true
		}

		selectedClientIndex := int(selectedClient[0]-'0') - 1

		if selectedClientIndex < 0 || selectedClientIndex >= len(availableClients) {
			fmt.Println("Invalid index")
			continue
		}

		return availableClients[selectedClientIndex], false
	}
}

type Receiver struct {
	Address string
	Connection *websocket.Conn
}

var currentReceiver *Receiver

func (receiver *Receiver) closeConnection() {
	receiver.Connection.Close()
}

func (receiver *Receiver) openConnection(message chan string) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://localhost:%s/connect", receiver.Address), nil)

	receiver.Connection = conn

	if err != nil {
		log.Println("Error connecting to receiver", err)
		return
	}

	go func() {
		fmt.Println("Enter message: ")
		for {

			enteredMessage := <-message

			if enteredMessage == "" {
				continue
			} else {
				message <- ""
			}

			err := conn.WriteMessage(websocket.TextMessage, []byte(enteredMessage))
			if err != nil {
				log.Println("Error writing message", err)
				return
			}
		}
	}()

}

var message = make(chan string)

func listenForInput(message chan string) {
	for {
		var enteredMessage string
		fmt.Scanln(&enteredMessage)
		fmt.Println("Message entered: ", enteredMessage)
		if enteredMessage == "switch" {
			fmt.Println("User wants to switch clients")
			switchReceiver()
		} else {
			message <- enteredMessage
		}
	}
}

func switchReceiver() {
	go func() {
		for {
		clients, err := getAvailableClients()

		if err != nil {
			log.Println("Error sending request", err)
			panic(err)
		}

		selectedClient, refresh := promptForReceiverClient(clients, myPort)
		if refresh {
			continue
		}

		if currentReceiver != nil {
			currentReceiver.closeConnection()
		}

		currentReceiver = &Receiver{Address: selectedClient}
		currentReceiver.openConnection(message)
		break
	}
	}()

}

func handleConnectionRequest(w http.ResponseWriter, r *http.Request) {
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
			log.Printf("Message from %s: %s", r.RemoteAddr, message)

		}
	}()

}

func registerWithCoordinator(listener net.Listener) error {
	myPort = listener.Addr().(*net.TCPAddr).Port
	requestBody := fmt.Sprintf("address=%d", myPort)

	_, err := http.Post("http://localhost:8080/register", "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(requestBody)))

	if err != nil {
		log.Println("Error sending request", err)
		return err
	}

	return nil
}

var myPort int

func main() {

	log.Println("Starting client")

	http.HandleFunc("/connect", handleConnectionRequest)

	server := &http.Server{
		Addr: ":0",
	}

	defer server.Close()

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	log.Println("Using port:", listener.Addr().(*net.TCPAddr).Port)

	if registerWithCoordinator(listener) != nil {
		log.Println("Error registering with coordinator")
		return
	}

	go switchReceiver()

	panic(http.Serve(listener, nil))
}
