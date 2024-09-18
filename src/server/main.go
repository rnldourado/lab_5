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

	welcome(conn)

	buf := make([]byte, 1024)
	n, err := conn.Read(buf[:])
	if err != nil {
		return
	}

	data := buf[:n]

	//fmt.Println(buf)
	tempHashes := make(map[string]string)
	err = json.Unmarshal(data, &tempHashes)
	if err != nil {
		fmt.Println("Erro ao decriptografar o arquivo: ", err)
		return
	}

	mutex.Lock()
	hashStorage(tempHashes, conn.RemoteAddr())
	mutex.Unlock()

	//fmt.Println(hashes["b8696ddcb191628675c8667cad61444fb8a367bdabed66053f06fc579ddc3804"])

	for {
		reader := bufio.NewReader(conn)

		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Cliente desconectado: ", conn.RemoteAddr())
			return
		}

		fmt.Println(message)

	}

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
