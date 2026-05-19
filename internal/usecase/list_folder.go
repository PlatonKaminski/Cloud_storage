package usecase

import (
	"cloud_storage/internal/dto"
	"context"
)

func (uc *UseCase) ListFolders(ctx context.Context, req dto.ListFoldersRequest) (dto.ListFoldersResponse, error) {
	folders, err := uc.folderRepo.ListByOwner(ctx, req.OwnerID)
	if err != nil {
		return dto.ListFoldersResponse{}, err
	}
	result := make([]dto.ListFolderResponse, len(folders))
	for i, f := range folders {
		resp := dto.ListFolderResponse{
			Id:       f.Id,
			Name:     f.Name,
			ParentId: f.ParentID,
		}

		if f.ParentID == "" {
			for _, child := range folders {
				if child.ParentID == f.Id {
				}
			}
		}

		result[i] = resp
	}

	return dto.ListFoldersResponse{Folders: result}, nil
}
