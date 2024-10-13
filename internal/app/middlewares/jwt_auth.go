package middlewares

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

var (
	auth0Domain = fmt.Sprintf("https://%s/", os.Getenv("AUTH0_DOMAIN"))
	audience    = os.Getenv("AUTH0_AUDIENCE")
)

// JWKS structure to hold the keys from Auth0
type Jwks struct {
	Keys []struct {
		Kty string `json:"kty"`
		Kid string `json:"kid"`
		Use string `json:"use"`
		N   string `json:"n"`
		E   string `json:"e"`
	} `json:"keys"`
}

// Fetches the JWKS from Auth0
func fetchJWKs() (*Jwks, error) {
	resp, err := http.Get(fmt.Sprintf("%s.well-known/jwks.json", auth0Domain))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jwks Jwks
	err = json.Unmarshal(body, &jwks)
	if err != nil {
		return nil, err
	}

	return &jwks, nil
}

// Get the RSA Public Key from the JWKS
func getPublicKey(token *jwt.Token) (*rsa.PublicKey, error) {
	jwks, err := fetchJWKs()
	if err != nil {
		return nil, err
	}

	for _, key := range jwks.Keys {
		if kid, ok := token.Header["kid"].(string); ok && key.Kid == kid {
			// Build the public key using the standard base64 decoding
			nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
			if err != nil {
				return nil, err
			}
			eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
			if err != nil {
				return nil, err
			}

			var eInt int
			if len(eBytes) == 3 {
				eInt = int(eBytes[0])<<16 | int(eBytes[1])<<8 | int(eBytes[2])
			} else if len(eBytes) == 1 {
				eInt = int(eBytes[0])
			} else {
				return nil, errors.New("invalid exponent length")
			}

			pubKey := &rsa.PublicKey{
				N: new(big.Int).SetBytes(nBytes),
				E: eInt,
			}
			return pubKey, nil
		}
	}
	return nil, errors.New("unable to find appropriate key")
}
