package cache

import (
	"reflect"
	"testing"
)

func TestRelatedAPIVersionKeys(t *testing.T) {
	cases := map[string][]string{
		"https://api.zsapi.net/zpa/mgmtconfig/v2/admin/customers/1/segmentGroup/9": {"https://api.zsapi.net/zpa/mgmtconfig/v1/admin/customers/1/segmentGroup/9"},
		"https://api.zsapi.net/zpa/mgmtconfig/v1/admin/customers/1/segmentGroup/9": {"https://api.zsapi.net/zpa/mgmtconfig/v2/admin/customers/1/segmentGroup/9"},
		"https://api.zsapi.net/zia/api/v1/urlFilteringRules/5":                     nil,
	}
	for key, want := range cases {
		if got := RelatedAPIVersionKeys(key); !reflect.DeepEqual(got, want) {
			t.Errorf("RelatedAPIVersionKeys(%q) = %v, want %v", key, got, want)
		}
	}
}
