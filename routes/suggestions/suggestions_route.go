package suggestions

import (
	"272-backend/library"
	"272-backend/pkg"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	route := pkg.App.Group("/suggestions")
	pkg.UseJWT(route)
	route.Post("/", createSuggestion)
	route.Get("/", getSuggestions)
	route.Put("/:id/star", starSuggestion)
	route.Put("/:id/upvote", upvoteSuggestion)
	route.Use(isAuthorized)
}

func isAuthorized(ctx *fiber.Ctx) error {
	user := ctx.Locals("user")
	if user == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not logged in",
		})
	}
	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	for _, role := range claims["roles"].([]interface{}) {
		if role == "admin" {
			return ctx.Next()
		}
	}
	return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "You are not authorized to access this route",
		"error":   "NOT_PERMITTED",
	})
}

func getSuggestions(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not logged in",
		})
	}
	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	userID := claims["username"].(string)
	/* userType := claims["user_type"].(string) */
	suggestions, err := library.GetSuggestions()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get suggestions",
			"error":   err.Error(),
		})
	}

	type suggest struct {
		ID         string   `json:"id"`
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		Author     string   `json:"author"`
		Upvotes    int      `json:"upvotes"`
		Stars      float64  `json:"stars"`
		Date       string   `json:"date"`
		Tags       []string `json:"tags"`
		Status     string   `json:"status"`
		Starred    float64  `json:"starred"`
		Voted      bool     `json:"voted"`
		Department int      `json:"department"`
	}
	var response []suggest
	for _, suggestion := range suggestions {
		starred := 0.00
		voted := false
		/* if userType == "teacher" { */
		for _, stars := range suggestion.Stars {
			if stars.UserID == userID {
				starred = stars.Star
				voted = true
				break
			}
		}
		/* } else { */
		for _, upvote := range suggestion.Upvotes {
			if upvote == userID {
				voted = true
				break
			}
		}
		/* } */
		response = append(response, suggest{
			ID:      suggestion.ID.Hex(),
			Title:   suggestion.Title,
			Content: suggestion.Content,
			Author:  suggestion.AuthorID,
			Upvotes: len(suggestion.Upvotes),
			Stars:   suggestion.CalculateAverageStars(),
			Date:    suggestion.Date,
			Tags:    suggestion.Tags,
			Status:  suggestion.Status,
			Starred: starred,
			Voted:   voted,
		})
	}
	return c.JSON(response)
}

func createSuggestion(c *fiber.Ctx) error {
	user := c.Locals("user")
	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	userID := claims["username"].(string)
	userType := claims["user_type"].(string)
	if userType != "student" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not authorized to create a suggestion",
		})
	}
	var suggestion library.Suggestion
	if err := c.BodyParser(&suggestion); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	suggestion.AuthorID = userID // TODO: change to user.ID
	if suggestion.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Title is required",
		})
	}
	if suggestion.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Content is required",
		})
	}
	if err := suggestion.InsertToDB(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create suggestion",
			"error":   err.Error(),
		})
	}
	return getSuggestions(c)
}

func upvoteSuggestion(c *fiber.Ctx) error {
	user := c.Locals("user")
	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	userID := claims["username"].(string)
	suggestion := library.Suggestion{}
	if suggestionID, err := primitive.ObjectIDFromHex(c.Params("id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid suggestion ID",
		})
	} else {
		suggestion.ID = suggestionID
	}
	if err := suggestion.GiveUpvote(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to upvote suggestion",
			"error":   err.Error(),
		})
	}
	return c.JSON(suggestion)
}

func starSuggestion(c *fiber.Ctx) error {
	type Params struct {
		Star int `json:"star"`
	}
	user := c.Locals("user")
	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	userID := claims["username"].(string)
	/* userType := claims["user_type"].(string)
	if userType == "student" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not authorized to star a suggestion",
		})
	} */
	var params Params
	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	suggestion := library.Suggestion{}
	if suggestionID, err := primitive.ObjectIDFromHex(c.Params("id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid suggestion ID",
		})
	} else {
		suggestion.ID = suggestionID
	}
	if err := suggestion.GiveStar(userID, params.Star); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to star suggestion",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(suggestion)
}

func GetDepartmentID(username string) int {
	chars := strings.Split(username[3:], "")
	id_slice := []string{}
	for i := 0; i < len(chars)-3; i++ {
		id_slice = append(id_slice, chars[i])
	}
	department, _ := strconv.Atoi(strings.Join(id_slice, ""))
	return department
}
