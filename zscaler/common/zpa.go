package common

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

func getMicrotenantIDFromBody(body interface{}) string {
	if body == nil {
		return ""
	}

	d, err := json.Marshal(body)
	if err != nil {
		return ""
	}
	dataMap := map[string]interface{}{}
	err = json.Unmarshal(d, &dataMap)
	if err != nil {
		return ""
	}
	if microTenantID, ok := dataMap["microtenantId"]; ok && microTenantID != nil && microTenantID != "" {
		return fmt.Sprintf("%v", microTenantID)
	}
	return ""
}

func getMicrotenantIDFromEnvVar() string {
	return os.Getenv("ZPA_MICROTENANT_ID")
}

func InjectMicrotentantID(body interface{}, q url.Values, microtenantIDFromConfig string) url.Values {
	if q.Has("microtenantId") && q.Get("microtenantId") != "" {
		return q
	}

	microTenantID := getMicrotenantIDFromBody(body)
	if microTenantID != "" {
		q.Add("microtenantId", microTenantID)
		return q
	}

	microTenantID = getMicrotenantIDFromEnvVar()
	if microTenantID != "" {
		q.Add("microtenantId", microTenantID)
		return q
	}

	microTenantID = microtenantIDFromConfig
	if microTenantID != "" {
		q.Add("microtenantId", microTenantID)
		return q
	}

	return q
}
