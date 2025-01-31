package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrfansi/go-api-boilerplate/internal/domain/entity"
	"github.com/mrfansi/go-api-boilerplate/internal/domain/errors"
	"github.com/mrfansi/go-api-boilerplate/internal/domain/repository"
)

type UserService interface {
	Create(ctx context.Context, email, password, name string) (*entity.User, error)
	Update(ctx context.Context, id uuid.UUID, name string) (*entity.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	List(ctx context.Context, page, limit int) ([]*entity.User, int64, error)
	ChangePassword(ctx context.Context, id uuid.UUID, oldPassword, newPassword string) error
	UpdateRole(ctx context.Context, id uuid.UUID, role string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Create(ctx context.Context, email, password, name string) (*entity.User, error) {
	// Check if user already exists
	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return nil, errors.ErrUserAlreadyExists
	}

	// Create new user
	user, err := entity.NewUser(email, password, name)
	if err != nil {
		return nil, err
	}

	// Save user to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Update(ctx context.Context, id uuid.UUID, name string) (*entity.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Update(name)

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

func (s *userService) List(ctx context.Context, page, limit int) ([]*entity.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.userRepo.List(ctx, page, limit)
}

func (s *userService) ChangePassword(ctx context.Context, id uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := user.ComparePassword(oldPassword); err != nil {
		return errors.ErrInvalidPassword
	}

	if err := user.UpdatePassword(newPassword); err != nil {
		return err
	}

	return s.userRepo.Update(ctx, user)
}

func (s *userService) UpdateRole(ctx context.Context, id uuid.UUID, role string) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Validate role
	switch role {
	case "admin", "user":
		user.SetRole(role)
	default:
		return errors.ErrInvalidRole
	}

	return s.userRepo.Update(ctx, user)
}
