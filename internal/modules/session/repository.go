package session

import "gorm.io/gorm"

type Repository interface {
	Create(session *Session) error
	FindByOne(tokenHash string) (*Session, error)
	Update(session *Session) error
	Delete(jti string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(session *Session) error {
	return r.db.Create(session).Error
}

func (r *repository) Delete(jti string) error {
	result := r.db.Delete(&Session{}, "access_jti = ?", jti)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *repository) FindByOne(data string) (*Session, error) {
	var session Session

	if err := r.db.Where("token_hash = ?", data).
		Or("access_jti = ?", data).
		First(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *repository) Update(session *Session) error {
	return r.db.Model(session).
		Select("token_hash", "access_jti", "expires_at").
		Updates(session).Error
}
