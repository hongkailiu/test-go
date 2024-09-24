package pkey

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/big"
	"os"
)

type Option struct {
	PubIn  bool
	PubOut bool
	Text   bool
	NoOut  bool
	In     string
	Out    string
}

func Run(o Option) error {
	log.WithField("pubout", o.PubOut).
		WithField("pubin", o.PubIn).
		WithField("noout", o.NoOut).
		WithField("text", o.Text).
		WithField("in", o.In).
		WithField("out", o.Out).
		Debug("pkey")

	if o.Out != "" && o.PubOut == true {
		bytes, err := os.ReadFile(o.In)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		block, _ := pem.Decode(bytes)
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse the private key: %w", err)
		}

		privateKey := key.(*rsa.PrivateKey)
		keyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
		if err != nil {
			return fmt.Errorf("failed to marshal the public key: %w", err)
		}

		publicKey := &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: keyBytes,
		}

		pubKeyBytes := pem.EncodeToMemory(publicKey)
		log.Debug("Public key generated")

		if err := os.WriteFile(o.Out, pubKeyBytes, 0644); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}

		log.WithField("file", o.Out).Info("Saved the public key")
	}
	if o.NoOut && o.Text {
		if o.PubIn {
			bytes, err := os.ReadFile(o.In)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			block, _ := pem.Decode(bytes)
			key, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				return fmt.Errorf("failed to parse the public key: %w", err)
			}
			publicKey := key.(*rsa.PublicKey)
			fmt.Printf("Private-Key: (%d bit)\n", publicKey.N.BitLen())
			printBigInt("modulus", *publicKey.N)

			fmt.Printf("publicExponent: %d (0x%x)\n", publicKey.E, publicKey.E)
		} else {
			bytes, err := os.ReadFile(o.In)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			block, _ := pem.Decode(bytes)
			key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return fmt.Errorf("failed to parse the private key: %w", err)
			}

			privateKey := key.(*rsa.PrivateKey)
			fmt.Printf("Private-Key: (%d bit, %d primes)\n", privateKey.N.BitLen(), len(privateKey.Primes))
			printBigInt("modulus", *privateKey.N)

			fmt.Printf("publicExponent: %d (0x%x)\n", privateKey.E, privateKey.E)

			printBigInt("privateExponent", *privateKey.D)

			for i, p := range privateKey.Primes {
				printBigInt(fmt.Sprintf("prime%d", i+1), *p)
			}

			privateKey.Precompute()
			printBigInt("exponent1", *privateKey.Precomputed.Dp)
			printBigInt("exponent2", *privateKey.Precomputed.Dq)
			printBigInt("coefficient", *privateKey.Precomputed.Qinv)
		}
	}
	return nil
}

func printBigInt(title string, n big.Int) {
	// TODO
	// determine when to insert 00 as the first byte
	// log.WithField("sign", n.Sign()).WithField("len(n.Bytes())", len(n.Bytes())).WithField("text", n.Text(16)).Debug("aaa===")
	fmt.Printf("%s:", title)

	for i, b := range n.Bytes() {
		if i%15 == 0 {
			fmt.Printf("\n    %02x", b)
		} else {
			fmt.Printf(":%02x", b)
		}
	}
	fmt.Println()
}
