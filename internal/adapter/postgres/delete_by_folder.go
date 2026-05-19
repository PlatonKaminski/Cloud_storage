package postgres

import (
	"context"
	"fmt"
)

func (p *Postgres) DeleteByFolder(ctx context.Context, folderID string) error {
	query := `DELETE FROM files WHERE folder_id = $1`

	result, err := p.pool.Exec(ctx, query, folderID)
	if err != nil {
		return fmt.Errorf("delete files by folder: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected > 0 {
		fmt.Printf("Deleted %d files from folder %s\n", rowsAffected, folderID)
	}

	return nil
}
