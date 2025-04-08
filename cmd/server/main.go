package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// .envの読み込み
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello from Docker-free Go! 🙌")
	})

	r.Run(":" + port)
}
