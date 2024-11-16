package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/database"
	"gorm.io/gorm"
)

func InstallRouter(app *fiber.App, db *gorm.DB, cfg *database.Config, middleware *Middleware) {
	setup(app, NewApiRouter(db, cfg, middleware), NewHttpRouter())
}
func setup(app *fiber.App, router ...Router) {
	for _, r := range router {
		r.InstallRouter(app)
	}
}
