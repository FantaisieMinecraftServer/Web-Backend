package lib

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// compress compresses the input string using zlib
func Compress(data string) (string, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	_, err := w.Write([]byte(data))
	if err != nil {
		return "", err
	}
	w.Close()

	return buf.String(), nil
}

// decompress decompresses the input string using zlib
func Decompress(data string) (string, error) {
	r, err := zlib.NewReader(bytes.NewReader([]byte(data)))
	if err != nil {
		return "", err
	}
	defer r.Close()

	var out bytes.Buffer
	_, err = io.Copy(&out, r)
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

// encrypt encrypts the input string with the given key using AES
func Encrypt(plaintext, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decrypt decrypts the base64 encoded ciphertext with the given key using AES
func Decrypt(ciphertext, key string) (string, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return string(ciphertextBytes), nil
}

// compressEncrypt compresses, then encrypts the input string
func CompressEncrypt(data, key string) (string, error) {
	compressedData, err := Compress(data)
	if err != nil {
		return "", err
	}

	encryptedData, err := Encrypt(compressedData, key)
	if err != nil {
		return "", err
	}

	return encryptedData, nil
}

// decryptDecompress decrypts, then decompresses the input string
func DecryptDecompress(data, key string) (string, error) {
	decryptedData, err := Decrypt(data, key)
	if err != nil {
		return "", err
	}

	decompressedData, err := Decompress(decryptedData)
	if err != nil {
		return "", err
	}

	return decompressedData, nil
}
