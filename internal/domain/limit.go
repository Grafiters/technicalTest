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
	CustomerID int64 `json:"customer_id"`
	Tenor      int   `json:"tenor"`
	Amount     int64 `json:"amount"`
}

type LimitResponse struct {
	ID        int64     `json:"id" gorm:"column:id;type:bigserial;primaryKey;autoIncrement"`
	Tenor     int       `json:"tenors"`
	Amount    int64     `json:"amount"`
	CreatedAT time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAT time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type LimitUsecase interface {
	Create(ctx fiber.Ctx, data *LimitInput) error
	GetByID(ctx fiber.Ctx, ID int64) (*Limit, error)
	GetByCustommerID(ctx fiber.Ctx, CustomerID int64) ([]*Limit, error)
	Update(
		ctx fiber.Ctx,
		CustomerID int64,
		oldData *Limit,
		newData *Limit,
	) (*Limit, error)
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
