package v1

import (
	"errors"
	"net/http"

	"cloud_storage/internal/controller/http/middleware"
	"cloud_storage/internal/domain"
	"cloud_storage/internal/dto"
	"cloud_storage/pkg/render"
)

const maxMultipartMemory = 32 << 20

func (h *Handler) uploadFile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxMultipartMemory); err != nil {
		render.Error(w, err, http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		render.Error(w, err, http.StatusBadRequest)
		return
	}
	defer file.Close()

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		render.Error(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	folderID := r.FormValue("folder_id")

	req := dto.UploadFileRequest{
		Name:      header.Filename,
		Mime:      header.Header.Get("Content-Type"),
		Size:      header.Size,
		Owner_id:  userID,
		Folder_id: folderID,
	}

	resp, err := h.uc.UploadFile(r.Context(), req, file)
	if err != nil {
		if errors.Is(err, domain.ErrFileTooLarge) {
			render.Error(w, err, http.StatusRequestEntityTooLarge)
			return
		}

		h.log.Error("upload file", "err", err)
		render.Error(w, errors.New("internal server error"), http.StatusInternalServerError)
		return
	}

	render.JSON(w, resp, http.StatusCreated)
}
