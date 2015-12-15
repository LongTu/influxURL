package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
)

var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

func main() {
	// Need to encrypt a string
	plaintext := []byte("My name is Astaxie")
	// If there is an incoming string of words to be encrypted, set plaintext to that incoming string
	if len(os.Args) > 1 {
		plaintext = []byte(os.Args[1])
	}

	// aes encryption string
	keyText := "astaxie12798akljzmknm.ahkjkljl;k"
	if len(os.Args) > 2 {
		keyText = os.Args[2]
	}

	// Create the aes encryption algorithm
	c, err := aes.NewCipher([]byte(keyText))
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(keyText), err)
		os.Exit(-1)
	}

	// Encrypted string
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	fmt.Printf("%s=>%x\n", plaintext, ciphertext)
}
