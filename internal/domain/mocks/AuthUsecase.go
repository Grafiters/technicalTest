package mocks

import (
	"github.com/Grafiters/archive/internal/domain"
	"github.com/stretchr/testify/mock"
)

type AuthUsecaseMock struct {
	mock.Mock
}

func (m *AuthUsecaseMock) Login(req *domain.AuthRequest) (*domain.Customer, error) {
	args := m.Called(req)
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *AuthUsecaseMock) Register(req *domain.RegisterRequest) (*domain.Customer, error) {
	args := m.Called(req)
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *AuthUsecaseMock) GetByEmail(email string) (*domain.Customer, error) {
	args := m.Called(email)
	return args.Get(0).(*domain.Customer), args.Error(1)
}
