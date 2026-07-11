package auth

import (
	"be-ecommerce/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockAuthRepository is a mock implementation of AuthRepository
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockAuthRepository) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}
