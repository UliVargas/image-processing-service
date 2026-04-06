package api

import (
	"encoding/json"
	"net/http"

	"image-processing-service/internal/api/middleware"
	"image-processing-service/internal/modules/auth"
	"image-processing-service/internal/modules/file"
	"image-processing-service/internal/modules/user"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
)

type pingResponse struct {
	Message string `json:"message"`
}

func NewRouter(
	authMW *middleware.AuthMiddleware,
	authHdl auth.Handler,
	userHdl user.Handler,
	fileHdl file.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chi_middleware.Logger)
	r.Use(chi_middleware.Recoverer)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pingResponse{Message: "pong"})
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1/auth", func(r chi.Router) {
			r.Post("/signup", authHdl.SignUp)
			r.Post("/signin", authHdl.SignIn)
			r.With(authMW.Authenticate).Post("/signout", authHdl.SignOut)
			r.Post("/renew-session", authHdl.RenewSession)
		})

		r.Route("/v1/users", func(r chi.Router) {
			r.Use(authMW.Authenticate)
			r.Get("/", userHdl.GetAll)
			r.Get("/{id}", userHdl.GetByID)
			r.Patch("/{id}", userHdl.Update)
			r.Patch("/change-password/me", userHdl.UpdatePassword)
			r.Delete("/{id}", userHdl.Delete)
		})

		r.Route("/v1/files", func(r chi.Router) {
			r.Use(authMW.Authenticate)
			r.Get("/", fileHdl.ListMine)
			r.Get("/*", fileHdl.GetOne)
			r.Post("/", fileHdl.Upload)
		})
	})

	return r
}
