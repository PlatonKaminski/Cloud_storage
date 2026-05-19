package usecase

import (
	"cloud_storage/internal/domain"
	"cloud_storage/internal/dto"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

const (
	MAX_SIZE = 1024 * 1024 * 100
)

func (u *UseCase) UploadFile(
	ctx context.Context,
	req dto.UploadFileRequest,
	r io.Reader,
) (dto.UploadFileResponse, error) {

	now := time.Now()

	path := filepath.Join(
		req.Owner_id,
		now.Format("2006"),
		now.Format("01"),
		req.Name,
	)

	if err := u.storage.Save(ctx, path, r); err != nil {
		return dto.UploadFileResponse{}, fmt.Errorf("save file: %w", err)
	}

	fileEntity := domain.File{
		ID:       uuid.New().String(),
		Name:     req.Name,
		Path:     path,
		Size:     req.Size,
		MimeType: req.Mime,
		OwnerID:  req.Owner_id,
		FolderID: req.Folder_id,
	}

	file, err := u.fileRepo.Save(ctx, fileEntity)
	if err != nil {
		return dto.UploadFileResponse{}, fmt.Errorf("save file metadata: %w", err)
	}

	return dto.UploadFileResponse{
		Id: file.ID,
	}, nil
}
