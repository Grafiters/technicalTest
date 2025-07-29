package mocks

import (
	"github.com/Grafiters/archive/internal/domain"
)

type MockCustomerRepository struct {
	CreateFunc       func(data *domain.CustomerInput) (*domain.Customer, error)
	GetFunc          func(filter *domain.CustomerFilter) ([]*domain.Customer, int64, error)
	GetByIDFunc      func(ID int64) (*domain.Customer, error)
	UpdateFunc       func(ID int64, newData *domain.CustomerUpdate) (*domain.Customer, error)
	UpdateSalaryFunc func(ID int64, salary domain.CustomerUpdateSalary) (*domain.Customer, error)
}

func (m *MockCustomerRepository) Create(data *domain.CustomerInput) (*domain.Customer, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(data)
	}
	return &domain.Customer{}, nil
}

func (m *MockCustomerRepository) Get(filter *domain.CustomerFilter) ([]*domain.Customer, int64, error) {
	if m.GetFunc != nil {
		return m.GetFunc(filter)
	}
	return []*domain.Customer{}, 0, nil
}

func (m *MockCustomerRepository) GetByID(ID int64) (*domain.Customer, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ID)
	}
	return &domain.Customer{}, nil
}

func (m *MockCustomerRepository) Update(ID int64, newData *domain.CustomerUpdate) (*domain.Customer, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ID, newData)
	}
	return &domain.Customer{}, nil
}

func (m *MockCustomerRepository) UpdateSalary(ID int64, salary domain.CustomerUpdateSalary) (*domain.Customer, error) {
	if m.UpdateSalaryFunc != nil {
		return m.UpdateSalaryFunc(ID, salary)
	}
	return &domain.Customer{}, nil
}
