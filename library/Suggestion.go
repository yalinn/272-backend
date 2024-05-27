package library

import (
	"272-backend/pkg"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Suggestions *mongo.Collection
	Rejections  *mongo.Collection
)

func init() {
	Suggestions = pkg.Mongo.Collection("suggestions")
	Rejections = pkg.Mongo.Collection("rejections")
}

type Suggestion struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title    string             `json:"title" bson:"title"`
	Content  string             `json:"content" bson:"content"`
	AuthorID string             `json:"author" bson:"author"`
	Upvotes  []string           `json:"upvotes" bson:"upvotes"`
	Stars    []struct {
		UserID string  `json:"userID" bson:"userID"`
		Star   float64 `json:"star" bson:"star"`
		Date   string  `json:"date" bson:"date"`
	} `json:"stars" bson:"stars"`
	Date   string   `json:"date" bson:"date"`
	Tags   []string `json:"tags" bson:"tags"`
	Status string   `json:"status" bson:"status"`
}

func (s *Suggestion) WithID(id string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	if err := Suggestions.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&s); err != nil {
		return err
	}
	return nil
}

func (s *Suggestion) InsertToDB() error {
	suggestion := bson.M{
		"title":   s.Title,
		"content": s.Content,
		"author":  s.AuthorID,
		"upvotes": []string{},
		"stars": []struct {
			UserID string  `json:"userID" bson:"userID"`
			Star   float64 `json:"star" bson:"star"`
			Date   string  `json:"date" bson:"date"`
		}{},
		"date":   time.Now().UTC().Format(time.RFC3339),
		"tags":   []string{},
		"status": "pending",
	}
	if _, err := Suggestions.InsertOne(context.TODO(), suggestion); err != nil {
		return err
	}
	if err := Suggestions.FindOne(context.TODO(), bson.M{"title": s.Title}).Decode(&s); err != nil {
		return err
	}
	return nil
}

func (s *Suggestion) GiveUpvote(userID string) error {
	for _, upvote := range s.Upvotes {
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
		"_id": s.ID,
	}
	if _, err := Suggestions.UpdateOne(context.TODO(), query, update); err != nil {
		return err
	}
	if err := Suggestions.FindOne(context.TODO(), bson.M{"_id": s.ID}).Decode(&s); err != nil {
		return err
	}
	return nil
}

func (s *Suggestion) GiveStar(userID string, star int) error {
	query := bson.M{
		"_id": s.ID,
		"stars": bson.M{
			"$elemMatch": bson.M{
				"userID": userID,
			},
		},
	}
	if err := Suggestions.FindOne(context.TODO(), query).Decode(&s); err == nil {
		/* log.Println("User already starred") */
		stars := []struct {
			UserID string  `json:"userID" bson:"userID"`
			Star   float64 `json:"star" bson:"star"`
			Date   string  `json:"date" bson:"date"`
		}{}
		for _, st := range s.Stars {
			/* log.Println(st.UserID) */
			if st.UserID == userID {
				st.Star = float64(star)
				st.Date = time.Now().UTC().Format(time.RFC3339)
				stars = append(stars, st)
			} else {
				stars = append(stars, st)
			}
		}
		/* log.Println("ctrl") */
		update := bson.M{
			"$set": bson.M{
				"stars": stars,
			},
		}
		query := bson.M{
			"_id": s.ID,
		}
		if _, err := Suggestions.UpdateOne(context.TODO(), query, update); err != nil {
			log.Println(err.Error())
			return err
		}
		if err := Suggestions.FindOne(context.TODO(), bson.M{"_id": s.ID}).Decode(&s); err != nil {
			log.Println(err.Error())
			return err
		}
		return nil
	} else {
		/* log.Println("User has not starred") */
		update := bson.M{
			"$addToSet": bson.M{
				"stars": bson.M{
					"userID": userID,
					"star":   float64(star),
					"date":   time.Now().UTC().Format(time.RFC3339),
				},
			},
		}
		query := bson.M{
			"_id": s.ID,
		}
		if _, err := Suggestions.UpdateOne(context.TODO(), query, update); err != nil {
			return err
		}
		if err := Suggestions.FindOne(context.TODO(), bson.M{"_id": s.ID}).Decode(&s); err != nil {
			return err
		}
		return nil
	}
}

type Rejection struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title      string             `json:"title" bson:"title"`
	Content    string             `json:"content" bson:"content"`
	AuthorID   string             `json:"author" bson:"author"`
	Upvotes    int                `json:"upvotes" bson:"upvotes"`
	Stars      float64            `json:"stars" bson:"stars"`
	Tags       []string           `json:"tags" bson:"tags"`
	Reasons    []string           `json:"reason" bson:"reason"`
	ExecutorID string             `json:"executor" bson:"executor"`
	Date       string             `json:"date" bson:"date"`
}

func (s *Suggestion) Reject(reasons []string, executorID string) (Rejection, error) {
	rejection := Rejection{
		Title:      s.Title,
		Content:    s.Content,
		AuthorID:   s.AuthorID,
		Upvotes:    len(s.Upvotes),
		Stars:      s.CalculateAverageStars(),
		Tags:       s.Tags,
		Reasons:    reasons,
		ExecutorID: executorID,
		Date:       time.Now().UTC().Format(time.RFC3339),
	}
	if _, err := Rejections.InsertOne(context.TODO(), rejection); err != nil {
		return Rejection{}, err
	}
	if err := Rejections.FindOne(context.TODO(), bson.M{"title": s.Title}).Decode(&rejection); err != nil {
		return Rejection{}, err
	}
	update := bson.M{
		"$set": bson.M{
			"status": "rejected",
		},
	}
	query := bson.M{
		"_id": s.ID,
	}
	if _, err := Suggestions.UpdateOne(context.TODO(), query, update); err != nil {
		return rejection, err
	}
	return rejection, nil
}

type Proposal struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title      string             `json:"title" bson:"title"`
	Content    string             `json:"content" bson:"content"`
	AuthorID   string             `json:"author" bson:"author"`
	Upvotes    int                `json:"upvotes" bson:"upvotes"`
	Stars      float64            `json:"stars" bson:"stars"`
	Tags       []string           `json:"tags" bson:"tags"`
	ExecutorID string             `json:"executor" bson:"executor"`
	Date       string             `json:"date" bson:"date"`
}

func (s *Suggestion) Approve(executorID string) error {
	update := bson.M{
		"$set": bson.M{
			"status": "approved",
		},
	}
	query := bson.M{
		"_id": s.ID,
	}
	if _, err := Suggestions.UpdateOne(context.TODO(), query, update); err != nil {
		return err
	}
	return nil
}

func (s *Suggestion) CalculateAverageStars() float64 {
	totalStars := 0.00
	starCount := 0.00
	for _, star := range s.Stars {
		totalStars += star.Star
		starCount++
	}
	average := 0.00
	if starCount != 0 {
		average = totalStars / starCount
	}
	return average
}

func GetSuggestions() ([]Suggestion, error) {
	var suggestions []Suggestion
	cursor, err := Suggestions.Find(context.TODO(), bson.M{})
	if err != nil {
		return []Suggestion{}, err
	}
	if err := cursor.All(context.TODO(), &suggestions); err != nil {
		return []Suggestion{}, err
	}
	return suggestions, nil
}
