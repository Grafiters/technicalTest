package repository

import (
	"fmt"

	"github.com/Grafiters/archive/configs"
	"github.com/Grafiters/archive/internal/domain"
	"github.com/Grafiters/archive/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type mysqlTransactionRepository struct {
	db     *gorm.DB
	logger *configs.LoggerFormat
}

func NewTranscationRepository(
	db *gorm.DB,
	logger *configs.LoggerFormat,
) domain.TransactionRepository {
	return &mysqlTransactionRepository{db, logger}
}

func (m *mysqlTransactionRepository) Create(UserID int64, data *domain.Transaction) (*domain.Transaction, error) {

	totalAmount := data.OTR + data.AdminFee + data.Installment

	err := m.db.Transaction(func(tx *gorm.DB) error {
		var limit domain.Limit
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&limit, "id = ?", data.LimitID).Error; err != nil {
			return fmt.Errorf("limit not found: %w", err)
		}

		if limit.Amount < totalAmount {
			return fmt.Errorf("insufficient limit: need %d, available %d", totalAmount, limit.Amount)
		}

		if err := tx.Create(data).Error; err != nil {
			return err
		}

		if err := tx.Exec(`UPDATE limits SET amount = amount - ? WHERE id = ?`, totalAmount, data.LimitID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (m *mysqlTransactionRepository) Get(UserId int64, filter *domain.TransactionFilter) ([]*domain.Transaction, int, error) {
	var (
		transaction      []*domain.Transaction
		totalSize        int64
		whereQuery, args = getWhereClause(filter)
	)

	query := m.db.Where("customer_id = ?", UserId)

	if whereQuery != "" {
		query = query.Where(whereQuery, args...)
	}

	err := query.Model(&domain.Transaction{}).Count(&totalSize)
	if err.Error != nil {
		m.logger.Error("failed to count query transaction, err: ", err)
		return nil, 0, err.Error
	}

	if totalSize == 0 {
		return nil, 0, fmt.Errorf(utils.DataNotFound)
	}

	// ambil data dengan pagination
	err = query.Limit(filter.PageSize).Offset(filter.Page).Find(&transaction)
	if err.Error != nil {
		m.logger.Error("failed to get query transaction, err: ", err)
		return nil, 0, err.Error
	}

	return transaction, int(totalSize), nil
}
func (m *mysqlTransactionRepository) GetByID(ID int64) (*domain.Transaction, error) {
	var (
		transaction *domain.Transaction
	)

	err := m.db.Where("id = ?", ID).First(&transaction)
	if err.Error != nil {
		m.logger.Error("failed to query transaction, err: %+v", err)
		return nil, err.Error
	}

	return transaction, nil
}

func (m *mysqlTransactionRepository) GetByCustomerID(ID int64) (*domain.Transaction, error) {
	var transaction *domain.Transaction
	query := m.db.Where("customer_id = ?")

	// ambil data dengan pagination
	if err := query.Find(&transaction).Error; err != nil {
		m.logger.Error("failed to query customer, err: %+v", err)
		return nil, err
	}

	return transaction, nil
}

func (m *mysqlTransactionRepository) GetInstallmentLogs(ID int64) ([]*domain.InstallmentLog, error) {
	var installmentLog []*domain.InstallmentLog
	query := m.db.Where("transaction_id = ?", ID)

	err := query.Find(&installmentLog)
	if err.Error != nil {
		m.logger.Error("failed to get data installment log, err: %+v", err)
		return nil, err.Error
	}

	return installmentLog, nil
}

func (m *mysqlTransactionRepository) PayOff(ID int64) (*domain.Transaction, error) {
	var transaction domain.Transaction

	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&transaction, "id = ? AND status = ?", ID, "active").Error; err != nil {
			return fmt.Errorf("transaction not found or already paid: %w", err)
		}

		if err := tx.Model(&transaction).Where("id = ?", transaction.ID).Update("status", "paid_off").Error; err != nil {
			return err
		}

		transaction.Status = "paid_off"
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (m *mysqlTransactionRepository) BulkInsertInstallment(ID int64, data *domain.BulkInstallmentInput) ([]*domain.InstallmentLog, error) {
	installmentLog := []*domain.InstallmentLog{}
	for _, val := range data.InstallmentInput {
		newInstallment := &domain.InstallmentLog{
			TransactionID: ID,
			Month:         val.Month,
			Amount:        val.Amount,
			DueDate:       val.DueDate,
		}

		err := m.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(newInstallment).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			m.logger.Error("failed to create installment log, err: %+v", err)
			return nil, err
		}

		installmentLog = append(installmentLog, newInstallment)
	}

	return installmentLog, nil
}

func (m *mysqlTransactionRepository) BulkUpdateInstallment(ID int64, data *domain.BulkInstallmentUpdate) ([]*domain.InstallmentLog, error) {
	for _, val := range data.InstallmentUpdae {
		err := m.db.Transaction(func(tx *gorm.DB) error {
			result := tx.Model(&domain.InstallmentLog{}).Where("transaction_id = ? AND month = ?", data.TransactionID, val.Month).Update("paid_at", val.PaidAt)

			if result.Error != nil {
				return result.Error
			}

			return nil
		})

		if err != nil {
			m.logger.Error("failed to update installment log, err: %+v", err)
			return nil, err
		}
	}

	installment, err := m.GetInstallmentLogs(ID)
	if err != nil {
		m.logger.Error("failed to get result of installment log, err: %+v", err)
		return nil, err
	}

	return installment, nil
}

// GetInstallmentsByDueDate implements domain.TransactionRepository.
func (m *mysqlTransactionRepository) GetInstallmentsByDueDate(date string) ([]*domain.InstallmentLog, error) {
	var installment []*domain.InstallmentLog
	err := m.db.Where("DATE(due_date) = ? AND paid_at IS NULL", date).Find(&installment)
	if err.Error != nil {
		return nil, err.Error
	}

	return installment, nil
}

func getWhereClause(filter *domain.TransactionFilter) (string, []interface{}) {
	var (
		whereClause string
		// whereClauses []string
		args []interface{}
	)

	return whereClause, args
}
