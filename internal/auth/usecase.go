package auth

import (
	"errors"

	"be-ecommerce/internal/customerror"
	"be-ecommerce/internal/models"
	"be-ecommerce/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authUseCase struct {
	repo AuthRepository
}

// NewAuthUseCase creates a new instance of AuthUseCase
func NewAuthUseCase(repo AuthRepository) AuthUseCase {
	return &authUseCase{
		repo: repo,
	}
}

func (u *authUseCase) Register(req RegisterRequest) (RegisterResponse, error) {
	// Check if email already exists
	_, err := u.repo.GetUserByEmail(req.Email)
	if err == nil {
		return RegisterResponse{}, customerror.NewConflictError("email already registered")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return RegisterResponse{}, customerror.NewInternalError("database error", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return RegisterResponse{}, customerror.NewInternalError("failed to hash password", err)
	}

	// Create user model
	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         "customer", // Default role
	}

	// Save to DB via repository
	if err := u.repo.CreateUser(&user); err != nil {
		return RegisterResponse{}, customerror.NewInternalError("failed to create user", err)
	}

	// Prepare response
	return RegisterResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (u *authUseCase) Login(req LoginRequest) (LoginResponse, error) {
	// Find user by email
	user, err := u.repo.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return LoginResponse{}, customerror.NewUnauthorizedError("invalid email or password")
		}
		return LoginResponse{}, customerror.NewInternalError("database error", err)
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return LoginResponse{}, customerror.NewUnauthorizedError("invalid email or password")
	}

	// Generate JWT Token
	token, err := utils.GenerateJWT(user.ID.String(), user.Role)
	if err != nil {
		return LoginResponse{}, customerror.NewInternalError("failed to generate token", err)
	}

	// Prepare response
	resp := LoginResponse{
		Token: token,
	}
	resp.User.ID = user.ID.String()
	resp.User.Name = user.Name
	resp.User.Email = user.Email
	resp.User.Role = user.Role

	return resp, nil
}
