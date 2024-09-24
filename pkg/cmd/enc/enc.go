package enc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Option struct {
	AES256CBC bool
	Pass      string
	P         bool
	D         bool
	Md        string
	In        string
	Out       string
	// TODO support -k key and -iv IV
	// It seems how generation of key/IV depends on the implementation
}

func Run(o Option) error {
	log.WithField("aes256cbc", o.AES256CBC).
		WithField("pass", o.Pass).
		WithField("p", o.P).
		WithField("d", o.D).
		WithField("md", o.Md).
		WithField("in", o.In).
		WithField("out", o.Out).
		Debug("enc")

	if !strings.HasPrefix(o.Pass, "file:") {
		return fmt.Errorf("failed to get the file from pass: %s", o.Pass)
	}
	file := strings.Replace(o.Pass, "file:", "", 1)

	bytes, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	key, err := base64.StdEncoding.DecodeString(string(bytes))

	d := sha256.New()
	_, err = d.Write(key)
	if err != nil {
		return fmt.Errorf("failed to hash: %w", err)
	}
	iv := d.Sum(nil)[:aes.BlockSize]
	if o.P {
		fmt.Printf("key=")
		for _, b := range key {
			fmt.Printf("%02X", b)
		}
		fmt.Println()
		fmt.Printf("iv =")
		for _, b := range iv {
			fmt.Printf("%02X", b)
		}
		fmt.Println()
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create aes cipher: %w", err)
	}

	inBytes, err := os.ReadFile(o.In)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	log.WithField("len(inBytes)", len(inBytes)).
		WithField("block.BlockSize()", block.BlockSize()).
		WithField("len(key)", len(key)).
		WithField("len(iv)", len(iv)).Debug("Print length")

	if o.D {
		if len(inBytes)%aes.BlockSize != 0 {
			return fmt.Errorf("valid input")
		}

		mode := cipher.NewCBCDecrypter(block, iv)
		plaintext := make([]byte, len(inBytes))
		mode.CryptBlocks(plaintext, inBytes)
		plaintext = unpad(plaintext)
		if err := os.WriteFile(o.Out, plaintext, 0644); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
		log.WithField("out", o.Out).Info("Saved output")
	} else {
		plaintext := pad(inBytes, aes.BlockSize)
		ciphertext := make([]byte, len(plaintext))
		mode := cipher.NewCBCEncrypter(block, iv)
		mode.CryptBlocks(ciphertext, plaintext)
		if err := os.WriteFile(o.Out, ciphertext, 0644); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
		log.WithField("out", o.Out).Info("Saved output")
	}
	return nil
}

func pad(data []byte, blockSize int) []byte {
	n := blockSize - len(data)%blockSize
	padding := bytes.Repeat([]byte{byte(n)}, n)
	return append(data, padding...)
}

func unpad(data []byte) []byte {
	n := int(data[len(data)-1])
	return data[:len(data)-n]
}
