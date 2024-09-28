package portal

type PostResponse struct {
	Message string `json:"message"`
}

type GetResponse struct {
	Message string `json:"message"`
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
