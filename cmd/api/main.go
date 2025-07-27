package main

import (
	"log"

	"github.com/Grafiters/archive/configs"
	"github.com/Grafiters/archive/route"

	_ "github.com/Grafiters/archive/cmd/api/docs"
)

// @title Fiber Swagger Example API
// @version 1.0
// @description This is a sample server for a Fiber application.
// @termsOfService http://swagger.io/terms/
// @contact.name Ryudelta Support
// @contact.url https://t.me/Grafiters
// @contact.email ryudelta7@gmail.com
// @Schemes http https
// @Host 127.0.0.1:8080
// @BasePath /
// @securityDefinitions.apikey Token
// @in header
// @name Authorization
// @license.name MIT
// @license.url https://github.com/Grafiters/go-archive/tree/main/license
func main() {
	if err := configs.Initialize(); err != nil {
		log.Fatal(err)
		return
	}
	r := route.SetupRouter()

	r.Listen(":8080")
}
