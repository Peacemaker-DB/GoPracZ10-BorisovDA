package repo

import (
	"errors"

	"example.com/goprac10-borisovda/internal/core"
	"golang.org/x/crypto/bcrypt"
)

type UserRecord struct {
	ID    int64
	Email string
	Role  string
	Hash  []byte
}

type UserMem struct {
	users map[string]UserRecord
}

func NewUserMem() *UserMem {
	hash := func(s string) []byte {
		h, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
		return h
	}

	return &UserMem{
		users: map[string]UserRecord{
			"admin@example.com": {ID: 1, Email: "admin@example.com", Role: "admin", Hash: hash("secret123")},
			"user@example.com":  {ID: 2, Email: "user@example.com", Role: "user", Hash: hash("secret123")},
		},
	}
}

var ErrNotFound = errors.New("user not found")
var ErrBadCreds = errors.New("bad credentials")

func (r *UserMem) CheckPassword(email, pass string) (core.UserInfo, error) {
	u, ok := r.users[email]
	if !ok {
		return core.UserInfo{}, ErrNotFound
	}

	if bcrypt.CompareHashAndPassword(u.Hash, []byte(pass)) != nil {
		return core.UserInfo{}, ErrBadCreds
	}

	return core.UserInfo{
		ID:    u.ID,
		Email: u.Email,
		Role:  u.Role,
	}, nil
}
