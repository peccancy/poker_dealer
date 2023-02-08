package service

import (
	"crypto/rsa"
	"github.com/peccancy/authorisation-service/jwt"
	"github.com/peccancy/authorisation-service/models"
)

func (s *Service) GetPublicKey() *rsa.PublicKey {
	return s.pubKey
}

func (s *Service) GetKeyID() string {
	return s.keyID
}

func (s *Service) BuildServiceToken(idTokenStr string) (string, error) {
	var idTokenClaims jwt.IDTokenClaims
	err := jwt.ParseSignedIDToken(idTokenStr, s.pubKey, &idTokenClaims)
	if err != nil {
		return "", models.ErrNotValidIDToken(err.Error())
	}

	err = idTokenClaims.Valid()
	if err != nil {
		return "", models.ErrNotValidIDToken(err.Error())
	}

	return s.exchanger.Exchange(idTokenClaims, s.key, s.keyID)
}
