package route

import (
	"html/template"

	"github.com/Grafiters/archive/configs"
	authHttp "github.com/Grafiters/archive/internal/auth/delivery"
	authMysql "github.com/Grafiters/archive/internal/auth/repository"
	authUsecase "github.com/Grafiters/archive/internal/auth/usecase"
	customerHttp "github.com/Grafiters/archive/internal/customer/delivery"
	customerMysql "github.com/Grafiters/archive/internal/customer/repository"
	customerUsecase "github.com/Grafiters/archive/internal/customer/usecase"
	limitMysql "github.com/Grafiters/archive/internal/limit/repository"
	transactionHttp "github.com/Grafiters/archive/internal/transaction/delivery"
	transactionMysql "github.com/Grafiters/archive/internal/transaction/repository"
	transactionUsecase "github.com/Grafiters/archive/internal/transaction/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

func SetupRouter() *fiber.App {
	app := fiber.New()
	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Get("/api/openapi/*", swagger.New(swagger.Config{
		Title:  "Skill Test Pertama - Bayu Grafit Nur Alfian",
		Layout: "BaseLayout",
		Plugins: []template.JS{
			template.JS(`SwaggerUIBundle.plugins.DownloadUrl`),
		},
		CustomStyle: template.CSS(`
			@import url('https://cdn.jsdelivr.net/npm/swagger-themes@1.4.3/themes/dark.css');
		`),
	}))

	app.Options("/*", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	api := app.Group("/api")

	limitRepo := limitMysql.NewLimitRepository(configs.DataBase, configs.Logger)

	authRepo := authMysql.NewAuthRepository(configs.DataBase, configs.Logger)
	authUsecase := authUsecase.NewAuthUsecase(authRepo, limitRepo, configs.Logger)
	authHttp.NewAuthHandler(api, configs.JwtConfig, authUsecase, configs.Logger)

	customerRepo := customerMysql.NewCustomerRepository(configs.DataBase, configs.Logger)
	customerUsecase := customerUsecase.NewCustomerUsecase(customerRepo, limitRepo, configs.Logger)
	customerHttp.NewCustomerHandler(api, customerUsecase, configs.Logger)

	transactionRepo := transactionMysql.NewTranscationRepository(configs.DataBase, configs.Logger)
	transactionUsecase := transactionUsecase.NewTransactionUsecase(transactionRepo, customerRepo, limitRepo, configs.Logger)
	transactionHttp.NewTransactionHandler(api, transactionUsecase, configs.Logger)

	return app
}
