package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
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

	//buf := make([]byte, 1024)
	//_, err := conn.Read(buf)
	buf, err := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Estou aqui!!!!")
	if err != nil {
		return
	}

	fmt.Println(buf)
	tempHashes := make(map[string][]net.Addr)
	err = json.Unmarshal([]byte(buf), &tempHashes)
	if err != nil {
		fmt.Println("Erro ao decriptografar o arquivo: ", err)
		//return
	}
	mutex.Lock()
	hashStorage(tempHashes, conn.RemoteAddr())
	mutex.Unlock()

	fmt.Println(hashes["870067c23728f0f86f0ca77091d1c60921a086caf23990c5ca38e8128b5e586e"])

	//message, _ := bufio.NewReader(conn).ReadString('\n')
	//fmt.Printf("Hash recebido: %v de: %v", string(message), conn.RemoteAddr())

	//hash := string(message)
	//hash = hash[:len(hash)-1]

	//mutex.Lock()
	//defer mutex.Unlock()

	//hashes[hash] = append(hashes[hash], conn.RemoteAddr())

	//ipList := hashes[hash]
	//response := "Endereços IPs que possuem o arquivo: "
	//for _, addr := range ipList {
	//	response += fmt.Sprintf("%v ", addr.String())
	//}

	//fmt.Printf("Hash recebido: %v de %v\n", hash, conn.RemoteAddr())
	//conn.Write([]byte(response + "\n"))
}

func hashStorage(tempHashes map[string][]net.Addr, addr net.Addr) {
	for hash, _ := range tempHashes {
		hashes[hash] = append(hashes[hash], addr)
	}
}
