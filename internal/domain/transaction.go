package domain

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type Transaction struct {
	ID          int64     `json:"id" gorm:"column:id;type:bigserial;primaryKey;autoIncrement"`
	CustomerID  int64     `json:"customer_id"`
	LimitID     int64     `json:"limit_id"`
	ContractNo  int64     `json:"contract_no"`
	OTR         int64     `json:"otr"`
	AdminFee    int64     `json:"admin_fee"`
	Installment int64     `json:"installment"`
	AssetName   string    `json:"asset_name"`
	Status      string    `json:"status"`
	CreatedAT   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAT   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type InstallmentLog struct {
	ID            int64      `json:"id"`
	TransactionID int64      `json:"transaction_id"`
	Month         int        `json:"month"`
	Amount        int64      `json:"amount"`
	DueDate       time.Time  `json:"due_date"`
	PaidAt        *time.Time `json:"paid_at"`
	CreatedAT     time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAT     time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

type InstallmentInput struct {
	Month   int       `json:"month"`
	Amount  int64     `json:"amount"`
	DueDate time.Time `json:"due_date"`
}

type BulkInstallmentInput struct {
	TransactionID    int64               `json:"transaction_id"`
	InstallmentInput []*InstallmentInput `json:"installment_input"`
}

type InstallmentUpdae struct {
	Month  int       `json:"month"`
	Amount int64     `json:"amount"`
	PaidAt time.Time `json:"paid_at"`
}

type BulkInstallmentUpdate struct {
	TransactionID    int64               `json:"transaction_id"`
	InstallmentUpdae []*InstallmentUpdae `json:"installment_update"`
}

type TransactionInput struct {
	Tenor       int64  `json:"tenor"`
	ContractNo  int64  `json:"contract_no"`
	OTR         int64  `json:"otr"`
	AdminFee    int64  `json:"admin_fee"`
	Installment int64  `json:"installment"`
	AssetName   string `json:"asset_name"`
}

type TransactionFilter struct {
	Page     int `json:"page" query:"page" validate:"omitempty"`
	PageSize int `json:"page_size" query:"pageSize" validate:"omitempty,oneof=10 20 50 100"`
	Sort     []SortObject
}

type TransactionResponse struct {
	ID          int64     `json:"id" gorm:"column:id;type:bigserial;primaryKey;autoIncrement"`
	ContractNo  int64     `json:"contract_no"`
	OTR         int64     `json:"otr"`
	AdminFee    int64     `json:"admin_fee"`
	Installment int64     `json:"installment"`
	AssetName   string    `json:"asset_name"`
	Status      string    `json:"status"`
	CreatedAT   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAT   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Customer        *CustomerResponse          `json:"customers,omitempty"`
	Limit           *LimitResponse             `json:"limits,omitempty"`
	InstallmentList []*InstallmentLogRespponse `json:"installment_list"`
}

type InstallmentLogRespponse struct {
	ID        int64      `json:"id"`
	Month     int        `json:"month"`
	Amount    int64      `json:"amount"`
	DueDate   time.Time  `json:"due_date"`
	PaidAt    *time.Time `json:"paid_at"`
	CreatedAT time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAT time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (t *Transaction) ToResponse() *TransactionResponse {
	return &TransactionResponse{
		ID:          t.ID,
		ContractNo:  t.ContractNo,
		OTR:         t.OTR,
		AdminFee:    t.AdminFee,
		Installment: t.Installment,
		AssetName:   t.AssetName,
		Status:      t.Status,
		CreatedAT:   t.CreatedAT,
		UpdatedAT:   t.UpdatedAT,
	}
}

func (il *InstallmentLog) ToRespose() *InstallmentLogRespponse {
	return &InstallmentLogRespponse{
		ID:        il.ID,
		Month:     il.Month,
		Amount:    il.Amount,
		DueDate:   il.DueDate,
		PaidAt:    il.PaidAt,
		CreatedAT: il.CreatedAT,
		UpdatedAT: il.UpdatedAT,
	}
}

type TransactionUsecase interface {
	Create(ctx *fiber.Ctx, UserID int64, data *TransactionInput) (*TransactionResponse, error)
	Get(ctx *fiber.Ctx, customerID int64, filter *TransactionFilter) ([]*TransactionResponse, int, error)
	GetByCustomerID(ctx *fiber.Ctx, ID int64) (*TransactionResponse, error)
	GetByID(ctx *fiber.Ctx, ID int64) (*TransactionResponse, error)
	BulkinsertInstallment(ID int64, data *BulkInstallmentInput) (int64, error)
	BulkUpdateInstallment(ctx *fiber.Ctx, ID int64, data *BulkInstallmentUpdate) (*TransactionResponse, error)
}

type TransactionRepository interface {
	Create(UserID int64, data *Transaction) (*Transaction, error)
	Get(UserId int64, filter *TransactionFilter) ([]*Transaction, int, error)
	GetByID(ID int64) (*Transaction, error)
	GetByCustomerID(ID int64) (*Transaction, error)
	PayOff(ID int64) (*Transaction, error)
	BulkInsertInstallment(ID int64, data *BulkInstallmentInput) ([]*InstallmentLog, error)
	BulkUpdateInstallment(ID int64, data *BulkInstallmentUpdate) ([]*InstallmentLog, error)
	GetInstallmentLogs(transactionID int64) ([]*InstallmentLog, error)
	GetInstallmentsByDueDate(date string) ([]*InstallmentLog, error)
}
