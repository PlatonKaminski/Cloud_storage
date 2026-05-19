package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"cloud_storage/internal/domain"
)

func (p *Postgres) ListByOwner(ctx context.Context, ownerID string) ([]domain.Folder, error) {
	query := `
		SELECT id, name, owner_id, parent_id, created_at
		FROM folders
		WHERE owner_id = $1
		ORDER BY created_at DESC
	`

	rows, err := p.pool.Query(ctx, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("query folders: %w", err)
	}
	defer rows.Close()

	var folders []domain.Folder
	for rows.Next() {
		var f domain.Folder
		var parentID sql.NullString

		if err := rows.Scan(
			&f.Id,
			&f.Name,
			&f.OwnerID,
			&parentID,
			&f.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan folder: %w", err)
		}

		if parentID.Valid {
			f.ParentID = parentID.String
		}

		folders = append(folders, f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if folders == nil {
		folders = []domain.Folder{}
	}

	return folders, nil
}
