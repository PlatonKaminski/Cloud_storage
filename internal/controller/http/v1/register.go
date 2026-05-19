package v1

import (
	"cloud_storage/internal/domain"
	"cloud_storage/internal/dto"
	"cloud_storage/pkg/render"
	"encoding/json"
	"errors"
	"net/http"
)

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("decode register request", "err", err)
		render.Error(w, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	resp, err := h.uc.Register(r.Context(), req)
	if err != nil {
		h.log.Error("register usecase", "err", err)
		if errors.Is(err, domain.ErrAlreadyExists) {
			render.Error(w, errors.New("user already exists"), http.StatusConflict)
			return
		}
		render.Error(w, errors.New("internal error"), http.StatusInternalServerError)
		return
	}

	h.log.Info("user registered", "id", resp.Id)
	render.JSON(w, resp, http.StatusCreated)
}
