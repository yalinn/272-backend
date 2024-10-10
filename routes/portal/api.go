package portal

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostResponse struct {
	Message string `json:"message"`
}

type GetResponse struct {
	Message string `json:"message"`
}
type fetchCirriculumParams struct {
	Password string `json:"password"`
}

type CurriculumObject struct {
	Code       string `json:"code"`
	CourseID   string `json:"course_id"`
	CourseCode string `json:"course_code"`
	CourseName string `json:"course_name"`
	Credits    string `json:"credits"`
	Grade      string `json:"grade"`
	Semester   int    `json:"semester"`
	Year       int    `json:"year"`
}

type CourseData struct {
	ID         string `json:"id,omitempty" bson:"_id,omitempty"`
	Code       string `json:"code" bson:"code"`
	CourseID   string `json:"course_id" bson:"course_id"`
	CourseCode string `json:"course_code" bson:"course_code"`
	CourseName string `json:"course_name" bson:"course_name"`
	Credits    string `json:"credits" bson:"credits"`
	Department int    `json:"department" bson:"department"`
	Semester   int    `json:"semester" bson:"semester"`
	Year       int    `json:"year" bson:"year"`
}

type Course struct {
	ID         string `json:"id,omitempty" bson:"_id,omitempty"`
	Code       string `json:"code" bson:"code"`
	CourseID   string `json:"course_id" bson:"course_id"`
	CourseCode string `json:"course_code" bson:"course_code"`
	CourseName string `json:"course_name" bson:"course_name"`
	Credit     string `json:"credit" bson:"credit"`
}

type CirriculumIndex struct {
	CourseID string `json:"course_id" bson:"course_id"`
	Semester int    `json:"semester" bson:"semester"`
}

type Cirriculum struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Department     int                `json:"department" bson:"department"`
	CirriculumName string             `json:"cirriculum_name" bson:"cirriculum_name"`
	Index          []CirriculumIndex  `json:"index" bson:"index"`
}

type UserCourse struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID   string             `json:"user_id" bson:"user_id"`
	CourseID string             `json:"course_id" bson:"course_id"`
}
