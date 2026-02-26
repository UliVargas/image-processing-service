package main

import (
	"encoding/json"
	"image-processing-service/internal/api/middleware"
	"image-processing-service/internal/modules/auth"
	"image-processing-service/internal/modules/session"
	"image-processing-service/internal/modules/user"
	tokenManager "image-processing-service/internal/shared/auth"
	"image-processing-service/internal/shared/config"
	"image-processing-service/internal/shared/database"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	// ==========================================
	// Carga de variables de entorno
	// ==========================================
	cfg := config.NewEnv()

	// ==========================================
	// Configuración de JWT
	// ==========================================
	m := tokenManager.NewTokenManager(cfg.SecretKey, time.Hour*24*7)
	// ==========================================
	// Conexión a base de datos
	// ==========================================
	db := database.NewConection(cfg.DatabaseURL)

	// ==========================================
	// Modulo de Usuario
	// ==========================================
	userRepo := user.NewRepository(db)
	userSvc := user.NewService(userRepo)
	userHdl := user.NewHandler(userSvc)

	// ==========================================
	// Modulo de Sesion
	// ==========================================
	sessionRepo := session.NewRepository(db)
	sessionSvc := session.NewService(sessionRepo)

	// ==========================================
	// Modulo de Auth
	// ==========================================
	authSvc := auth.NewService(userRepo, sessionSvc, m)
	authHdl := auth.NewHandler(authSvc)

	// ==========================================
	// Se crear una instancia del enrutador
	// ==========================================
	r := chi.NewRouter()

	// ==========================================
	// Se dan de alta middlewares para autenticación, la impresión de rutas visitadas y el recuperador.
	// Este último por si hay un error no se detenga el servidor
	// ==========================================
	r.Use(chi_middleware.Logger)
	r.Use(chi_middleware.Recoverer)
	authMiddleware := middleware.NewAuthMiddleware(m, sessionSvc)

	// ==========================================
	// Se registra la ruta
	// ==========================================
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		res := Response{Message: "pong"}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	})

	// ==========================================
	// Registro de rutas de usuarios
	// ==========================================
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1/auth", func(router chi.Router) {
			router.Post("/signup", authHdl.SignUp)
			router.Post("/signin", authHdl.SignIn)
			router.With(authMiddleware.Authenticate).Post("/signout", authHdl.SignOut)
			router.Post("/renew-session", authHdl.RenewSession)
		})

		r.Route("/v1/users", func(router chi.Router) {
			router.Use(authMiddleware.Authenticate)
			router.Get("/", userHdl.GetAll)

			router.Get("/{id}", userHdl.GetByID)
			router.Patch("/{id}", userHdl.Update)
			router.Patch("/change-password/me", userHdl.UpdatePassword)
			router.Delete("/{id}", userHdl.Delete)
		})
	})

	// ==========================================
	// Se inicia servidor en el puerto 3000
	// ==========================================
	addr := ":" + cfg.Port
	log.Printf("Iniciando servidor en el puerto %s", cfg.Port)
	http.ListenAndServe(addr, r)
}
