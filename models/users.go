package models

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	//Id           uint    `gorm:"primary_key"`
	Name     string    `gorm:"varchar(255);not null"`
	Address  string    `gorm:"varchar(255);not null"`
	Contact  string    `gorm:"column:contact"`
	Email    string    `gorm:"column:email;unique_index"`
	Password string    `gorm:"column:password;not null"`
	Created  time.Time `gorm:"column:created_at"`

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

func (user *User) GenerateJwtToken() string {

	jwt_token := jwt.New(jwt.SigningMethodHS512)

	jwt_token.Claims = jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Email,
		"exp":      time.Now().Add(time.Hour * 24 * 90).Unix(),
	}
	// Sign and get the complete encoded token as a string
	token, _ := jwt_token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return token
}