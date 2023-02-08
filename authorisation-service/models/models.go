package models

type Permissions []string

type TokenType string

const (
	UserToken     TokenType = "user"
	ClientToken   TokenType = "client"
	InternalToken TokenType = "internal"
)
