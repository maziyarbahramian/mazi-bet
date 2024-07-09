package utils

import (
	"errors"
	"mazi-bet/models"
	"net/mail"
	"strings"

	"gorm.io/gorm"
)

func ValidateJsonFormat(jsonBody map[string]interface{}, fields ...string) (string, error) {
	msg := "OK"
	for _, field := range fields {
		if _, ok := jsonBody[field]; !ok {
			msg = "Input Json doesn't include " + field
			break
		}
	}
	if msg != "OK" {
		return msg, errors.New("")
	}
	return msg, nil
}

// this function used to check user properties validation
func ValidateUser(jsonBody map[string]interface{}) (string, models.User, error) {
	msg := "OK"
	// Create User Object
	var user models.User
	user.FirstName = jsonBody["firstname"].(string)
	user.LastName = jsonBody["lastname"].(string)
	user.Email = jsonBody["email"].(string)
	user.Phone = jsonBody["phone"].(string)

	// Check FirstName Validation
	if len(strings.TrimSpace(user.FirstName)) == 0 {
		msg = "First Name can't be empty"
		return msg, models.User{}, errors.New("")
	}

	// Check LastName Validation
	if len(strings.TrimSpace(user.LastName)) == 0 {
		msg = "Last Name can't be empty"
		return msg, models.User{}, errors.New("")
	}

	// Check Phone Number Validation
	if !ValidatePhone(user.Phone) {
		msg = "Invalid Phone Number"
		return msg, models.User{}, errors.New("")
	}

	// Check Email Validation
	if !ValidateEmail(user.Email) {
		msg = "Invalid Email Address"
		return msg, models.User{}, errors.New("")
	}

	return msg, user, nil
}

// This Function Validates Input Email.
func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// This Function Validates Input Phone Number.
func ValidatePhone(phone string) bool {
	hasCharacter := false
	for _, digit := range phone {
		if digit < 48 || digit > 57 {
			hasCharacter = true
			break
		}
	}
	if hasCharacter {
		return false
	}
	return strings.HasPrefix(phone, "09") && len(phone) == 11
}

func CheckUnique(user models.User, username string, db *gorm.DB) (string, error) {
	msg := "OK"
	// Is Input Phone Number Unique or Not
	var existingUser models.User
	db.Where("phone = ?", user.Phone).First(&existingUser)
	if existingUser.ID != 0 {
		msg = "Input Phone Number has already been registered"
		return msg, errors.New("")
	}

	// Is Input Email Address Unique or Not
	db.Where("email = ?", user.Email).First(&existingUser)
	if existingUser.ID != 0 {
		msg = "Input Email Address has already been registered"
		return msg, errors.New("")
	}

	// Is Input Username Unique or Not
	var existingAccount models.Account
	db.Where("username = ?", username).First(&existingAccount)
	if existingAccount.ID != 0 {
		msg = "Input Username has already been registered"
		return msg, errors.New("")
	}
	return msg, nil
}
