package user

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
	GetByID(id string) (*User, error)
	GetAll() ([]*User, error)
	Update(user *User) error
	UpdatePassword(user *User) error
	Delete(id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *repository) GetByEmail(email string) (*User, error) {
	var user User

	if err := r.db.Unscoped().
		Where("email = ?", email).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *repository) GetByID(id string) (*User, error) {
	var user User

	if err := r.db.Where("ID = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) GetAll() ([]*User, error) {
	var users []*User

	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repository) Update(user *User) error {
	return r.db.Model(user).
		Omit("password", "created_at", "id").
		Save(user).Error
}

func (r *repository) UpdatePassword(user *User) error {
	return r.db.Model(user).
		Select("password").
		Save(user).Error
}

func (r *repository) Delete(id string) error {
	result := r.db.Delete(&User{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
