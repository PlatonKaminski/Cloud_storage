package v1

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"cloud_storage/internal/controller/http/middleware"
	"cloud_storage/internal/domain"
	"cloud_storage/internal/dto"
	"cloud_storage/pkg/render"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) downloadFile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		render.Error(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	fileID := chi.URLParam(r, "id")
	if fileID == "" {
		render.Error(w, errors.New("file id is required"), http.StatusBadRequest)
		return
	}

	meta, err := h.uc.DownloadFile(r.Context(), dto.DownloadFileRequest{
		File_id:  fileID,
		Owner_id: userID,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFileNotFound):
			render.Error(w, err, http.StatusNotFound)
		case errors.Is(err, domain.ErrFileAccessDenied):
			render.Error(w, err, http.StatusForbidden)
		default:
			h.log.Error("download file", "err", err)
			render.Error(w, errors.New("internal server error"), http.StatusInternalServerError)
		}
		return
	}

	rc, err := h.storage.Read(r.Context(), meta.Path)
	if err != nil {
		if errors.Is(err, domain.ErrFileNotFound) {
			render.Error(w, err, http.StatusNotFound)
			return
		}
		h.log.Error("open file from storage", "err", err)
		render.Error(w, errors.New("internal server error"), http.StatusInternalServerError)
		return
	}
	defer rc.Close()

	w.Header().Set("Content-Type", meta.Mime_type)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, meta.Name))
	w.WriteHeader(http.StatusOK)

	if _, err = io.Copy(w, rc); err != nil {
		h.log.Error("stream file", "err", err)
	}
}
