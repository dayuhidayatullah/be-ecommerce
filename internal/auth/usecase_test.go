package auth

import (
	"os"
	"testing"

	"be-ecommerce/internal/customerror"
	"be-ecommerce/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	// Set dummy JWT secret for testing utils.GenerateJWT inside Login
	os.Setenv("JWT_SECRET", "dummy_secret_for_testing")
	os.Exit(m.Run())
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name          string
		req           RegisterRequest
		mockSetup     func(repo *MockAuthRepository)
		expectedError *customerror.AppError
	}{
		{
			name: "Success Registration",
			req: RegisterRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(repo *MockAuthRepository) {
				// Mock finding email: returns ErrRecordNotFound (meaning email is available)
				repo.On("GetUserByEmail", "test@example.com").Return(nil, gorm.ErrRecordNotFound)
				// Mock creating user: returns no error
				repo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)
			},
			expectedError: nil, // Should succeed
		},
		{
			name: "Failed - Email Already Exists",
			req: RegisterRequest{
				Name:     "Test User 2",
				Email:    "exist@example.com",
				Password: "password123",
			},
			mockSetup: func(repo *MockAuthRepository) {
				// Mock finding email: returns a user (meaning email is taken)
				repo.On("GetUserByEmail", "exist@example.com").Return(&models.User{}, nil)
			},
			expectedError: customerror.NewConflictError("email already registered"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockAuthRepository)
			tt.mockSetup(mockRepo)
			useCase := NewAuthUseCase(mockRepo)

			// Act
			resp, err := useCase.Register(tt.req)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				// Check if it's our custom AppError
				appErr, ok := err.(*customerror.AppError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedError.Type, appErr.Type)
				assert.Equal(t, tt.expectedError.Message, appErr.Message)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, resp.ID)
				assert.Equal(t, tt.req.Email, resp.Email)
			}

			mockRepo.AssertExpectations(t) // Verify that all expected mocks were called
		})
	}
}

func TestLogin(t *testing.T) {
	// Create a dummy hashed password for testing
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	
	dummyUser := &models.User{
		ID:           uuid.New(),
		Name:         "John Doe",
		Email:        "john@example.com",
		PasswordHash: string(hashedPassword),
		Role:         "customer",
	}

	tests := []struct {
		name          string
		req           LoginRequest
		mockSetup     func(repo *MockAuthRepository)
		expectedError *customerror.AppError
	}{
		{
			name: "Success Login",
			req: LoginRequest{
				Email:    "john@example.com",
				Password: "password123",
			},
			mockSetup: func(repo *MockAuthRepository) {
				// Mock finding email: returns dummyUser
				repo.On("GetUserByEmail", "john@example.com").Return(dummyUser, nil)
			},
			expectedError: nil,
		},
		{
			name: "Failed - Invalid Email (Not Found)",
			req: LoginRequest{
				Email:    "wrong@example.com",
				Password: "password123",
			},
			mockSetup: func(repo *MockAuthRepository) {
				// Mock finding email: returns ErrRecordNotFound
				repo.On("GetUserByEmail", "wrong@example.com").Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError: customerror.NewUnauthorizedError("invalid email or password"),
		},
		{
			name: "Failed - Invalid Password",
			req: LoginRequest{
				Email:    "john@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(repo *MockAuthRepository) {
				// Mock finding email: returns dummyUser (but the password won't match)
				repo.On("GetUserByEmail", "john@example.com").Return(dummyUser, nil)
			},
			expectedError: customerror.NewUnauthorizedError("invalid email or password"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockAuthRepository)
			tt.mockSetup(mockRepo)
			useCase := NewAuthUseCase(mockRepo)

			// Act
			resp, err := useCase.Login(tt.req)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				appErr, ok := err.(*customerror.AppError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedError.Type, appErr.Type)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, resp.Token)
				assert.Equal(t, dummyUser.Email, resp.User.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
