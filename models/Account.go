package models

type Account struct {
	ID       uint   `gorm:"primary_key"`
	UserID   uint   `gorm:"not null"`
	Username string `gorm:"type:varchar(255);unique;not null"`
	Password string `gorm:"type:varchar(255)"`
	Token    string `gorm:"not null"`
	IsActive bool   `gorm:"default:true"`
	IsAdmin  bool   `gorm:"default:false"`
}