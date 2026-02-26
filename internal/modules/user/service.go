package user

import (
	"errors"
	"image-processing-service/internal/shared/utils"

	"gorm.io/gorm"
)

var (
	ErrNotFound        = utils.NewError(404, "USER_NOT_FOUND", "Usuario no encontrado", nil)
	ErrAlreadyExists   = utils.NewError(409, "USER_ALREADY_EXISTS", "No es posible registrar este correo. Intenta con otro.", nil)
	ErrInvalidPassword = utils.NewError(403, "INVALID_PASSWORD", "La contraseña actual no concide", nil)
)

type Service interface {
	GetByID(id string) (*User, error)
	GetAll() ([]*User, error)
	Update(id string, req UpdateUserRequest) (*User, error)
	UpdatePassword(id string, req UpdatePasswordUserRequest) (*User, error)
	Delete(id string) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetByID(id string) (*User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *service) GetAll() ([]*User, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *service) Update(id string, req UpdateUserRequest) (*User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if req.Email != nil && *req.Email != user.Email {
		existingUser, _ := s.repo.GetByEmail(*req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, ErrAlreadyExists
		}

		user.Email = *req.Email
	}

	if req.Name != nil {
		user.Name = *req.Name
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) UpdatePassword(id string, req UpdatePasswordUserRequest) (*User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if !utils.CheckPasswordHash(req.CurrentPassword, user.Password) {
		return nil, ErrInvalidPassword
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	if err := s.repo.UpdatePassword(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}
