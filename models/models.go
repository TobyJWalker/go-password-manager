package models

import (
	"time"

	"gorm.io/gorm"
)

type Master struct {
	gorm.Model
	Password string
	LastLogin time.Time
}

type Credential struct {
	ID uint `gorm:"primary_key"`
	Service string
	Username string
	Password string
	Key string
}