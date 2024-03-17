package suggestions

import (
	"272-backend/library"
	"272-backend/package/app"
	jwts "272-backend/package/jwt"
	"context"
	"log"

	db "272-backend/package/database"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	route := app.App.Group("/suggestions")
	jwts.UseJWT(route)
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
	userType := claims["user_type"].(string)
	var suggestions []library.Suggestion
	cursor, err := db.Suggestions.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get suggestions",
			"error":   err.Error(),
		})
	}
	if err := cursor.All(context.TODO(), &suggestions); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get suggestions",
			"error":   err.Error(),
		})
	}
	type suggest struct {
		ID      string   `json:"id"`
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Author  string   `json:"author"`
		Upvotes int      `json:"upvotes"`
		Stars   float64  `json:"stars"`
		Date    string   `json:"date"`
		Tags    []string `json:"tags"`
		Status  string   `json:"status"`
		Starred float64  `json:"starred"`
		Voted   bool     `json:"voted"`
	}
	var response []suggest
	for _, suggestion := range suggestions {
		starred := 0.00
		voted := false
		if userType == "teacher" {
			for _, stars := range suggestion.Stars {
				if stars.UserID == userID {
					starred = stars.Star
					voted = true
					break
				}
			}
		} else {
			for _, upvote := range suggestion.Upvotes {
				if upvote == userID {
					voted = true
					break
				}
			}
		}
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
	sess, err := app.SessionStore.Get(c)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get session",
		})
	}
	token := sess.Get("token")
	if token == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	var user library.User
	if err := user.InitToken(token.(string)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   "Invalid token",
			"context": err.Error(),
		})
	}
	var suggestion library.Suggestion
	if err := c.BodyParser(&suggestion); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	suggestion.AuthorID = user.Username // TODO: change to user.ID
	if suggestion.Title == "" || suggestion.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
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

	return nil
}
