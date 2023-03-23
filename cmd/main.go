package main

import (
	"github.com/dissurender/hn-news/api"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	api.InitializeCache()

	r.GET("/api", api.HandleAPIRequest)
	r.GET("/api/:item", api.HandleItemRequest)

	r.Run() // listen and at "localhost:8080"
}
