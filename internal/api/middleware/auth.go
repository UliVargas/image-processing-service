package middleware

import (
	"context"
	"image-processing-service/internal/modules/session"
	"image-processing-service/internal/shared/auth"
	"image-processing-service/internal/shared/utils"
	"net/http"
	"strings"
)

var (
	ErrInvalidToken    = utils.NewError(401, "INVALID_TOKEN", "Token invalido", nil)
	ErrInvalidIDFormat = utils.NewError(400, "INVALID_ID", "El formato del identificador proporcionado es incorrecto", nil)
)

type authMiddleware struct {
	tokenManager *auth.TokenManager
	sessionSrv   session.Service
}

func NewAuthMiddleware(tm *auth.TokenManager, sess session.Service) *authMiddleware {
	return &authMiddleware{
		tokenManager: tm,
		sessionSrv:   sess,
	}
}

func (m *authMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := m.extractTokenFromHeader(r)
		if err != nil {
			utils.HandleError(w, err)
			return
		}

		claims, err := m.tokenManager.Validate(token)
		if err != nil {
			utils.HandleError(w, err)
			return
		}

		_, err = m.sessionSrv.IsValid(claims.ID, session.Access)
		if err != nil {
			utils.HandleError(w, ErrInvalidToken)
			return
		}

		userID := claims.UserID
		if !utils.IsValidID(userID) {
			utils.HandleError(w, ErrInvalidIDFormat)
			return
		}
		jti := claims.ID
		if !utils.IsValidID(jti) {
			utils.HandleError(w, ErrInvalidIDFormat)
			return
		}

		authUser := auth.AuthenticatedUser{
			UserID: userID,
			JTI:    jti,
		}

		ctx := context.WithValue(r.Context(), auth.AuthKey, authUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *authMiddleware) extractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrInvalidToken
	}

	token := strings.Split(authHeader, " ")
	if len(token) != 2 || token[0] != "Bearer" {
		return "", ErrInvalidToken
	}

	return token[1], nil
}
