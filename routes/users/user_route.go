package users

import (
	"272-backend/library"
	"272-backend/pkg"

	"github.com/gofiber/fiber/v2"
)

func init() {
	userRoutes := pkg.App.Group("/users")
	pkg.UseJWT(userRoutes)
	userRoutes.Get("/", getUsers)
	userRoutes.Get("/:id", getUser)
	// userRoutes.Post("/", register)
}

func getUsers(c *fiber.Ctx) error {
	users, err := library.GetUsers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to get users",
		})
	}

	return c.JSON(users)
}

func getUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := library.GetUser(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to get user",
		})
	}
	return c.JSON(user)
}
