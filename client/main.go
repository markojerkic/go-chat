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

// func getAvailableClients() ([]string, error) {
// 	response, err := http.Get("http://localhost:8080/clients")
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	clients, err := io.ReadAll(response.Body)
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var availableClients []string
//
// 	err = json.Unmarshal(clients, &availableClients)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return availableClients, nil
// }
//
// func promptForReceiverClient(availableClients []string, myPort int) (string, bool) {
// 	var selectedClient string
//
// 	for {
// 		isLocked := <-stdLock
// 		if isLocked {
// 			continue
// 		}
// 		stdLock <- true
//
// 		for index, client := range availableClients {
// 			if client == fmt.Sprintf("%d", myPort) {
// 				continue
// 			}
//
// 			fmt.Println(index+1, ": ", client)
// 		}
//
// 		fmt.Print("Select client: ")
// 		fmt.Scanln(&selectedClient)
//
// 		if selectedClient == "refresh" {
// 			return "", true
// 		}
//
// 		selectedClientIndex := int(selectedClient[0]-'0') - 1
//
// 		if selectedClientIndex < 0 || selectedClientIndex >= len(availableClients) {
// 			fmt.Println("Invalid index")
// 			continue
// 		}
//
// 		stdLock <- false
//
// 		return availableClients[selectedClientIndex], false
// 	}
// }
//
// type Receiver struct {
// 	Address    string
// 	Connection *websocket.Conn
// }
//
// var currentReceiver *Receiver
//
// func (receiver *Receiver) closeConnection() {
// 	receiver.Connection.Close()
// }
//
// func (receiver *Receiver) openConnection(message chan string) {
// 	log.Println("Opening connection to receiver", receiver.Address)
// 	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://localhost:%s/connect", receiver.Address), nil)
//
// 	receiver.Connection = conn
//
// 	if err != nil {
// 		log.Println("Error connecting to receiver", err)
// 		return
// 	}
//
// 	go func() {
// 		fmt.Println("Enter message: ")
// 		for {
//
// 			enteredMessage := <-message
//
// 			if enteredMessage == "" {
// 				continue
// 			} else {
// 				message <- ""
// 			}
//
// 			err := conn.WriteMessage(websocket.TextMessage, []byte(enteredMessage))
// 			if err != nil {
// 				log.Println("Error writing message", err)
// 				return
// 			}
// 		}
// 	}()
//
// }
//
// func listenForInput(message chan string) {
// 	for {
// 		isLocked := <-stdLock
// 		if isLocked {
// 			continue
// 		}
//
// 		var enteredMessage string
//
// 		stdLock <- true
// 		fmt.Scanln(&enteredMessage)
// 		fmt.Println("Message entered: ", enteredMessage)
// 		if enteredMessage == "switch" {
// 			fmt.Println("User wants to switch clients")
// 			stdLock <- false
// 			switchReceiver()
// 		} else {
// 			message <- enteredMessage
// 		}
// 		stdLock <- false
// 	}
// }
//
// func switchReceiver() {
// 	for {
// 		log.Println("Switching receiver")
// 		clients, err := getAvailableClients()
//
// 		if err != nil {
// 			log.Println("Error sending request", err)
// 			panic(err)
// 		}
//
// 		selectedClient, refresh := promptForReceiverClient(clients, myPort)
// 		if refresh {
// 			continue
// 		}
//
// 		if currentReceiver != nil {
// 			log.Println("Closing connection to current receiver")
// 			currentReceiver.closeConnection()
// 		}
//
// 		currentReceiver = &Receiver{Address: selectedClient}
// 		currentReceiver.openConnection(message)
// 		break
// 	}
//
// }
//
// func handleConnectionRequest(w http.ResponseWriter, r *http.Request) {
// 	upgrader := websocket.Upgrader{
// 		ReadBufferSize:  1024,
// 		WriteBufferSize: 1024,
// 	}
//
// 	conn, err := upgrader.Upgrade(w, r, nil)
//
// 	if err != nil {
// 		log.Println("Error upgrading connection", err)
// 		return
// 	}
//
// 	log.Println("Connection established with", r.RemoteAddr)
//
// 	go func() {
// 		defer conn.Close()
// 		for {
// 			_, message, err := conn.ReadMessage()
// 			if err != nil {
// 				log.Println("Error reading message", err)
// 				return
// 			}
// 			log.Printf("Message from %s: %s", r.RemoteAddr, message)
//
// 		}
// 	}()
//
// }
//
// func registerWithCoordinator(listener net.Listener) error {
// 	myPort = listener.Addr().(*net.TCPAddr).Port
// 	requestBody := fmt.Sprintf("address=%d", myPort)
//
// 	_, err := http.Post("http://localhost:8080/register", "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(requestBody)))
//
// 	if err != nil {
// 		log.Println("Error sending request", err)
// 		return err
// 	}
//
// 	return nil
// }
//
// var myPort int
// var message chan string
// var stdLock chan bool
//
// func main() {
//
// 	log.Println("Starting client")
//
// 	http.HandleFunc("/connect", handleConnectionRequest)
//
// 	server := &http.Server{
// 		Addr: ":0",
// 	}
//
// 	defer server.Close()
//
// 	listener, err := net.Listen("tcp", ":0")
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	log.Println("Using port:", listener.Addr().(*net.TCPAddr).Port)
//
// 	if registerWithCoordinator(listener) != nil {
// 		log.Println("Error registering with coordinator")
// 		return
// 	}
//
// 	stdLock = make(chan bool, 1)
// 	message = make(chan string, 1)
//
// 	stdLock <- false
// 	log.Println("Registered with coordinator")
// 	go switchReceiver()
// 	go listenForInput(message)
//
// 	panic(http.Serve(listener, nil))
// }

type Peer struct {
	id   string
	port int
}

func getAvailablePeers() ([]Peer, error) {
	response, err := http.Get("http://localhost:8080/clients")

	if err != nil {
		return nil, err
	}

	clients, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var availableClients []Peer

	err = json.Unmarshal(clients, &availableClients)
	if err != nil {
		return nil, err
	}

	return availableClients, nil
}

func selectPeer() (Peer, error) {
	for {
		peers, err := getAvailablePeers()
		if err != nil {
			return Peer{}, err
		}

		fmt.Println("Select a peer to connect to:")
		for index, peer := range peers {
			fmt.Printf("%d: %s\n", index+1, peer.id)
		}
		fmt.Print("Enter the number of the peer: ")
		var selectedPeerIndex int
		fmt.Scanln(&selectedPeerIndex)

		if selectedPeerIndex < 1 || selectedPeerIndex > len(peers) {
			fmt.Println("Invalid selection, try again")
			continue
		}

		return peers[selectedPeerIndex-1], nil
	}

}

func getOrCreatePeerConnection(clients chan map[string]Connection, currentPeer Peer) (Connection, error) {
	connectedClients := <-clients

	savedConnection, hasSavedConnection := connectedClients[currentPeer.id]
	if !hasSavedConnection {

		connection := Connection{currentPeer.id, nil, make(chan string)}
		err := connection.openConnection(currentPeer.port)
		if err != nil {
			log.Println("Error opening connection to peer", err)
			return Connection{}, err
		}

		connectedClients[currentPeer.id] = connection

		return connection, nil
	}

	return savedConnection, nil

}

func watchStdIn(clients chan map[string]Connection) {
	var currentPeerId *Peer

	for {

		var message string
		fmt.Scanln(&message)

		if currentPeerId == nil {
			peer, err := selectPeer()
			if err != nil {
				log.Println("Error selecting peer", err)
				continue
			}
			currentPeerId = &peer
		}

		peerConnection, err := getOrCreatePeerConnection(clients, *currentPeerId)
		if err != nil {
			log.Println("Error getting or creating peer connection", err)
			panic(err)
		}

		if message == "switch" {
			log.Println("User wants to switch peers")
		} else {
			peerConnection.sendMessage(message)
		}

	}
}

func main() {
}
