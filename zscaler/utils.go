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
	var mapData map[string]interface{}
	err = json.Unmarshal(data, &mapData)
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
	// removes url prefix from oneapi to legacy api (/zia, /zpa, /zcc)
	if strings.HasPrefix(endpoint, "/zia") {
		return strings.TrimPrefix(endpoint, "/zia")
	}
	if strings.HasPrefix(endpoint, "/zpa") {
		return strings.TrimPrefix(endpoint, "/zpa")
	}
	if strings.HasPrefix(endpoint, "/zcc") {
		return strings.TrimPrefix(endpoint, "/zcc")
	}
	return endpoint
}
