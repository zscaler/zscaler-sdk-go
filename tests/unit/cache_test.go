package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/cache"
	"github.com/stretchr/testify/assert"
)

func Test_cache_key_can_be_created_from_request_object(t *testing.T) {
	var buff io.ReadWriter
	request, _ := http.NewRequest("GET", "https://example.com/sample/cache-key/test+test@test."+
		"com?with=a&query=string",
		buff)

	cacheKey := cache.CreateCacheKey(request)

	assert.Equal(t, "https://example.com/sample/cache-key/test+test@test.com?with=a&query=string", cacheKey,
		"The cache key was not created correctly.")
}

func Test_an_item_can_be_stored_in_cache(t *testing.T) {
	log.Println("Starting Test_an_item_can_be_stored_in_cache")

	url := "https://example.com/sample/cache-key/"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	log.Println("Created HTTP request")

	cacheKey := cache.CreateCacheKey(request)
	log.Printf("Generated cache key: %s", cacheKey)

	myCache, err := cache.NewCache(30*time.Minute, 30*time.Minute, 100)
	if err != nil {
		log.Fatalf("Failed to create cache: %v", err)
	}
	log.Println("Created cache instance")

	if myCache.Get(cacheKey) != nil {
		t.Errorf("Item already existed in cache")
	}
	log.Println("Checked cache for existing item")

	toCache := "test Item"
	resp := httptest.NewRecorder().Result()
	resp.Body = io.NopCloser(strings.NewReader(toCache))
	log.Println("Created mock HTTP response")

	myCache.Set(cacheKey, resp)
	log.Println("Cached the response")

	pulledFromCache := myCache.Get(cacheKey)
	if pulledFromCache == nil {
		t.Fatalf("Item was not found in cache after setting")
	}
	log.Println("Retrieved item from cache")

	cachedBody, err := ioutil.ReadAll(pulledFromCache.Body)
	if err != nil {
		t.Fatalf("Failed to read body from cached response: %v", err)
	}
	log.Printf("Read body from cached response: %s", string(cachedBody))

	if string(cachedBody) != toCache {
		t.Errorf("Item pulled from cache was not correct, expected '%s', got '%s'", toCache, string(cachedBody))
	}

	log.Println("Completed Test_an_item_can_be_stored_in_cache")
}

func Test_an_item_can_be_deleted_from_cache(t *testing.T) {
	url := "https://example.com/sample/cache-key/delete"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	cacheKey := cache.CreateCacheKey(request)

	myCache, err := cache.NewCache(30*time.Minute, 30*time.Minute, 100)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	toCache := "test Item"
	resp := httptest.NewRecorder().Result()
	resp.Body = io.NopCloser(bytes.NewBufferString(toCache))

	myCache.Set(cacheKey, resp)

	if myCache.Get(cacheKey) == nil {
		t.Errorf("Item does not exist in cache after setting")
	}

	myCache.Delete(cacheKey)

	if myCache.Get(cacheKey) != nil {
		t.Errorf("Item was not deleted from cache")
	}
}

func Test_cache_can_be_cleared(t *testing.T) {
	url := "https://example.com/sample/cache-key/clear"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	cacheKey := cache.CreateCacheKey(request)

	// Creating cache instance
	myCache, err := cache.NewCache(30*time.Minute, 30*time.Minute, 100)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	toCache := "test Item"
	resp := httptest.NewRecorder().Result()
	resp.Body = io.NopCloser(bytes.NewBufferString(toCache))

	myCache.Set(cacheKey, resp)

	if myCache.Get(cacheKey) == nil {
		t.Errorf("Item does not exist in cache after setting")
	}

	myCache.Clear()

	if myCache.Get(cacheKey) != nil {
		t.Errorf("Cache was not cleared")
	}
}
