package dgst

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type Option struct {
	Sha1      bool
	Sign      string
	Verify    string
	Signature string
	Out       string
}

func Run(o Option, key string) error {
	log.WithField("sha1", o.Sha1).
		WithField("sign", o.Sign).
		WithField("verify", o.Verify).
		WithField("signature", o.Signature).
		WithField("out", o.Out).
		WithField("key", key).
		Debug("pkeyutl")

	keyBytes, err := os.ReadFile(key)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	msgHash := sha1.New()
	_, err = msgHash.Write(keyBytes)
	if err != nil {
		return fmt.Errorf("failed to hash: %w", err)
	}
	msgHashSum := msgHash.Sum(nil)

	if o.Sign != "" {
		bytes, err := os.ReadFile(o.Sign)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		block, _ := pem.Decode(bytes)
		privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse the private key: %w", err)
		}

		signature, err := rsa.SignPSS(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA1, msgHashSum, nil)
		if err != nil {
			return fmt.Errorf("failed to sign: %w", err)
		}
		if err := os.WriteFile(o.Out, signature, 0644); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}
	if o.Verify != "" {
		bytes, err := os.ReadFile(o.Verify)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		block, _ := pem.Decode(bytes)
		publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse the public key: %w", err)
		}

		signature, err := os.ReadFile(o.Signature)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		err = rsa.VerifyPSS(publicKey.(*rsa.PublicKey), crypto.SHA1, msgHashSum, signature, nil)
		if err != nil {
			return fmt.Errorf("failed to verify: %w", err)
		}
		log.Info("Signature is verified")
	}
	return nil
}
