package models

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
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

// check if encryption keys are generated
func CheckEncryptionKeysExist(db *gorm.DB) bool {
	key := os.Getenv("GO_PWM_KEY")
	iv := os.Getenv("GO_PWM_IV")
	if key == "" || iv == "" {
		return false
	} else {return true}
}

// generate random string with a specified length
func randomString(length int) string {
    b := make([]byte, length+2)
    rand.Read(b)
    return fmt.Sprintf("%x", b)[2 : length+2]
}

// generate encryption keys
func GenerateEncryptionKeys(db *gorm.DB) {
	key := randomString(32)
	iv := randomString(16)
	os.Setenv("GO_PWM_KEY", key)
	os.Setenv("GO_PWM_IV", iv)
}

// get encryption keys
func getEncryptionKeys() (string, string) {
	key := os.Getenv("GO_PWM_KEY")
	iv := os.Getenv("GO_PWM_IV")
	return key, iv
}

// define padding function
func PKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// encrypt password
func encryptPassword(password string) string {
	key, iv := getEncryptionKeys()
	
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	bPlaintext := PKCS5Padding([]byte(password), aes.BlockSize, len(password))
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, bPlaintext)
	str := hex.EncodeToString(ciphertext)

	return str
}

// decrypt password
func decryptPassword(password string) string {
	key, iv := getEncryptionKeys()
	
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	bCiphertext, _ := hex.DecodeString(password)
	plaintext := make([]byte, len(bCiphertext))
	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(plaintext, bCiphertext)

	return string(plaintext)
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
	enc_pwd := encryptPassword(password)
	db.Create(&Credential{Service: service, Username: username, Password: enc_pwd})
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