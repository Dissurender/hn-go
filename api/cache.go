package api

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var c *cache.Cache

func InitializeCache() {
	// Initialize the cache with a 5 minute default expiration time and a 10 minute cleanup interval
	c = cache.New(5*time.Minute, 10*time.Minute)
}

func GetFromCache(key string) (interface{}, bool) {
	return c.Get(key)
}

func AddToCache(key string, value interface{}) {
	c.Set(key, value, cache.DefaultExpiration)
}

func AddToCacheWithExpiration(key string, value interface{}, expiration time.Duration) {
	c.Set(key, value, expiration)
}
