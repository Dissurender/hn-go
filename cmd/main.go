package main

// "https://hacker-news.firebaseio.com/v0/topstories.json"
// "https://hacker-news.firebaseio.com/v0/item/%d.json", id

import (
	"github.com/dissurender/hn-news/api"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	api.InitializeCache()

	r.GET("/api", api.HandleAPIRequest)
	r.GET("/api/:item", api.HandleItemRequest)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
