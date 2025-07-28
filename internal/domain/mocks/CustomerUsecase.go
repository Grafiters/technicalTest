package mocks

import (
	"github.com/Grafiters/archive/internal/domain"
	"github.com/stretchr/testify/mock"
)

type CustomerUsecaseMock struct {
	mock.Mock
}

func (m *CustomerUsecaseMock) GetByID(id int64) (*domain.Customer, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *CustomerUsecaseMock) Update(id int64, data *domain.Customer) (*domain.Customer, error) {
	args := m.Called(id, data)
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *CustomerUsecaseMock) UpdateSalary(id int64, salary domain.CustomerUpdateSalary) (*domain.Customer, error) {
	args := m.Called(id, salary)
	return args.Get(0).(*domain.Customer), args.Error(1)
}
