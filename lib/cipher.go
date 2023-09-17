package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"os"
)

var secretKey string
var iv string

func Encrypt(plaintext []byte) string {
	secretKey = os.Getenv("SECRET_KEY")
	iv = os.Getenv("IV")
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		panic(err)
	}
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	encrypted := make([]byte, len(plaintext))
	mode.CryptBlocks(encrypted, plaintext)
	encoded := base64.StdEncoding.EncodeToString(encrypted)
	return encoded
}

func Decrypt(plaintext string) []byte {
	// Decode the data field
	decoded, err := base64.StdEncoding.DecodeString(plaintext)
	if err != nil {
		panic(err)
	}

	// Decrypt the response
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		panic(err)
	}
	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	decrypted := make([]byte, len(decoded))
	mode.CryptBlocks(decrypted, decoded)
	return decrypted
}
