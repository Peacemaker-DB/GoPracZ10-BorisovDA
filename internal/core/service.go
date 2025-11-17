package core

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"example.com/goprac10-borisovda/internal/http/middleware"
	"example.com/goprac10-borisovda/internal/platform/jwt"
	"github.com/go-chi/chi/v5"
)

type UserInfo struct {
	ID    int64
	Email string
	Role  string
}

type userRepo interface {
	CheckPassword(email, pass string) (UserInfo, error)
}

type jwtSigner interface {
	Sign(userID int64, email, role string) (string, error)
}

type Service struct {
	repo      userRepo
	jwt       jwtSigner
	refresh   *jwt.RefreshManager
	blacklist map[string]int64 // refresh → exp
}

func NewService(r userRepo, j jwtSigner) *Service {
	refreshTTL := 7 * 24 * time.Hour
	return &Service{
		repo:      r,
		jwt:       j,
		refresh:   jwt.NewRefresh(j.(*jwt.HS256).Secret(), refreshTTL),
		blacklist: make(map[string]int64),
	}
}

func (s *Service) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var in struct{ Email, Password string }
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpError(w, 400, "invalid_credentials")
		return
	}

	u, err := s.repo.CheckPassword(in.Email, in.Password)
	if err != nil {
		httpError(w, 401, "unauthorized")
		return
	}

	access, err := s.jwt.Sign(u.ID, u.Email, u.Role)
	if err != nil {
		httpError(w, 500, "token_error")
		return
	}

	refresh, err := s.refresh.Sign(u.ID)
	if err != nil {
		httpError(w, 500, "token_error")
		return
	}

	jsonOK(w, map[string]any{
		"access":  access,
		"refresh": refresh,
	})
}

func (s *Service) MeHandler(w http.ResponseWriter, r *http.Request) {
	// клеймы положим в контекст в AuthN-мидлваре
	claims := r.Context().Value(middleware.CtxClaimsKey).(map[string]any)
	jsonOK(w, map[string]any{
		"id": claims["sub"], "email": claims["email"], "role": claims["role"],
	})
}

func (s *Service) AdminStats(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, map[string]any{"users": 2, "version": "1.0"})
}

// утилиты и ключ для контекста — экспортируем из middleware
type ctxClaims struct{}

func jsonOK(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
func httpError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (s *Service) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var in struct{ Refresh string }
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		httpError(w, 400, "invalid_refresh")
		return
	}

	// Проверяем blacklist
	if exp, ok := s.blacklist[in.Refresh]; ok && exp > time.Now().Unix() {
		httpError(w, 401, "revoked_refresh")
		return
	}

	claims, err := s.refresh.Parse(in.Refresh)
	if err != nil {
		httpError(w, 401, "bad_refresh")
		return
	}

	userID := int64(claims["sub"].(float64))

	// Кладём старый refresh в blacklist
	s.blacklist[in.Refresh] = int64(claims["exp"].(float64))

	// Создаём новую пару токенов
	access, _ := s.jwt.Sign(userID, "placeholder@example.com", "user")
	refresh, _ := s.refresh.Sign(userID)

	jsonOK(w, map[string]any{
		"access":  access,
		"refresh": refresh,
	})
}

var mockUsers = map[int64]map[string]any{
	1: {"id": 1, "email": "admin@example.com"},
	2: {"id": 2, "email": "user@example.com"},
}

func (s *Service) UserByID(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.CtxClaimsKey).(map[string]any)
	role := claims["role"].(string)
	sub := int64(claims["sub"].(float64))

	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if role == "user" && id != sub {
		httpError(w, 403, "forbidden_abac")
		return
	}

	u, ok := mockUsers[id]
	if !ok {
		httpError(w, 404, "not_found")
		return
	}

	jsonOK(w, u)
}
