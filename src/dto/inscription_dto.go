package dto

type EnrollRequestResponseDto struct {
	CourseId string `json:"course_id"`
	UserId   string `json:"user_id"`
}

type Student struct {
	UserId string `json:"user_id"`
}

type MyCourse struct {
	CourseId string `json:"course_id"`
}

type StudentsInCourse []Student
type MyCourses []MyCourse
