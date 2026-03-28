package file

import "gorm.io/gorm"

type Repository interface {
	Create(file *File) error
	FindOne(storageKey string) (*File, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(file *File) error {
	return r.db.Create(file).Error
}

func (r *repository) FindOne(storageKey string) (*File, error) {
	var file File

	if err := r.db.Where("storage_key = ?", storageKey).First(&file).Error; err != nil {
		return nil, err
	}

	return &file, nil
}
