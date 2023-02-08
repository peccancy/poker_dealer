package srvgql

type Service interface {
	BuildServiceToken(idTokenStr string) (string, error)
}
