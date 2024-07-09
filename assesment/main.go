package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	r := gin.Default()
	r.GET("/hotels", func(c *gin.Context) {
		GetHotels(c)
	})

	r.Run(":8080")
}