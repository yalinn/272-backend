package events

import (
	"272-backend/library"
	"272-backend/pkg"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostEventParams struct {
	Title     string `json:"title"`
	StartTime string `json:"start_time" bson:"start_time"`
}

func init() {
	router := pkg.App.Group("/events")
	pkg.UseJWT(router)
	router.Get("/", getEvents)
	router.Post("/", postEvent)
	router.Use(isAdmin)
	router.Get("/pending", getPendingEvents)
	router.Patch("/:id", approveEvent)
	router.Delete("/:id", deleteEvent)
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
		if role == "haysev_admin" {
			return ctx.Next()
		}
	}
	return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "You are not authorized to access this route",
		"error":   "NOT_PERMITTED",
	})
}

// getEvents godoc
// @Summary Get events
// @Description Get events
// @Tags events
// @Accept json
// @Produce json
// @Security Bearer
// @Failure 401 {object} library.ErrorPayload
// @Failure 500 {object} library.ErrorPayload
// @Router /events [get]
func getEvents(c *fiber.Ctx) error {
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
	events, err := library.GetAllEvents()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to get events",
		})
	}
	return c.JSON(events)
}

// postEvent godoc
// @Summary Post event
// @Description Post event
// @Tags events
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body PostEventParams true "Body"
// @Failure 400 {object} library.ErrorPayload
// @Failure 401 {object} library.ErrorPayload
// @Failure 500 {object} library.ErrorPayload
// @Router /events [post]
func postEvent(c *fiber.Ctx) error {
	var params PostEventParams
	if err := c.BodyParser(&params); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
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
	startTime, err := time.Parse(time.RFC3339, params.StartTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid start_time format",
		})
	}
	event := library.Event{
		Title:       "Köpek Besleme",
		StartTime:   primitive.NewDateTimeFromTime(startTime),
		Type:        "haysev",
		OrganizerID: userID,
	}
	if err := event.CreateEvent(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to create event",
		})
	}
	return c.JSON(event)
}

// getPendingEvents godoc
// @Summary Get pending events
// @Description Get pending events
// @Tags events
// @Accept json
// @Produce json
// @Security Bearer
// @Failure 401 {object} library.ErrorPayload
// @Failure 500 {object} library.ErrorPayload
// @Router /events/pending [get]
func getPendingEvents(c *fiber.Ctx) error {
	events, err := library.GetPendingEvents()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to get events",
		})
	}
	return c.JSON(events)
}

// approveEvent godoc
// @Summary Approve event
// @Description Approve event
// @Tags events
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Event ID"
// @Failure 400 {object} library.ErrorPayload
// @Failure 500 {object} library.ErrorPayload
// @Router /events/{id} [patch]
func approveEvent(c *fiber.Ctx) error {
	event := library.Event{}
	if eventID, err := primitive.ObjectIDFromHex(c.Params("id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid suggestion ID",
		})
	} else {
		event.ID = eventID
	}
	if err := event.ApproveEvent(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to approve event",
		})
	}
	return c.JSON(event)
}

// deleteEvent godoc
// @Summary Delete event
// @Description Delete event
// @Tags events
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Event ID"
// @Failure 400 {object} library.ErrorPayload
// @Failure 500 {object} library.ErrorPayload
// @Router /events/{id} [delete]
func deleteEvent(c *fiber.Ctx) error {
	event := library.Event{}
	if eventID, err := primitive.ObjectIDFromHex(c.Params("id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid suggestion ID",
		})
	} else {
		event.ID = eventID
	}
	if err := event.RemoveEvent(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to reject event",
		})
	}
	return c.JSON(event)
}
