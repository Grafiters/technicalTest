package domain

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type Customer struct {
	ID             int64     `json:"id" gorm:"column:id;type:bigserial;primaryKey;autoIncrement"`
	Email          string    `json:"email" gorm:"unique" validate:"required"`
	Password       string    `json:"password"`
	NIK            int       `json:"bik"`
	FullName       string    `json:"full_name"`
	LegalName      string    `json:"legal_name"`
	BirthPlace     string    `json:"birth_place"`
	BirthDate      time.Time `json:"birth_date"`
	Salary         int64     `json:"salary"`
	KTPImageUrl    string    `json:"ktp_image_url"`
	SelfieImageUrl string    `json:"selfie_image_url"`
	CreatedAT      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAT      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CustomerInput struct {
	Email       string    `json:"email" validate:"required" example:"test@gmail.com"`
	NIK         int       `json:"nik" validate:"required,len=16,numeric"`
	FullName    string    `json:"full_name" validate:"required"`
	LegalName   string    `json:"legal_name" validate:"required"`
	BirthPlace  string    `json:"birth_place" validate:"required"`
	BirthDate   time.Time `json:"birth_date" validate:"required,datetime=2006-01-02"`
	Salary      int64     `json:"salary" validate:"required, gt=0"`
	KTPImage    string    `json:"ktp_image" validate:"required"`
	SelfieImage string    `json:"selfie_image" validate:"required"`
}

type CustomerUpdate struct {
	Email       string    `json:"email" validate:"required" example:"test@gmail.com"`
	FullName    string    `json:"full_name" validate:"required"`
	LegalName   string    `json:"legal_name" validate:"required"`
	BirthPlace  string    `json:"birth_place" validate:"required"`
	BirthDate   time.Time `json:"birth_date" validate:"required,datetime=2006-01-02"`
	KTPImage    string    `json:"ktp_image" validate:"required"`
	SelfieImage string    `json:"selfie_image" validate:"required"`
}

type CustomerUpdateSalary struct {
	Salary int `json:"salary" validate:"required, gt=0"`
}

type CustomerResponse struct {
	ID             int64     `json:"id"`
	Email          string    `json:"email"`
	NIK            int       `json:"bik"`
	FullName       string    `json:"full_name"`
	LegalName      string    `json:"legal_name"`
	BirthPlace     string    `json:"birth_place"`
	BirthDate      time.Time `json:"birth_date"`
	Salary         int64     `json:"salary"`
	KTPImageUrl    string    `json:"ktp_image_url"`
	SelfieImageUrl string    `json:"selfie_image_url"`
	CreatedAT      time.Time `json:"created_at"`

	Limit []*LimitResponse `json:"limit,omitempty"`
}

type CustomerFilter struct {
	Page     int `json:"page" query:"page" validate:"omitempty"`
	PageSize int `json:"page_size" query:"pageSize" validate:"omitempty,oneof=10 20 50 100"`
	Sort     []SortObject
}

func (c *Customer) ToCustomerResponse() *CustomerResponse {
	return &CustomerResponse{
		ID:             c.ID,
		Email:          c.Email,
		FullName:       c.FullName,
		LegalName:      c.LegalName,
		BirthPlace:     c.BirthPlace,
		BirthDate:      c.BirthDate,
		Salary:         c.Salary,
		KTPImageUrl:    c.KTPImageUrl,
		SelfieImageUrl: c.SelfieImageUrl,
		CreatedAT:      c.CreatedAT,
	}
}

type CustomerUsecase interface {
	Create(ctx *fiber.Ctx, data *CustomerInput) (*CustomerResponse, error)
	Get(ctx *fiber.Ctx, ID int64) (*CustomerResponse, error)
	Update(ctx *fiber.Ctx, ID int64, data *CustomerUpdate) (*CustomerResponse, error)
	UpdateSalary(ctx *fiber.Ctx, ID int64, salary *CustomerUpdateSalary) (*CustomerResponse, error)
}

type CustomerRepository interface {
	Create(data *CustomerInput) (*Customer, error)
	Get(filter *CustomerFilter) ([]*Customer, int64, error)
	GetByID(ID int64) (*Customer, error)
	Update(
		ID int64,
		newData *CustomerUpdate,
	) (*Customer, error)
	UpdateSalary(ID int64, salary CustomerUpdateSalary) (*Customer, error)
}
