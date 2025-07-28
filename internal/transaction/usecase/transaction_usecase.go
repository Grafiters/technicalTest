package usecase

import (
	"fmt"
	"time"

	"github.com/Grafiters/archive/configs"
	"github.com/Grafiters/archive/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type transactionUsecase struct {
	transactionRepo domain.TransactionRepository
	limitRepo       domain.LimitRepository
	customerRepo    domain.CustomerRepository
	logger          *configs.LoggerFormat
}

func NewTransactionUsecase(
	transactionRepo domain.TransactionRepository,
	customerRepo domain.CustomerRepository,
	limitRepo domain.LimitRepository,
	logger *configs.LoggerFormat,
) domain.TransactionUsecase {
	return &transactionUsecase{
		transactionRepo: transactionRepo,
		limitRepo:       limitRepo,
		logger:          logger,
	}
}

func (t *transactionUsecase) Create(ctx *fiber.Ctx, UserID int64, data *domain.TransactionInput) (*domain.TransactionResponse, error) {
	err := t.validateLimit(data)
	if err != nil {
		return nil, err
	}

	transaction, err := t.transactionRepo.Create(UserID, data)
	if err != nil {
		t.logger.Error("failed to create transaction, err: %+v", err)
		return nil, err
	}

	transactionToResponse := transaction.ToResponse()

	customer, err := t.customerRepo.GetByID(transaction.CustomerID)
	if err != nil {
		t.logger.Error("failed to get customer data, err: %+v", err)
		return nil, err
	}

	transactionToResponse.Customer = customer.ToCustomerResponse()

	limit, err := t.limitRepo.GetByID(transaction.LimitID)
	if err != nil {
		t.logger.Error("failed to get limit customer data, err: %+v", err)
		return nil, err
	}

	transactionToResponse.Limit = limit.ToLimitResponse()

	return transactionToResponse, nil
}

func (t *transactionUsecase) Get(ctx *fiber.Ctx, customerID int64, filter *domain.TransactionFilter) ([]*domain.TransactionResponse, int, error) {
	transaction, totalSize, err := t.transactionRepo.Get(customerID, filter)
	if err != nil {
		t.logger.Error("failed to get transaction customer data, err: %+v", err)
		return nil, 0, err
	}

	limits := []*domain.Limit{}
	customers := []*domain.Customer{}
	transactionResponse := []*domain.TransactionResponse{}
	for _, val := range transaction {
		txResponse := val.ToResponse()

		limit, err := t.limitRepo.GetByID(val.LimitID)
		if err != nil {
			t.logger.Error("failed to get limit customer data, err: %+v", err)
			return nil, 0, err
		}

		limitResponse := limit.ToLimitResponse()
		txResponse.Limit = limitResponse

		limits[val.LimitID] = limit

		customer, err := t.customerRepo.GetByID(val.CustomerID)
		if err != nil {
			t.logger.Error("failed to get customer data, err: %+v", err)
			return nil, 0, err
		}

		customerResponse := customer.ToCustomerResponse()
		txResponse.Customer = customerResponse
		customers[val.CustomerID] = customer

		transactionResponse = append(transactionResponse, txResponse)
	}

	return transactionResponse, totalSize, nil
}

func (t *transactionUsecase) GetByID(ctx *fiber.Ctx, CustomerID int64) (*domain.TransactionResponse, error) {
	transaction, err := t.transactionRepo.GetByID(CustomerID)
	if err != nil {
		t.logger.Error("failed to get transaction customer data, err: %+v", err)
		return nil, err
	}

	transactionResponse := transaction.ToResponse()

	limit, err := t.limitRepo.GetByID(transaction.LimitID)
	if err != nil {
		t.logger.Error("failed to get limit customer data, err: %+v", err)
		return nil, err
	}

	limitResponse := limit.ToLimitResponse()
	transactionResponse.Limit = limitResponse

	customer, err := t.customerRepo.GetByID(transaction.CustomerID)
	if err != nil {
		t.logger.Error("failed to get customer data, err: %+v", err)
		return nil, err
	}

	customerResponse := customer.ToCustomerResponse()
	transactionResponse.Customer = customerResponse

	return transactionResponse, nil
}

func (t *transactionUsecase) GetByCustomerID(ctx *fiber.Ctx, CustomerID int64) (*domain.TransactionResponse, error) {
	transaction, err := t.transactionRepo.GetByCustomerID(CustomerID)
	if err != nil {
		t.logger.Error("failed to get transaction customer data, err: %+v", err)
		return nil, err
	}

	transactionResponse := transaction.ToResponse()

	limit, err := t.limitRepo.GetByID(transaction.LimitID)
	if err != nil {
		t.logger.Error("failed to get limit customer data, err: %+v", err)
		return nil, err
	}

	limitResponse := limit.ToLimitResponse()
	transactionResponse.Limit = limitResponse

	customer, err := t.customerRepo.GetByID(transaction.CustomerID)
	if err != nil {
		t.logger.Error("failed to get customer data, err: %+v", err)
		return nil, err
	}

	customerResponse := customer.ToCustomerResponse()
	transactionResponse.Customer = customerResponse

	return transactionResponse, nil
}

func (t *transactionUsecase) BulkUpdateInstallment(c *fiber.Ctx, ID int64, data *domain.BulkInstallmentUpdate) (*domain.TransactionResponse, error) {
	// Ambil transaksi
	tx, err := t.transactionRepo.GetByID(ID)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	// Ambil semua installment log yang belum dibayar
	existingLogs, err := t.transactionRepo.GetInstallmentLogs(ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch existing installment logs: %w", err)
	}

	var unpaidTotal int64
	for _, l := range existingLogs {
		if l.PaidAt == nil {
			unpaidTotal += l.Amount
		}
	}

	var inputTotal int64
	for _, v := range data.InstallmentUpdae {
		inputTotal += v.Amount
	}

	if inputTotal != unpaidTotal {
		return nil, fmt.Errorf("installment total mismatch: expected %d, got %d", unpaidTotal, inputTotal)
	}
	logs, err := t.transactionRepo.BulkUpdateInstallment(ID, data)
	if err != nil {
		return nil, fmt.Errorf("failed to update installment logs: %w", err)
	}

	_, err = t.handleUpdateTransaction(tx.ID, "paid_off")
	if err != nil {
		fmt.Errorf("failed to update transaction status: %w", err)
		return nil, fmt.Errorf("failed to update transaction data")
	}

	response := tx.ToResponse()
	installmentResponse := []*domain.InstallmentLogRespponse{}
	for _, val := range logs {
		installmentResponse = append(installmentResponse, val.ToRespose())
	}

	response.InstallmentList = installmentResponse

	return response, nil
}

func (t *transactionUsecase) BulkinsertInstallment(ID int64, data *domain.BulkInstallmentInput) (*domain.TransactionResponse, error) {
	// Ambil transaksi
	tx, err := t.transactionRepo.GetByID(ID)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	// Ambil semua installment log yang belum dibayar
	existingLogs, err := t.transactionRepo.GetInstallmentLogs(ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch existing installment logs: %w", err)
	}

	var unpaidTotal int64
	for _, l := range existingLogs {
		if l.PaidAt == nil {
			unpaidTotal += l.Amount
		}
	}

	var inputTotal int64
	for _, v := range data.InstallmentInput {
		inputTotal += v.Amount
	}

	if inputTotal != unpaidTotal {
		return nil, fmt.Errorf("installment total mismatch: expected %d, got %d", unpaidTotal, inputTotal)
	}
	logs, err := t.transactionRepo.BulkInsertInstallment(ID, data)
	if err != nil {
		return nil, fmt.Errorf("failed to insert installment logs: %w", err)
	}

	response := tx.ToResponse()
	installmentResponse := []*domain.InstallmentLogRespponse{}
	for _, val := range logs {
		installmentResponse = append(installmentResponse, val.ToRespose())
	}

	response.InstallmentList = installmentResponse

	return response, nil
}

func (t *transactionUsecase) mappinInstallmentData(ID int64, data *domain.Transaction) error {
	limit, err := t.limitRepo.GetByID(data.LimitID)
	if err != nil {
		t.logger.Error("failed to get limit customer data, err: %+v", err)
		return err
	}

	installmentInput := []*domain.InstallmentInput{}
	installmentPerMonth := data.Installment / int64(limit.Tenor)
	startDate := time.Now().AddDate(0, 1, 0)
	for i := 1; i <= limit.Tenor; i++ {
		dueDate := time.Date(startDate.Year(), startDate.Month()+time.Month(i-1), 1, 0, 0, 0, 0, time.Local)
		installmentInput = append(installmentInput, &domain.InstallmentInput{
			Month:   i,
			Amount:  installmentPerMonth,
			DueDate: dueDate,
		})
	}

	builkInstallmentInput := &domain.BulkInstallmentInput{
		TransactionID:    ID,
		InstallmentInput: installmentInput,
	}

	_, err = t.BulkinsertInstallment(ID, builkInstallmentInput)
	if err != nil {
		t.logger.Error("failed to generate installment logs, err: %+v", err)
		return err
	}

	return nil
}

func (t *transactionUsecase) handleUpdateTransaction(ID int64, status string) (*domain.Transaction, error) {
	installment, err := t.transactionRepo.GetInstallmentLogs(ID)
	if err != nil {
		fmt.Errorf("failed to fetch existing installment logs: %w", err)
		return nil, fmt.Errorf("failed to get installment log data")
	}

	allPaid := true
	for _, val := range installment {
		if val.PaidAt != nil {
			allPaid = false
			break
		}
	}

	if allPaid {
		tx, err := t.transactionRepo.PayOff(ID)
		if err != nil {
			fmt.Errorf("failed to update transaction status: %w", err)
			return nil, fmt.Errorf("failed to update transaction data")
		}

		return tx, nil
	}

	return nil, fmt.Errorf("failed to update transacton data")
}

func (t *transactionUsecase) validateLimit(data *domain.TransactionInput) error {
	limit, err := t.limitRepo.GetByID(data.LimitID)
	if err != nil {
		t.logger.Error("failed to get limit customer, err: %+v", err)
		return err
	}

	totalAmount := data.OTR + data.AdminFee + data.Installment
	if limit.Amount < totalAmount {
		t.logger.Error("Transaction is to higher then limit customer, err: ", limit.Amount)
		return fmt.Errorf("your limit is to low then your transaction")
	}

	return nil
}
