package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const DB_PATH = "data"

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	err = os.MkdirAll(DB_PATH, os.ModePerm)
	if err != nil {
		log.Fatalf(err.Error())
	}

	router := gin.Default()

	user := router.Group("/user")
	{
		user.POST("/create", SignUp)
		user.POST("/login", Login)
	}
	router.GET("/btcRate", AuthMiddleware(), btcRate)
	router.Run(":8081")
}
