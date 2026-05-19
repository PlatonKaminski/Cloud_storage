package usecase

import (
	"cloud_storage/internal/domain"
	"cloud_storage/internal/dto"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func (u *UseCase) Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.RegisterResponse{}, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now()

	user := domain.User{
		ID:           uuid.New().String(),
		Name:         req.Name,
		Email:        req.Email,
		CreatedAt:    now,
		PasswordHash: string(hash),
	}
	created, err := u.userRepo.Create(ctx, user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return dto.RegisterResponse{}, domain.ErrUserAlreadyExists
		}
		return dto.RegisterResponse{}, fmt.Errorf("create user: %w", err)
	}

	return dto.RegisterResponse{Id: created.ID}, nil

}
