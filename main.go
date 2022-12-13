package main

import (
	"fmt"
	"os"
	"shopping_cart/controllers"
	"shopping_cart/database"
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
	gin.SetMode(gin.ReleaseMode)
	makeMigration(db)
	app := gin.Default()
	app.Use(cors.Default())
	apiRoutes := app.Group("/api")

	controllers.UsersRoutes(apiRoutes.Group("/users"))
	app.Run()
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
