package postgres

import (
	"context"
	"fmt"

	"cloud_storage/internal/domain"
)

func (p *Postgres) Create(ctx context.Context, user domain.User) (domain.User, error) {
	query := `
		INSERT INTO users (id, name, email, password_hash, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`
	_, err := p.pool.Exec(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("insert user: %w", err)
	}

	return user, nil
}
