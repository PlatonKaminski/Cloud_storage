package postgres

import (
	"cloud_storage/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	_ usecase.FileRepository = (*Postgres)(nil)
	_ usecase.UserRepository = (*Postgres)(nil)
	//_ usecase.FolderRepository = (*Postgres)(nil)
)

type Postgres struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Postgres {
	return &Postgres{pool: pool}
}
