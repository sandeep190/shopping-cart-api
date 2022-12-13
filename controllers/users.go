package controllers

import (
	"errors"
	"net/http"
	"shopping_cart/database"
	"shopping_cart/dtobjects"
	"shopping_cart/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func UsersRoutes(routes *gin.RouterGroup) {
	routes.POST("/signup", Registration)
	routes.POST("/login", Login)
}

func Registration(c *gin.Context) {
	var json dtobjects.SignupRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dtobjects.BadRequestDto(err))
		return
	}

	pass, _ := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)
	db := database.GetConnection()

	result := db.Create(&models.User{
		Name:     json.Name,
		Password: string(pass),
		Address:  json.Address,
		Contact:  json.Contact,
		Email:    json.Email,
	})

	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, dtobjects.DetailedErrors("database", result.Error))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":  true,
		"messages": []string{"User created successfully"},
	})
}

func Login(c *gin.Context) {
	var request dtobjects.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dtobjects.BadRequestDto(err))
		return
	}
	db := database.GetConnection()
	var user models.User
	result := db.Where("email = ?", request.Email).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusForbidden, dtobjects.DetailedErrors("login_error", result.Error))
		return
	}

	bytePassword := []byte(request.Password)
	byteHashedPassword := []byte(user.Password)
	err := bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
	if err != nil {
		c.JSON(http.StatusForbidden, dtobjects.DetailedErrors("login", errors.New("invalid credential")))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   user.GenerateJwtToken(),
		"user_id": user.ID,
	})

}
