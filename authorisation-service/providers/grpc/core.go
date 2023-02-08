package grpc

import "github.com/pkg/errors"

type Config struct {
	Target string
}

var ErrNotSupportedTokenType = errors.New("not supported token type")

type ErrInit string

func (e ErrInit) Error() string {
	return "initialisation error: " + string(e)
}

type ErrConn string

func (e ErrConn) Error() string {
	return "connection error: " + string(e)
}
