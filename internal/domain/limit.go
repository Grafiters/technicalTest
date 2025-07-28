package domain

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type Limit struct {
	ID         int64     `json:"id" gorm:"column:id;type:bigserial;primaryKey;autoIncrement"`
	CustomerID int64     `json:"customer_id"`
	Tenor      int       `json:"tenor"`
	Amount     int64     `json:"amount"`
	CreatedAT  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAT  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type LimitInput struct {
	Tenor  int   `json:"tenor"`
	Amount int64 `json:"amount"`
}

type BulkLimitInput struct {
	CustomerID int64         `json:"customer_id"`
	LimitTenor []*LimitInput `json:"limit_tenor"`
}

type LimitResponse struct {
	ID        int64     `json:"id" gorm:"column:id;type:bigserial;primaryKey;autoIncrement"`
	Tenor     int       `json:"tenors"`
	Amount    int64     `json:"amount"`
	CreatedAT time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAT time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type LimitFilter struct {
	ID         int64 `json:"id"`
	CustomerID int64 `json:"customer_id"`
	Tenor      int64 `json:"tenor"`
	Amount     int64 `json:"amount"`
}

type LimitUsecase interface {
	BulkCreateLimit(ctx *fiber.Ctx, data *BulkLimitInput) error
	GetByCustommerID(ctx *fiber.Ctx, CustomerID int64) ([]*Limit, error)
	Update(
		ctx *fiber.Ctx,
		newData *BulkLimitInput,
	) ([]*Limit, error)
}

type LimitRepository interface {
	Get(filter *LimitFilter) ([]*Limit, error)
	GetByCustommerID(CustomerID int64) ([]*Limit, error)
	GetByID(ID int64) (*Limit, error)
	GetByTenor(CustomerID int64, tenor int64) (*Limit, error)
	BulkCreateLimit(data *BulkLimitInput) ([]*int64, error)
	BulkUpdateLimit(data *BulkLimitInput) ([]*int64, error)
}

func (l *Limit) ToLimitResponse() *LimitResponse {
	return &LimitResponse{
		ID:        l.ID,
		Tenor:     l.Tenor,
		Amount:    l.Amount,
		CreatedAT: l.CreatedAT,
		UpdatedAT: l.UpdatedAT,
	}
}

func BuildTenorFactor(salary int64) []*LimitInput {
	tenorFactor := map[int]float64{
		1: 0.4,
		2: 0.6,
		3: 0.75,
		6: 1.0,
	}

	limits := []*LimitInput{}
	for tenor, factor := range tenorFactor {
		limit := &LimitInput{
			Tenor:  tenor,
			Amount: int64(float64(salary) * factor),
		}
		limits = append(limits, limit)
	}
	return limits
}
