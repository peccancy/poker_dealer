package srvhttp

import (
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"strings"

	"github.com/peccancy/authorisation-service/models"
	log "github.com/sirupsen/logrus"
)

func safeEncode(p []byte) string {
	data := base64.URLEncoding.EncodeToString(p)
	return strings.TrimRight(data, "=")
}

// /.well-known/jwks.json
func (s *HTTPSrv) getPublicKey(w http.ResponseWriter, r *http.Request) {
	type Key struct {
		Kid string `json:"kid"`
		Kty string `json:"kty"`
		N   string `json:"n"`
		E   string `json:"e"`
	}
	type keysSet struct {
		Keys []*Key `json:"keys"`
		KTY  string `json:"kty"`
	}

	pKey := s.service.GetPublicKey()

	key := &Key{
		Kid: s.service.GetKeyID(),
		Kty: "RSA",
		N:   safeEncode(pKey.N.Bytes()),
		E:   safeEncode(big.NewInt(int64(pKey.E)).Bytes()),
	}

	res := keysSet{Keys: []*Key{key}}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Error("get public key err:", err.Error())
	}
}

// /.well-known/jwks.json
func (s *HTTPSrv) getServiceToken(w http.ResponseWriter, r *http.Request) {
	token, err := s.service.BuildServiceToken(getIDToken(r))
	if err != nil {
		switch err.(type) {
		case models.ErrNotValidIDToken:
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Error("gqlsrv: failed to build service token:", err)
			return
		}
	}
	w.Write([]byte(token)) //nolint
}

func getIDToken(r *http.Request) string {
	const header = "Authorization"
	parts := strings.Split(r.Header.Get(header), " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}
