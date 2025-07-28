package repository

import (
	"fmt"

	"github.com/Grafiters/archive/configs"
	"github.com/Grafiters/archive/internal/domain"
	"github.com/Grafiters/archive/utils"
	"gorm.io/gorm"
)

type mysqlCustomerRepository struct {
	db     *gorm.DB
	logger *configs.LoggerFormat
}

func NewCustomerRepository(db *gorm.DB, logger *configs.LoggerFormat) domain.CustomerRepository {
	return &mysqlCustomerRepository{db, logger}
}

func (u *mysqlCustomerRepository) Create(data *domain.CustomerInput) (*domain.Customer, error) {
	newCustomer := &domain.Customer{
		Email:          data.Email,
		NIK:            data.NIK,
		FullName:       data.FullName,
		LegalName:      data.LegalName,
		BirthPlace:     data.BirthPlace,
		BirthDate:      data.BirthDate,
		Salary:         data.Salary,
		KTPImageUrl:    data.KTPImage,
		SelfieImageUrl: data.SelfieImage,
	}

	err := u.db.Transaction(func(tx *gorm.DB) error {
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

func (u *mysqlCustomerRepository) Get(filter *domain.CustomerFilter) ([]*domain.Customer, int64, error) {
	var (
		customer         []*domain.Customer
		totalSize        int64
		whereQuery, args = getWhereClause(filter)
		orderClause      = domain.GetOrderClause(filter.Sort)
	)

	query := u.db.Where(whereQuery, args...).Order(orderClause)

	if err := query.Model(&domain.Customer{}).Count(&totalSize).Error; err != nil {
		u.logger.Error("failed to count query customer, err: %+v", err)
		return nil, 0, err
	}

	// jika tidak ada data
	if totalSize == 0 {
		return []*domain.Customer{}, 0, fmt.Errorf(utils.DataNotFound)
	}

	// ambil data dengan pagination
	if err := query.Limit(filter.PageSize).Offset(filter.Page).Find(&customer).Error; err != nil {
		u.logger.Error("failed to query customer, err: %+v", err)
		return nil, 0, err
	}

	return customer, totalSize, nil
}

func (u *mysqlCustomerRepository) GetByID(ID int64) (*domain.Customer, error) {
	var customer *domain.Customer
	result := u.db.Where("id = ?", ID).First(&customer)
	if result.Error != nil {
		u.logger.Error("failed to get customer, err: ", result.Error)
		return &domain.Customer{}, result.Error
	}
	return customer, nil
}

func (u *mysqlCustomerRepository) Update(
	ID int64,
	newData *domain.CustomerUpdate,
) (*domain.Customer, error) {
	customer, err := u.GetByID(ID)
	if err != nil {
		u.logger.Error("failed get data customer, err: %+v", err)
		return &domain.Customer{}, fmt.Errorf(utils.DataNotFound)
	}

	fields, values := domain.MapToKeyValueArrays(newData)

	if len(fields) == 0 || len(values) == 0 {
		return &domain.Customer{}, fmt.Errorf("no update data for customer")
	}

	updateMap := make(map[string]interface{})
	for i := range fields {
		updateMap[fields[i]] = values[i]
	}

	if err := u.db.Model(&customer).Updates(updateMap).Error; err != nil {
		u.logger.Error("faild to update data customer, err: %+v", err)
		return nil, err
	}

	u.db.Where("id = ?", ID).First(&customer)
	return customer, nil
}

func (u *mysqlCustomerRepository) UpdateSalary(ID int64, salary domain.CustomerUpdateSalary) (*domain.Customer, error) {
	customer, err := u.GetByID(ID)
	if err != nil {
		u.logger.Error("failed get data customer, err: %+v", err)
		return &domain.Customer{}, fmt.Errorf(utils.DataNotFound)
	}

	if err := u.db.Model(&customer).Where("id = ?", ID).Update("salary", salary.Salary).Error; err != nil {
		u.logger.Error("faild to update data customer, err: %+v", err)
		return nil, err
	}

	u.db.Where("id = ?", ID).First(&customer)
	return customer, nil
}

func getWhereClause(filter *domain.CustomerFilter) (string, []interface{}) {
	var (
		whereClause string
		// whereClauses []string
		args []interface{}
	)

	return whereClause, args
}
