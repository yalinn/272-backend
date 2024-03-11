package users

import (
	"272-backend/package/app"
	db "272-backend/package/database"
	jwts "272-backend/package/jwt"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
)

func init() {
	userRoutes := app.App.Group("/users")
	jwts.UseJWT(userRoutes)
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
	sess, err := app.SessionStore.Get(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to get session",
		})
	}
	token := sess.Get("token")
	rolesString, err := db.Redis.Get(token.(string))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to get user",
		})
	}
	var roles []string
	err = json.Unmarshal([]byte(rolesString), &roles)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to parse roles",
			"error":   err.Error(),
		})
	}
	for _, role := range roles {
		log.Println(role)
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
