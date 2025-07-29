package delivery

import (
	"github.com/Grafiters/archive/configs"
	"github.com/Grafiters/archive/internal/domain"
	"github.com/Grafiters/archive/middleware"
	"github.com/Grafiters/archive/utils"
	"github.com/gofiber/fiber/v2"
)

type customerHandler struct {
	customerUsecase domain.CustomerUsecase
	logger          *configs.LoggerFormat
}

func NewCustomerHandler(
	router fiber.Router,
	cs domain.CustomerUsecase,
	logger *configs.LoggerFormat,
) {
	h := &customerHandler{
		customerUsecase: cs,
		logger:          logger,
	}

	router.Get("/customer/get", middleware.Authenticate, h.Get)
	router.Put("/customer/update", middleware.Authenticate, h.Update)
	router.Put("/customer/update/salary", middleware.Authenticate, h.UpdateSalary)
}

// Customers
// @Router /api/customer/get [get]
// @Summary Get Customer
// @Description Get Customer data
// @Security Token
// @Tags Customers
// @Accept  json
// @Produce  json
// @Success 200 {object} domain.SingleResponse{data=domain.CustomerResponse}
// @Failure 422 {object} domain.SingleResponse
// @Failure 500 {object} domain.SingleResponse
func (ch *customerHandler) Get(c *fiber.Ctx) error {
	CurrentUser := c.Locals("CurrentUser").(*domain.Customer)

	customer, err := ch.customerUsecase.Get(c, CurrentUser.ID)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: utils.BodyParamMissing,
		})
	}
	return c.Status(fiber.StatusOK).JSON(domain.SingleResponse{
		Code:    fiber.StatusOK,
		Data:    customer,
		Message: "successfully created user",
	})
}

// Customers
// @Router /api/customer/update [put]
// @Summary Update Customer
// @Description Update Customer data for non crucial data
// @Tags Customers
// @Security Token
// @Accept  json
// @Produce  json
// @Param customer body domain.CustomerUpdate true "customer to create"
// @Success 200 {object} domain.SingleResponse{data=domain.CustomerResponse}
// @Failure 422 {object} domain.SingleResponse
// @Failure 500 {object} domain.SingleResponse
func (ch *customerHandler) Update(c *fiber.Ctx) error {
	CurrentUser := c.Locals("CurrentUser").(*domain.Customer)
	payload := new(domain.CustomerUpdate)

	if len(c.Body()) <= 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: utils.BodyParamMissing,
		})
	}

	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: utils.InvalidMessageBody,
		})
	}

	customer, err := ch.customerUsecase.Update(c, CurrentUser.ID, payload)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: utils.BodyParamMissing,
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.SingleResponse{
		Code:    fiber.StatusUnprocessableEntity,
		Data:    customer,
		Message: "successfully created user",
	})
}

// Customers
// @Router /api/customer/update/salary [put]
// @Summary Update Customer
// @Description Update Customer data for non crucial data
// @Tags Customers
// @Security Token
// @Accept  json
// @Produce  json
// @Param customer body domain.CustomerUpdateSalary true "customer to create"
// @Success 200 {object} domain.SingleResponse{data=domain.CustomerResponse}
// @Failure 422 {object} domain.SingleResponse
// @Failure 500 {object} domain.SingleResponse
func (ch *customerHandler) UpdateSalary(c *fiber.Ctx) error {
	CurrentUser := c.Locals("CurrentUser").(*domain.Customer)
	payload := new(domain.CustomerUpdateSalary)

	if len(c.Body()) <= 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: utils.BodyParamMissing,
		})
	}

	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: utils.InvalidMessageBody,
		})
	}
	customer, err := ch.customerUsecase.UpdateSalary(c, CurrentUser.ID, payload)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: utils.BodyParamMissing,
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.SingleResponse{
		Code:    fiber.StatusUnprocessableEntity,
		Data:    customer,
		Message: "successfully created user",
	})
	return nil
}
