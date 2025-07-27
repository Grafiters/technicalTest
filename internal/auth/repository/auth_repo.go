package repository

import (
	"github.com/Grafiters/archive/configs"
	"github.com/Grafiters/archive/internal/domain"
	"gorm.io/gorm"
)

type mysqlAuthRepository struct {
	db     *gorm.DB
	logger *configs.LoggerFormat
}

func NewAuthRepository(
	db *gorm.DB,
	logger *configs.LoggerFormat,
) domain.AuthRepository {
	return &mysqlAuthRepository{db, logger}
}

// Login implements domain.AuthRepository.
func (m *mysqlAuthRepository) Login(auth *domain.AuthRequest) (*domain.Customer, error) {
	var customer *domain.Customer
	err := m.db.Where("email = ?", auth.Email).First(&customer)
	if err != nil {
		m.logger.Error("failed to get customer data for login, err: %+v", err)
		return &domain.Customer{}, err.Error
	}

	return customer, nil
}

// Registeer implements domain.AuthRepository.
func (m *mysqlAuthRepository) Register(data *domain.RegisterRequest) (*domain.Customer, error) {
	newCustomer := &domain.Customer{
		Email:    data.Email,
		Password: data.Password,
	}

	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(newCustomer).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return &domain.Customer{}, err
	}

	return newCustomer, nil
}

func (m *mysqlAuthRepository) GetByEmail(email string) (*domain.Customer, error) {
	var customer *domain.Customer
	err := m.db.Where("email = ?", email).First(&customer)
	if err != nil {
		m.logger.Error("failed to get customer data for login, err: %+v", err)
		return &domain.Customer{}, err.Error
	}

	return customer, nil
}
