package middleware

import (
	"strings"

	"github.com/Grafiters/archive/configs"
	"github.com/Grafiters/archive/configs/response"
	"github.com/Grafiters/archive/internal/domain"
	"github.com/Grafiters/archive/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type Auth struct {
	IDF string `json:"idf"`
	jwt.StandardClaims
}

func Authenticate(c *fiber.Ctx) error {
	var (
		oauthData Auth
		member    *domain.Customer
	)

	token := c.Get("Authorization")
	if len(token) == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Errors{
			Code:   fiber.StatusUnauthorized,
			Errors: []string{utils.AuthorizationBearerMissing},
		})
	}

	signature := c.Get("X-Signature")
	if len(signature) == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Errors{
			Code:   fiber.StatusUnauthorized,
			Errors: []string{utils.SignatureMissing},
		})
	}
	token = strings.Replace(token, configs.Prefix, "", -1)
	err := configs.JwtConfig.DecodeTokenSession(token, oauthData)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"errors": []string{utils.JwtDecodeAndVerify},
		})
	}

	configs.DataBase.Where("uid = ?", oauthData.IDF).First(&member)
	if member == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&response.Errors{
			Code:   fiber.StatusUnauthorized,
			Errors: []string{utils.JwtDecodeAndVerify},
		})
	}
	c.Locals("CurrentUser", member)

	return c.Next()
}
