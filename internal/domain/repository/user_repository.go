package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrfansi/go-api-boilerplate/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	List(ctx context.Context, page, limit int) ([]*entity.User, int64, error)
}
