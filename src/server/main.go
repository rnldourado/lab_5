package main

import (
	"bufio"
	"fmt"
	"net"
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
	fmt.Println("New connection from:", conn.RemoteAddr())

	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Mensage received: ", string(message))

	conn.Write([]byte("Mensage received: " + message))
}
