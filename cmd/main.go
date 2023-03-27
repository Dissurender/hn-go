package main

import (
	"time"

	"github.com/dissurender/hn-news/api"
	"github.com/gin-gonic/gin"

	"github.com/itsjamie/gin-cors"
)

func main() {
	r := gin.Default()

	// Lock this to host machine in Prod
	r.Use(cors.Middleware(cors.Config{
		Origins:         "http://localhost",
		Methods:         "GET",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	api.InitializeCache()

	r.GET("/api/best", api.HandleAPIRequest)

	// Sort cached items and send to client
	r.GET("/api/new", api.HandleAPIRequest)
	r.GET("/api/top", api.HandleAPIRequest)

	r.GET("/api/:item", api.HandleItemRequest)

	r.Run() // listen and at "localhost:8080"
}
