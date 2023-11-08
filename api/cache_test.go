package api

import (
	"reflect"
	"testing"
	"time"
)

func TestAddToCacheWithExpiration(t *testing.T) {
	type args struct {
		key        string
		value      interface{}
		expiration time.Duration
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddToCacheWithExpiration(tt.args.key, tt.args.value, tt.args.expiration)
		})
	}
}

func TestGetFromCache(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		args  args
		want  interface{}
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetFromCache(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFromCache() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetFromCache() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestAddToCache(t *testing.T) {
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddToCache(tt.args.key, tt.args.value)
		})
	}
}

func TestInitializeCache(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitializeCache()
		})
	}
}
