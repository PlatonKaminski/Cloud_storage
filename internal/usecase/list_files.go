package usecase

import (
	"cloud_storage/internal/dto"
	"context"
	"fmt"
)

const DEAULT_LIMIT = 20
const MAX_LIMIT = 100

func (u *UseCase) ListByOwner(ctx context.Context, req dto.ListFileRequest) (dto.ListFileResponse, error) {
	if req.Limit <= 0 {
		req.Limit = DEAULT_LIMIT
	}
	if req.Limit > MAX_LIMIT {
		req.Limit = MAX_LIMIT
	}
	files, total, err := u.fileRepo.List(ctx, req.Owner_id, req.Limit, req.Offset)
	if err != nil {
		return dto.ListFileResponse{}, fmt.Errorf("list files %w", err)
	}

	return dto.ListFileResponse{
		Files: files,
		Total: total,
	}, nil

}
