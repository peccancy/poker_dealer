package jwt

import (
	"github.com/peccancy/authorisation-service/models"

	jwtgo "github.com/dgrijalva/jwt-go"
)

type IDTokenClaims struct {
	jwtgo.StandardClaims
	UserID    string           `json:"userId"`
	Account   string           `json:"account"`
	TokenType models.TokenType `json:"type"`
	Context   string           `json:"context"`
	// next fields are only for user token id case
	GivenName  string `json:"given_name,omitempty"`
	FamilyName string `json:"family_name,omitempty"`
	Email      string `json:"email,omitempty"`
	Source     string `json:"source,omitempty"`
}

type ServiceTokenClaims struct {
	UserID      string             `json:"userId"`
	Account     string             `json:"account"`
	TokenType   models.TokenType   `json:"type"`
	Permissions models.Permissions `json:"permissions"`
	Context     string             `json:"context"`
	jwtgo.StandardClaims
}
