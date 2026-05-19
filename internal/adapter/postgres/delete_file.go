package postgres

import (
	"context"
	"fmt"

	"cloud_storage/internal/domain"
)

func (p *Postgres) Delete(ctx context.Context, id string) error {
	tag, err := p.pool.Exec(ctx,
		`DELETE FROM files WHERE id = $1`, id,
	)
	if err != nil {
		return fmt.Errorf("delete file: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrFileNotFound
	}

	return nil
}
