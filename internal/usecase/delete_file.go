package usecase

import (
	"context"
	"errors"
	"fmt"

	"cloud_storage/internal/domain"
	"cloud_storage/internal/dto"
)

func (uc *UseCase) DeleteFile(ctx context.Context, req dto.DeleteFileRequest) error {
	file, err := uc.fileRepo.GetByID(ctx, req.File_id)
	if err != nil {
		if errors.Is(err, domain.ErrFileNotFound) {
			return domain.ErrFileNotFound
		}
		return fmt.Errorf("get file metadata: %w", err)
	}

	if file.OwnerID != req.Owner_id {
		return domain.ErrFileAccessDenied
	}

	if err = uc.storage.Delete(ctx, file.Path); err != nil {
		return fmt.Errorf("delete file from storage: %w", err)
	}

	if err = uc.fileRepo.Delete(ctx, req.File_id); err != nil {
		return fmt.Errorf("delete file metadata: %w", err)
	}

	return nil
}
