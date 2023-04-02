package api

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHandleAPIRequest(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleAPIRequestBest(tt.args.c)
		})
	}
}

func TestHandleItemRequest(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleItemRequest(tt.args.c)
		})
	}
}
