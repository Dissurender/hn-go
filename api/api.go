package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// "https://hacker-news.firebaseio.com/v0/topstories.json"
// "https://hacker-news.firebaseio.com/v0/item/%d.json", id

func HandleAPIRequestBest(c *gin.Context) {

	// Check if the results are already cached
	cacheKey := "results"
	cachedResult, found := GetFromCache(cacheKey)
	if found {
		// If the results are cached, return them
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

	/*
	 * Look into refactoring the primary GET request into a helper function
	 * for use throughout API
	 */

	// Make a request to HN API
	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Unmarshal the response body into a slice of int
	var data []int
	err = json.Unmarshal(body, &data)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Make additional requests to the HN API with each integer as an ID parameter concurrently
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
				// If the result is not cached, make the HN API request and cache the result
				url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
				resp, err := http.Get(url)

				if err != nil {
					fmt.Println("Error making request to API:", err)
					return
				}
				defer resp.Body.Close()

				// Read the response body
				body, err := io.ReadAll(resp.Body)

				if err != nil {
					fmt.Println("Error reading response body:", err)
					return
				}

				// Unmarshal the response body into the Base model
				var responseData Base
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Unmarshal the response body into an interface{}
	var responseData Base
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	buildComments := retrieveKids(c, responseData.Kids)

	responseDataWithKids := BaseWithKids{
		ID:          responseData.ID,
		Type:        responseData.Type,
		By:          responseData.By,
		Time:        responseData.Time,
		Kids:        buildComments,
		Dead:        responseData.Dead,
		Deleted:     responseData.Deleted,
		Descendants: responseData.Descendants,
		Score:       responseData.Score,
		Title:       responseData.Title,
		URL:         responseData.URL,
	}

	// Cache the result with a 5 minute expiration time
	AddToCacheWithExpiration(cacheKey, responseData, 5*time.Minute)

	// Write the result as the response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, responseDataWithKids)
}

func retrieveKids(c *gin.Context, kids []int) []interface{} {

	comments := make([]interface{}, len(kids))

	// Check if the result for this item is already cached
	// Make additional requests to the HN API with each integer as an ID parameter concurrently
	var wg sync.WaitGroup
	for i, kid := range kids {
		wg.Add(1)
		go func(i int, kid int) {
			defer wg.Done()

			// Check if the result for this ID is already cached
			cacheKey := fmt.Sprintf("story-%v", kid)
			cachedResult, found := GetFromCache(cacheKey)
			if found {
				// If the result is cached, add it to the result slice
				comments[i] = cachedResult
			} else {
				// If the result is not cached, make the HN API request and cache the result
				url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", kid)
				resp, err := http.Get(url)

				if err != nil {
					fmt.Println("Error making request to API:", err)
					return
				}
				defer resp.Body.Close()

				// Read the response body
				body, err := io.ReadAll(resp.Body)

				if err != nil {
					fmt.Println("Error reading response body:", err)
					return
				}

				// Unmarshal the response body into the Base model
				var responseData Comment
				err = json.Unmarshal(body, &responseData)

				if err != nil {
					fmt.Println("Error unmarshalling response body:", err)
					return
				}

				comments[i] = responseData
			}
		}(i, kid)
	}
	wg.Wait()

	return comments
}
