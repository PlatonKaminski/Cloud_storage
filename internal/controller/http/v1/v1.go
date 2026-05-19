package v1

import (
	"context"
	"io"
	"log/slog"

	"cloud_storage/internal/dto"

	"github.com/go-chi/chi/v5"
)

type UseCase interface {
	// Файлы
	UploadFile(ctx context.Context, req dto.UploadFileRequest, r io.Reader) (dto.UploadFileResponse, error)
	ListByOwner(ctx context.Context, req dto.ListFileRequest) (dto.ListFileResponse, error)
	DownloadFile(ctx context.Context, req dto.DownloadFileRequest) (dto.DownloadFileResponse, error)
	DeleteFile(ctx context.Context, req dto.DeleteFileRequest) error

	// Авторизация
	Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)

	// Папки
	CreateFolder(ctx context.Context, req dto.CreateFolderRequest) (dto.CreateFolderResponse, error)
	GetFolder(ctx context.Context, req dto.GetFolderRequest) (dto.GetFolderResponse, error)
	DeleteFolder(ctx context.Context, req dto.DeleteFolderRequest) error
	ListFolders(ctx context.Context, req dto.ListFoldersRequest) (dto.ListFoldersResponse, error)
}

type FileOpener interface {
	Read(ctx context.Context, path string) (io.ReadCloser, error)
}

type Handler struct {
	uc      UseCase
	storage FileOpener
	log     *slog.Logger
}

func New(uc UseCase, storage FileOpener, log *slog.Logger) *Handler {
	return &Handler{
		uc:      uc,
		storage: storage,
		log:     log,
	}
}

func (h *Handler) RegisterAuth(r chi.Router) {
	r.Post("/register", h.register)
	r.Post("/login", h.login)
}

func (h *Handler) RegisterAPI(r chi.Router) {
	// Файлы
	r.Post("/files", h.uploadFile)
	r.Get("/files", h.listFiles)
	r.Get("/files/{id}", h.downloadFile)
	r.Delete("/files/{id}", h.deleteFile)

	// Папки
	r.Post("/folders", h.createFolder)
	r.Get("/folders", h.listFolders)
	r.Get("/folders/{id}", h.getFolder)
	r.Delete("/folders/{id}", h.deleteFolder)
}
