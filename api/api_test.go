package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHandleAPIRequest(t *testing.T) {
	// Create a new Gin router
	router := gin.Default()

	// Add the HandleAPIRequest function to the router
	router.GET("/api", HandleAPIRequest)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP response recorder
	w := httptest.NewRecorder()

	// Serve the HTTP request using the Gin router
	router.ServeHTTP(w, req)

	// Check the status code of the response
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the content type of the response
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	// Check the body of the response
	var responseBody []interface{}
	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Len(t, responseBody, 10)
}

func TestHandleItemRequest(t *testing.T) {
	// Create a new Gin router
	router := gin.Default()

	// Add the HandleItemRequest function to the router
	router.GET("/api/:item", HandleItemRequest)

	// Create a new HTTP request with the "item" parameter set to "12345"
	req, err := http.NewRequest("GET", "/api/35276032", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP response recorder
	w := httptest.NewRecorder()

	// Serve the HTTP request using the Gin router
	router.ServeHTTP(w, req)

	// Check the status code of the response
	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d but got %d", http.StatusOK, w.Code)
	}

	// Check the content type of the response
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content type %s but got %s", "application/json", contentType)
	}

	// Check the body of the response
	expectedBody := `{"by":"maze-le","descendants":0,"id":35276032,"score":3,"time":1679584159,"title":"How to Polis, 101, Part IIa: Politeia in the Polis","type":"story","url":"https://acoup.blog/2023/03/17/collections-how-to-polis-part-iia-politeia-in-the-polis/"}`
	if w.Body.String() != expectedBody {
		t.Errorf("expected body %s but got %s", expectedBody, w.Body.String())
	}
}
