package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"
	"shopping_cart/database"
	"shopping_cart/dtobjects"
	"shopping_cart/models"
	"strconv"

	ginsession "shopping_cart/middleware"

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
	routes.GET("/login", SignIn)
	routes.POST("/login", SignIn)
	routes.GET("/pages/about", AboutUs)
	routes.GET("/users/carts", Carts)
	routes.POST("/users/addtoCarts/:productid", addCards)
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
	store := ginsession.FromContext(c)
	userId, ok := store.Get("userId")
	if !ok {
		c.Redirect(302, "/login")
	}
	log.Println("session ", userId, ok)

	data["title"] = "users Carts "
	data["content"] = "users carts details"

	c.HTML(http.StatusOK, "cart.html", gin.H{
		"content": data,
	})
}

func SignIn(c *gin.Context) {
	log.Println("methods===>", c.Request.Method)
	data := make(map[string]string)
	data["title"] = "Sigin/SignUp "
	data["content"] = "This Is the Login Page"

	store := ginsession.FromContext(c)
	sessemail, ok := store.Get("email")
	log.Println("emails sessison===", sessemail, ok)
	if ok {
		c.Redirect(302, "/")
	}

	if c.Request.Method == "POST" {
		var request dtobjects.LoginRequest
		var user models.User

		err2 := c.ShouldBindJSON(&request)
		if err2 == nil {
			log.Println("errror", err2)
		}
		Pass := GetMD5Hash(request.Password)
		database.DB.Select("id,email,name,contact").Where("email", request.Email).Where("password", Pass).First(&user)

		if request.Email == user.Email {
			store.Set("email", user.Email)
			store.Set("userId", user.ID)
			err := store.Save()
			if err != nil {
				c.AbortWithError(500, err)
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Login successfull",
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":  false,
				"message": "Email/password do not match!",
			})
			return
		}
	} else {

		c.HTML(http.StatusOK, "login.html", gin.H{
			"content": data,
		})
	}
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func addCards(ctx *gin.Context) {
	productId, _ := strconv.Atoi(ctx.Param("productid"))
	log.Println("product id is ", productId)

	store := ginsession.FromContext(ctx)
	userId, ok := store.Get("userId")
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  false,
			"message": "please sign-In for add to card",
		})
	}

	log.Println("session ", userId, ok)
	var usersCarts models.UsersCarts
	result := database.DB.Table("user_carts").Where("user_id", userId).Where("product_id", productId).First(&usersCarts)
	log.Printf("users cards===>#%v", usersCarts)
	log.Println("affected row====>", result.RowsAffected)

	if result.RowsAffected == 0 {
		insertCarts := models.UsersCarts{
			ProductId: productId,
			UserId:    userId.(string),
		}

		database.DB.Table("user_carts").Save(&insertCarts)
	} else {
		log.Println("update case")
	}
}
