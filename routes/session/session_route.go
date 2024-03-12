package session

import (
	"272-backend/library"
	"272-backend/package/app"
	db "272-backend/package/database"

	"github.com/gofiber/fiber/v2"
)

func init() {
	sessionRoutes := app.App.Group("/session")
	sessionRoutes.Get("/", getSession)
	sessionRoutes.Post("/", login)
	sessionRoutes.Delete("/", logout)
}

func getSession(c *fiber.Ctx) error {
	sess, err := app.SessionStore.Get(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to get session",
		})
	}
	token := sess.Get("token")
	if token == nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	return c.JSON(fiber.Map{
		"token": token,
	})
}

type loginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserType string `json:"user_type"`
}

func login(c *fiber.Ctx) error {
	var form loginForm
	if err := c.BodyParser(&form); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	user := library.User{
		Username: form.Username,
		UserType: form.UserType,
	}

	if err := user.LoginByEmail(form.Password); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	session, err := app.SessionStore.Get(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to get session",
		})
	}
	session.Set("token", user.Token)
	if err := db.Redis.Set(user.Token, user.Stringify()); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to save token to redis",
		})
	}
	if err := session.Save(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to save session",
		})
	}

	return c.JSON(fiber.Map{
		"token": user.Token,
		"user":  user.Stringify(),
	})
}

func logout(c *fiber.Ctx) error {
	session, err := app.SessionStore.Get(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to get session",
		})
	}
	session.Delete("token")
	if err := session.Save(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to save session",
		})
	}
	return c.SendStatus(204)
}
