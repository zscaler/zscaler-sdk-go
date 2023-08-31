package cache

import (
	"bufio"
	"bytes"
	"context"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/allegro/bigcache/v3"
)

type cache struct {
	bcache *bigcache.BigCache
}

func NewCache(ttl, cleanWindow time.Duration, maxCacheSizeMB int) (Cache, error) {
	bCache, err := bigcache.New(context.Background(), bigcache.Config{
		Shards:           16,
		LifeWindow:       ttl,
		CleanWindow:      cleanWindow,
		HardMaxCacheSize: maxCacheSizeMB,
		Verbose:          false,
		StatsEnabled:     false,
	})
	if err != nil {
		return nil, err
	}
	return &cache{bCache}, nil
}

func (c cache) Get(key string) *http.Response {
	item, err := c.bcache.Get(key)
	if err == nil && item != nil {
		r := bufio.NewReader(bytes.NewReader(item))
		resp, _ := http.ReadResponse(r, nil)
		return resp
	}

	return nil
}

func (c cache) Set(key string, value *http.Response) {
	cacheableResponse, _ := httputil.DumpResponse(value, true)

	c.bcache.Set(key, cacheableResponse)
}

func (c cache) Delete(key string) {
	c.bcache.Delete(key)
}

func (c cache) Clear() {
	c.bcache.Reset()
}

func (c cache) ClearAllKeysWithPrefix(prefix string) {
	it := c.bcache.Iterator()
	for it.SetNext() {
		e, err := it.Value()
		if err != nil {
			continue
		}
		if strings.HasPrefix(e.Key(), prefix) {
			c.bcache.Delete(e.Key())
		}
	}
}
