package DB

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql/driver"
	"encoding/base64"
	"fmt"
	"github.com/lambda-platform/lambda/config"
	"golang.org/x/crypto/pbkdf2"
	"io"
)

type SecureString string

// Value implements the driver.Valuer interface, encrypting the string value.
func (s SecureString) Value() (driver.Value, error) {
	encryptedString := Encrypt(string(s), config.Config.JWT.Secret)

	return encryptedString, nil
}

// Scan implements the sql.Scanner interface, decrypting the string value.
func (s *SecureString) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	decryptedString := Decrypt(value.(string), config.Config.JWT.Secret)

	fmt.Println("decryptedString")
	fmt.Println("decryptedString")
	fmt.Println("decryptedString")
	fmt.Println("decryptedString")
	fmt.Println("decryptedString")

	*s = SecureString(decryptedString)
	return nil
}

func Encrypt(stringToEncrypt string, passphrase string) (encryptedString string) {
	key := generateAESKey(passphrase, 32)
	plaintext := []byte(stringToEncrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return base64.URLEncoding.EncodeToString(ciphertext)
}

func Decrypt(encryptedString string, passphrase string) (decryptedString string) {
	key := generateAESKey(passphrase, 32)

	enc, err := base64.URLEncoding.DecodeString(encryptedString)
	if err != nil {
		panic(err.Error())
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return string(plaintext)
}

func generateAESKey(passphrase string, keySize int) []byte {
	salt := []byte(config.Config.JWT.Secret) // Use a proper random salt in production
	iterations := 4096

	key := pbkdf2.Key([]byte(passphrase), salt, iterations, keySize, sha256.New)
	return key
}
