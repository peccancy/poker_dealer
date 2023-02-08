package service

import (
	"crypto/rsa"
	"github.com/peccancy/authorisation-service/jwt"
	"time"
)

type (
	exchanger interface {
		Exchange(idTokenClaims jwt.IDTokenClaims, key *rsa.PrivateKey, keyID string) (string, error)
	}

	Service struct {
		exchanger  exchanger
		issuer     string
		key        *rsa.PrivateKey
		pubKey     *rsa.PublicKey
		keyID      string
		defaultTTL time.Duration
		readiness  bool
	}
)

func New(
	ex exchanger,
	issuer string,
	key *rsa.PrivateKey,
	keyID string,
	defaultTTL time.Duration,
) *Service {
	// build service
	return &Service{
		exchanger:  ex,
		issuer:     issuer,
		key:        key,
		pubKey:     &key.PublicKey,
		keyID:      keyID,
		defaultTTL: defaultTTL,
		readiness:  true,
	}
}
