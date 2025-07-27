package domain

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type Transaction struct {
	ID          int64     `json:"id" gorm:"column:id;type:bigserial;primaryKey;autoIncrement"`
	UserID      int64     `json:"user_id"`
	CustonerID  int64     `json:"customer_id"`
	ContractNo  int64     `json:"contract_no"`
	OTR         int64     `json:"otr"`
	AdminFee    int64     `json:"admin_fee"`
	Installment int64     `json:"installment"`
	AssetName   int64     `json:"asset_name"`
	CreatedAT   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAT   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type TransactionInput struct {
	CustonerID  int64 `json:"customer_id"`
	ContractNo  int64 `json:"contract_no"`
	OTR         int64 `json:"otr"`
	AdminFee    int64 `json:"admin_fee"`
	Installment int64 `json:"installment"`
	AssetName   int64 `json:"asset_name"`
}

type TransactionFilter struct {
}

type TransactionResponse struct {
	ID          int64     `json:"id" gorm:"column:id;type:bigserial;primaryKey;autoIncrement"`
	Customer    *Customer `json:"customer"`
	ContractNo  int64     `json:"contract_no"`
	OTR         int64     `json:"otr"`
	AdminFee    int64     `json:"admin_fee"`
	Installment int64     `json:"installment"`
	AssetName   int64     `json:"asset_name"`
	CreatedAT   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAT   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (t *Transaction) ToResponse() *TransactionResponse {
	return &TransactionResponse{
		ID:          t.ID,
		ContractNo:  t.ContractNo,
		OTR:         t.OTR,
		AdminFee:    t.AdminFee,
		Installment: t.Installment,
		AssetName:   t.AssetName,
		CreatedAT:   t.CreatedAT,
		UpdatedAT:   t.UpdatedAT,
	}
}

type TransactionUsecase interface {
	Create(ctx *fiber.Ctx, data *TransactionInput) error
	Get(ctx *fiber.Ctx, filter *TransactionFilter) ([]*TransactionResponse, int, error)
	GetByCustomerID(ctx *fiber.Ctx, ID int64) (*TransactionResponse, error)
}

type TransactionRepository interface {
	Create(ctx *fiber.Ctx, data *TransactionInput) error
	Get(ctx *fiber.Ctx, filter *TransactionFilter) ([]*Transaction, int, error)
	GetByCustomerID(ctx *fiber.Ctx, ID int64) (*Transaction, error)
	Update(
		ctx *fiber.Ctx,
		ID int64,
		oldData *Transaction,
		newData *Transaction,
	) (*Transaction, error)
}
