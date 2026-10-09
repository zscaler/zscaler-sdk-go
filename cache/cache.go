package cache

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

type Cache interface {
	Get(key string) *http.Response
	Set(key string, value *http.Response)
	Delete(key string)
	Clear()
	ClearAllKeysWithPrefix(prefix string)
	Close()
}

func CreateCacheKey(req *http.Request) string {
	s := req.URL.Scheme + "://" + req.URL.Host + req.URL.RequestURI()
	return s
}

// zpaMgmtConfigVersions are the ZPA management API path segments that address
// the same objects through different API versions.
var zpaMgmtConfigVersions = [][2]string{
	{"/mgmtconfig/v1/", "/mgmtconfig/v2/"},
	{"/mgmtconfig/v2/", "/mgmtconfig/v1/"},
}

// RelatedAPIVersionKeys returns the cache key prefixes of the same ZPA
// management API object under its other API version. Some objects are read
// through one version and written through another (e.g. segment groups: GET
// /mgmtconfig/v1/..., PUT /mgmtconfig/v2/...), so invalidating only the path
// that was written would leave the other version's cached reads stale.
func RelatedAPIVersionKeys(key string) []string {
	var keys []string
	for _, v := range zpaMgmtConfigVersions {
		if strings.Contains(key, v[0]) {
			keys = append(keys, strings.Replace(key, v[0], v[1], 1))
		}
	}
	return keys
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
