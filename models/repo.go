package models

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
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

// generate random string with a specified length
func randomString(length int) string {
    b := make([]byte, length+2)
    rand.Read(b)
    return fmt.Sprintf("%x", b)[2 : length+2]
}

// define padding function
func PKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// encrypt password
func encryptPassword(password string) (string, string, string) {
	key := randomString(32)
	iv := randomString(16)
	
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	bPlaintext := PKCS5Padding([]byte(password), aes.BlockSize, len(password))
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, bPlaintext)
	str := hex.EncodeToString(ciphertext)

	return str, key, iv
}

// decrypt password
func decryptPassword(password string, enc_key string) string {
	
	decoded_key, err := base64.StdEncoding.DecodeString(enc_key)
	if err != nil {
		fmt.Println("Error decoding key.")
		os.Exit(1)
	}

	key := string(decoded_key[:32])
	iv := string(decoded_key[32:])	
	
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	bCiphertext, _ := hex.DecodeString(password)
	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks([]byte(bCiphertext), []byte(bCiphertext))

	return string(bCiphertext)
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

// save credentials
func SaveCredentials(service string, username string, password string, db *gorm.DB) {
	enc_pwd, key, iv := encryptPassword(password)
	enc_key := base64.StdEncoding.EncodeToString([]byte(key + iv))
	db.Create(&Credential{Service: service, Username: username, Password: enc_pwd, Key: enc_key})
}

// get all services
func GetServices(db *gorm.DB) []string {
	var credentials []Credential
	var services []string
	db.Find(&credentials)
	for _, credential := range credentials {
		services = append(services, credential.Service)
	}
	return services
}

// get credentials for a service
func GetCredentials(service string, db *gorm.DB) (string, string) {
	var credential Credential
	err := db.Where("service = ?", service).First(&credential).Error
	if err != nil {
		fmt.Println("Service not found.")
		os.Exit(0)
	}
	return credential.Username, decryptPassword(credential.Password, credential.Key)
}

// get entire entry
func GetEntry(service string, db *gorm.DB) Credential {
	var credential Credential
	err := db.Where("service = ?", service).First(&credential).Error
	if err != nil {
		fmt.Println("Service not found.")
		os.Exit(0)
	}
	return credential
}

// delete credentials for a service
func DeleteCredentials(service string, db *gorm.DB) {
	var credential Credential
	result := db.Where("service = ?", service).First(&credential)
	if result.Error != nil {
		fmt.Println("Service not found.")
		os.Exit(0)
	} else {
		db.Delete(&credential)
		fmt.Printf("Deleted credentials for %s.\n", service)
	}
}

// edit credentials
func EditCredentials(service string, username string, password string, db *gorm.DB) {
	var credential Credential
	result := db.Where("service = ?", service).First(&credential)
	if result.Error != nil {
		fmt.Println("Service not found.")
		os.Exit(0)
	} else {
		enc_pwd, key, iv := encryptPassword(password)
		enc_key := base64.StdEncoding.EncodeToString([]byte(key + iv))
		credential.Username = username
		credential.Password = enc_pwd
		credential.Key = enc_key
		db.Save(&credential)
	}
}

// check if service exists
func CheckServiceExists(service string, db *gorm.DB) bool {
	var credential Credential
	if err := db.Where("service = ?", service).First(&credential).Error; err != nil {
		return false
	} else {return true}
}
