package domain

import (
	"time"
)

type File struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	MimeType  string    `json:"mimeType"`
	Path      string    `json:"path"`
	OwnerID   string    `json:"ownerID"`
	FolderID  string    `json:"folder_id"`
	CreatedAt time.Time `json:"created_at"`
}
