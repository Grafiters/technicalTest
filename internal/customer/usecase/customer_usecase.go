package usecase

import (
	"fmt"

	"github.com/Grafiters/archive/configs"
	"github.com/Grafiters/archive/internal/domain"
	"github.com/Grafiters/archive/utils"
	"github.com/gofiber/fiber/v2"
)

type customerUsecase struct {
	customerRepo domain.CustomerRepository
	limitRepo    domain.LimitRepository
	logger       *configs.LoggerFormat
}

func NewCustomerUsecase(
	customerRepo domain.CustomerRepository,
	limitRepo domain.LimitRepository,
	logger *configs.LoggerFormat,
) domain.CustomerUsecase {
	return &customerUsecase{
		customerRepo: customerRepo,
		limitRepo:    limitRepo,
		logger:       logger,
	}
}

func (c *customerUsecase) Create(ctx *fiber.Ctx, data *domain.CustomerInput) (*domain.CustomerResponse, error) {
	customer, err := c.customerRepo.Create(data)
	if err != nil {
		c.logger.Error("failed to create customer, err: %+v", err)
		return &domain.CustomerResponse{}, fmt.Errorf(utils.ProsessError)
	}

	customerResponse := customer.ToCustomerResponse()
	return customerResponse, nil
}

func (c *customerUsecase) Get(ctx *fiber.Ctx, ID int64) (*domain.CustomerResponse, error) {
	customer, err := c.customerRepo.GetByID(ID)
	if err != nil {
		c.logger.Error("failed to get customer, err: %+v", err)
		return &domain.CustomerResponse{}, err
	}

	customerResponse := customer.ToCustomerResponse()
	return customerResponse, nil
}

func (c *customerUsecase) Update(ctx *fiber.Ctx, ID int64, data *domain.CustomerUpdate) (*domain.CustomerResponse, error) {
	customer, err := c.customerRepo.Update(ID, data)
	if err != nil {
		c.logger.Error("failed to update customer, err: %+v", err)
		return &domain.CustomerResponse{}, fmt.Errorf(utils.ProsessError)
	}

	customerResponse := customer.ToCustomerResponse()
	return customerResponse, nil
}

func (c *customerUsecase) UpdateSalary(ctx *fiber.Ctx, ID int64, salary *domain.CustomerUpdateSalary) (*domain.CustomerResponse, error) {
	customer, err := c.customerRepo.UpdateSalary(ID, *salary)
	if err != nil {
		c.logger.Error("failed to update customer, err: %+v", err)
		return &domain.CustomerResponse{}, fmt.Errorf(utils.ProsessError)
	}

	err = c.handleCerateLimit(customer)

	customerResponse := customer.ToCustomerResponse()
	return customerResponse, nil
}

func (c *customerUsecase) handleCerateLimit(data *domain.Customer) error {
	tenorLimit := domain.BuildTenorFactor(data.Salary)
	limitInput := &domain.BulkLimitInput{
		CustomerID: data.ID,
		LimitTenor: tenorLimit,
	}

	_, err := c.limitRepo.BulkCreateLimit(limitInput)
	if err != nil {
		c.logger.Error("failed to bulk create limit, err: %+v", err)
		return err
	}

	return nil

}
