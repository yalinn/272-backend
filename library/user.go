package library

import (
	db "272-backend/package/database"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"`
	UserType string             `json:"userType" bson:"userType"`
	Roles    []string           `json:"roles" bson:"roles"`
	Token    string             `json:"token" bson:"token"`
}

func (u *User) InsertToDB() error {
	user := bson.M{
		"username": u.Username,
		"password": u.Password,
		"userType": u.UserType,
		"roles":    u.Roles,
		"token":    u.Token,
	}
	if _, err := db.Users.InsertOne(context.TODO(), user); err != nil {
		return err
	}
	if err := db.Users.FindOne(context.TODO(), bson.M{"username": u.Username}).Decode(&u); err != nil {
		return err
	}
	return nil
}

func (u *User) PutRoles(roles []string) error {
	u.Roles = append(u.Roles, roles...)
	update := bson.M{
		"$set": bson.M{
			"roles": u.Roles,
		},
	}
	query := bson.M{
		"_id": u.ID,
	}
	if _, err := db.Users.UpdateOne(context.TODO(), query, update); err != nil {
		return err
	}
	return nil
}

func GetUsers() ([]User, error) {
	var users []User
	cursor, err := db.Users.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

func GetUser() (*User, error) {
	var user User
	query := bson.M{
		"username": "admin",
	}
	if err := db.Users.FindOne(context.TODO(), query).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
