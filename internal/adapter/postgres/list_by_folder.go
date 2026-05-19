package postgres

import (
	"context"
	"fmt"

	"cloud_storage/internal/domain"
)

func (p *Postgres) ListByFolder(ctx context.Context, folderID string) ([]domain.File, int, error) {
	var total int
	err := p.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM files WHERE folder_id = $1`,
		folderID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count files in folder: %w", err)
	}

	rows, err := p.pool.Query(ctx, `
		SELECT id, name, size, mime_type, path, owner_id, folder_id, created_at
		FROM files
		WHERE folder_id = $1
		ORDER BY created_at DESC
	`,
		folderID,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list files in folder: %w", err)
	}
	defer rows.Close()

	var files []domain.File
	for rows.Next() {
		var f domain.File
		if err = rows.Scan(
			&f.ID, &f.Name, &f.Size, &f.MimeType,
			&f.Path, &f.OwnerID, &f.FolderID, &f.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan file row: %w", err)
		}
		files = append(files, f)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return files, total, nil
}
