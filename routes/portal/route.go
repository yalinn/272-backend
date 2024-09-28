package portal

import (
	"272-backend/config"
	"272-backend/pkg"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func init() {
	router := pkg.App.Group("/portal")
	pkg.UseJWT(router)
	router.Get("/curriculum", getCurriculum)
	router.Post("/curriculum", postCurriculum)
}

// getCurriculum godoc
// @Summary Get curriculum
// @Description Get curriculum
// @Tags portal
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} GetResponse
// @Router /portal/curriculum [get]
func getCurriculum(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"message": "Curriculum",
	})
}

type postBody struct {
	Password string `json:"password"`
}

// postCurriculum godoc
// @Summary Post curriculum
// @Description Post curriculum
// @Tags portal
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body postBody true "Body"
// @Success 200 {object} PostResponse
// @Router /portal/curriculum [post]
func postCurriculum(c *fiber.Ctx) error {
	data := new(postBody)
	if err := c.BodyParser(data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	user := c.Locals("user")
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not logged in",
		})
	}
	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	userID := claims["username"].(string)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Authentication token is invalid or expired",
		})
	}
	pwd := data.Password
	postData := map[string]string{
		"username": userID,
		"password": pwd,
	}
	jsonData, err := json.Marshal(postData)
	if err != nil {
		return err
	}
	response, err := http.Post(config.BSL_URI+"/curriculum", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var info []CurriculumObject
	err = json.Unmarshal(body, &info)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(info)
}
