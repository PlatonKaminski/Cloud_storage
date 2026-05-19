package usecase

import (
	"cloud_storage/internal/domain"
	"cloud_storage/internal/dto"
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (u *UseCase) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return dto.LoginResponse{}, domain.ErrInvalidCredentials
		}
		return dto.LoginResponse{}, fmt.Errorf("get user: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return dto.LoginResponse{}, domain.ErrInvalidCredentials
	}

	token, err := u.jwt.Generate(user.ID)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{Token: token}, nil
}
