package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Digite o IP do server: ")
	serverURI, _ := reader.ReadString('\n')

	conn, err := net.Dial("tcp", strings.Trim(serverURI, "\n"))

	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return
	}
	defer conn.Close()

	fmt.Print("Digite sua mensagem: ")
	text, _ := reader.ReadString('\n')

	fmt.Fprintf(conn, text)

	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Resposta do servidor: " + message)
}
