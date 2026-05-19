package v1

import (
	"cloud_storage/internal/dto"
	"errors"
	"net/http"

	"cloud_storage/internal/controller/http/middleware"
	"cloud_storage/pkg/render"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) deleteFile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		render.Error(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	fileID := chi.URLParam(r, "id")
	if fileID == "" {
		render.Error(w, errors.New("file id required"), http.StatusBadRequest)
		return
	}
	req := dto.DeleteFileRequest{
		File_id:  fileID,
		Owner_id: userID,
	}

	err := h.uc.DeleteFile(r.Context(), req)
	if err != nil {
		h.log.Error("delete file", "err", err)
		render.Error(w, errors.New("internal server error"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
