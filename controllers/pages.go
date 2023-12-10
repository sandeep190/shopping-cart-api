package controllers

import (
	"html/template"
	"net/http"
	"shopping_cart/database"
	"shopping_cart/models"

	"github.com/gin-gonic/gin"
)

type Pages struct {
	Title   string        `json:"title"`
	Content template.HTML `json:"content"`
}

type result struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Filename string `json:"filename"`
	FilePath string `json:"file_path"`
}

func PagesRoutes(routes *gin.RouterGroup) {
	routes.GET("", Index)
	routes.GET("/pages/about", AboutUs)
}

func AboutUs(c *gin.Context) {
	data := make(map[string]string)
	data["title"] = "About Us page "
	data["content"] = "This Is the About Us Page"
	c.HTML(http.StatusOK, "about.html", gin.H{
		"content": data,
	})
}

func Index(c *gin.Context) {
	data := make(map[string]interface{})
	data["title"] = "index page "
	data["content"] = "this is the index page"
	var CategoryList []models.CatagoryList

	var result []result
	database.DB.Table("categories").Select("id,name").Find(&CategoryList)
	database.DB.Model(&CategoryList).Select("categories.id,categories.name, f.filename,f.file_path").
		Joins("left join file_uploads as f on f.category_id = categories.id").
		Where("f.default_image", 1).Where("parent_id", 0).
		Group("categories.id").Find(&result)

	data["categoryList"] = result

	c.HTML(http.StatusOK, "index.html", gin.H{
		"content": data,
	})
}
