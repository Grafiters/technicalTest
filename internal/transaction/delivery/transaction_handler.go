package delivery

import (
	"strconv"

	"github.com/Grafiters/archive/configs"
	"github.com/Grafiters/archive/internal/domain"
	"github.com/Grafiters/archive/middleware"
	"github.com/Grafiters/archive/utils"
	"github.com/gofiber/fiber/v2"
)

type transactionHandler struct {
	transactionUsecase domain.TransactionUsecase
	logger             *configs.LoggerFormat
}

func NewTransactionHandler(
	router fiber.Router,
	tu domain.TransactionUsecase,
	logger *configs.LoggerFormat,
) {
	h := &transactionHandler{
		transactionUsecase: tu,
		logger:             logger,
	}

	router.Get("/transaction/get", middleware.Authenticate, h.Get)
	router.Get("/transaction/get/:id", middleware.Authenticate, h.GetByID)
	router.Post("/transaction/create", middleware.Authenticate, h.Create)
	router.Post("/transaction/installment/pay", middleware.Authenticate, h.PayOff)
}

// Transactions
// @Router /api/transaction/get [get]
// @Summary Get Transaction
// @Description Get Transaction data
// @Tags Transactions
// @Security Token
// @Accept  json
// @Produce  json
// @Params transaction query domain.TransactionFilter false "filter for transaction"
// @Success 200 {object} domain.PaginationResponse{data=domain.TransactionResponse}
// @Failure 422 {object} domain.SingleResponse
// @Failure 500 {object} domain.SingleResponse
func (th *transactionHandler) Get(c *fiber.Ctx) error {
	CurrentUser := c.Locals("CurrentUser").(*domain.Customer)
	params := new(domain.TransactionFilter)
	if err := c.QueryParser(params); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: utils.BodyParamMissing,
		})
	}

	if params.PageSize == 0 {
		params.PageSize = 10
	}

	if params.Page == 0 {
		params.Page = 0
	}

	transaction, totalSize, err := th.transactionUsecase.Get(c, CurrentUser.ID, params)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: utils.BodyParamMissing,
		})
	}

	response := domain.PaginationResponse{
		Code:    fiber.StatusOK,
		Data:    convertToPaginationData(transaction, params.Page, params.PageSize, totalSize, params.Sort),
		Message: "successfully get transaction data",
	}

	return c.Status(fiber.StatusOK).JSON(response)

}

// Transactions
// @Router /api/transaction/create [post]
// @Summary Transaction
// @Description Create Transaction data for non crucial data
// @Tags Transactions
// @Security Token
// @Accept  json
// @Produce  json
// @Param transaction body domain.TransactionInput true "transaction to create"
// @Success 200 {object} domain.SingleResponse{data=domain.TransactionResponse}
// @Failure 422 {object} domain.SingleResponse
// @Failure 500 {object} domain.SingleResponse
func (th *transactionHandler) Create(c *fiber.Ctx) error {
	CurrentUser := c.Locals("CurrentUser").(*domain.Customer)
	payload := new(domain.TransactionInput)

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

	transaction, err := th.transactionUsecase.Create(c, CurrentUser.ID, payload)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.SingleResponse{
		Code:    fiber.StatusCreated,
		Data:    transaction,
		Message: "successfully create data transaction",
	})
}

// Transactions
// @Router /api/transaction/installment/pay [post]
// @Summary Transaction
// @Description Create Transaction data for non crucial data
// @Tags Transactions
// @Security Token
// @Accept  json
// @Produce  json
// @Param transaction body domain.BulkInstallmentUpdate true "transaction to pay installment"
// @Success 200 {object} domain.SingleResponse{data=domain.TransactionResponse}
// @Failure 422 {object} domain.SingleResponse
// @Failure 500 {object} domain.SingleResponse
func (th *transactionHandler) PayOff(c *fiber.Ctx) error {
	CurrentUser := c.Locals("CurrentUser").(*domain.Customer)
	payload := new(domain.BulkInstallmentUpdate)

	if len(c.Body()) <= 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: utils.BodyParamMissing,
		})
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: utils.InvalidMessageBody,
		})
	}

	transaction, err := th.transactionUsecase.BulkUpdateInstallment(c, CurrentUser.ID, payload)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.SingleResponse{
		Code:    fiber.StatusCreated,
		Data:    transaction,
		Message: "successfully update data installment",
	})
}

// Transactions
// @Router /api/transaction/get/{id} [get]
// @Summary Get Transaction
// @Description Get Transaction data
// @Tags Transactions
// @Security Token
// @Accept  json
// @Produce  json
// @Param id path int true "Transaction ID"
// @Success 200 {object} domain.PaginationResponse{data=domain.TransactionResponse}
// @Failure 422 {object} domain.SingleResponse
// @Failure 500 {object} domain.SingleResponse
func (th *transactionHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		th.logger.Error("failed query parsing data, err: %+v", err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: utils.BodyParamMissing,
		})
	}

	transaction, err := th.transactionUsecase.GetByID(c, int64(id))
	if err != nil {
		th.logger.Error("failed query parsing data, err: %+v", err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: utils.BodyParamMissing,
		})
	}
	return c.Status(fiber.StatusOK).JSON(domain.SingleResponse{
		Code:    fiber.StatusOK,
		Data:    transaction,
		Message: "successfully Get data transaction",
	})
}

func convertToPaginationData(transaction []*domain.TransactionResponse, page, pageSize, totalSize int, sort []domain.SortObject) domain.PaginationData {
	return domain.PaginationData{
		Content:    transaction,
		First:      page == 0,
		Last:       pageSize > 0 && totalSize <= pageSize*(page+1),
		Page:       page,
		PageSize:   pageSize,
		TotalSize:  totalSize,
		TotalPages: (totalSize + pageSize - 1) / pageSize,
		Sort:       sort,
	}
}
