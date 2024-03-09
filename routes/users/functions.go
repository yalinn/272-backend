package users

import (
	"272-backend/library"
	db "272-backend/package/database"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUsers() ([]library.User, error) {
	var users []library.User
	cursor, err := db.Users.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

func GetUser(id string) (library.User, error) {
	var user library.User
	objID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{
		"_id": objID,
	}
	if err := db.Users.FindOne(context.TODO(), query).Decode(&user); err != nil {
		return user, err
	}
	return user, nil
}
