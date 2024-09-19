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
		fmt.Print("> ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)
		parts := strings.Fields(command)

		switch parts[0] {
		case "files":
			files()
		case "search":
			search(parts[1], *conn)
		case "exit":
			fmt.Println("Até Logo!")
			return
		}
	}
}

func join(serverURI string) *net.Conn {

	//conn, err := net.Dial("tcp", IP)
	conn, err := net.Dial("tcp", strings.Trim(serverURI, "\n"))

	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		os.Exit(1)
	}

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Erro ao ler mensagem do servidor:", err)
		return nil
	}

	fmt.Print(message)

	return &conn
}

func search(hash string, conn net.Conn) {

	_, err := conn.Write([]byte(fmt.Sprintf("search %s\n", hash)))
	if err != nil {
		fmt.Println("Erro ao enviar o comando de busca:", err)
		return
	}

	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Erro ao ler a resposta do servidor:", err)
	}

	var addresses []string
	err = json.Unmarshal([]byte(response), &addresses)
	if err != nil {
		fmt.Println("Erro ao desserializar a resposta:", err)
	}

	fmt.Println("Endereços com o arquivo:", addresses)

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
	_, err = conn.Write(jsonData)
	fmt.Fprintf(conn, "\n")
	if err != nil {
		fmt.Println("Erro ao enviar ao servidor: ", err)
		return err
	}

	return nil
}

func files() {
	for hash, name := range fileHashes {
		fmt.Println("Arquivo: ", name)
		fmt.Printf("Hash: %s\n", hash)
	}
}
