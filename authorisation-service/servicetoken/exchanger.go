package servicetoken

import (
	"time"

	"crypto/rsa"

	jwtgo "github.com/dgrijalva/jwt-go"
	core "github.com/peccancy/authorisation-service/providers/grpc"

	"github.com/peccancy/authorisation-service/jwt"
	"github.com/peccancy/authorisation-service/models"
)

type (
	Exchanger struct {
		boClient    boClient
		usersClient usersClient
		tk          tokenData
		skipUpsert  bool
	}

	boClient interface {
		GetTenantMachineTokenPermissions(accountID, clientID string) (models.Permissions, error)
		GetInternalMachineTokenPermissions(clientID string) (models.Permissions, error)
	}

	usersClient interface {
		GetUserPermissions(accountID, userID, email string) (models.Permissions, string, error)
		Upsert(accountID, userID, givenName, familyName, email, source string) error
		GetTokenContext(accountID, userID, email string) (string, error)
	}

	tokenData struct {
		Issuer     string
		DefaultTTL time.Duration
		PubKey     *rsa.PublicKey
	}
)

func NewExchanger(b boClient, u usersClient, issuer string, defaultTTL time.Duration,
	pubKey *rsa.PublicKey, skipUpsert bool) Exchanger {
	return Exchanger{
		boClient:    b,
		usersClient: u,
		tk: tokenData{
			Issuer:     issuer,
			DefaultTTL: defaultTTL,
			PubKey:     pubKey,
		},
		skipUpsert: skipUpsert,
	}
}

func (e Exchanger) Exchange(idTokenClaims jwt.IDTokenClaims, key *rsa.PrivateKey, keyID string) (string, error) {
	switch idTokenClaims.TokenType {
	case models.UserToken:
		return userExchange(e.skipUpsert, e.usersClient, e.tk, idTokenClaims, key, keyID)
	case models.ClientToken:
		return clientExchange(e.boClient, e.tk, idTokenClaims, key, keyID)
	case models.InternalToken:
		return internalExchange(e.boClient, e.tk, idTokenClaims, key, keyID)
	}
	return "", core.ErrNotSupportedTokenType
}

func createTokenClaims(idTokenClaims jwt.IDTokenClaims, permissions models.Permissions, context string, issuer string,
	defaultTTL time.Duration, customUserID *string) jwt.ServiceTokenClaims {

	now := time.Now().UTC()
	userID := idTokenClaims.UserID
	if customUserID != nil {
		userID = *customUserID
	}

	if permissions == nil {
		permissions = make([]string, 0)
	}

	// Create the JWT claims, which includes the username and expiry time
	return jwt.ServiceTokenClaims{
		UserID:      userID,
		Account:     idTokenClaims.Account,
		Context:     context,
		Permissions: permissions,
		TokenType:   idTokenClaims.TokenType,
		StandardClaims: jwtgo.StandardClaims{
			Issuer:   issuer,
			IssuedAt: now.Unix(),
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: time.Now().Add(defaultTTL).Unix(),
		},
	}
}
