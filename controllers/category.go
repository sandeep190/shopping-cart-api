package controllers

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
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
	err := database.Preload("Images", "category_id IS NOT NULL").Find(&categories).Error
	//err := database.Find(&categories).Error
	if err != nil {
		c.JSON(http.StatusNotFound, dtobjects.DetailedErrors("fetch_error", err))
		return
	}
	c.JSON(http.StatusOK, dtobjects.CategoryListDto(categories))
}

func CreateCategory(c *gin.Context) {
	name := c.PostForm("name")
	description := c.PostForm("description")
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	files := form.File["image"]

	var categoryImages = make([]models.FileUpload, len(files))
	for index, file := range files {
		fileName := name + "__" + randomString(16) + ".png"

		dirPath := filepath.Join(".", "static", "images", "categories")
		filePath := filepath.Join(dirPath, fileName)
		// Create directory if does not exist
		if _, err = os.Stat(dirPath); os.IsNotExist(err) {
			err = os.MkdirAll(dirPath, os.ModeDir)
			if err != nil {
				c.JSON(http.StatusInternalServerError, dtobjects.DetailedErrors("io_error", err))
				return
			}
		}
		// Create file that will hold the image
		outputFile, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer outputFile.Close()

		// Open the temporary file that contains the uploaded image
		inputFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusOK, dtobjects.DetailedErrors("io_error", err))
		}
		defer inputFile.Close()

		// Copy the temporary image to the permanent location outputFile
		_, err = io.Copy(outputFile, inputFile)
		if err != nil {
			log.Fatal(err)
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		fileSize := (uint)(file.Size)
		categoryImages[index] = models.FileUpload{Filename: fileName, FilePath: string(filepath.Separator) + filePath, FileSize: fileSize}
	}

	database := database.GetConnection()
	category := models.Category{Name: name, Description: description, Images: categoryImages}
	err = database.Create(&category).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtobjects.DetailedErrors("db_error", err))
	}
	c.JSON(http.StatusOK, dtobjects.CreateCategoryDto(category))
}

func randomString(length int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
