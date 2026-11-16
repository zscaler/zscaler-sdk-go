package zscaler

import (
	"bytes"
	"encoding/json"
	"html"
	"strings"
)

func decodeJSON(respData []byte, v interface{}) error {
	return json.NewDecoder(bytes.NewBuffer(respData)).Decode(&v)
}

func unescapeHTML(entity interface{}) {
	if entity == nil {
		return
	}
	data, err := json.Marshal(entity)
	if err != nil {
		return
	}
	// UseNumber keeps numbers exact: as float64, integers above 2^53 (such as
	// 18-digit IDs) would be rounded when the entity is decoded again below.
	var mapData map[string]interface{}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	err = dec.Decode(&mapData)
	if err != nil {
		return
	}
	for _, field := range []string{"name", "description"} {
		if v, ok := mapData[field]; ok && v != nil {
			str, ok := v.(string)
			if ok {
				mapData[field] = html.UnescapeString(html.UnescapeString(str))
			}
		}
	}
	data, err = json.Marshal(mapData)
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, entity)
}

func removeOneApiEndpointPrefix(endpoint string) string {
	// removes url prefix from oneapi to legacy api (/zia, /zpa, /zcc, /zdx)
	if strings.HasPrefix(endpoint, "/zia") {
		return strings.TrimPrefix(endpoint, "/zia")
	}
	if strings.HasPrefix(endpoint, "/ztw/api/v1") {
		return strings.TrimPrefix(endpoint, "/ztw/api/v1")
	}
	if strings.HasPrefix(endpoint, "/zpa") {
		return strings.TrimPrefix(endpoint, "/zpa")
	}
	if strings.HasPrefix(endpoint, "/zcc/papi") {
		return strings.TrimPrefix(endpoint, "/zcc/papi")
	}
	if strings.HasPrefix(endpoint, "/zcc") {
		return strings.TrimPrefix(endpoint, "/zcc")
	}
	if strings.HasPrefix(endpoint, "/zdx") {
		return strings.TrimPrefix(endpoint, "/zdx")
	}
	return endpoint
}

func Difference(slice1 []string, slice2 []string) []string {
	var diff []string
	for _, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if s1 == s2 {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, s1)
		}
	}
	return diff
}
