package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/vviveksharma/lorebound/internal/logger"
	"github.com/vviveksharma/lorebound/internal/queue"
	"github.com/vviveksharma/lorebound/internal/utils"
	"gorm.io/gorm"
)

// Deps are runtime dependencies injected into HTTP handlers.
type Deps struct {
	DB       *gorm.DB
	Producer *queue.Producer
}

// New builds the Fiber app with middleware and routes.
func New(deps Deps) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "lorebound",
	})

	allowOrigins := utils.Getenv("CORS_ALLOW_ORIGINS", "*")

	app.Use(recover.New())
	app.Use(fiberlogger.New(fiberlogger.Config{
		Output: logger.Writer(),
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: allowOrigins != "*",
	}))

	registerRoutes(app, deps)

	return app
}

func registerRoutes(app *fiber.App, deps Deps) {
	_ = deps

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})
}
