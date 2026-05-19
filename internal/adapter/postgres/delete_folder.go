package postgres

import (
	"cloud_storage/internal/domain"
	"context"
	"fmt"
)

func (p *Postgres) DeleteFolder(ctx context.Context, id string) error {

	tag, err := p.pool.Exec(ctx,
		`DELETE FROM folders WHERE id = $1`, id,
	)
	if err != nil {
		return fmt.Errorf("delete folder: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrFolderNotFound
	}

	return nil
}
