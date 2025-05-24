package routes

import (
	"project_restfulApi_go/controllers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserRoutes(app *fiber.App, db *gorm.DB) {
	user := app.Group("/users")
	user.Get("/:id", func(c *fiber.Ctx) error {
		return controllers.GetUserHandler(c, db)
	})
	user.Get("/", func(c *fiber.Ctx) error {
		return controllers.GetUsersHandler(c, db)
	})
	user.Post("/", func(c *fiber.Ctx) error {
		return controllers.CreateUserHandler(c, db)
	})
	user.Patch("/:id", func(c *fiber.Ctx) error {
		return controllers.UpdateUserHandler(c, db)
	})
	user.Delete("/:id", func(c *fiber.Ctx) error {
		return controllers.DeleteUserHandler(c, db)
	})
}
