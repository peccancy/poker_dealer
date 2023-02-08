package servicetoken

import (
	"crypto/rsa"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/peccancy/authorisation-service/jwt"
)

func userExchange(skipUpsert bool, client usersClient, tk tokenData, idTokenClaims jwt.IDTokenClaims, key *rsa.PrivateKey,
	keyID string) (string, error) {
	if !skipUpsert {
		err := client.Upsert(idTokenClaims.Account, idTokenClaims.UserID, idTokenClaims.GivenName,
			idTokenClaims.FamilyName, idTokenClaims.Email, idTokenClaims.Source)
		if err != nil {
			log.WithError(err).Error("failed to upsert the user")
		}
	}

	permissions, swapUserID, err := client.GetUserPermissions(
		idTokenClaims.Account,
		idTokenClaims.UserID,
		idTokenClaims.Email,
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to read permissions")
	}

	tokenCtx, err := client.GetTokenContext(
		idTokenClaims.Account,
		idTokenClaims.UserID,
		idTokenClaims.Email,
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to read token context")
	}

	if tokenCtx == "" {
		tokenCtx = idTokenClaims.Context
	}

	serviceClaims := createTokenClaims(idTokenClaims, permissions, tokenCtx, tk.Issuer, tk.DefaultTTL, &swapUserID)

	// Create the JWT string
	return jwt.BuildSignedServiceToken(key, keyID, serviceClaims)
}
