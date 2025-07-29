package delivery

import (
	"github.com/Grafiters/archive/configs"
	"github.com/Grafiters/archive/internal/domain"
	"github.com/Grafiters/archive/utils"
	"github.com/gofiber/fiber/v2"
)

type authHandler struct {
	authUsecase domain.AuthUsecase
	jwtConfigs  *configs.JwtService
	logger      *configs.LoggerFormat
}

func NewAuthHandler(
	router fiber.Router,
	jwtConfigs *configs.JwtService,
	au domain.AuthUsecase,
	logger *configs.LoggerFormat,
) {
	h := &authHandler{
		authUsecase: au,
		jwtConfigs:  jwtConfigs,
		logger:      logger,
	}

	router.Post("/auth", h.Login)
	router.Post("/register", h.Register)
}

// Auth
// @Router /api/auth [post]
// @Summary Session
// @Description Generate access token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param auth body domain.AuthRequest true "generate access token"
// @Success 200 {object} domain.SingleResponse{data=domain.TokenResponse}
// @Failure 422 {object} domain.SingleResponse
// @Failure 500 {object} domain.SingleResponse
func (ah *authHandler) Login(c *fiber.Ctx) error {
	payload := new(domain.AuthRequest)

	if len(c.Body()) <= 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: utils.BodyParamMissing,
		})
	}

	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: utils.InvalidMessageBody,
		})
	}

	customer, err := ah.authUsecase.Login(c, payload)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: err.Error(),
		})
	}

	accessToken, err := ah.jwtConfigs.GenerateTokenSession(customer.Email)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: utils.ProsessError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.SingleResponse{
		Code: fiber.StatusCreated,
		Data: domain.TokenResponse{
			AccessToken: accessToken,
		},
		Message: "login success",
	})
}

// Auth
// @Router /api/register [post]
// @Summary Register
// @Description Register to system
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param auth body domain.RegisterRequest true "register to system"
// @Success 200 {object} domain.SingleResponse{data=domain.CustomerResponse}
// @Failure 422 {object} domain.SingleResponse
// @Failure 500 {object} domain.SingleResponse
func (ah *authHandler) Register(c *fiber.Ctx) error {
	payload := new(domain.RegisterRequest)
	if len(c.Body()) <= 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: utils.BodyParamMissing,
		})
	}

	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: utils.InvalidMessageBody,
		})
	}

	customer, err := ah.authUsecase.Register(c, payload)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.SingleResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Data:    nil,
			Message: err.Error(),
		})
	}

	customerResponse := customer.ToCustomerResponse()
	return c.Status(fiber.StatusOK).JSON(domain.SingleResponse{
		Code:    fiber.StatusUnprocessableEntity,
		Data:    customerResponse,
		Message: "success",
	})
}
