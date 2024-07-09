package models

type User struct {
	ID        uint   `gorm:"primary_key"`
	FirstName string `gorm:"type:varchar(255);not null"`
	LastName  string `gorm:"type:varchar(255);not null"`
	Phone     string `gorm:"type:varchar(255);unique;not null"`
	Email     string `gorm:"type:varchar(255);unique;not null"`
	Balance   float64 `gorm:"type:decimal(10,2);default:0"`
}
