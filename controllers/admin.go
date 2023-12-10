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
	routes.GET("/admin/master/category", AdminCategoryList)
	routes.POST("/admin/master/savecategory", saveCategory)
	routes.POST("/admin/master/del_category", DelCategory)
}

func AdminIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_index.html", gin.H{
		"content": "Admin page",
	})
}

func AdminCategoryList(c *gin.Context) {
	database := database.GetConnection()
	edit_id, _ := strconv.Atoi(c.Request.URL.Query().Get("edit"))
	rows, _ := database.Raw("select cat.id, cat.name, cat.description,cat.parent_id, f.filename,f.file_path, pcat.name as parent from categories as cat left join file_uploads as f on f.category_id = cat.id left join categories as pcat on pcat.id = cat.parent_id where f.default_image=1 group by cat.id ").Rows()

	var cat []models.CatagoryList
	var id int
	for rows.Next() {
		database.ScanRows(rows, &cat)
		rows.Scan(&id)
	}

	var selectedCategory models.CatagoryList
	for _, value := range cat {
		if value.ID == edit_id {
			selectedCategory = value
		}
	}
	c.HTML(http.StatusOK, "admin_category.html", gin.H{
		"title":            "Admin - Category Details",
		"category":         dtobjects.CategoryListAdminDto(cat),
		"endpoint":         Geturl(c),
		"selectedCategory": selectedCategory,
		"id":               edit_id,
	})
}

func saveCategory(c *gin.Context) {
	database := database.GetConnection()
	name := c.PostForm("name")
	description := c.PostForm("description")
	parent, _ := strconv.Atoi(c.PostForm("parent_id"))
	ID, _ := strconv.Atoi(c.PostForm("id"))
	form, err := c.MultipartForm()
	log.Println("parent id for save ", parent)
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
			return
		}
		defer inputFile.Close()

		// Copy the temporary image to the permanent location outputFile
		_, err = io.Copy(outputFile, inputFile)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusUnprocessableEntity, dtobjects.DetailedErrors("database", err))
			return
		}

		fileSize := (uint)(file.Size)
		categoryImages[index] = models.FileUpload{Filename: fileName, FilePath: string(filepath.Separator) + filePath, FileSize: fileSize}
	}

	if ID == 0 {
		category := models.Category{Name: name, Description: description, Images: categoryImages, ParentId: parent}
		err = database.Create(&category).Error
	} else {
		var updateData = map[string]interface{}{
			"name":        name,
			"description": description,
			"parent_id":   parent,
			"id":          ID,
		}
		err = database.Table("categories").Where("id=?", ID).Updates(&updateData).Error
	}
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusUnprocessableEntity, dtobjects.DetailedErrors("database", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"messages": "Data save successfully",
	})

}

func Geturl(c *gin.Context) string {
	return c.Request.URL.Path
}

func DelCategory(c *gin.Context) {
	db := database.GetConnection()
	delid, _ := strconv.Atoi(c.Request.URL.Query().Get("delid"))
	img := c.Request.URL.Query().Get("img")
	path, _ := os.Getwd()

	os.Remove(path + img)
	category := models.Category{ID: delid}
	db.Delete(&category)
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"messages": "Data Deleted successfully",
	})
}
