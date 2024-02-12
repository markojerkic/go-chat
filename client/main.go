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

type Peer struct {
	ID   string "json:id"
	Port int    "json:port"
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
			fmt.Printf("%d: %s\n", index+1, peer.ID)
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

	savedConnection, hasSavedConnection := connectedClients[currentPeer.ID]
	if !hasSavedConnection {

		connection := Connection{currentPeer.ID, nil, make(chan string)}
		err := connection.openConnection(currentPeer.Port)
		if err != nil {
			log.Println("Error opening connection to peer", err)
			return Connection{}, err
		}

		connectedClients[currentPeer.ID] = connection

		return connection, nil
	}

	return savedConnection, nil

}

func watchStdIn(clients chan map[string]Connection) {
	var currentPeerId *Peer

	for {

		if currentPeerId == nil {
			peer, err := selectPeer()
			if err != nil {
				log.Println("Error selecting peer", err)
				continue
			}
			currentPeerId = &peer
		}

		var message string
		fmt.Scanln(&message)

		if currentPeerId == nil || message == "switch" {
			peer, err := selectPeer()
			if err != nil {
				log.Println("Error selecting peer", err)
				continue
			}
			currentPeerId = &peer

			log.Println("User wants to switch peers")

		} else {
			peerConnection, err := getOrCreatePeerConnection(clients, *currentPeerId)
			if err != nil {
				log.Println("Error getting or creating peer connection", err)
				panic(err)
			}
			go peerConnection.sendMessage(message)
		}

	}
}

func registerWithCoordinator(listener net.Listener) (int, error) {

	var myId string
	fmt.Print("Enter your id: ")
	fmt.Scanln(&myId)

	myPort := listener.Addr().(*net.TCPAddr).Port
	requestBody := fmt.Sprintf("address=%d&id=%s", myPort, myId)

	_, err := http.Post("http://localhost:8080/register", "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(requestBody)))

	if err != nil {
		log.Println("Error sending request", err)
		return 0, err
	}

	return myPort, nil
}

func main() {
	connections := make(chan map[string]Connection)

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	registerWithCoordinator(listener)

	go watchStdIn(connections)

	panic(http.Serve(listener, nil))

}
