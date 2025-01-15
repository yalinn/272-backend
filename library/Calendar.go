package library

import (
	"272-backend/pkg"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Events *mongo.Collection

func init() {
	Events = pkg.Mongo.Collection("events")
}

type Event struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	StartTime   primitive.DateTime `json:"start_time" bson:"start_time"`
	EndTime     primitive.DateTime `json:"end_time" bson:"end_time"`
	Location    string             `json:"location" bson:"location"`
	OrganizerID string             `json:"organizer_id" bson:"organizer_id"`
	Author      string             `json:"author" bson:"author"`
	CreatedAt   primitive.DateTime `json:"created_at" bson:"created_at"`
	Tags        []string           `json:"tags" bson:"tags"`
	Status      string             `json:"status" bson:"status"`
	Type        string             `json:"type" bson:"type"`
}

func (e *Event) CreateEvent() error {
	if e.Title == "" {
		return errors.New("INVALID_EVENT")
	}
	if e.OrganizerID == "" {
		return errors.New("INVALID_ORGANIZER")
	}
	if e.Author == "" {
		return errors.New("INVALID_ORGANIZER_NAME")
	}
	e.Status = "pending"
	date, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}
	e.CreatedAt = primitive.NewDateTimeFromTime(date)
	if e.StartTime.Time().IsZero() {
		return errors.New("INVALID_START_TIME")
	}
	if e.Type == "" {
		return errors.New("INVALID_EVENT_TYPE")
	}
	doc, err := Events.InsertOne(context.Background(), e)
	if err != nil {
		return err
	}
	e.ID = doc.InsertedID.(primitive.ObjectID)
	if err := Events.FindOne(context.Background(), bson.D{{Key: "_id", Value: e.ID}}).Decode(&e); err != nil {
		return err
	}
	return nil
}

func (u *User) GetEvents() ([]Event, error) {
	cursor, err := Events.Find(context.Background(), bson.D{{Key: "organizer_id", Value: u.Username}})
	if err != nil {
		return nil, err
	}
	var events []Event
	if err := cursor.All(context.Background(), &events); err != nil {
		return nil, err
	}
	return events, nil
}

func (e *Event) ApproveEvent() error {
	_id, err := Events.UpdateOne(context.Background(), bson.D{{Key: "_id", Value: e.ID}}, bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: "approved"}}}})
	if err != nil {
		return err
	} else if _id.MatchedCount == 0 {
		return errors.New("EVENT_NOT_FOUND")
	}
	if err := Events.FindOne(context.Background(), bson.D{{Key: "_id", Value: e.ID}}).Decode(&e); err != nil {
		return err
	}
	return nil
}

func (e *Event) RemoveEvent() error {
	if err := Events.FindOne(context.Background(), bson.D{{Key: "_id", Value: e.ID}}).Decode(&e); err != nil {
		return err
	}
	_, err := Events.DeleteOne(context.Background(), bson.D{{Key: "_id", Value: e.ID}})
	return err
}

func GetAllEvents() ([]Event, error) {
	cursor, err := Events.Find(context.Background(), bson.D{{Key: "status", Value: "approved"}})
	if err != nil {
		return nil, err
	}
	var events []Event
	if err := cursor.All(context.Background(), &events); err != nil {
		return nil, err
	}
	return events, nil
}

func GetPendingEvents() ([]Event, error) {
	cursor, err := Events.Find(context.Background(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	var events []Event
	if err := cursor.All(context.Background(), &events); err != nil {
		return nil, err
	}
	return events, nil
}
