package suggestions

import (
	"272-backend/library"
	"272-backend/pkg"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	route := pkg.App.Group("/suggestions")
	pkg.UseJWT(route)
	route.Get("/", getApprovedSuggestions)
	route.Get("/:id", getSuggestion)
	route.Post("/", createSuggestion)
	route.Put("/:id/upvote", upvoteSuggestion)
	route.Get("/rejected", getRejectedSuggestions)
	route.Use(isAdmin)
	route.Get("/pending", getPendingSuggestions)
	route.Get("/reported", getReportedSuggestions)
	route.Put("/:id/star", starSuggestion)
	route.Patch("/:id/approve", approveSuggestion)
	route.Patch("/:id/reject", rejectSuggestion)
	route.Patch("/:id/report", reportSuggestion)
}

func isAdmin(ctx *fiber.Ctx) error {
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

// getApprovedSuggestions godoc
// @Summary Get Approved Suggestions
// @Description Get Approved suggestions
// @Tags suggestions
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} library.SuggestionResponse
// @Failure 401 {object} library.ErrorPayload
// @Failure 500 {object} library.ErrorPayload
// @Router /suggestions [get]
func getApprovedSuggestions(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not logged in",
		})
	}
	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	userID := claims["username"].(string)
	/* userType := claims["user_type"].(string) */
	suggestions, err := library.GetApprovedSuggestions()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get suggestions",
			"error":   err.Error(),
		})
	}
	var response []library.SuggestionResponse
	for _, suggestion := range suggestions {
		response = append(response, suggestion.ToResponse(userID))
	}
	return c.JSON(response)
}

// getPendingSuggestions godoc
// @Summary Get Pending Suggestions
// @Description Get all pending suggestions
// @Tags suggestions
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} library.SuggestionResponse
// @Failure 401 {object} library.ErrorPayload
// @Failure 500 {object} library.ErrorPayload
// @Router /suggestions/pending [get]
func getPendingSuggestions(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not logged in",
		})
	}
	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	userID := claims["username"].(string)
	suggestions, err := library.GetPendingSuggestions()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get suggestions",
			"error":   err.Error(),
		})
	}
	var response []library.SuggestionResponse
	for _, suggestion := range suggestions {
		response = append(response, suggestion.ToResponse(userID))
	}
	return c.JSON(response)
}

// getReportedSuggestions godoc
// @Summary Get Reported Suggestions
// @Description Get all reported suggestions
// @Tags suggestions
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} library.SuggestionResponse
// @Failure 401 {object} library.ErrorPayload
// @Failure 500 {object} library.ErrorPayload
// @Router /suggestions/reported [get]
func getReportedSuggestions(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not logged in",
		})
	}
	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	userID := claims["username"].(string)
	suggestions, err := library.GetReportedSuggestions()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get suggestions",
			"error":   err.Error(),
		})
	}
	var response []library.SuggestionResponse
	for _, suggestion := range suggestions {
		response = append(response, suggestion.ToResponse(userID))
	}
	return c.JSON(response)
}

// getSuggestion godoc
// @Summary Get Suggestion
// @Description Get a suggestion
// @Tags suggestions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Suggestion ID"
// @Success 200 {object} library.SuggestionResponse
// @Failure 401 {object} library.ErrorPayload
// @Failure 500 {object} library.ErrorPayload
// @Router /suggestions/{id} [get]
func getSuggestion(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not logged in",
		})
	}
	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	userID := claims["username"].(string)
	suggestion := library.Suggestion{}
	if c.Params("id") == "pending" || c.Params("id") == "reported" || c.Params("id") == "rejected" {
		return c.Next()
	}
	if err := suggestion.WithID(c.Params("id")); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get suggestion",
			"error":   err.Error(),
		})
	}
	response := suggestion.ToResponse(userID)
	return c.JSON(response)
}

// createSuggestion godoc
// @Summary Create Suggestion
// @Description Create a suggestion
// @Tags suggestions
// @Accept json
// @Produce json
// @Security Bearer
// @Param suggestion body library.CreateSuggestionParams true "Suggestion"
// @Success 200 {object} library.SuggestionResponse
// @Router /suggestions [post]
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
	return c.JSON(suggestion.ToResponse(userID))
}

// getRejectedSuggestions godoc
// @Summary Get Rejected Suggestions
// @Description Get all rejected suggestions
// @Tags suggestions
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} library.SuggestionResponse
// @Failure 401 {object} library.ErrorPayload
// @Failure 500 {object} library.ErrorPayload
// @Router /suggestions/rejected [get]
func getRejectedSuggestions(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not logged in",
		})
	}
	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	userID := claims["username"].(string)
	suggestions, err := library.GetRejectedSuggestions()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get suggestions",
			"error":   err.Error(),
		})
	}
	var response []library.SuggestionResponse
	for _, suggestion := range suggestions {
		response = append(response, suggestion.ToResponse(userID))
	}
	return c.JSON(response)
}

// upvoteSuggestion godoc
// @Summary Upvote Suggestion
// @Description Upvote a suggestion
// @Tags suggestions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Suggestion ID"
// @Success 200 {object} library.SuggestionResponse
// @Failure 400 {object} library.ErrorPayload
// @Failure 500 {object} library.ErrorPayload
// @Router /suggestions/{id}/upvote [put]
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
	return c.JSON(suggestion.ToResponse(userID))
}

// starSuggestion godoc
// @Summary Star Suggestion
// @Description Star a suggestion
// @Tags suggestions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Suggestion ID"
// @Param star body library.StarSuggestionParams true "Star"
// @Success 200 {object} library.SuggestionResponse
// @Failure 400 {object} library.ErrorPayload
// @Failure 500 {object} library.ErrorPayload
// @Router /suggestions/{id}/star [put]
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

	return c.Status(fiber.StatusOK).JSON(suggestion.ToResponse(userID))
}

// approveSuggestion godoc
// @Summary Approve Suggestion
// @Description Approve a suggestion
// @Tags suggestions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Suggestion ID"
// @Success 200 {object} library.SuggestionResponse
// @Failure 400 {object} library.ErrorPayload
// @Failure 500 {object} library.ErrorPayload
// @Router /suggestions/{id}/approve [patch]
func approveSuggestion(c *fiber.Ctx) error {
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
	if err := suggestion.Approve(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to approve suggestion",
			"error":   err.Error(),
		})
	}
	return c.JSON(suggestion.ToResponse(userID))
}

// rejectSuggestion godoc
// @Summary Reject Suggestion
// @Description Reject a suggestion
// @Tags suggestions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Suggestion ID"
// @Param reason body library.WithReasonParams true "Reason"
// @Success 200 {object} library.SuggestionResponse
// @Failure 400 {object} library.ErrorPayload
// @Failure 500 {object} library.ErrorPayload
// @Router /suggestions/{id}/reject [patch]
func rejectSuggestion(c *fiber.Ctx) error {
	user := c.Locals("user")
	claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	userID := claims["username"].(string)
	suggestion := library.Suggestion{}
	var params struct {
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	if params.Reason == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Reason is required",
		})
	}
	if suggestionID, err := primitive.ObjectIDFromHex(c.Params("id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid suggestion ID",
		})
	} else {
		suggestion.ID = suggestionID
	}
	if err := suggestion.Reject(userID, params.Reason); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to reject suggestion",
			"error":   err.Error(),
		})
	}
	return c.JSON(suggestion.ToResponse(userID))
}

// reportSuggestion godoc
// @Summary Report Suggestion
// @Description Report a suggestion
// @Tags suggestions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Suggestion ID"
// @Success 200 {object} library.SuggestionResponse
// @Failure 400 {object} library.ErrorPayload
// @Failure 500 {object} library.ErrorPayload
// @Router /suggestions/{id}/report [patch]
func reportSuggestion(c *fiber.Ctx) error {
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
	if err := suggestion.Report(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to report suggestion",
			"error":   err.Error(),
		})
	}
	return c.JSON(suggestion.ToResponse(userID))
}
