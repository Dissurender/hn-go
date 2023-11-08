package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/dissurender/hn-go/utils"

	"github.com/gin-gonic/gin"
)

// "https://hacker-news.firebaseio.com/v0/topstories.json"
// "https://hacker-news.firebaseio.com/v0/item/%d.json", id

func HandleAPIRequestBest(c *gin.Context) {
	cacheKey := "results"
	results, found := getResultsFromCache(cacheKey)
	if found {
		utils.Logger("Results found.")
		c.JSON(http.StatusOK, results)
		return
	}

	data, err := fetchTopStories()
	if err != nil {
		utils.Logger(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	results = retrieveKids(data)
	addResultsToCache(cacheKey, results)
	c.JSON(http.StatusOK, results)
}

// HandleItemRequest handles the request for an individual item.
func HandleItemRequest(c *gin.Context) {
	itemID := c.Param("item")
	cacheKey := fmt.Sprintf("story-%v", itemID)

	utils.Logger(fmt.Sprintf("GET: %v", cacheKey))

	responseData, found := getStoryFromCache(cacheKey)
	if found {
		c.JSON(http.StatusOK, responseData)
		return
	}

	responseDataWithKids, err := fetchStoryWithKids(itemID)
	if err != nil {
		utils.Logger(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	addStoryToCacheWithExpiration(cacheKey, responseDataWithKids, 5*time.Minute)
	c.JSON(http.StatusOK, responseDataWithKids)
}

/*

	--------------------------------------
	Helper Functions
	--------------------------------------
	Below are the helper functions used in handling
	API requests.

*/

func retrieveKids(data []int) []interface{} {
	var wg sync.WaitGroup
	results := make([]interface{}, len(data))

	// Mutex to protect concurrent writes to results.
	mu := new(sync.Mutex)

	for i, id := range data {
		wg.Add(1)
		go func(i, id int) {
			defer wg.Done()
			item, err := fetchOrRetrieveFromCache(id)
			if err != nil {
				return
			}
			// claim the mutex while mutating
			mu.Lock()
			results[i] = item
			mu.Unlock()
		}(i, id)
	}
	wg.Wait()

	utils.Logger(fmt.Sprintf("Number of results found: %v", len(results)))

	return results
}

// fetchOrRetrieveFromCache tries to get the item from the cache, or fetches from the API if not cached.
func fetchOrRetrieveFromCache(id int) (interface{}, error) {
	cacheKey := fmt.Sprintf("story-%v", id)
	if cachedResult, found := GetFromCache(cacheKey); found {
		return cachedResult, nil
	}

	// If not in cache, fetch from the API.
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseData Base
	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, err
	}

	// Cache the newly fetched result.
	AddToCache(cacheKey, responseData)
	return responseData, nil
}

func fetchTopStories() ([]int, error) {
	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data []int
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func getResultsFromCache(cacheKey string) ([]interface{}, bool) {
	cachedResult, found := GetFromCache(cacheKey)
	if !found {
		return nil, false
	}

	result, ok := cachedResult.([]interface{})
	if !ok {
		// TODO: Handle the situation where the result is not a []interface{}
		return nil, false
	}

	return result, true
}

func getStoryFromCache(cacheKey string) (BaseWithKids, bool) {
	cachedResult, found := GetFromCache(cacheKey)
	if !found {
		return BaseWithKids{}, false
	}

	responseData, ok := cachedResult.(BaseWithKids)
	if !ok {
		// TODO: result is not a BaseWithKids
		return BaseWithKids{}, false
	}

	return responseData, true
}

func addResultsToCache(cacheKey string, results []interface{}) {
	AddToCache(cacheKey, results)
}

func addStoryToCacheWithExpiration(cacheKey string, story BaseWithKids, expiration time.Duration) {
	AddToCacheWithExpiration(cacheKey, story, expiration)
}

// buildBaseWithKids constructs a BaseWithKids object from the provided Base and kids data.
func buildBaseWithKids(responseData Base, kidsData []interface{}) BaseWithKids {
	return BaseWithKids{
		ID:          responseData.ID,
		Type:        responseData.Type,
		By:          responseData.By,
		Time:        responseData.Time,
		Kids:        kidsData,
		Dead:        responseData.Dead,
		Deleted:     responseData.Deleted,
		Descendants: responseData.Descendants,
		Score:       responseData.Score,
		Title:       responseData.Title,
		URL:         responseData.URL,
	}
}

// fetchItem makes a GET request to the provided URL and unmarshals the response into Base.
func fetchItem(url string) (Base, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Base{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Base{}, err
	}

	var responseData Base
	if err := json.Unmarshal(body, &responseData); err != nil {
		return Base{}, err
	}

	return responseData, nil
}

// fetchStoryWithKids fetches a story and its kids from the HN API.
func fetchStoryWithKids(itemID string) (BaseWithKids, error) {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%v.json", itemID)
	responseData, err := fetchItem(url)
	if err != nil {
		return BaseWithKids{}, err
	}

	kidsData := retrieveKids(responseData.Kids)
	return buildBaseWithKids(responseData, kidsData), nil
}
