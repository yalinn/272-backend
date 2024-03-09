package users

import (
	"272-backend/package/app"

	"github.com/gofiber/fiber/v2"
)

func init() {
	userRoutes := app.App.Group("/users")
	userRoutes.Get("/", getUsers)
	userRoutes.Get("/:id", getUser)
	// userRoutes.Post("/", register)
}

func getUsers(c *fiber.Ctx) error {
	users, err := GetUsers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to get users",
		})
	}
	return c.JSON(users)
}

func getUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := GetUser(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to get user",
		})
	}
	return c.JSON(user)
}
