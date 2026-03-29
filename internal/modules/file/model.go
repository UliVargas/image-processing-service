package file

import "time"

type File struct {
	ID                  string    `gorm:"primaryKey;size=24" json:"id"`
	FileName            string    `gorm:"not null" json:"file_name"`
	StorageKey          string    `gorm:"unique not null" json:"storage_key"`
	ThumbnailStorageKey string    `gorm:"uniqueIndex" json:"thumbnail_storage_key"`
	MimeType            string    `gorm:"not null" json:"mime_type"`
	FileSize            int64     `gorm:"not null" json:"file_size"`
	UserID              string    `gorm:"index" json:"user_id"`
	Format              string    `gorm:"not null" json:"format"`
	Width               int64     `gorm:"not null" json:"width"`
	Height              int64     `gorm:"not null" json:"height"`
	CreatedAt           time.Time `json:"created_at"`
}

type FileUploadRequest struct {
	FileName string
	MimeType string
	FileSize int64
	UserID   string
	Format   string
	Width    int64
	Height   int64
}
