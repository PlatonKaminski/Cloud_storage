package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"cloud_storage/internal/domain"
)

func (p *Postgres) Save(ctx context.Context, file domain.File) (domain.File, error) {
	query := `
		INSERT INTO files (
			id,
			name,
			path,
			size,
			mime_type,
			owner_id,
			folder_id,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
	`

	var folderID sql.NullString
	if file.FolderID != "" {
		folderID = sql.NullString{String: file.FolderID, Valid: true}
	} else {
		folderID = sql.NullString{Valid: false}
	}

	_, err := p.pool.Exec(
		ctx,
		query,
		file.ID,
		file.Name,
		file.Path,
		file.Size,
		file.MimeType,
		file.OwnerID,
		folderID,
	)
	if err != nil {
		return domain.File{}, fmt.Errorf("insert file: %w", err)
	}

	return file, nil
}
