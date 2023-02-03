package common

import (
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v1/zia"
)

const pageSize = 1000

type IDNameExtensions struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type IDExtensions struct {
	ID         int                    `json:"id,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type UserGroups struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	IdpID    int    `json:"idp_id,omitempty"`
	Comments string `json:"comments,omitempty"`
}

type UserDepartment struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	IdpID    int    `json:"idp_id,omitempty"`
	Comments string `json:"comments,omitempty"`
	Deleted  bool   `json:"deleted,omitempty"`
}

type DeviceGroups struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Devices struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func ReadAllPages[T any](client *zia.Client, endpoint string, list *[]T) error {
	if list == nil {
		return nil
	}
	page := 1
	if !strings.Contains(endpoint, "?") {
		endpoint += "?"
	}

	for {
		pageItems := []T{}
		err := client.Read(fmt.Sprintf("%s&pageSize=%d&page=%d", endpoint, pageSize, page), &pageItems)
		if err != nil {
			return err
		}
		*list = append(*list, pageItems...)
		if len(pageItems) < pageSize {
			break
		}
		page++
	}
	return nil
}
