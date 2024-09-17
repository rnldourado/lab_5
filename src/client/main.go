package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/rnldourado/lab_5/src/util"
)

const IP = "127.0.0.1:13000"

var fileHashes = make(map[string]string)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Digite o IP do server: ")
	serverURI, _ := reader.ReadString('\n')

	conn := join(strings.Trim(serverURI, "\n"))
	if conn != nil {
	}

	loadHashesFromDataset()
	sendHashToServer(fileHashes, *conn)
	//defer conn.Close()

	for {
		readed, _ := reader.ReadString('\n')
		command := strings.Split(readed, " ")

		switch command[0] {
		case "search":
			search(command[1], conn)
		case "exit":
			os.Exit(1)
		}
	}

	//fmt.Fprintf(conn, text)

	//message, _ := bufio.NewReader(conn).ReadString('\n')
	//fmt.Print("Resposta do servidor: " + message)
}

/*
	args := os.Args

	if len(args) < 2 {
		help()
		return
	}

	if args[1] == "help" {
		help()
	} else if args[1] == "join" && args[3] == "search" {
		conn := join(args[2])
		fmt.Print(search(args[4], conn))
		// (*conn).Close()
	} else {
		help()
	}
*/

func help() {
	fmt.Println("Uso: go run main.go join [serverip] seach [file hash]")
	return
}

func join(serverURI string) *net.Conn {

	conn, err := net.Dial("tcp", IP)
	//conn, err := net.Dial("tcp", strings.Trim(serverURI, "\n"))

	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		os.Exit(1)
	}

	return &conn
}

func search(hash string, conn *net.Conn) string {
	defer (*conn).Close()
	fmt.Fprintf(*conn, hash+"\n")

	message, _ := bufio.NewReader(*conn).ReadString('\n')

	return "Resposta do servidor: \n" + message
}

func loadHashesFromDataset() {
	datasetPath := "/tmp/dataset"
	files, err := os.ReadDir(datasetPath)
	if err != nil {
		fmt.Println("Erro ao ler o diretório:", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			hash := calculateHash(file.Name())
			fileHashes[hash] = datasetPath + "/" + file.Name()
		}
	}
}

func calculateHash(fileName string) string {
	filepath := "/tmp/dataset/" + fileName
	hash, err := util.CalculateHash(filepath)

	if err != nil {
		fmt.Println("Error ao calcular o hash do arquivo: ", err)
		return ""
	}

	return hash
}

func sendHashToServer(fileHashes map[string]string, conn net.Conn) error {
	jsonData, err := json.Marshal(fileHashes)
	if err != nil {
		fmt.Println("Erro ao serializar: ", err)
		return err
	}
	fmt.Println(jsonData)
	_, err = conn.Write(jsonData)
	fmt.Fprintf(conn, "\n")
	if err != nil {
		fmt.Println("Erro ao enviar ao servidor: ", err)
		return err
	}
	//fmt.Println("Cheguei aqui")

	return nil
}
