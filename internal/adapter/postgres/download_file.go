package postgres

import (
	"context"
	"errors"
	"fmt"

	"cloud_storage/internal/domain"

	"github.com/jackc/pgx/v5"
)

func (p *Postgres) GetByIDU(ctx context.Context, id string) (domain.File, error) {
	query := `
		SELECT id, name, size, mime_type, path, owner_id, created_at
		FROM files
		WHERE id = $1
	`

	var file domain.File
	err := p.pool.QueryRow(ctx, query, id).Scan(
		&file.ID,
		&file.Name,
		&file.Size,
		&file.MimeType,
		&file.Path,
		&file.OwnerID,
		&file.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.File{}, domain.ErrFileNotFound
		}
		return domain.File{}, fmt.Errorf("get file by id: %w", err)
	}

	return file, nil
}
