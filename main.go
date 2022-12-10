package main

import (
	"fmt"

	"shopping_cart/database"

	"github.com/joho/godotenv"
)

func main() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	db := database.Connection()
	fmt.Println(db)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
}
