package models

import (
	"time"

	"gorm.io/gorm"
)

// function to check if logged in recently
func CheckSessionValid(db *gorm.DB) bool {
	var master Master
	db.First(&master)
	diff := time.Since(master.LastLogin)
	return diff.Minutes() < 5
}

// check if master password already exists
func CheckMasterExists(db *gorm.DB) bool {
	var master Master
	if err := db.Where("id = ?", 1).First(&master).Error; err != nil {
		return false
	} else {return true}
}

// save master password
func SaveMasterPassword(hash string, db *gorm.DB) {
	db.Create(&Master{Password: hash, LastLogin: time.Now()})
}

// update login time
func UpdateLoginTime(db *gorm.DB) {
	var master Master
	db.First(&master)
	db.Model(&master).Update("LastLogin", time.Now())
}