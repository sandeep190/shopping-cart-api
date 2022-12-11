package main

import (
	"fmt"
	"shopping_cart/controllers"
	"shopping_cart/database"

	"github.com/gin-contrib/cors"

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
	app := gin.Default()
	app.Use(cors.Default())
	apiRoutes := app.Group("/api")

	controllers.UsersRoutes(apiRoutes.Group("/users"))
	app.Run()
}
