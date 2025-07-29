package domain

import (
	"github.com/gofiber/fiber/v2"
)

type AuthRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (r *TokenResponse) ToTokenResponse() TokenResponse {
	return TokenResponse{
		AccessToken: r.AccessToken,
	}
}

type AuthUsecase interface {
	Login(ctx *fiber.Ctx, auth *AuthRequest) (*Customer, error)
	Register(ctx *fiber.Ctx, data *RegisterRequest) (*Customer, error)
}

type AuthRepository interface {
	Login(auth *AuthRequest) (*Customer, error)
	Register(data *RegisterRequest) (*Customer, error)
	GetByEmail(email string) (*Customer, error)
}
