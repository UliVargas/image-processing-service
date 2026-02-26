package session

import (
	"errors"
	"image-processing-service/internal/shared/utils"
	"time"

	"gorm.io/gorm"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
)

var (
	ErrNotFound       = utils.NewError(404, "SESSION_NOT_FOUND", "Session no encontrado", nil)
	ErrInvalidSession = utils.NewError(401, "INVALID_SESSION", "Session Inválida", nil)
)

type Service interface {
	Create(req CreateSessionRequest) (*Session, error)
	Delete(jti string) error
	IsValid(req string, t TokenType) (*Session, error)
	RenewSession(req UpdateSessionRequest) (*Session, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Create(req CreateSessionRequest) (*Session, error) {
	tokenHash := utils.GenerateSHA256(req.TokenHash)
	newSession := &Session{
		ID:        utils.GenerateID(),
		TokenHash: tokenHash,
		AccessJti: req.AccessJti,
		UserID:    req.UserID,
		ExpiresAt: req.ExpiresAt,
	}

	err := s.repo.Create(newSession)
	if err != nil {
		return nil, err
	}

	return newSession, nil
}

func (s *service) Delete(jti string) error {
	return s.repo.Delete(jti)
}

func (s *service) IsValid(req string, t TokenType) (*Session, error) {
	if t == Refresh {
		hashedToken := utils.GenerateSHA256(req)
		sess, err := s.repo.FindByOne(hashedToken)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		if time.Now().After(sess.ExpiresAt) {
			return nil, ErrInvalidSession
		}

		return sess, nil
	}

	sess, err := s.repo.FindByOne(req)
	if err != nil {
		return nil, err
	}

	if time.Now().After(sess.ExpiresAt) {
		return nil, ErrInvalidSession
	}

	return sess, nil
}

func (s *service) RenewSession(req UpdateSessionRequest) (*Session, error) {
	newHashedToken := utils.GenerateSHA256(req.NewTokenHash)

	updatedSession := &Session{
		ID:        req.SessionID,
		TokenHash: newHashedToken,
		AccessJti: req.NewAccessJti,
		ExpiresAt: req.ExpiresAt,
	}

	err := s.repo.Update(updatedSession)
	if err != nil {
		return nil, err
	}

	return updatedSession, nil
}
