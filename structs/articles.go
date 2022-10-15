package structs

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Id                     int
	Title, Anons, FullText string
}
