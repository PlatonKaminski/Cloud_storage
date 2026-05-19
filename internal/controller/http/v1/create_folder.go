package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"cloud_storage/internal/controller/http/middleware"
	"cloud_storage/internal/domain"
	"cloud_storage/internal/dto"
	"cloud_storage/pkg/render"
)

func (h *Handler) createFolder(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		render.Error(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	var req dto.CreateFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	req.OwnerId = userID

	resp, err := h.uc.CreateFolder(r.Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrFolderNotFound) {
			render.Error(w, errors.New("parent folder not found"), http.StatusNotFound)
			return
		}
		h.log.Error("create folder", "err", err)
		render.Error(w, errors.New("internal server error"), http.StatusInternalServerError)
		return
	}
	
	render.JSON(w, resp, http.StatusCreated)
}
