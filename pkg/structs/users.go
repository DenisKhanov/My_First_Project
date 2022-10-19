package structs

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	Login, Password string
}
