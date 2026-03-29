package file

import "gorm.io/gorm"

type Repository interface {
	Create(file *File) error
	FindOne(storageKey string) (*File, error)
	FindOneByUserID(storageKey string, userID string) (*File, error)
	FindOneByAnyKeyAndUserID(objectKey string, userID string) (*File, error)
	FindByUserID(userID string) ([]File, error)
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

func (r *repository) FindOneByUserID(storageKey string, userID string) (*File, error) {
	var file File

	if err := r.db.Where("storage_key = ? AND user_id = ?", storageKey, userID).First(&file).Error; err != nil {
		return nil, err
	}

	return &file, nil
}

func (r *repository) FindOneByAnyKeyAndUserID(objectKey string, userID string) (*File, error) {
	var file File

	if err := r.db.Where("(storage_key = ? OR thumbnail_storage_key = ?) AND user_id = ?", objectKey, objectKey, userID).First(&file).Error; err != nil {
		return nil, err
	}

	return &file, nil
}

func (r *repository) FindByUserID(userID string) ([]File, error) {
	var files []File

	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&files).Error; err != nil {
		return nil, err
	}

	return files, nil
}
