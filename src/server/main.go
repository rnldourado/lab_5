package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
)

var (
	hashes = make(map[string][]net.Addr)
	mutex  sync.Mutex
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
	fmt.Println("Nova conexão de: ", conn.RemoteAddr())

	welcome(conn)

	buf := make([]byte, 1024)
	n, err := conn.Read(buf[:])
	if err != nil {
		return
	}

	data := buf[:n]

	tempHashes := make(map[string]string)
	err = json.Unmarshal(data, &tempHashes)
	if err != nil {
		fmt.Println("Erro ao decriptografar o arquivo: ", err)
		return
	}

	mutex.Lock()
	hashStorage(tempHashes, conn.RemoteAddr())
	mutex.Unlock()

	for {
		reader := bufio.NewReader(conn)

		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Cliente desconectado: ", conn.RemoteAddr())
			removeClientFiles(conn.RemoteAddr())
			return
		}

		message = strings.TrimSpace(message)
		command := strings.Fields(message)

		switch command[0] {
		case "search":
			search(command[1], conn)
		}

		fmt.Println(message)

	}

}

func hashStorage(tempHashes map[string]string, addr net.Addr) {
	for hash, _ := range tempHashes {
		hashes[hash] = append(hashes[hash], addr)
	}
}

func welcome(conn net.Conn) {
	mensagem := "Bem-vindo ao servidor!\n"
	_, err := conn.Write([]byte(mensagem))
	if err != nil {
		fmt.Println("Erro ao enviar mensagem ao cliente:", err)
		return
	}
}

func search(hash string, conn net.Conn) {
	addresses, ok := hashes[hash]
	if !ok {
	} else {

		var sAddr []string
		for _, adrr := range addresses {
			sAddr = append(sAddr, adrr.String())
		}

		jsonResponse, err := json.Marshal(sAddr)
		if err != nil {
			fmt.Println("Erro ao serializar: ", err)
			return
		}
		_, err = conn.Write(append(jsonResponse, '\n'))
		fmt.Fprintf(conn, "\n")
		if err != nil {
			fmt.Println("Erro ao enviar resposta ao cliente: ", err)
			return
		}
	}
}

func removeClientFiles(addr net.Addr) {

	for hash, owners := range hashes {

		var updatedOwners []net.Addr
		for _, owner := range owners {
			if owner.String() != addr.String() {
				updatedOwners = append(updatedOwners, owner)
			}
		}

		if len(updatedOwners) > 0 {
			hashes[hash] = updatedOwners
		} else {
			delete(hashes, hash)
		}
	}

	fmt.Printf("Tabela atualizada! Registros de %s excluidos!\n", addr.String())
}
