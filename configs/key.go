package configs

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
)

type JwtService struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func LoadKeys() (*JwtService, error) {
	privateKey, err := LoadPrevKey()
	if err != nil {
		fmt.Printf("Failed to load Private Key : %s", err)
		return nil, err
	}
	publicKey, err := LoadPublicKey()
	if err != nil {
		fmt.Printf("Failed to load Private Key : %s", err)
		return nil, err
	}

	ks := &JwtService{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}

	return ks, nil

}

func LoadPublicKey() (*rsa.PublicKey, error) {
	publicKeyBase64 := os.Getenv("PUBLIC_KEY_BASE64")
	if publicKeyBase64 == "" {
		log.Fatalf("PUBLIC_KEY_BASE64 not found")
		return nil, nil
	}
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

func LoadPrevKey() (*rsa.PrivateKey, error) {
	privKeyBase64 := os.Getenv("PRIVATE_KEY_BASE64")
	if privKeyBase64 == "" {
		log.Fatalf("PRIVATE_KEY_BASE64 not found")
		return nil, nil
	}
	privKeyBytes, err := base64.StdEncoding.DecodeString(privKeyBase64)
	if err != nil {
		return nil, err
	}

	// Parse kunci privat dari PEM
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privKeyBytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
