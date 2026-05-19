package v1

import (
	"errors"
	"net/http"
	"strconv"

	"cloud_storage/internal/controller/http/middleware"
	"cloud_storage/internal/dto"
	"cloud_storage/pkg/render"
)

func (h *Handler) listFiles(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		render.Error(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	limit := 100
	offset := 0

	if v := r.URL.Query().Get("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			limit = parsed
		}
	}

	if v := r.URL.Query().Get("offset"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			offset = parsed
		}
	}

	resp, err := h.uc.ListByOwner(r.Context(), dto.ListFileRequest{
		Owner_id: userID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		h.log.Error("list files", "err", err)
		render.Error(w, errors.New("internal server error"), http.StatusInternalServerError)
		return
	}

	render.JSON(w, resp, http.StatusOK)
}
