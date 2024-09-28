package middlewares

import (
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

// Mock JWKS response for testing
var mockJWKS = `{
	"keys": [
		{
			"kty": "RSA",
			"kid": "1234",
			"use": "sig",
			"n": "sXchxkWIpIEV5MBRka8Zgp8FvXjhMQHzS5IdwOYO_JBfBo9Itiu1FrKN3dNmO6dZ9EXac3mfHlg50ztFLFttjZNdOSLDVseWtP0sXQqNGfr-J8VVcx0LIPZLrsw3FCrw58EMaxSmyRjWJYeV56CUPiBpRxz3A0EpAyIfUkNodLM",
			"e": "AQAB"
		}
	]
}`

// Mock HTTP server for JWKS
func setupJWKSMockServer(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
}

// TestGetPublicKey_Success tests successful retrieval of the RSA public key
func TestGetPublicKey_Success(t *testing.T) {
	// Mock token with header containing kid
	token := &jwt.Token{
		Header: map[string]interface{}{
			"kid": "1234",
		},
	}

	// Set up mock JWKS server
	server := setupJWKSMockServer(mockJWKS)
	defer server.Close()

	// Override the auth0Domain to point to the mock server with a trailing slash
	auth0Domain = server.URL + "/"

	// Call the getPublicKey function
	pubKey, err := getPublicKey(token)
	assert.NoError(t, err)
	assert.NotNil(t, pubKey)
	assert.IsType(t, &rsa.PublicKey{}, pubKey)
}

// TestGetPublicKey_KidNotFound tests the case where the kid is not found in the JWKS
func TestGetPublicKey_KidNotFound(t *testing.T) {
	// Mock token with a kid that does not exist
	token := &jwt.Token{
		Header: map[string]interface{}{
			"kid": "5678",
		},
	}

	// Set up mock JWKS server
	server := setupJWKSMockServer(mockJWKS)
	defer server.Close()

	// Override the auth0Domain to point to the mock server with a trailing slash
	auth0Domain = server.URL + "/"

	// Call the getPublicKey function
	pubKey, err := getPublicKey(token)
	assert.Error(t, err)
	assert.Nil(t, pubKey)
	assert.EqualError(t, err, "unable to find appropriate key")
}

// TestGetPublicKey_InvalidExponent tests the case where the exponent is invalid
func TestGetPublicKey_InvalidExponent(t *testing.T) {
	// Modify mock JWKS to have an invalid exponent (not valid base64url encoding)
	invalidJWKS := `{
		"keys": [
			{
				"kty": "RSA",
				"kid": "1234",
				"use": "sig",
				"n": "sXchxkWIpIEV5MBRka8Zgp8FvXjhMQHzS5IdwOYO_JBfBo9Itiu1FrKN3dNmO6dZ9EXac3mfHlg50ztFLFttjZNdOSLDVseWtP0sXQqNGfr-J8VVcx0LIPZLrsw3FCrw58EMaxSmyRjWJYeV56CUPiBpRxz3A0EpAyIfUkNodLM",
				"e": "!!invalid!!"
			}
		]
	}`

	// Mock token with header containing kid
	token := &jwt.Token{
		Header: map[string]interface{}{
			"kid": "1234",
		},
	}

	// Set up mock JWKS server
	server := setupJWKSMockServer(invalidJWKS)
	defer server.Close()

	// Override the auth0Domain to point to the mock server
	auth0Domain = server.URL + "/"

	// Call the getPublicKey function
	pubKey, err := getPublicKey(token)
	assert.Error(t, err)
	assert.Nil(t, pubKey)
	assert.Contains(t, err.Error(), "illegal base64 data")
}
