package suggestions

import (
	"272-backend/library"
	db "272-backend/package/database"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetSuggestion(id string) (library.Suggestion, error) {
	var suggestion library.Suggestion
	objID, _ := primitive.ObjectIDFromHex(id)
	if err := db.Suggestions.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&suggestion); err != nil {
		return library.Suggestion{}, err
	}
	return suggestion, nil
}

func GetSuggestions() ([]library.Suggestion, error) {
	var suggestions []library.Suggestion
	cursor, err := db.Suggestions.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.Background(), &suggestions); err != nil {
		return nil, err
	}
	return suggestions, nil
}

func GetUsersSuggestions(userID string) ([]library.Suggestion, error) {
	var suggestions []library.Suggestion
	cursor, err := db.Suggestions.Find(context.Background(), bson.M{"userID": userID})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.Background(), &suggestions); err != nil {
		return nil, err
	}
	return suggestions, nil
}
