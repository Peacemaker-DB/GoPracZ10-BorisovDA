package router

import (
	"net/http"
	"time"

	"example.com/goprac10-borisovda/internal/core"
	"example.com/goprac10-borisovda/internal/http/middleware"
	"example.com/goprac10-borisovda/internal/platform/config"
	"example.com/goprac10-borisovda/internal/platform/jwt"
	"example.com/goprac10-borisovda/internal/repo"
	"github.com/go-chi/chi/v5"
)

func Build(cfg config.Config) http.Handler {
	r := chi.NewRouter()

	userRepo := repo.NewUserMem()
	jwtv := jwt.NewRS256(15 * time.Minute)
	svc := core.NewService(userRepo, jwtv, cfg)

	r.Post("/api/v1/login", svc.LoginHandler)

	r.Group(func(priv chi.Router) {
		priv.Use(middleware.AuthN(jwtv))
		priv.Use(middleware.AuthZRoles("admin", "user"))
		priv.Get("/api/v1/me", svc.MeHandler)
		priv.Get("/api/v1/users/{id}", svc.UserByID)
	})

	r.Group(func(admin chi.Router) {
		admin.Use(middleware.AuthN(jwtv))
		admin.Use(middleware.AuthZRoles("admin"))
		admin.Get("/api/v1/admin/stats", svc.AdminStats)
	})

	r.Post("/api/v1/refresh", svc.RefreshHandler)
	r.With(middleware.RateLimitLogin).Post("/api/v1/login", svc.LoginHandler)
	return r
}
