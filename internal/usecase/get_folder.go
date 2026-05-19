package usecase

import (
	"context"
	"errors"
	"fmt"

	"cloud_storage/internal/domain"
	"cloud_storage/internal/dto"
)

func (u *UseCase) GetFolder(ctx context.Context, req dto.GetFolderRequest) (dto.GetFolderResponse, error) {
	folder, err := u.folderRepo.GetByIDFolder(ctx, req.FolderId)
	if err != nil {
		if errors.Is(err, domain.ErrFolderNotFound) {
			return dto.GetFolderResponse{}, domain.ErrFolderNotFound
		}
		return dto.GetFolderResponse{}, fmt.Errorf("get folder: %w", err)
	}

	if folder.OwnerID != req.OwnerId { // ← OwnerId (как в твоём dto)
		return dto.GetFolderResponse{}, domain.ErrFolderNotFound
	}
	files, total, err := u.fileRepo.ListByFolder(ctx, folder.Id)

	if err != nil {
		return dto.GetFolderResponse{}, fmt.Errorf("list files in folder: %w", err)
	}

	return dto.GetFolderResponse{
		Folder: folder,
		Files:  files,
		Total:  total,
	}, nil
}
