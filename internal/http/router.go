package router

import (
	"net/http"

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
	jwtv := jwt.NewHS256(cfg.JWTSecret, cfg.JWTTTL)
	svc := core.NewService(userRepo, jwtv)

	r.Post("/api/v1/login", svc.LoginHandler)

	r.Group(func(priv chi.Router) {
		priv.Use(middleware.AuthN(jwtv))
		priv.Use(middleware.AuthZRoles("admin", "user"))
		priv.Get("/api/v1/me", svc.MeHandler)
	})

	r.Group(func(admin chi.Router) {
		admin.Use(middleware.AuthN(jwtv))
		admin.Use(middleware.AuthZRoles("admin"))
		admin.Get("/api/v1/admin/stats", svc.AdminStats)
	})

	return r
}
