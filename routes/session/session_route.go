package session

import (
	"272-backend/library"
	"272-backend/package/app"

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
	Email    string `json:"email"`
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
		Email:    form.Email,
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
	if err := session.Save(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to save session",
		})
	}

	return c.JSON(fiber.Map{
		"token": user.Token,
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
