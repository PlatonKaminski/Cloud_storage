package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cloud_storage/internal/domain"
)

func (p *Postgres) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `
		SELECT id, name, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`
	var u domain.User
	err := p.pool.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("get user by email: %w", err)
	}
	return u, nil
}
