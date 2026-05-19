// internal/usecase/usecase.go
package usecase

import (
	pkgjwt "cloud_storage/pkg/jwt"
	"context"
	"io"

	"cloud_storage/internal/domain"
)

//go:generate mockery --name=FileRepository
type FileRepository interface {
	Save(ctx context.Context, file domain.File) (domain.File, error)
	GetByID(ctx context.Context, id string) (domain.File, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, ownerID string, limit, offset int) ([]domain.File, int, error)
	ListByFolder(ctx context.Context, folderID string) ([]domain.File, int, error)
	DeleteByFolder(ctx context.Context, folderID string) error
}

//go:generate mockery --name=FileStorage
type FileStorage interface {
	Save(ctx context.Context, path string, r io.Reader) error
	Read(ctx context.Context, path string) (io.ReadCloser, error)
	Delete(ctx context.Context, path string) error
}

//go:generate mockery --name=UserRepository
type UserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByIDUser(ctx context.Context, id string) (domain.User, error)
}

//go:generate mockery --name=FolderRepository
type FolderRepository interface {
	CreateFolder(ctx context.Context, folder domain.Folder) (domain.Folder, error)
	GetByIDFolder(ctx context.Context, id string) (domain.Folder, error)
	DeleteFolder(ctx context.Context, id string) error
	ListByOwner(ctx context.Context, id string) ([]domain.Folder, error)
}

type UseCase struct {
	fileRepo   FileRepository
	userRepo   UserRepository
	folderRepo FolderRepository
	storage    FileStorage
	jwt        *pkgjwt.Manager
}

func New(
	fileRepo FileRepository,
	userRepo UserRepository,
	folderRepo FolderRepository,
	storage FileStorage,
	jwt *pkgjwt.Manager,
) *UseCase {
	return &UseCase{
		fileRepo:   fileRepo,
		userRepo:   userRepo,
		folderRepo: folderRepo,
		storage:    storage,
		jwt:        jwt,
	}
}
