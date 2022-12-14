package controllers

import (
	"net/http"
	"shopping_cart/database"
	"shopping_cart/dtobjects"
	"shopping_cart/middleware"
	"shopping_cart/models"

	"github.com/gin-gonic/gin"
)

func CategoryRoutes(router *gin.RouterGroup) {
	router.GET("", CategoryList)
	router.Use(middleware.AuthMiddleware())
	{
		router.POST("/create", CreateCategory)
	}
}

func CategoryList(c *gin.Context) {
	database := database.GetConnection()
	var categories []models.Category
	// err := database.Preload("Images", "category_id IS NOT NULL").Find(&categories).Error
	err := database.Find(&categories).Error
	if err != nil {
		c.JSON(http.StatusNotFound, dtobjects.DetailedErrors("fetch_error", err))
		return
	}
	c.JSON(http.StatusOK, dtobjects.CategoryListDto(categories))
}

func CreateCategory(c *gin.Context) {

}
