package srvgrpc

type Service interface {
	BuildServiceToken(idTokenStr string) (string, error)
}
