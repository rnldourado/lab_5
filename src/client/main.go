package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		help()
		return
	}


	if (args[1] == "help"){
		help()
	}else if (args[1] == "join" && args[3] == "search"){
		conn := join(args[2])
		fmt.Print(search(args[4], conn))
		// (*conn).Close()
	}else{
		help()
	}

}

func help() {
	fmt.Println("Uso: go run main.go join [serverip] seach [file hash]")
	return
}

func join(serverURI string) *net.Conn {

	conn, err := net.Dial("tcp", strings.Trim(serverURI, "\n"))

	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return nil
	}

	return &conn
}

func search(hash string, conn *net.Conn) string {
	defer (*conn).Close()
	fmt.Fprintf(*conn, hash + "\n")
	
	message, _ := bufio.NewReader(*conn).ReadString('\n')

	return "Resposta do servidor: \n" + message
}