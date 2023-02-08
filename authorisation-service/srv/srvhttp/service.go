package srvhttp

import "crypto/rsa"

type Service interface {
	GetKeyID() string
	GetPublicKey() *rsa.PublicKey
	BuildServiceToken(idTokenStr string) (string, error)
}
