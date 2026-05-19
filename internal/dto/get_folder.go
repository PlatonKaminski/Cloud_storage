package dto

import "cloud_storage/internal/domain"

type GetFolderRequest struct {
	FolderId string
	OwnerId  string
}

type GetFolderResponse struct {
	Folder domain.Folder `json:"folder"`
	Files  []domain.File `json:"files"`
	Total  int           `json:"total"`
}
