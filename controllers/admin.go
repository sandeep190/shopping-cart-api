package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"shopping_cart/database"
	"shopping_cart/dtobjects"
	"shopping_cart/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(routes *gin.RouterGroup) {
	routes.GET("admin", AdminIndex)
	routes.Any("/admin/master/category", AdminCategoryList)
}

func AdminIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_index.html", gin.H{
		"content": "Admin page",
	})
}

func AdminCategoryList(c *gin.Context) {
	log.Println("details ---->", c.Request.Method)
	database := database.GetConnection()
	if c.Request.Method == "POST" {
		name := c.PostForm("name")
		description := c.PostForm("description")
		parent, _ := strconv.Atoi(c.PostForm("parent_id"))
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
				return
			}

			fileSize := (uint)(file.Size)
			categoryImages[index] = models.FileUpload{Filename: fileName, FilePath: string(filepath.Separator) + filePath, FileSize: fileSize}
		}

		category := models.Category{Name: name, Description: description, Images: categoryImages, ParentId: parent}
		err = database.Create(&category).Error
		if err != nil {
			log.Fatal(err)
		}
	}
	var categories []models.Category
	err := database.Preload("Images", "category_id IS NOT NULL").Find(&categories).Error
	if err != nil {
		log.Println(err)
	}
	c.HTML(http.StatusOK, "admin_category.html", gin.H{
		"content":  "Category Details",
		"category": dtobjects.CategoryListDto(categories),
	})
}
