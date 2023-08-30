package cache

import (
	"net/http"
)

type nopCache struct {
}

func NewNopCache() Cache {
	return &nopCache{}
}

func (c nopCache) Get(key string) *http.Response {
	return nil
}

func (c nopCache) Set(key string, value *http.Response) {
}

func (c nopCache) Delete(key string) {
}

func (c nopCache) Clear() {
}

func (c nopCache) ClearAllKeysWithPrefix(prefix string) {
}
