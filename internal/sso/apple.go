package sso

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrEmptyKey   = errors.New("empty key")
	ErrInvalidKey = errors.New("invalid key")
)

func AppleSecret(clientID, teamID, keyID, keySecret string) (string, error) {
	// Decode a pem encoded PKCS8 PrivateKey
	block, _ := pem.Decode([]byte(keySecret))
	if block == nil {
		return "", ErrEmptyKey
	}

	pkUntyped, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	pk, ok := pkUntyped.(*ecdsa.PrivateKey)
	if !ok {
		return "", ErrInvalidKey
	}

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": teamID,
		"iat": now.Unix(),
		"exp": now.Add(time.Hour*24*30*6 - time.Minute).Unix(),
		"aud": "https://appleid.apple.com",
		"sub": clientID,
	})

	token.Header["kid"] = keyID

	return token.SignedString(pk)
}
