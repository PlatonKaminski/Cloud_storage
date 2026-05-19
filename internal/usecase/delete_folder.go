package usecase

import (
	"cloud_storage/internal/dto"
	"context"
	"errors"
	"fmt"
)

func (u *UseCase) DeleteFolder(ctx context.Context, req dto.DeleteFolderRequest) error {

	folder, err := u.folderRepo.GetByIDFolder(ctx, req.FolderId)
	if err != nil {
		return fmt.Errorf("folder not found: %w", err)
	}

	if folder.OwnerID != req.OwnerId {
		return errors.New("folder not owned by user")
	}

	err = u.fileRepo.DeleteByFolder(ctx, req.FolderId)
	if err != nil {
		return fmt.Errorf("delete files in folder: %w", err)
	}

	err = u.folderRepo.DeleteFolder(ctx, req.FolderId)
	if err != nil {
		return fmt.Errorf("delete folder: %w", err)
	}

	return nil
}
