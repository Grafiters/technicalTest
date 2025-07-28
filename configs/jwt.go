package configs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	tokens string
	Prefix string = "Bearer "
	err    error
)

func NewJwtConfig() (*JwtService, error) {
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

func (js *JwtService) GenerateTokenSession(idf string) (string, error) {
	claims := jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().UTC().Add(time.Hour).Unix(),
		"sub": "auth",
		"iss": "challenge",
		"idf": idf,
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokens, err = rawToken.SignedString(js.PrivateKey)
	if err != nil {
		Logger.Error("Error Signed : %s", err)
		return "", err
	}

	return tokens, nil
}

func (js *JwtService) DecodeTokenSession(tokenString string, target interface{}) error {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return js.PublicKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		recordJSON, err := json.Marshal(claims)
		if err != nil {
			return err
		}
		err = json.Unmarshal(recordJSON, &target)
		if err != nil {
			return err
		}

		return nil
	} else {
		return err
	}
}
