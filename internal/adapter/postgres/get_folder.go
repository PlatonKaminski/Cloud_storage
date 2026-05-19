package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cloud_storage/internal/domain"
)

func (p *Postgres) GetByIDFolder(ctx context.Context, id string) (domain.Folder, error) {
	query := `
		SELECT id, name, owner_id, parent_id, created_at
		FROM folders
		WHERE id = $1
	`

	var f domain.Folder
	var parentID sql.NullString

	err := p.pool.QueryRow(ctx, query, id).Scan(
		&f.Id,
		&f.Name,
		&f.OwnerID,
		&parentID,
		&f.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Folder{}, domain.ErrFolderNotFound
		}
		return domain.Folder{}, fmt.Errorf("scan folder: %w", err)
	}

	if parentID.Valid {
		f.ParentID = parentID.String
	}

	return f, nil
}
