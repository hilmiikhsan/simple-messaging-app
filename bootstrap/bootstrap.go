package bootstrap

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/hilmiikhsan/simple-messaging-app/app/repository/user_session"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/database"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/env"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/router"
)

func NewApplication() *fiber.App {
	env.SetupEnvFile()
	db, cfg := database.SetupDatabase()

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Get("/dashboard", monitor.New())

	userSessionRepo := user_session.NewRepository(db)
	middleware := router.NewMiddleware(userSessionRepo)

	apiRouter := router.NewApiRouter(db, cfg, middleware)
	apiRouter.InstallRouter(app)

	return app
}
