package mocks

import (
	"github.com/Grafiters/archive/internal/domain"
)

type MockAuthRepository struct {
	LoginFunc      func(auth *domain.AuthRequest) (*domain.Customer, error)
	RegisterFunc   func(data *domain.RegisterRequest) (*domain.Customer, error)
	GetByEmailFunc func(email string) (*domain.Customer, error)
}

func (m *MockAuthRepository) Login(auth *domain.AuthRequest) (*domain.Customer, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(auth)
	}
	return &domain.Customer{}, nil
}

func (m *MockAuthRepository) Register(data *domain.RegisterRequest) (*domain.Customer, error) {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(data)
	}
	return &domain.Customer{}, nil
}

func (m *MockAuthRepository) GetByEmail(email string) (*domain.Customer, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(email)
	}
	return &domain.Customer{}, nil
}
