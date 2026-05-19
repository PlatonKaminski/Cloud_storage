package postgres

import (
	"context"
	"fmt"

	"cloud_storage/internal/domain"
)

func (p *Postgres) List(ctx context.Context, ownerID string, limit, offset int) ([]domain.File, int, error) {
	var total int
	err := p.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM files WHERE owner_id = $1`, ownerID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count files: %w", err)
	}

	rows, err := p.pool.Query(ctx, `
		SELECT 
			f.id, f.name, f.size, f.mime_type, f.path, f.owner_id, 
			f.folder_id, f.created_at
		FROM files f
		LEFT JOIN folders fo ON f.folder_id = fo.id
		WHERE f.owner_id = $1
		ORDER BY 
			CASE WHEN f.folder_id IS NULL THEN 0 ELSE 1 END,
			fo.name,
			f.created_at DESC
		LIMIT $2 OFFSET $3
	`, ownerID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list files: %w", err)
	}
	defer rows.Close()

	var files []domain.File
	for rows.Next() {
		var f domain.File
		var folderID *string
		if err = rows.Scan(
			&f.ID, &f.Name, &f.Size, &f.MimeType,
			&f.Path, &f.OwnerID, &folderID, &f.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan file row: %w", err)
		}
		if folderID != nil {
			f.FolderID = *folderID
		}
		files = append(files, f)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return files, total, nil
}
