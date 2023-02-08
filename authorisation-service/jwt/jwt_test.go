package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dgrijalva/jwt-go"
)

func Test_BuildTokenFlow(t *testing.T) {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	claims := ServiceTokenClaims{
		UserID:  "user id",
		Account: "acc name",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(rsaPrivateKey)
	assert.NoError(t, err)

	publicKeyBytes, err := ExportRsaPublicKeyAsPemStr(&rsaPrivateKey.PublicKey)
	assert.NoError(t, err)

	var claims2 ServiceTokenClaims
	tkn, err := jwt.ParseWithClaims(tokenString, &claims2, func(token *jwt.Token) (interface{}, error) {
		return ParseRsaPublicKeyFromPemStr(publicKeyBytes)
	})
	assert.NoError(t, err)
	assert.True(t, tkn.Valid)
	assert.Equal(t, claims, claims2)
}

func Test_BuildParseFlow(t *testing.T) {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	claims := ServiceTokenClaims{
		UserID:  "user id",
		Account: "acc id",
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
	}

	tokenString, err := BuildSignedServiceToken(rsaPrivateKey, "keyID", claims)
	assert.NoError(t, err)

	var parsedClaims IDTokenClaims
	err = ParseSignedIDToken(tokenString, &rsaPrivateKey.PublicKey, &parsedClaims)
	assert.NoError(t, err)
	assert.NoError(t, parsedClaims.Valid())

	assert.Equal(t, claims.UserID, parsedClaims.UserID)
	assert.Equal(t, claims.Account, parsedClaims.Account)
}

func Test_BuildParseUserIDTokenClaims(t *testing.T) {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	claims := IDTokenClaims{
		UserID:  "user id",
		Account: "acc id",
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
		GivenName:  "Dale",
		FamilyName: "Cooper",
		Email:      "black@lodge.com",
		Source:     "twinpeaks",
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "keyID"

	// Create the JWT string
	signedToken, err := token.SignedString(rsaPrivateKey)
	assert.NoError(t, err)

	var parsedClaims IDTokenClaims
	err = ParseSignedIDToken(signedToken, &rsaPrivateKey.PublicKey, &parsedClaims)
	assert.NoError(t, err)
	assert.NoError(t, parsedClaims.Valid())
	assert.Equal(t, claims, parsedClaims)
}
