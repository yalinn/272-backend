package session

import (
	"272-backend/library"
	"272-backend/package/app"
	db "272-backend/package/database"
	jwts "272-backend/package/jwt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func init() {
	sessionRoutes := app.App.Group("/session")
	sessionRoutes.Post("/", login)
	jwts.UseJWT(sessionRoutes)
	sessionRoutes.Get("/", getSession)
	sessionRoutes.Delete("/", logout)
}

func getSession(c *fiber.Ctx) error {
	auth := c.Locals("user")
	if auth == nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "You are not logged in",
		})
	}
	claims := auth.(*jwt.Token).Claims.(jwt.MapClaims)
	if claims["username"] == nil || claims["user_type"] == nil || claims["roles"] == nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Authentication token is invalid or expired",
		})
	}
	user := library.User{
		Username: claims["username"].(string),
		UserType: claims["user_type"].(string),
	}
	if err := user.FindUser(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to get user",
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"username":  user.Username,
			"user_type": user.UserType,
			"roles":     user.Roles,
		},
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
		Username: strings.Split(form.Username, "@")[0],
		UserType: form.UserType,
	}

	if err := user.LoginByEmail(form.Password); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   err.Error(),
		})
	}
	/* claims := &jwt.MapClaims{
		"username":  user.Username,
		"user_type": user.UserType,
		"roles":     user.Roles,
	} */

	token := jwts.CreateToken(jwt.MapClaims{
		"username":  user.Username,
		"user_type": user.UserType,
		"roles":     user.Roles,
	})
	session, err := app.SessionStore.Get(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to get session",
		})
	}
	session.Set("token", token)
	if err := db.Redis.Set(token, user.Stringify()); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to save token to redis",
			"error":   err.Error(),
		})
	}

	if err := session.Save(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to save session",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": fiber.Map{
			"username":  user.Username,
			"user_type": user.UserType,
			"roles":     user.Roles,
		},
		"token": token,
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
