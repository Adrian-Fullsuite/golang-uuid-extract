package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

func encrypt(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}

func main() {
	godotenv.Load()
	uuid := exec.Command("wmic", "csproduct", "get", "uuid")
	serial := exec.Command("wmic", "diskdrive", "get", "serialnumber")

	uuid_output, err := uuid.CombinedOutput()
	if err != nil {
		panic(err)
	}

	serial_output, err := serial.CombinedOutput()
	if err != nil {
		panic(err)
	}


	stdout_uuid := string(uuid_output)
	stdout_serial := string(serial_output)
	stdout_serial = strings.ReplaceAll(stdout_serial, "SerialNumber", "")
	stdout_uuid = strings.ReplaceAll(stdout_uuid, "UUID", "")
	
	re := regexp.MustCompile(`[\s\n]`)
	stdout_serial = re.ReplaceAllString(stdout_serial, "")
	stdout_uuid = re.ReplaceAllString(stdout_uuid, "")

	output := stdout_uuid + stdout_serial

	text := []byte(output)
	key := []byte(os.Getenv("SECRET_KEY"))

	hash, _ := encrypt(text, key)
	
	filename := "serial.spi"

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	
	_, err = writer.WriteString(hash)

	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	// Flush the writer to ensure that all data is written to the file
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing writer:", err)
		return
	}

}