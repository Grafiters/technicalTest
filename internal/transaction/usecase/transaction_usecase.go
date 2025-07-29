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
		customerRepo:    customerRepo,
		limitRepo:       limitRepo,
		logger:          logger,
	}
}

func (t *transactionUsecase) Create(ctx *fiber.Ctx, UserID int64, data *domain.TransactionInput) (*domain.TransactionResponse, error) {
	limit, err := t.validateLimit(UserID, data)
	if err != nil {
		t.logger.Error("failed to failed limit, err: ", err)
		return nil, err
	}

	transactionData := &domain.Transaction{
		CustomerID:  UserID,
		LimitID:     limit.ID,
		ContractNo:  data.ContractNo,
		OTR:         data.OTR,
		AdminFee:    data.AdminFee,
		Installment: data.Installment,
		AssetName:   data.AssetName,
		Status:      "active",
	}

	transaction, err := t.transactionRepo.Create(UserID, transactionData)
	if err != nil {
		t.logger.Error("failed to create transaction, err: ", err)
		return nil, err
	}

	transactionToResponse := transaction.ToResponse()

	customer, err := t.customerRepo.GetByID(transaction.CustomerID)
	if err != nil {
		t.logger.Error("failed to get customer data, err: ", err)
		return nil, err
	}

	t.mappinInstallmentData(UserID, transaction)

	transactionToResponse.Customer = customer.ToCustomerResponse()

	transactionToResponse.Limit = limit.ToLimitResponse()

	installment, err := t.transactionRepo.GetInstallmentLogs(transactionToResponse.ID)
	if err != nil {
		t.logger.Error("failed to get installment log data, err: ", err)
		return nil, err
	}

	installmentResponse := []*domain.InstallmentLogRespponse{}
	for _, i := range installment {
		installmentResponse = append(installmentResponse, i.ToRespose())
	}

	transactionToResponse.InstallmentList = installmentResponse

	return transactionToResponse, nil
}

func (t *transactionUsecase) Get(ctx *fiber.Ctx, customerID int64, filter *domain.TransactionFilter) ([]*domain.TransactionResponse, int, error) {
	transaction, totalSize, err := t.transactionRepo.Get(customerID, filter)
	if err != nil {
		t.logger.Error("failed to get transaction customer data, err: ", err)
		return nil, 0, err
	}

	limits := make(map[int64]*domain.Limit)
	customers := make(map[int64]*domain.Customer)
	transactionResponse := []*domain.TransactionResponse{}
	for _, val := range transaction {
		txResponse := val.ToResponse()

		limit, err := t.limitRepo.GetByID(val.LimitID)
		if err != nil {
			t.logger.Error("failed to get limit customer data, err: ", err)
			return nil, 0, err
		}

		limitResponse := limit.ToLimitResponse()
		txResponse.Limit = limitResponse

		limits[val.LimitID] = limit

		customer, err := t.customerRepo.GetByID(val.CustomerID)
		if err != nil {
			t.logger.Error("failed to get customer data, err: ", err)
			return nil, 0, err
		}

		customerResponse := customer.ToCustomerResponse()
		txResponse.Customer = customerResponse
		customers[val.CustomerID] = customer

		installment, err := t.transactionRepo.GetInstallmentLogs(val.ID)
		if err != nil {
			t.logger.Error("failed to get installment log data, err: ", err)
			return nil, 0, err
		}

		installmentResponse := []*domain.InstallmentLogRespponse{}
		for _, i := range installment {
			installmentResponse = append(installmentResponse, i.ToRespose())
		}

		txResponse.InstallmentList = installmentResponse

		transactionResponse = append(transactionResponse, txResponse)
	}

	return transactionResponse, totalSize, nil
}

func (t *transactionUsecase) GetByID(ctx *fiber.Ctx, CustomerID int64) (*domain.TransactionResponse, error) {
	transaction, err := t.transactionRepo.GetByID(CustomerID)
	if err != nil {
		t.logger.Error("failed to get transaction customer data, err: ", err)
		return nil, err
	}

	transactionResponse := transaction.ToResponse()

	limit, err := t.limitRepo.GetByID(transaction.LimitID)
	if err != nil {
		t.logger.Error("failed to get limit customer data, err: ", err)
		return nil, err
	}

	limitResponse := limit.ToLimitResponse()
	transactionResponse.Limit = limitResponse

	customer, err := t.customerRepo.GetByID(transaction.CustomerID)
	if err != nil {
		t.logger.Error("failed to get customer data, err: ", err)
		return nil, err
	}

	customerResponse := customer.ToCustomerResponse()
	transactionResponse.Customer = customerResponse

	installment, err := t.transactionRepo.GetInstallmentLogs(transaction.ID)
	if err != nil {
		t.logger.Error("failed to get installment log data, err: ", err)
		return nil, err
	}

	installmentResponse := []*domain.InstallmentLogRespponse{}
	for _, i := range installment {
		installmentResponse = append(installmentResponse, i.ToRespose())
	}

	transactionResponse.InstallmentList = installmentResponse

	return transactionResponse, nil
}

func (t *transactionUsecase) GetByCustomerID(ctx *fiber.Ctx, CustomerID int64) (*domain.TransactionResponse, error) {
	transaction, err := t.transactionRepo.GetByCustomerID(CustomerID)
	if err != nil {
		t.logger.Error("failed to get transaction customer data, err: ", err)
		return nil, err
	}

	transactionResponse := transaction.ToResponse()

	limit, err := t.limitRepo.GetByID(transaction.LimitID)
	if err != nil {
		t.logger.Error("failed to get limit customer data, err: ", err)
		return nil, err
	}

	limitResponse := limit.ToLimitResponse()
	transactionResponse.Limit = limitResponse

	customer, err := t.customerRepo.GetByID(transaction.CustomerID)
	if err != nil {
		t.logger.Error("failed to get customer data, err: ", err)
		return nil, err
	}

	customerResponse := customer.ToCustomerResponse()
	transactionResponse.Customer = customerResponse

	return transactionResponse, nil
}

func (t *transactionUsecase) BulkUpdateInstallment(c *fiber.Ctx, ID int64, data *domain.BulkInstallmentUpdate) (*domain.TransactionResponse, error) {
	// Ambil transaksi
	tx, err := t.transactionRepo.GetByID(data.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	existingLogs, err := t.transactionRepo.GetInstallmentLogs(tx.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch existing installment logs: %w", err)
	}
	var month []int64
	for _, item := range data.InstallmentUpdae {
		month = append(month, int64(item.Month))
	}

	var unpaidTotal int64
	for _, l := range existingLogs {
		isMonthInclude := IsMonthIncluded(int64(l.Month), month)
		if l.PaidAt == nil && isMonthInclude {
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

func (t *transactionUsecase) BulkinsertInstallment(ID int64, data *domain.BulkInstallmentInput) (int64, error) {
	// Ambil semua installment log yang belum dibayar
	installment, err := t.transactionRepo.BulkInsertInstallment(ID, data)
	if err != nil {
		t.logger.Error("failed to bulk insert installment, err: ", err)
		return 0, err
	}

	return int64(len(installment)), nil
}

func (t *transactionUsecase) mappinInstallmentData(ID int64, data *domain.Transaction) error {
	limit, err := t.limitRepo.GetByID(data.LimitID)
	if err != nil {
		t.logger.Error("failed to get limit customer data, err: ", err)
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

	_, err = t.BulkinsertInstallment(data.ID, builkInstallmentInput)
	if err != nil {
		t.logger.Error("failed to generate installment logs, err: ", err)
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

	return nil, nil
}

func (t *transactionUsecase) validateLimit(customerID int64, data *domain.TransactionInput) (*domain.Limit, error) {
	limit, err := t.limitRepo.GetByTenor(customerID, data.Tenor)
	if err != nil {
		t.logger.Error("failed to get limit customer, err: ", err)
		return &domain.Limit{}, err
	}

	totalAmount := data.OTR + data.AdminFee + data.Installment
	if limit.Amount < totalAmount {
		t.logger.Error("Transaction is to higher then limit customer, err: ", limit.Amount)
		return &domain.Limit{}, fmt.Errorf("your limit is to low then your transaction")
	}

	return limit, nil
}

func IsMonthIncluded(amount int64, list []int64) bool {
	for _, v := range list {
		if v == amount {
			return true
		}
	}
	return false
}
