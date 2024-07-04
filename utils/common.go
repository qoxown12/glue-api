package utils

import (
	"Glue-API/model"
	"encoding/json"
	"log"
	"os"
	"runtime"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func HandleError(err error) (b bool) {
	if err != nil {
		// notice that we're using 1, so it will actually log where
		// the error happened, 0 = this function, we don't want that.
		_, filename, line, _ := runtime.Caller(1)
		log.Printf("[error] %s:%d %v", filename, line, err)
		b = true
	}
	return
}

// FancyHandleError this logs the function name as well.
func FancyHandleError(err error) (b bool) {
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		pc, filename, line, _ := runtime.Caller(1)

		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), filename, line, err)

		b = true
	}
	return
}

// Read the settings file.
func ReadConfFile() (settings model.Settings, err error) {
	content, err := os.ReadFile("./conf.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
		return
	}

	err = json.Unmarshal(content, &settings)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
		return
	}
	return
}

// 패딩을 추가하는 함수
func pkcs7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

// 패딩을 제거하는 함수
func pkcs7UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

// 텍스트를 암호화하는 함수
func encrypt(plainText, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plainText = pkcs7Padding(plainText, block.BlockSize())
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// 암호화된 텍스트를 복호화하는 함수
func decrypt(cipherText string, key []byte) (string, error) {
	cipherData, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(cipherData) < aes.BlockSize {
		return "", fmt.Errorf("cipherText too short")
	}

	iv := cipherData[:aes.BlockSize]
	cipherData = cipherData[aes.BlockSize:]

	if len(cipherData)%aes.BlockSize != 0 {
		return "", fmt.Errorf("cipherText is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherData, cipherData)

	plainText := pkcs7UnPadding(cipherData)
	return string(plainText), nil
}

// Password encryption method
func PasswordEncryption(plainText string) (encryptedText string, err error) {
	// 32바이트(256비트) 길이의 키 (AES-256)
	key := []byte("ablecloudablecloudablecloud12345")
	// 암호화
	encryptedText, err = encrypt([]byte(plainText), key)
	if err != nil {
		return "", err
	}
	return encryptedText, nil
}

// Password decryption method
func PasswordDecryption(encryptedText string) (decryptedText string, err error) {
	// 32바이트(256비트) 길이의 키 (AES-256)
	key := []byte("ablecloudablecloudablecloud12345")
	//복호화
	decryptedText, err = decrypt(encryptedText, key)
	if err != nil {
		return "", err
	}
	return decryptedText, nil
}