package utils

import (
	"errors"
	"log"
	"mazi-bet/models"

	"gorm.io/gorm"
)

func Withdraw(db *gorm.DB, userId uint, amount float64,description string) (models.User, error) {
	// start new transaction
	tx := db.Begin()
	var user models.User
	if tx.Error != nil {
		log.Fatalf("Failed to begin transaction %v", tx.Error)

		return user, errors.New("Failed to begin transaction")
	}

	// select balance for update
	if err := tx.Raw("SELECT * FROM users WHERE id = ? FOR UPDATE", userId).Scan(&user).Error; err != nil {
		tx.Rollback()
		log.Fatalf("failed to lock row: %v", err)
		return user, errors.New("Failed to lock row")
	}

	if user.Balance < amount {
		return user, errors.New("insufficient balance")
	}
	user.Balance -= amount
	if err := tx.Model(&models.User{}).Where("id = ?", user.ID).Update("balance", user.Balance).Error; err != nil {
		tx.Rollback()
		log.Fatalf("failed to update balance: %v", err)
		return user, errors.New("Failed to update balance")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Fatalf("failed to commit transaction: %v", err)
		return user, errors.New("Failed to commit transaction")
	}

	return user, nil
}

func Deposit(db *gorm.DB, userId uint, amount float64, description string) (models.User, error) {
	var user models.User

	if amount <= 0 {
		return user, errors.New("amount should be greater than 0")
	}

	// start new transaction
	tx := db.Begin()

	if tx.Error != nil {
		log.Fatalf("Failed to begin transaction %v", tx.Error)

		return user, errors.New("Failed to begin transaction")
	}

	// select balance for update
	if err := tx.Raw("SELECT * FROM users WHERE id = ? FOR UPDATE", userId).Scan(&user).Error; err != nil {
		tx.Rollback()
		log.Fatalf("failed to lock row: %v", err)
		return user, errors.New("Failed to lock row")
	}
	user.Balance += amount
	if err := tx.Model(&models.User{}).Where("id = ?", user.ID).Update("balance", user.Balance).Error; err != nil {
		tx.Rollback()
		log.Fatalf("failed to update balance: %v", err)
		return user, errors.New("Failed to update balance")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Fatalf("failed to commit transaction: %v", err)
		return user, errors.New("Failed to commit transaction")
	}

	return user, nil
}
