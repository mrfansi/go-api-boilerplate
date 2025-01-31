package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/mrfansi/go-api-boilerplate/internal/domain/entity"
	"github.com/mrfansi/go-api-boilerplate/internal/domain/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of repository.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context, page, limit int) ([]*entity.User, int64, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]*entity.User), args.Get(1).(int64), args.Error(2)
}

func TestUserService_Create(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		email := "test@example.com"
		password := "password123"
		name := "Test User"

		mockRepo.On("FindByEmail", ctx, email).Return(nil, errors.ErrUserNotFound)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

		user, err := service.Create(ctx, email, password, name)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, name, user.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserAlreadyExists", func(t *testing.T) {
		email := "existing@example.com"
		password := "password123"
		name := "Existing User"

		existingUser, _ := entity.NewUser(email, password, name)
		mockRepo.On("FindByEmail", ctx, email).Return(existingUser, nil)

		user, err := service.Create(ctx, email, password, name)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrUserAlreadyExists, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		id := uuid.New()
		expectedUser, _ := entity.NewUser("test@example.com", "password123", "Test User")
		expectedUser.ID = id

		mockRepo.On("FindByID", ctx, id).Return(expectedUser, nil)

		user, err := service.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, id, user.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		id := uuid.New()
		mockRepo.On("FindByID", ctx, id).Return(nil, errors.ErrUserNotFound)

		user, err := service.GetByID(ctx, id)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrUserNotFound, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}
