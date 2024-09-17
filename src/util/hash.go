package util

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func CalculateHash(filePath string) (string, error) {
	// Abre o arquivo
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Cria um objeto para calcular o hash SHA-256
	hash := sha256.New()

	// Copia o conteúdo do arquivo para o objeto hash
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// Converte o hash para um slice de bytes
	hashBytes := hash.Sum(nil)

	// Converte o slice de bytes para uma string hexadecimal
	hashString := fmt.Sprintf("%x", hashBytes)

	return hashString, nil
}