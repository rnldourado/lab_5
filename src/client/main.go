package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/rnldourado/lab_5/src/util"
)

var fileHashes = make(map[string]string)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Uso: go run client.go [serverip:port]")
		return
	}

	serverAddress := os.Args[1]

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return
	}
	defer conn.Close()

	hashes := loadHashesFromDataset()
	fmt.Fprintf(conn, strings.Join(hashes, ",")+"\n")

	fmt.Println("Conectado ao servidor.")

	for {
		fmt.Print("('ls' para listar os seus aqruivos ou 'exit' para sair).\nDigite o hash para buscar: ")
		reader := bufio.NewReader(os.Stdin)
		hash, _ := reader.ReadString('\n')
		hash = strings.TrimSpace(hash)

		if hash == "exit" {
			fmt.Println("Encerrando conexão.")
			break
		} else if hash == "ls" {
			files()
		} else {
			fmt.Fprintf(conn, hash+"\n")

			scanner := bufio.NewScanner(conn)

			for scanner.Scan() {
				line := scanner.Text()
				if strings.Trim(line, "EOF") == "" {
					break
				}
				fmt.Println(line)
			}

		}
	}

}

func loadHashesFromDataset() []string {
	var fileHashes []string
	datasetPath := "/tmp/dataset"

	files, err := os.ReadDir(datasetPath)
	if err != nil {
		fmt.Println("Erro ao ler o diretório:", err)
		return fileHashes
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := datasetPath + "/" + file.Name()
			hash, err := util.CalculateHash(filePath)
			if err == nil {
				fileHashes = append(fileHashes, hash)
			}
		}
	}
	return fileHashes
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

func files() {
	datasetPath := "/tmp/dataset"

	files, err := os.ReadDir(datasetPath)
	if err != nil {
		fmt.Println("Erro ao ler o diretório:", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := datasetPath + "/" + file.Name()
			hash, err := util.CalculateHash(filePath)
			if err == nil {
				fmt.Println("Arquivo: ", file.Name())
				fmt.Printf("Hash: %s\n", hash)
			} else {
				fmt.Println("Erro ao calcular o hash para o arquivo:", file.Name())
			}
		}
	}
}
