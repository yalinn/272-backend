package library

import (
	"272-backend/pkg"
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Suggestions *mongo.Collection
	Rejections  *mongo.Collection
	Approvals   *mongo.Collection
	Reports     *mongo.Collection
)

func init() {
	Suggestions = pkg.Mongo.Collection("suggestions")
	Rejections = pkg.Mongo.Collection("rejections")
	Approvals = pkg.Mongo.Collection("approvals")
	Reports = pkg.Mongo.Collection("reports")
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
	s.Date = time.Now().UTC().Format(time.RFC3339)
	s.Status = "pending"
	s.Upvotes = []string{}
	s.Stars = []struct {
		UserID string  `json:"userID" bson:"userID"`
		Star   float64 `json:"star" bson:"star"`
		Date   string  `json:"date" bson:"date"`
	}{}
	s.Tags = []string{}
	res, err := Suggestions.InsertOne(context.TODO(), s)
	if err != nil {
		return err
	}
	if err := Suggestions.FindOne(context.TODO(), bson.M{"_id": res.InsertedID}).Decode(&s); err != nil {
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
	Reason     string             `json:"reason" bson:"reason"`
	ExecutorID string             `json:"executor" bson:"executor"`
	Date       string             `json:"date" bson:"date"`
}

func (s *Suggestion) Reject(executorID string, reason string) error {
	rejection := Rejection{
		ID:         s.ID,
		Reason:     reason,
		ExecutorID: executorID,
		Date:       time.Now().UTC().Format(time.RFC3339),
	}
	if _, err := Rejections.InsertOne(context.TODO(), rejection); err != nil {
		return err
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
		return err
	}
	if err := Suggestions.FindOne(context.TODO(), bson.M{"_id": s.ID}).Decode(&s); err != nil {
		return err
	}
	return nil
}

type Approval struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ExecutorID string             `json:"executor" bson:"executor"`
	Date       string             `json:"date" bson:"date"`
}

func (s *Suggestion) Approve(executorID string) error {
	approval := Approval{
		ID:         s.ID,
		ExecutorID: executorID,
		Date:       time.Now().UTC().Format(time.RFC3339),
	}
	if _, err := Approvals.InsertOne(context.TODO(), approval); err != nil {
		return err
	}
	update := bson.M{
		"$set": bson.M{
			"status": "approved",
		},
	}
	query := bson.M{
		"_id": s.ID,
	}
	res, err := Suggestions.UpdateOne(context.TODO(), query, update)
	if err != nil {
		return err
	}
	if err := Suggestions.FindOne(context.TODO(), bson.M{"_id": res.UpsertedID}).Decode(&s); err != nil {
		return err
	}
	return nil
}

type Report struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ExecutorID string             `json:"executor" bson:"executor"`
	Date       string             `json:"date" bson:"date"`
}

func (s *Suggestion) Report(executorID string) error {
	report := Report{
		ID:         s.ID,
		ExecutorID: executorID,
		Date:       time.Now().UTC().Format(time.RFC3339),
	}
	if _, err := Reports.InsertOne(context.TODO(), report); err != nil {
		return err
	}
	update := bson.M{
		"$set": bson.M{
			"status": "reported",
		},
	}
	query := bson.M{
		"_id": s.ID,
	}
	res, err := Suggestions.UpdateOne(context.TODO(), query, update)
	if err != nil {
		return err
	}
	if err := Suggestions.FindOne(context.TODO(), bson.M{"_id": res.UpsertedID}).Decode(&s); err != nil {
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

func GetRejectedSuggestions() ([]Suggestion, error) {
	var suggestions []Suggestion
	cursor, err := Suggestions.Find(context.TODO(), bson.M{"status": "rejected"})
	if err != nil {
		return []Suggestion{}, err
	}
	if err := cursor.All(context.TODO(), &suggestions); err != nil {
		return []Suggestion{}, err
	}
	return suggestions, nil
}

func GetApprovedSuggestions() ([]Suggestion, error) {
	var suggestions []Suggestion
	cursor, err := Suggestions.Find(context.TODO(), bson.M{"status": "approved"})
	if err != nil {
		return []Suggestion{}, err
	}
	if err := cursor.All(context.TODO(), &suggestions); err != nil {
		return []Suggestion{}, err
	}
	return suggestions, nil
}

func GetAllSuggestions() ([]Suggestion, error) {
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

func GetPendingSuggestions() ([]Suggestion, error) {
	var suggestions []Suggestion
	cursor, err := Suggestions.Find(context.TODO(), bson.M{"status": "pending"})
	if err != nil {
		return []Suggestion{}, err
	}
	if err := cursor.All(context.TODO(), &suggestions); err != nil {
		return []Suggestion{}, err
	}
	return suggestions, nil
}

func GetReportedSuggestions() ([]Suggestion, error) {
	var suggestions []Suggestion
	cursor, err := Suggestions.Find(context.TODO(), bson.M{"status": "reported"})
	if err != nil {
		return []Suggestion{}, err
	}
	if err := cursor.All(context.TODO(), &suggestions); err != nil {
		return []Suggestion{}, err
	}
	return suggestions, nil
}

type SuggestionResponse struct {
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

func (s *Suggestion) ToResponse(userID string) SuggestionResponse {
	starred := 0.00
	voted := false
	for _, stars := range s.Stars {
		if stars.UserID == userID {
			starred = stars.Star
			voted = true
			break
		}
	}
	response := SuggestionResponse{
		ID:         s.ID.Hex(),
		Title:      s.Title,
		Content:    s.Content,
		Author:     s.AuthorID,
		Upvotes:    len(s.Upvotes),
		Stars:      s.CalculateAverageStars(),
		Date:       s.Date,
		Tags:       s.Tags,
		Status:     s.Status,
		Starred:    starred,
		Voted:      voted,
		Department: GetDepartmentID(s.AuthorID),
	}
	return response
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
