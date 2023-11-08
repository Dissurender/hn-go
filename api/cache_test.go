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
			t.Skip("Test cases not set yet")
			AddToCacheWithExpiration(tt.args.key, tt.args.value, tt.args.expiration)
		})
	}
}

func TestGetFromCache(t *testing.T) {
	InitializeCache()

	AddToCache("existingKey", "existingValue")
	AddToCacheWithExpiration("expiredKey", "expiredValue", 1*time.Second)

	time.Sleep(2 * time.Second)

	type args struct {
		key string
	}

	tests := []struct {
		name  string
		args  args
		want  interface{}
		want1 bool
	}{
		{
			name:  "Should retrieve existing item from the cache",
			args:  args{key: "existingKey"},
			want:  "existingValue",
			want1: true,
		},
		{
			name:  "Should fail to retrieve nil item from the cache",
			args:  args{key: "nonExistingKey"},
			want:  nil,
			want1: false,
		},
		{
			name:  "Should fail to retrieve exprired item from the cache",
			args:  args{key: "expiredKey"},
			want:  nil,
			want1: false,
		},
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
		{
			name: "",
			args: args{key: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("Test cases not set yet")
			AddToCache(tt.args.key, tt.args.value)
		})
	}
}
