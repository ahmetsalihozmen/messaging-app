package db

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Content  string
	Receiver string
	Status   string
}
