package auth

import (
	"errors"
	"image-processing-service/internal/modules/session"
	"image-processing-service/internal/modules/user"
	"time"

	"image-processing-service/internal/shared/auth"
	"image-processing-service/internal/shared/utils"

	"gorm.io/gorm"
)

var (
	ErrNotFound           = utils.NewError(404, "SESSION_NOT_FOUND", "Sesión no encontrada", nil)
	ErrAlreadyExists      = utils.NewError(409, "USER_ALREADY_EXISTS", "No es posible registrar este correo. Intenta con otro.", nil)
	ErrInvalidCredentials = utils.NewError(401, "INVALID_CREDENTIALS", "Credenciales inválidas", nil)
)

type Service interface {
	SignUp(req RegisterRequest) (*user.User, error)
	SignIn(req LoginRequest) (*Auth, error)
	SignOut(jti string) error
	RenewSession(refreshToken string) (*Auth, error)
}

type service struct {
	userRepo     user.Repository
	sessionSrv   session.Service
	tokenManager *auth.TokenManager
}

func NewService(r user.Repository, sSrv session.Service, m *auth.TokenManager) Service {
	return &service{userRepo: r, sessionSrv: sSrv, tokenManager: m}
}

func (s *service) SignUp(req RegisterRequest) (*user.User, error) {
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, ErrAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	newUser := &user.User{
		ID:       utils.GenerateID(),
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *service) SignIn(req LoginRequest) (*Auth, error) {
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser == nil || existingUser.DeletedAt.Valid {
		return nil, ErrInvalidCredentials
	}

	isValid := utils.CheckPasswordHash(req.Password, existingUser.Password)
	if !isValid {
		return nil, ErrInvalidCredentials
	}

	result, err := s.tokenManager.GeneratePair(existingUser.ID)
	if err != nil {
		return nil, err
	}

	_, err = s.sessionSrv.Create(session.CreateSessionRequest{
		TokenHash: result.RefreshToken,
		AccessJti: result.JTI,
		UserID:    existingUser.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	})

	return &Auth{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}

func (s *service) SignOut(jti string) error {
	if err := s.sessionSrv.Delete(jti); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func (s *service) RenewSession(refreshToken string) (*Auth, error) {
	sess, err := s.sessionSrv.IsValid(refreshToken, session.Refresh)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	result, err := s.tokenManager.GeneratePair(sess.UserID)
	if err != nil {
		return nil, err
	}

	_, err = s.sessionSrv.RenewSession(session.UpdateSessionRequest{
		SessionID:    sess.ID,
		NewTokenHash: result.RefreshToken,
		NewAccessJti: result.JTI,
		ExpiresAt:    time.Now().Add(time.Hour * 24 * 7),
	})
	if err != nil {
		return nil, err
	}

	return &Auth{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}
