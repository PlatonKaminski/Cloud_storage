package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"cloud_storage/internal/domain"
)

func (p *Postgres) CreateFolder(ctx context.Context, folder domain.Folder) (domain.Folder, error) {
	query := `
		INSERT INTO folders (
			id,
			name,
			owner_id,
			parent_id,
			created_at
		)
		VALUES ($1, $2, $3, $4, NOW())
	`
	
	var parentID sql.NullString

	if folder.ParentID != "" {
		parentID = sql.NullString{
			String: folder.ParentID,
			Valid:  true,
		}
	} else {
		parentID = sql.NullString{
			Valid: false,
		}
	}

	_, err := p.pool.Exec(
		ctx,
		query,
		folder.Id,
		folder.Name,
		folder.OwnerID,
		parentID,
	)

	if err != nil {
		return domain.Folder{}, fmt.Errorf("insert folder: %w", err)
	}

	return folder, nil
}
