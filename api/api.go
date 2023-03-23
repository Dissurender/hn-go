package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleAPIRequest(c *gin.Context) {

	// Check if the results are already cached
	cacheKey := "results"
	cachedResult, found := GetFromCache(cacheKey)
	if found {
		// If the results are cached, return them directly
		result, ok := cachedResult.([]interface{})

		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cached result is of invalid type"})
			return
		}

		// Write the cached results as the response
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, result)
		return
	}

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

			// Check if the result for this ID is already cached
			cacheKey := fmt.Sprintf("story-%v", id)
			cachedResult, found := GetFromCache(cacheKey)
			if found {
				// If the result is cached, add it to the result slice
				results[i] = cachedResult
			} else {
				// If the result is not cached, make the API request and cache the result
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

				// Cache the result
				fmt.Println("Added key to cache:", cacheKey)
				AddToCache(cacheKey, responseData)

				results[i] = responseData
			}
		}(i, id)
	}
	wg.Wait()

	// Cache the results and return them as the response
	AddToCache(cacheKey, results)
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, results)
}

func HandleItemRequest(c *gin.Context) {
	item := c.Param("item")

	// Check if the result for this item is already cached
	cacheKey := fmt.Sprintf("story-%v", item)
	cachedResult, found := GetFromCache(cacheKey)
	if found {
		// If the result is cached, return it directly
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, cachedResult)
		return
	}

	// Make a request to the API with the single item parameter
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%v.json", item)
	resp, err := http.Get(url)
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

	// Unmarshal the response body into an interface{}
	var responseData interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Cache the result with a 5 minute expiration time
	AddToCacheWithExpiration(cacheKey, responseData, 5*time.Minute)

	// Write the result as the response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, responseData)
}
