package jwts

import (
	"272-backend/config"
	"log"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func UseJWT(route fiber.Router) {
	route.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.JWT_SECRET_KEY)},
		ContextKey: "user",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Println(err.Error())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		},
	}))
}

func CreateToken(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.JWT_SECRET_KEY))
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	return t
}
