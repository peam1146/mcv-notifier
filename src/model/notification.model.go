package model

import "gorm.io/gorm"

type Notification struct {
	*gorm.Model
	Link       string
	CourseImg  string
	CourseName string
	Subject    string
	Date       string
}
