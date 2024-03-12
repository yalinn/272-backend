package library

import (
	"272-backend/config"
	db "272-backend/package/database"
	jwts "272-backend/package/jwt"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"

	"github.com/emersion/go-imap/client"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	UserType string             `json:"user_type" bson:"user_type"`
	Roles    []string           `json:"roles" bson:"roles"`
	Token    string             `json:"token" bson:"token"`
}

func (u *User) InsertToDB() error {
	if u.Username == "" || u.UserType == "" {
		return errors.New("INVALID_USER")
	}
	token := jwts.CreateToken(jwt.MapClaims{
		"username": u.Username,
		"role":     u.Roles,
	})
	user := bson.M{
		"username":  u.Username,
		"user_type": u.UserType,
		"token":     token,
		"roles":     []string{u.UserType},
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

func (u *User) FindUser() (*User, error) {
	var user User
	query := bson.M{
		"username":  u.Username,
		"user_type": u.UserType,
	}
	if err := db.Users.FindOne(context.TODO(), query).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) GetByUsername() error {
	query := bson.M{
		"username":  u.Username,
		"user_type": u.UserType,
	}
	if err := db.Users.FindOne(context.TODO(), query).Decode(&u); err != nil {
		return err
	}
	return nil
}

func (u *User) LoginByEmail(pwd string) error {
	host := config.IMAP_S_HOST
	if u.UserType == "teacher" {
		host = config.IMAP_T_HOST
	}
	Imap, err := client.DialTLS(host+":"+config.IMAP_PORT, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	defer Imap.Logout()
	username := u.Username
	if err := Imap.Login(username, pwd); err != nil {
		return err
	}
	defer Imap.Logout()
	u.Username = username
	if err := u.GetByUsername(); err != nil {
		u.InsertToDB()
	}
	token := jwts.CreateToken(jwt.MapClaims{
		"username": username,
		"roles":    u.Roles,
	})
	if err := u.SetToken(token); err != nil {
		return err
	}
	return nil
}

func (u *User) GetToken() string {
	return u.Token
}

func (u *User) SetToken(token string) error {
	u.Token = token
	update := bson.M{
		"$set": bson.M{
			"token": token,
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

func (u *User) Stringify() string {
	out, _ := json.Marshal(u.Roles)
	return string(out)
}
