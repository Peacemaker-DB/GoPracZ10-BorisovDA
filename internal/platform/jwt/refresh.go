package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RefreshManager struct {
	secret []byte
	ttl    time.Duration
}

func NewRefresh(secret []byte, ttl time.Duration) *RefreshManager {
	return &RefreshManager{secret: secret, ttl: ttl}
}

func (m *RefreshManager) Sign(userID int64) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  userID,
		"iat":  now.Unix(),
		"exp":  now.Add(m.ttl).Unix(),
		"type": "refresh",
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(m.secret)
}

func (m *RefreshManager) Parse(token string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return m.secret, nil
	})
	if err != nil || !t.Valid {
		return nil, err
	}

	claims := t.Claims.(jwt.MapClaims)

	// refresh токен должен иметь type=refresh
	if claims["type"] != "refresh" {
		return nil, jwt.ErrTokenMalformed
	}

	return claims, nil
}
