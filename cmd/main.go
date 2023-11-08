package main

import (
	"time"

	"github.com/dissurender/hn-go/api"
	"github.com/gin-gonic/gin"

	"github.com/itsjamie/gin-cors"
)

func main() {
	r := gin.Default()

	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	api.InitializeCache()

	r.GET("/api/best", api.HandleAPIRequestBest)
	r.GET("/api/:item", api.HandleItemRequest)

	// Sort cached items and send to client
	// r.GET("/api/new", api.HandleAPIRequest)
	// r.GET("/api/top", api.HandleAPIRequest)

	r.Run(":8888") // listen and at "localhost:8080"
}
