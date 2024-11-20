package bootstrap

import (
	"io"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/hilmiikhsan/simple-messaging-app/app/controllers"
	"github.com/hilmiikhsan/simple-messaging-app/app/repository/message"
	"github.com/hilmiikhsan/simple-messaging-app/app/repository/user_session"
	"github.com/hilmiikhsan/simple-messaging-app/app/ws"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/database"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/env"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/router"
	"go.elastic.co/apm"
)

func NewApplication() *fiber.App {
	env.SetupEnvFile()
	SetupLogfile()

	db, cfg := database.SetupDatabase()

	database.SetupMongoDB()

	apm.DefaultTracer.Service.Name = "simple-mesaging-app"

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Get("/dashboard", monitor.New())
	app.Get("/", controllers.RenderUI)

	userSessionRepo := user_session.NewRepository(db)
	messageRepo := message.NewRepository(db)

	ws := ws.NewService(cfg, messageRepo)

	go ws.ServeWSMessaging(app)

	middleware := router.NewMiddleware(userSessionRepo)

	apiRouter := router.NewApiRouter(db, cfg, middleware)
	apiRouter.InstallRouter(app)

	return app
}

func SetupLogfile() {
	logFile, err := os.OpenFile("./logs/simple_messaging_app.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}
