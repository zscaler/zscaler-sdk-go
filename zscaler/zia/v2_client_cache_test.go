package zia

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/cache"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
)

// Same interleaving as the OneAPI client test
// (TestExecuteRequest_ReadOverlappingWriteIsNotCached): a list read that
// reaches the API while a child PUT is being processed returns the list from
// before the PUT, and must not be cached after the PUT invalidated the cache.
func TestLegacyDo_ReadOverlappingWriteIsNotCached(t *testing.T) {
	const list = "/api/v1/urlCategories"
	const item = "/api/v1/urlCategories/CUSTOM_14"
	var (
		gets       int32
		updated    atomic.Bool
		putArrived = make(chan struct{})
		listServed = make(chan struct{})
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == list:
			atomic.AddInt32(&gets, 1)
			if updated.Load() {
				_, _ = w.Write([]byte(`[{"id":"CUSTOM_14","urls":[]}]`))
			} else {
				_, _ = w.Write([]byte(`[{"id":"CUSTOM_14","urls":["before.example.com"]}]`))
			}
		case r.Method == http.MethodPut && r.URL.Path == item:
			close(putArrived)
			<-listServed
			updated.Store(true)
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	cch, err := cache.NewCache(10*time.Minute, 8*time.Minute, 1)
	if err != nil {
		t.Fatal(err)
	}
	c := &Client{
		session:          &Session{JSessionID: "test"},
		sessionRefreshed: time.Now(),
		sessionTimeout:   time.Hour,
		HTTPClient:       srv.Client(),
		Logger:           logger.GetDefaultLogger("zia-cache-test: "),
		cacheEnabled:     true,
		cache:            cch,
	}
	do := func(method, path string) string {
		t.Helper()
		req, err := http.NewRequest(method, srv.URL+path, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := c.do(req, time.Now(), "test")
		if err != nil {
			t.Fatalf("%s %s: %v", method, path, err)
		}
		body, _ := io.ReadAll(resp.Body)
		return string(body)
	}

	putDone := make(chan struct{})
	go func() {
		defer close(putDone)
		do(http.MethodPut, item)
	}()
	<-putArrived
	if body := do(http.MethodGet, list); !strings.Contains(body, "before.example.com") {
		t.Fatalf("overlapping read = %s, want the list from before the PUT", body)
	}
	close(listServed)
	<-putDone

	if body := do(http.MethodGet, list); strings.Contains(body, "before.example.com") {
		t.Fatalf("read after the PUT was served the overlapping read's response: %s", body)
	}
	if n := atomic.LoadInt32(&gets); n != 2 {
		t.Fatalf("list requests = %d, want 2 (the overlapping read must not be cached)", n)
	}

	// A read that does not overlap a write is cached.
	do(http.MethodGet, list)
	if n := atomic.LoadInt32(&gets); n != 2 {
		t.Fatalf("list requests = %d, want 2 (the non-overlapping read must be cached)", n)
	}
}
