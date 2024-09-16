package session

import (
	"272-backend/library"
	"272-backend/pkg"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func init() {
	sessionRoutes := pkg.App.Group("/session")
	sessionRoutes.Post("/", login)
	pkg.UseJWT(sessionRoutes)
	sessionRoutes.Get("/", getSession)
	sessionRoutes.Delete("/", logout)
}

// getSession godoc
// @Summary Get user session
// @Description Get user session
// @Tags session
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} library.User
// @Failure 401 {object} library.ErrorPayload
// @Failure 500 {object} library.ErrorPayload
// @Router /session [get]
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
	return c.Status(200).JSON(fiber.Map{
		"user": fiber.Map{
			"username":   user.Username,
			"user_type":  user.UserType,
			"roles":      user.Roles,
			"department": user.GetDepartmentID(),
		},
	})
}

type loginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserType string `json:"user_type"`
}

// login godoc
// @Summary Login
// @Description Login
// @Tags session
// @Accept json
// @Produce json
// @Param body body loginForm true "Login Form"
// @Success 200 {object} library.User
// @Failure 400 {object} library.ErrorPayload
// @Failure 401 {object} library.ErrorPayload
// @Failure 500 {object} library.ErrorPayload
// @Router /session [post]
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

	token := pkg.CreateToken(jwt.MapClaims{
		"username":  user.Username,
		"user_type": user.UserType,
		"roles":     user.Roles,
	})

	if err := pkg.Redis.Set(token, user.Stringify()); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to save token to redis",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": fiber.Map{
			"id":         user.Username,
			"user_type":  user.UserType,
			"department": user.GetDepartmentID(),
			"full_name":  user.FullName,
			"roles":      user.Roles,
		},
		"token": token,
	})
}

// logout godoc
// @Summary Logout
// @Description Logout
// @Tags session
// @Accept json
// @Produce json
// @Security Bearer
// @Success 204
// @Failure 401 {object} library.ErrorPayload
// @Failure 500 {object} library.ErrorPayload
// @Router /session [delete]
func logout(c *fiber.Ctx) error {
	auth := c.Locals("user")
	if auth == nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "You are not logged in",
		})
	}
	claims := auth.(*jwt.Token).Raw
	log.Println(claims)
	if err := pkg.Redis.Del(claims); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to delete token from redis",
			"error":   err.Error(),
		})
	}
	return c.SendStatus(204)
}
