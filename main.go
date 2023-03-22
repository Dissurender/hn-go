package main

// "https://hacker-news.firebaseio.com/v0/topstories.json"
// "https://hacker-news.firebaseio.com/v0/item/%d.json", id

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/api", func(c *gin.Context) {
		// Make a request to another API
		resp, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Unmarshal the response body into a slice of integers
		var data []int
		err = json.Unmarshal(body, &data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Make additional requests to the same API with each integer as an ID parameter concurrently
		var wg sync.WaitGroup
		results := make([]interface{}, len(data))
		for i, id := range data {
			wg.Add(1)
			go func(i int, id int) {
				defer wg.Done()

				url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println("Error making request to API:", err)
					return
				}
				defer resp.Body.Close()

				// Read the response body
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("Error reading response body:", err)
					return
				}

				// Unmarshal the response body into an interface{}
				var responseData interface{}
				err = json.Unmarshal(body, &responseData)
				if err != nil {
					fmt.Println("Error unmarshalling response body:", err)
					return
				}

				results[i] = responseData
			}(i, id)
		}
		wg.Wait()

		// Write the mutated results slice to a local file
		file, err := ioutil.TempFile("", "results-*.json")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		err = json.NewEncoder(file).Encode(results)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Write the mutated results slice as the response
		c.Header("Content-Type", "application/json")
		c.File(file.Name())
	})

	// Start the server
	fmt.Println("Listening on :8080")
	r.Run(":8080")
}
