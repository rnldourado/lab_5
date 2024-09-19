package util

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func CalculateHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashBytes := hash.Sum(nil)

	hashString := fmt.Sprintf("%x", hashBytes)

	return hashString, nil
}
