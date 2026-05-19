package usecase

import (
	"cloud_storage/internal/domain"
	"cloud_storage/internal/dto"
	"context"
	"errors"
	"fmt"
)

func (u *UseCase) DownloadFile(ctx context.Context, req dto.DownloadFileRequest) (dto.DownloadFileResponse, error) {

	file, err := u.fileRepo.GetByID(ctx, req.File_id)
	if err != nil {
		if errors.Is(err, domain.ErrFileNotFound) {
			return dto.DownloadFileResponse{}, domain.ErrFileNotFound
		}
		return dto.DownloadFileResponse{}, fmt.Errorf("get file metadata: %w", err)
	}

	if file.OwnerID != req.Owner_id {
		return dto.DownloadFileResponse{}, domain.ErrFileAccessDenied
	}

	return dto.DownloadFileResponse{
		Path:      file.Path,
		Mime_type: file.MimeType,
		Name:      file.Name,
	}, nil
}
