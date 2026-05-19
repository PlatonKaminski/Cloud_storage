package v1

import (
	"errors"
	"net/http"

	"cloud_storage/internal/controller/http/middleware"
	"cloud_storage/internal/domain"
	"cloud_storage/internal/dto"
	"cloud_storage/pkg/render"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) deleteFolder(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		render.Error(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	folderID := chi.URLParam(r, "id")
	if folderID == "" {
		render.Error(w, errors.New("folder id is required"), http.StatusBadRequest)
		return
	}

	err := h.uc.DeleteFolder(r.Context(), dto.DeleteFolderRequest{
		FolderId: folderID,
		OwnerId:  userID,
	})
	if err != nil {
		if errors.Is(err, domain.ErrFolderNotFound) {
			render.Error(w, errors.New("folder not found"), http.StatusNotFound)
			return
		}
		h.log.Error("delete folder", "err", err)
		render.Error(w, errors.New("internal server error"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
