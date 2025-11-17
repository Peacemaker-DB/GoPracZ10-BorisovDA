package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RS256 struct {
	privateKeys map[string]*rsa.PrivateKey
	publicKeys  map[string]*rsa.PublicKey
	ttl         time.Duration
	activeKID   string
}

func loadPrivate(path string) *rsa.PrivateKey {
	b, _ := os.ReadFile(path)
	k, _ := jwt.ParseRSAPrivateKeyFromPEM(b)
	return k
}

func loadPublic(path string) *rsa.PublicKey {
	b, _ := os.ReadFile(path)
	k, _ := jwt.ParseRSAPublicKeyFromPEM(b)
	return k
}

func NewRS256(ttl time.Duration) *RS256 {
	return &RS256{
		privateKeys: map[string]*rsa.PrivateKey{
			"k1": loadPrivate("private1.pem"),
			"k2": loadPrivate("private2.pem"),
		},
		publicKeys: map[string]*rsa.PublicKey{
			"k1": loadPublic("public1.pem"),
			"k2": loadPublic("public2.pem"),
		},
		ttl:       ttl,
		activeKID: "k2",
	}
}

func (r *RS256) Sign(userID int64, email, role string) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"role":  role,
		"iat":   now.Unix(),
		"exp":   now.Add(r.ttl).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	t.Header["kid"] = r.activeKID

	priv := r.privateKeys[r.activeKID]

	return t.SignedString(priv)
}

func (r *RS256) Parse(tokenStr string) (jwt.MapClaims, error) {
	parser := jwt.NewParser()
	token, err := parser.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, errors.New("no kid header")
		}

		key := r.publicKeys[kid]
		if key == nil {
			return nil, fmt.Errorf("unknown kid: %s", kid)
		}

		return key, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}
