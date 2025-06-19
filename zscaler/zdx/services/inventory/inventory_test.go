package inventory

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
)

func TestGetSoftware(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Define a time filter for the last 2 hours
	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix()
	to := now.Unix()
	filters := GetSoftwareFilters{
		GetFromToFilters: common.GetFromToFilters{
			From: int(from),
			To:   int(to),
		},
	}

	// Call GetSoftware with the filters
	softwareList, nextOffset, resp, err := GetSoftware(context.Background(), service, filters)
	if err != nil {
		t.Fatalf("Error getting software: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(softwareList) == 0 {
		t.Log("No software found.")
	} else {
		t.Logf("Retrieved %d software entries", len(softwareList))
		for _, software := range softwareList {
			t.Logf("Software Key: %s, Software Name: %s, Vendor: %s", software.SoftwareKey, software.SoftwareName, software.Vendor)
		}
	}

	t.Logf("Next Offset: %s", nextOffset)
}

func TestGetSoftwareKey(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Define a time filter for the last 2 hours
	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix()
	to := now.Unix()
	filters := GetSoftwareFilters{
		GetFromToFilters: common.GetFromToFilters{
			From: int(from),
			To:   int(to),
		},
	}

	// Call GetSoftware with the filters to get the software key
	softwareList, _, _, err := GetSoftware(context.Background(), service, filters)
	if err != nil {
		t.Fatalf("Error getting software: %v", err)
	}

	if len(softwareList) == 0 {
		t.Log("No software found, skipping GetSoftwareKey test.")
		return
	}

	softwareKey := softwareList[0].SoftwareKey

	// Call GetSoftwareKey with the software key
	softwareKeyList, nextOffset, resp, err := GetSoftwareKey(context.Background(), service, softwareKey, filters)
	if err != nil {
		t.Fatalf("Error getting software key: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(softwareKeyList) == 0 {
		t.Log("No software key entries found.")
	} else {
		t.Logf("Retrieved %d software key entries", len(softwareKeyList))
		for _, software := range softwareKeyList {
			t.Logf("Software Key: %s, Software Name: %s, Vendor: %s, User ID: %d", software.SoftwareKey, software.SoftwareName, software.Vendor, software.UserID)
		}
	}

	t.Logf("Next Offset: %s", nextOffset)
}
