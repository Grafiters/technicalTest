package usecase

import (
	"fmt"

	"github.com/Grafiters/archive/configs"
	"github.com/Grafiters/archive/internal/domain"
	"github.com/Grafiters/archive/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	authRepo domain.AuthRepository
	logger   *configs.LoggerFormat
}

func NewAuthUsecase(
	authRepo domain.AuthRepository,
	logger *configs.LoggerFormat,
) domain.AuthUsecase {
	return &authUsecase{authRepo, logger}
}

// Login implements domain.AuthUsecase.
func (a *authUsecase) Login(ctx *fiber.Ctx, auth *domain.AuthRequest) (*domain.Customer, error) {
	customer, err := a.authRepo.Login(auth)
	if err != nil {
		a.logger.Error("failed to get customer with email, err: %+v", err)
		return nil, fmt.Errorf(utils.DataNotFound)
	}

	err = utils.VerifyPassword(customer.Password, auth.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, fmt.Errorf(utils.PasswordInvalid)
	}

	return customer, nil
}

// Register implements domain.AuthUsecase.
func (a *authUsecase) Register(ctx *fiber.Ctx, data *domain.RegisterRequest) (*domain.Customer, error) {
	_, err := a.authRepo.GetByEmail(data.Email)
	if err == nil {
		a.logger.Error("failed to validate customer email, err: %+v", err)
		return nil, fmt.Errorf("email already in use")
	}

	password, err := utils.Hash(data.Password)
	if err != nil {
		a.logger.Error("failed to hash password, err: %+v", err)
		return nil, fmt.Errorf(utils.ProsessError)
	}

	data.Password = password

	customer, err := a.authRepo.Register(data)
	if err != nil {
		a.logger.Error("failed to register customer, err: %+v", err)
		return nil, fmt.Errorf(utils.ProsessError)
	}

	return customer, nil
}
