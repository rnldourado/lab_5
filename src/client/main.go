package main

import (
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

	var conn *net.Conn

	if (args[1] == "help"){
		help()
	}else if (args[1] == "join" && args[3] == "search"){
		conn = join(args[2])
		fmt.Println(search(args[2]))
	}else{
		help()
	}

	(*conn).Close()
}

func help() {
	fmt.Println("Uso: go run main.go join [serverip] seach [file hash]")
	return
}

func join(serverURI string) *net.Conn {
	// reader := bufio.NewReader(os.Stdin)

	// fmt.Print("Digite o IP do server: ")
	// serverURI, _ := reader.ReadString('\n')

	conn, err := net.Dial("tcp", strings.Trim(serverURI, "\n"))

	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return nil
	}

	return &conn

	// fmt.Print("Digite sua mensagem: ")
	// text, _ := reader.ReadString('\n')
	//
	// fmt.Fprintf(conn, text)
	//
	// message, _ := bufio.NewReader(conn).ReadString('\n')
	// fmt.Print("Resposta do servidor: " + message)
}

func search(hash string) []byte {
	return []byte(hash)
}
