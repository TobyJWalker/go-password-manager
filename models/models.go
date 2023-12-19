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
	gorm.Model
	Service string
	Username string
	Password string
}