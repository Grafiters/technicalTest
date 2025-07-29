package mocks

import "github.com/Grafiters/archive/internal/domain"

type MockLimitRepository struct {
	GetFunc              func(filter *domain.LimitFilter) ([]*domain.Limit, error)
	GetByCustommerIDFunc func(customerID int64) ([]*domain.Limit, error)
	GetByIDFunc          func(ID int64) (*domain.Limit, error)
	GetByTenorFunc       func(customerID int64, tenor int64) (*domain.Limit, error)
	BulkCreateLimitFunc  func(data *domain.BulkLimitInput) ([]*int64, error)
	BulkUpdateLimitFunc  func(data *domain.BulkLimitInput) ([]*int64, error)
}

func (m *MockLimitRepository) Get(filter *domain.LimitFilter) ([]*domain.Limit, error) {
	return m.GetFunc(filter)
}

func (m *MockLimitRepository) GetByCustommerID(customerID int64) ([]*domain.Limit, error) {
	return m.GetByCustommerIDFunc(customerID)
}

func (m *MockLimitRepository) GetByID(ID int64) (*domain.Limit, error) {
	return m.GetByIDFunc(ID)
}

func (m *MockLimitRepository) GetByTenor(customerID int64, tenor int64) (*domain.Limit, error) {
	return m.GetByTenorFunc(customerID, tenor)
}

func (m *MockLimitRepository) BulkCreateLimit(data *domain.BulkLimitInput) ([]*int64, error) {
	return m.BulkCreateLimitFunc(data)
}

func (m *MockLimitRepository) BulkUpdateLimit(data *domain.BulkLimitInput) ([]*int64, error) {
	return m.BulkUpdateLimitFunc(data)
}
