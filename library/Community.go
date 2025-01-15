package library

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Community struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Tags        []string           `json:"tags" bson:"tags"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

type Role struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Allow []string           `json:"allow" bson:"allow"`
	Deny  []string           `json:"deny" bson:"deny"`
}

type CommunityMember struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	MemberID string `json:"member_id" bson:"member_id"`
	Roles    []Role `json:"roles" bson:"roles"`
}
