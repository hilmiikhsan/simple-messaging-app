package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	userController "github.com/hilmiikhsan/simple-messaging-app/app/controllers/user"
	userRepository "github.com/hilmiikhsan/simple-messaging-app/app/repository/user"
	userSessionRepository "github.com/hilmiikhsan/simple-messaging-app/app/repository/user_session"
	userService "github.com/hilmiikhsan/simple-messaging-app/app/service/user"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/database"
	"gorm.io/gorm"
)

type ApiRouter struct {
	db         *gorm.DB
	cfg        *database.Config
	middleware *Middleware
}

func (h ApiRouter) InstallRouter(app *fiber.App) {
	api := app.Group("/api", limiter.New())
	api.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Hello from api",
		})
	})

	userRepo := userRepository.NewRepository(h.db)
	userSessionRepo := userSessionRepository.NewRepository(h.db)

	userService := userService.NewService(h.cfg, userRepo, userSessionRepo)

	userController := userController.NewController(app, userService)

	userGroup := app.Group("/user")
	userV1Group := userGroup.Group("/v1")

	userV1Group.Post("/register", userController.Register)
	userV1Group.Post("/login", userController.Login)
	userV1Group.Delete("/logout", h.middleware.MiddlewareValidateAuth, userController.Logout)
	userV1Group.Put("/refresh-token", h.middleware.MiddlewareRefreshToken, userController.RefreshToken)
}

func NewApiRouter(db *gorm.DB, cfg *database.Config, middleware *Middleware) *ApiRouter {
	return &ApiRouter{
		db:         db,
		cfg:        cfg,
		middleware: middleware,
	}
}