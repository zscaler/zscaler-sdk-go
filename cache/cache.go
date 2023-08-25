package cache

import (
	"bytes"
	"io"
	"net/http"
)

type Cache interface {
	Get(key string) *http.Response
	Set(key string, value *http.Response)
	Delete(key string)
	Clear()
	ClearAllKeysWithPrefix(prefix string)
}

func CreateCacheKey(req *http.Request) string {
	s := req.URL.Scheme + "://" + req.URL.Host + req.URL.RequestURI()
	return s
}

func CopyResponse(resp *http.Response) *http.Response {
	c := *resp
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
	c.Body = io.NopCloser(bytes.NewBuffer(respBody))

	return &c
}
