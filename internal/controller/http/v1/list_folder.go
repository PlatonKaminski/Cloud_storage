package v1

import (
	"errors"
	"net/http"

	"cloud_storage/internal/controller/http/middleware"
	"cloud_storage/internal/dto"
	"cloud_storage/pkg/render"
)

func (h *Handler) listFolders(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		render.Error(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	resp, err := h.uc.ListFolders(r.Context(), dto.ListFoldersRequest{OwnerID: userID})
	if err != nil {
		h.log.Error("list folders", "err", err)
		render.Error(w, errors.New("internal server error"), http.StatusInternalServerError)
		return
	}

	render.JSON(w, resp, http.StatusOK)
}
