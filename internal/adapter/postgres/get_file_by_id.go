package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cloud_storage/internal/domain"
)

func (p *Postgres) GetByID(ctx context.Context, id string) (domain.File, error) {
	query := `
        SELECT id, name, path, size, mime_type, owner_id, folder_id, created_at
        FROM files
        WHERE id = $1
    `

	var f domain.File
	var folderID sql.NullString

	err := p.pool.QueryRow(ctx, query, id).Scan(
		&f.ID,
		&f.Name,
		&f.Path,
		&f.Size,
		&f.MimeType,
		&f.OwnerID,
		&folderID,
		&f.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.File{}, domain.ErrFileNotFound
		}
		return domain.File{}, fmt.Errorf("scan file: %w", err)
	}

	if folderID.Valid {
		f.FolderID = folderID.String
	}

	return f, nil
}
