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

var Projects *mongo.Collection

func init() {
	Projects = pkg.Mongo.Collection("projects")
}

type Project struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title    string             `json:"title" bson:"title"`
	Content  string             `json:"content" bson:"content"`
	AuthorID string             `json:"author" bson:"author"`
	Date     string             `json:"date" bson:"date"`
	Team     []struct {
		UserID   string `json:"userID" bson:"userID"`
		Role     string `json:"role" bson:"role"`
		JoinedAt string `json:"joinedAt" bson:"joinedAt"`
	} `json:"team" bson:"team"`
	AdvisorID string `json:"advisor" bson:"advisor"`
	Stars     []struct {
		UserID string  `json:"userID" bson:"userID"`
		Star   float64 `json:"star" bson:"star"`
		Date   string  `json:"date" bson:"date"`
	} `json:"stars" bson:"stars"`
	Tags    []string `json:"tags" bson:"tags"`
	Upvotes []string `json:"upvotes" bson:"upvotes"`
}

func (p *Project) CreateFrom(s Suggestion) error {
	if s.Status != "approved" {
		return errors.New("SUGGESTION_NOT_APPROVED")
	}
	if p.AdvisorID == "" {
		return errors.New("ADVISOR_NOT_ASSIGNED")
	}
	p.ID = s.ID
	p.Title = s.Title
	p.Content = s.Content
	p.AuthorID = s.AuthorID
	p.Date = time.Now().UTC().Format(time.RFC3339)
	p.Team = []struct {
		UserID   string `json:"userID" bson:"userID"`
		Role     string `json:"role" bson:"role"`
		JoinedAt string `json:"joinedAt" bson:"joinedAt"`
	}{
		{
			UserID:   s.AuthorID,
			Role:     "leader",
			JoinedAt: p.Date,
		},
	}
	p.Stars = s.Stars
	p.Tags = s.Tags
	p.Upvotes = []string{}
	if err := p.insertToDB(); err != nil {
		return err
	}
	return nil
}

func (p *Project) insertToDB() error {
	res, err := Projects.InsertOne(context.TODO(), p)
	if err != nil {
		return err
	}
	if err := Projects.FindOne(context.TODO(), bson.M{"_id": res.InsertedID}).Decode(&p); err != nil {
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

func GetAllProjects() ([]Project, error) {
	projects := []Project{}
	cursor, err := Projects.Find(context.TODO(), bson.M{})
	if err != nil {
		return projects, err
	}
	if err := cursor.All(context.TODO(), &projects); err != nil {
		return projects, err
	}
	return projects, nil
}
