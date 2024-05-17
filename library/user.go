package library

import (
	"272-backend/config"
	db "272-backend/package/database"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/emersion/go-imap/client"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Username   string   `json:"id,omitempty" bson:"_id,omitempty"`
	UserType   string   `json:"user_type" bson:"user_type"`
	Roles      []string `json:"roles" bson:"roles"`
	Department int      `json:"department" bson:"department"`
}

func (u *User) GetDepartmentID() int {
	chars := strings.Split(u.Username[3:], "")
	id_slice := []string{}
	for i := 0; i < len(chars)-3; i++ {
		id_slice = append(id_slice, chars[i])
	}
	department, _ := strconv.Atoi(strings.Join(id_slice, ""))
	return department
}

func (u *User) InsertToDB() error {
	if u.Username == "" || u.UserType == "" {
		return errors.New("INVALID_USER")
	}
	user := bson.M{
		"_id":        u.Username,
		"user_type":  u.UserType,
		"roles":      []string{u.UserType},
		"department": u.GetDepartmentID(),
	}
	if _, err := db.Users.InsertOne(context.TODO(), user); err != nil {
		return err
	}
	if err := db.Users.FindOne(context.TODO(), bson.M{"_id": u.Username}).Decode(&u); err != nil {
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
		"_id": u.Username,
	}
	if _, err := db.Users.UpdateOne(context.TODO(), query, update); err != nil {
		return err
	}
	return nil
}

func (u *User) FindUser() error {
	query := bson.M{
		"_id":       u.Username,
		"user_type": u.UserType,
	}
	if err := db.Users.FindOne(context.TODO(), query).Decode(&u); err != nil {
		return err
	}
	return nil
}

func (u *User) GetByUsername() error {
	query := bson.M{
		"_id":       u.Username,
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
	return nil
}

func (u *User) Stringify() string {
	out, _ := json.Marshal(u)
	return string(out)
}

func (u *User) InitToken(token string) error {
	us, err := db.Redis.Get(token)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(us), &u); err != nil {
		return err
	}
	return nil
}
