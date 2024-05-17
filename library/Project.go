package library

import (
	db "272-backend/package/database"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Projects *mongo.Collection

func init() {
	Projects = db.Mongo.Collection("projects")
}

type Project struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title      string             `json:"title" bson:"title"`
	Content    string             `json:"content" bson:"content"`
	AuthorID   string             `json:"author" bson:"author"`
	Date       string             `json:"date" bson:"date"`
	Status     string             `json:"status" bson:"status"`
	ApproverID string             `json:"approver" bson:"approver"`
	FormedFrom string             `json:"formedFrom" bson:"formedFrom"`
	Team       []struct {
		UserID   string `json:"userID" bson:"userID"`
		Role     string `json:"role" bson:"role"`
		JoinedAt string `json:"joinedAt" bson:"joinedAt"`
	} `json:"team" bson:"team"`
	AdvisorID string `json:"advisor" bson:"advisor"`
	Comments  []struct {
		UserID  string `json:"userID" bson:"userID"`
		Comment string `json:"comment" bson:"comment"`
		Date    string `json:"date" bson:"date"`
	} `json:"comments" bson:"comments"`
	Stars []struct {
		UserID string  `json:"userID" bson:"userID"`
		Star   float64 `json:"star" bson:"star"`
		Date   string  `json:"date" bson:"date"`
	} `json:"stars" bson:"stars"`
	Tags    []string `json:"tags" bson:"tags"`
	Upvotes []string `json:"upvotes" bson:"upvotes"`
}

func (p *Project) InsertToDB() error {
	project := bson.M{
		"title":      p.Title,
		"content":    p.Content,
		"author":     p.AuthorID,
		"date":       p.Date,
		"status":     p.Status,
		"approver":   p.ApproverID,
		"formedFrom": p.FormedFrom,
		"team": []struct {
			UserID   string `json:"userID" bson:"userID"`
			Role     string `json:"role" bson:"role"`
			JoinedAt string `json:"joinedAt" bson:"joinedAt"`
		}{},
		"advisor": "",
		"comments": []struct {
			UserID  string `json:"userID" bson:"userID"`
			Comment string `json:"comment" bson:"comment"`
			Date    string `json:"date" bson:"date"`
		}{},
		"stars": []struct {
			UserID string  `json:"userID" bson:"userID"`
			Star   float64 `json:"star" bson:"star"`
			Date   string  `json:"date" bson:"date"`
		}{},
		"tags":    []string{},
		"upvotes": []string{},
	}
	if _, err := Projects.InsertOne(context.TODO(), project); err != nil {
		return err
	}
	if err := Projects.FindOne(context.TODO(), bson.M{"title": p.Title}).Decode(&p); err != nil {
		return err
	}
	return nil
}

func (p *Project) GiveUpvote(userID string) error {
	for _, upvote := range p.Upvotes {
		if upvote == userID {
			return nil
		}
	}
	update := bson.M{
		"$addToSet": bson.M{
			"upvotes": userID,
		},
	}
	query := bson.M{
		"_id": p.ID,
	}
	if _, err := Projects.UpdateOne(context.TODO(), query, update); err != nil {
		return err
	}
	if err := Projects.FindOne(context.TODO(), bson.M{"_id": p.ID}).Decode(&p); err != nil {
		return err
	}
	return nil
}

func (p *Project) AddComment(userID, comment string) error {
	commentObj := struct {
		UserID  string `json:"userID" bson:"userID"`
		Comment string `json:"comment" bson:"comment"`
		Date    string `json:"date" bson:"date"`
	}{
		UserID:  userID,
		Comment: comment,
		Date:    p.Date,
	}
	update := bson.M{
		"$addToSet": bson.M{
			"comments": commentObj,
		},
	}
	query := bson.M{
		"_id": p.ID,
	}
	if _, err := Projects.UpdateOne(context.TODO(), query, update); err != nil {
		return err
	}
	if err := Projects.FindOne(context.TODO(), bson.M{"_id": p.ID}).Decode(&p); err != nil {
		return err
	}
	return nil
}

func (p *Project) AddStar(userID string, star float64) error {
	starObj := struct {
		UserID string  `json:"userID" bson:"userID"`
		Star   float64 `json:"star" bson:"star"`
		Date   string  `json:"date" bson:"date"`
	}{
		UserID: userID,
		Star:   star,
		Date:   p.Date,
	}
	update := bson.M{
		"$addToSet": bson.M{
			"stars": starObj,
		},
	}
	query := bson.M{
		"_id": p.ID,
	}
	if _, err := Projects.UpdateOne(context.TODO(), query, update); err != nil {
		return err
	}
	if err := Projects.FindOne(context.TODO(), bson.M{"_id": p.ID}).Decode(&p); err != nil {
		return err
	}
	return nil
}

func (p *Project) GetProject() error {
	query := bson.M{
		"_id": p.ID,
	}
	if err := Projects.FindOne(context.TODO(), query).Decode(&p); err != nil {
		return err
	}
	return nil
}
