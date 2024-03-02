package controllers

import (
	"html/template"
	"log"
	"net/http"
	"shopping_cart/database"
	"shopping_cart/models"

	"github.com/gin-gonic/gin"
)

type Pages struct {
	Title   string        `json:"title"`
	Content template.HTML `json:"content"`
}

type Result struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Filename string `json:"filename"`
	FilePath string `json:"file_path"`
}

type CategoryMenu struct {
	Id     int
	Name   string
	Parent int
	Child  []CategoryMenu
}

func PagesRoutes(routes *gin.RouterGroup) {
	routes.GET("", Index)
	routes.GET("/pages/about", AboutUs)
	routes.GET("/users/carts", Carts)
	routes.GET("/products/:categoryid", Shop)
	routes.GET("/products/details/:productid", ShopDetails)
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

	result := make(map[int]Result)
	catMenuList := make(map[int]CategoryMenu)
	database.DB.Table("categories").Select("categories.id,categories.name,categories.parent_id,f.filename,f.file_path").
		Joins("left join file_uploads as f on f.category_id = categories.id").
		Group("categories.id").Find(&CategoryList)
	//log.Printf("category data===>%#v", CategoryList)
	for _, val := range CategoryList {
		catMenuList[val.ID] = CategoryMenu{Id: val.ID, Name: val.Name, Parent: val.ParentId}
		if val.ParentId == 0 {
			result[val.ID] = Result{Id: val.ID, Name: val.Name, Filename: val.Filename, FilePath: val.FilePath}
		}
	}
	var featureProd []models.ProductList
	database.DB.Table("products").Limit(8).Scan(&featureProd)

	data["categoryList"] = result
	data["catMenuList"] = catMenuList
	data["featureProd"] = featureProd

	c.HTML(http.StatusOK, "index.html", gin.H{
		"content": data,
	})
}

func Shop(c *gin.Context) {
	data := make(map[string]interface{})
	data["title"] = "products list "
	data["content"] = "this is the products list"

	db := database.DB
	var prodLists []models.ProductList
	db.Table("products").Scan(&prodLists)
	c.HTML(http.StatusOK, "shop.html", gin.H{
		"content":   data,
		"prodLists": prodLists,
	})
}

func ShopDetails(c *gin.Context) {
	data := make(map[string]interface{})
	productId := c.Param("productid")
	log.Printf("product id===%#v", productId)

	var prodDetails models.ProductList
	var RelatedProduct []models.ProductList

	db := database.DB
	db.Table("products").Where("id", productId).Scan(&prodDetails)
	data["title"] = prodDetails.Title
	data["content"] = "this is the descriptions"

	db.Table("products").Select("id,price,title,images").Where("cat_id", prodDetails.CatID).Scan(&RelatedProduct)

	log.Printf("category ==> %d and subcateogry==>%d", prodDetails.CatID, prodDetails.SubcatID)

	c.HTML(http.StatusOK, "detail.html", gin.H{
		"content":        data,
		"prodDetails":    prodDetails,
		"relatedProduct": RelatedProduct,
	})
}

func Carts(c *gin.Context) {
	data := make(map[string]interface{})
	data["title"] = "users Carts "
	data["content"] = "users carts details"
	c.HTML(http.StatusOK, "cart.html", gin.H{
		"content": data,
	})
}
