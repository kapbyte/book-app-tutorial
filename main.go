package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	routes "github.com/kapbyte/book-app-tutorial/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	routes.BookRoutes(router)

	router.Run(":" + port)
}
