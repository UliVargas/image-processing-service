package file

import "time"

type File struct {
	ID         string    `gorm:"primaryKey;size=24" json:"id"`
	FileName   string    `gorm:"not null" json:"file_name"`
	StorageKey string    `gorm:"unique not null" json:"storage_key"`
	MimeType   string    `gorm:"not null" json:"mime_type"`
	FileSize   int64     `gorm:"not null" json:"file_size"`
	CreatedAt  time.Time `json:"created_at"`
}

type FileUploadRequest struct {
	FileName   string
	StorageKey string
	MimeType   string
	FileSize   int64
}
