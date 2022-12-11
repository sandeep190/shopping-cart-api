package controllers

import (
	"log"
	"net/http"
	"shopping_cart/database"
	"shopping_cart/dtobjects"
	"shopping_cart/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func UsersRoutes(routes *gin.RouterGroup) {
	routes.POST("/signup", Registration)
}

func Registration(c *gin.Context) {
	var json dtobjects.SignupRequestDto
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dtobjects.BadRequestDto(err))
		return
	}

	pass, _ := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)
	log.Println(pass)
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
