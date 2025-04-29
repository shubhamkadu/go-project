package main

import (
	"GoProject/config"
	"GoProject/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	config.ConnectDB()
	config.ConnectRedis()

	r := gin.Default()
	r.Use(cors.Default())

	routes.SetupRoutes(r)

	log.Println("Server started at :8080")
	r.Run(":8080")
}
