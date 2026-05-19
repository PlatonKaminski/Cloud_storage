package dto

import (
	"cloud_storage/internal/domain"
)

type ListFileRequest struct {
	Owner_id string
	Limit    int
	Offset   int
}

type ListFileResponse struct {
	Files []domain.File `json:"files"`
	Total int           `json:"total"`
}
