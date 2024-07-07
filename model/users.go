package model

import (
	"go-api-kits/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"html"
	"strings"
)

type Users struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Email    string `gorm:"size:255;not null;unique" json:"email"`
	Password string `gorm:"size:500;not null;" json:"-"`
}

func (user *Users) Save() (*Users, error) {
	err := database.Database.Create(&user).Error
	if err != nil {
		return &Users{}, err
	}
	return user, nil
}

func (user *Users) BeforeSave(*gorm.DB) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	return nil
}

func (user *Users) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func FindUserByUsername(username string) (Users, error) {
	var users Users
	err := database.Database.Where("username=?", username).Find(&users).Error
	if err != nil {
		return Users{}, err
	}
	return users, nil
}
