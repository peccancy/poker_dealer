package servicetoken

import (
	"crypto/rsa"
	"github.com/peccancy/authorisation-service/jwt"
	"github.com/pkg/errors"
)

func clientExchange(client boClient, tk tokenData, idTokenClaims jwt.IDTokenClaims,
	key *rsa.PrivateKey, keyID string) (string, error) {
	permissions, err := client.GetTenantMachineTokenPermissions(idTokenClaims.Account, idTokenClaims.UserID)
	if err != nil {
		return "", errors.Wrap(err, "failed to read permissions")
	}

	claims := createTokenClaims(idTokenClaims, permissions, idTokenClaims.Context, tk.Issuer, tk.DefaultTTL, nil)

	// Create the JWT string
	return jwt.BuildSignedServiceToken(key, keyID, claims)
}
