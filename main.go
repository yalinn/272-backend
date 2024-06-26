package main

import (
	"github.com/gofiber/fiber/v2/log"

	"272-backend/config"
	_ "272-backend/docs"
	"272-backend/pkg"
	_ "272-backend/routes"
)

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host api-probee.yalin.app
// @BasePath /
func main() {
	if err := pkg.App.Listen(config.PORT); err != nil {
		log.Fatal("Oops... Server is not running! Reason: %v", err)
	} else {
		log.Info("Server is running on port 3000...")
	}
}
