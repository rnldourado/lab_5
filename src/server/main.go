package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var (
	clients = make(map[string][]string)
	mutex   sync.Mutex
)

func main() {
	ln, err := net.Listen("tcp", ":13000")
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer ln.Close()
	fmt.Println("Server is running: 13000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr().String()

	fmt.Println("Nova conexão de: ", conn.RemoteAddr())

	hashes, _ := bufio.NewReader(conn).ReadString('\n')
	hashes = strings.TrimSpace(hashes)
	clientHashes := strings.Split(hashes, ",")

	mutex.Lock()
	clients[clientAddr] = clientHashes
	mutex.Unlock()

	for {
		requestedHash, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("Cliente %s desconectado\n", clientAddr)
			mutex.Lock()
			delete(clients, clientAddr)
			mutex.Unlock()
			return
		}
		requestedHash = strings.TrimSpace(requestedHash)

		clientsWithFile := findClientsWithHash(requestedHash)

		var response string

		if len(clientsWithFile) > 0 {
			response = "Clientes com o arquivo:\n" + strings.Join(clientsWithFile, "\n")
			fmt.Println(response)
		} else {
			response = "Nenhum cliente possui o arquivo com o hash solicitado."
		}

		conn.Write([]byte(response + "\nEOF\n"))
	}
}

func findClientsWithHash(hash string) []string {
	mutex.Lock()
	defer mutex.Unlock()

	var clientsWithFile []string
	for client, hashes := range clients {
		for _, clientHash := range hashes {
			if clientHash == hash {
				clientsWithFile = append(clientsWithFile, client)
			}
		}
	}
	return clientsWithFile
}
