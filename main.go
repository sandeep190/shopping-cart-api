package main

import (
	"fmt"
	"os"
	"shopping_cart/controllers"
	"shopping_cart/database"
	"shopping_cart/middleware"
	"shopping_cart/models"

	"github.com/gin-contrib/cors"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	db := database.Connection()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	gin.SetMode(os.Getenv("GIN_MODE"))
	makeMigration(db)
	app := gin.Default()
	app.Use(cors.Default())
	app.Use(middleware.JwtTokenVerify())
	apiRoutes := app.Group("/api")

	controllers.UsersRoutes(apiRoutes.Group("/users"))
	controllers.CategoryRoutes(apiRoutes.Group("category"))
	app.Static("/assets", "./assets")
	app.Static("/static", "./static")
	app.LoadHTMLGlob("templates/*/*.html")

	controllers.PagesRoutes(&app.RouterGroup)
	controllers.AdminRoutes(&app.RouterGroup)
	app.Run(":" + os.Getenv("APPPORT"))
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.Comment{})
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.ProductCategory{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.FileUpload{})
}

func makeMigration(db *gorm.DB) {
	cmd := os.Args
	if len(cmd) > 1 {
		if cmd[1] == "make_migration" {
			migrate(db)
		}
	}
}
