package usecase_test

import (
	"errors"
	"testing"

	"github.com/Grafiters/archive/configs"
	"github.com/Grafiters/archive/internal/auth/usecase"
	"github.com/Grafiters/archive/internal/domain"
	"github.com/Grafiters/archive/internal/domain/mocks"
	"github.com/Grafiters/archive/utils"
	"github.com/stretchr/testify/assert"
)

func TestRegister_Success(t *testing.T) {
	mockAuthRepo := &mocks.MockAuthRepository{}
	mockLimitRepo := &mocks.MockLimitRepository{}
	mockLogger := &configs.LoggerFormat{}

	req := &domain.RegisterRequest{
		Email:    "test@example.com",
		Password: "securepassword",
	}

	mockCustomer := &domain.Customer{
		ID:       1,
		Email:    req.Email,
		Password: req.Password,
	}

	mockAuthRepo.GetByEmailFunc = func(email string) (*domain.Customer, error) {
		return nil, errors.New("not found")
	}

	mockAuthRepo.RegisterFunc = func(data *domain.RegisterRequest) (*domain.Customer, error) {
		return mockCustomer, nil
	}

	mockLimitRepo.BulkCreateLimitFunc = func(data *domain.BulkLimitInput) ([]*int64, error) {
		return []*int64{utils.Int64Ptr(1), utils.Int64Ptr(2)}, nil
	}

	uc := usecase.NewAuthUsecase(mockAuthRepo, mockLimitRepo, mockLogger)

	result, err := uc.Register(nil, req)

	assert.NoError(t, err)
	assert.Equal(t, req.Email, result.Email)
	assert.NotEmpty(t, result.ID)
}

func TestRegister_EmailAlreadyUsed(t *testing.T) {
	mockAuthRepo := &mocks.MockAuthRepository{
		GetByEmailFunc: func(email string) (*domain.Customer, error) {
			return &domain.Customer{Email: email}, nil
		},
	}
	mockLimitRepo := &mocks.MockLimitRepository{}
	mockLogger := &configs.LoggerFormat{}

	req := &domain.RegisterRequest{
		Email:    "test@example.com",
		Password: "securepassword",
	}

	uc := usecase.NewAuthUsecase(mockAuthRepo, mockLimitRepo, mockLogger)

	_, err := uc.Register(nil, req)

	assert.Error(t, err)
	assert.Equal(t, "email already in use", err.Error())
}

func TestLogin_Success(t *testing.T) {
	mockAuthRepo := &mocks.MockAuthRepository{}
	mockLogger := &configs.LoggerFormat{}
	mockLimitRepo := &mocks.MockLimitRepository{}

	hashedPassword, _ := utils.Hash("securepassword")

	authReq := &domain.AuthRequest{
		Email:    "test@example.com",
		Password: "securepassword",
	}

	mockCustomer := &domain.Customer{
		ID:       1,
		Email:    authReq.Email,
		Password: hashedPassword,
	}

	mockAuthRepo.LoginFunc = func(auth *domain.AuthRequest) (*domain.Customer, error) {
		return mockCustomer, nil
	}

	uc := usecase.NewAuthUsecase(mockAuthRepo, mockLimitRepo, mockLogger)

	result, err := uc.Login(nil, authReq)

	assert.NoError(t, err)
	assert.Equal(t, authReq.Email, result.Email)
}
