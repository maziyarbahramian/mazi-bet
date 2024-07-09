package utils

import (
	"errors"
	"mazi-bet/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// this function used to create an account and insert it into database
func CreateAccount(user_id int, username string, is_admin bool, password string, db *gorm.DB) (string, models.Account, error) {
	msg := "OK"
	// Instantiating Account Object
	var account models.Account
	account.UserID = uint(user_id)
	account.Username = username
	account.IsAdmin = is_admin
	account.Token = ""

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		msg = "Failed to Hashing Password"
		return msg, models.Account{}, errors.New("")
	}
	account.Password = string(hash)

	// insert account into database
	createdAccount := db.Create(&account)
	if createdAccount.Error != nil {
		msg = "Failed to Create Account"
		return msg, models.Account{}, errors.New("")
	}

	// generate token
	var token *jwt.Token
	if is_admin {
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":    account.ID,
			"exp":   time.Now().Add(time.Hour).Unix(),
			"admin": true,
		})
	} else {
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":    account.ID,
			"exp":   time.Now().Add(time.Hour).Unix(),
			"admin": false,
		})
	}
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		msg = "Failed To Create Token"
		return msg, models.Account{}, errors.New("")
	}
	account.Token = tokenString

	// update account
	db.Save(&account)

	return msg, account, nil
}

func Login(username, passwrod string, db *gorm.DB) (string, models.Account, error) {
	msg := "OK"

	// find account based on input username
	var account models.Account
	db.Where("username = ?", username).First(&account)

	// Account Not Found
	if account.ID == 0 {
		msg = "Invalid Username"
		return msg, models.Account{}, errors.New("")
	}

	// Incorrect Password
	err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(passwrod))
	if err != nil {
		msg = "Wrong Password"
		return msg, models.Account{}, errors.New((""))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  account.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		msg = "Failed To Create Token"
		return msg, models.Account{}, errors.New("")
	}

	// Update Account's Token In Database
	account.Token = tokenString
	db.Save(&account)
	return msg, account, nil
}
