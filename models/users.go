package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	//Id           uint    `gorm:"primary_key"`
	FirstName string `gorm:"varchar(255);not null"`
	LastName  string `gorm:"varchar(255);not null"`
	Username  string `gorm:"column:username"`
	Email     string `gorm:"column:email;unique_index"`
	Password  string `gorm:"column:password;not null"`

	// Comments []Comment `gorm:"foreignkey:UserId"`

	// Roles     []Role     `gorm:"many2many:users_roles;"`
	// UserRoles []UserRole `gorm:"foreignkey:UserId"`
}

func (u *User) SetPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty")
	}
	bytePassword := []byte(password)
	// Make sure the second param `bcrypt generator cost` between [4, 32)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.Password = string(passwordHash)
	return nil
}
