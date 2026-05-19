package usecase

import (
	"cloud_storage/internal/domain"
	"cloud_storage/internal/dto"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (uc *UseCase) CreateFolder(ctx context.Context, req dto.CreateFolderRequest) (dto.CreateFolderResponse, error) {
	if req.ParentId != "" {
		parent, err := uc.folderRepo.GetByIDFolder(ctx, req.ParentId)
		if err != nil {
			return dto.CreateFolderResponse{}, fmt.Errorf("get parent folder: %w", err)
		}
		if parent.OwnerID != req.OwnerId {
			return dto.CreateFolderResponse{}, domain.ErrFolderAccessDenied
		}
	}

	folder := domain.Folder{
		Id:       uuid.New().String(),
		Name:     req.Name,
		OwnerID:  req.OwnerId,
		ParentID: req.ParentId,
	}

	created, err := uc.folderRepo.CreateFolder(ctx, folder)
	if err != nil {
		return dto.CreateFolderResponse{}, fmt.Errorf("create folder: %w", err)
	}

	return dto.CreateFolderResponse{Id: created.Id}, nil
}
