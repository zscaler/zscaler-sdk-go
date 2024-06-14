package lssconfigcontroller

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
)

func TestGetAllFormats(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Failed to create ZPA client: %v", err)
	}
	service := services.New(client)

	// List of logTypes to be tested
	logTypes := []string{
		"zpn_ast_comprehensive_stats",
		"zpn_auth_log",
		"zpn_pbroker_comprehensive_stats",
		"zpn_ast_auth_log",
		"zpn_audit_log",
		"zpn_trans_log",
		"zpn_http_trans_log",
		"zpn_waf_http_exchanges_log",
	}

	// Expected formats
	expectedFormats := map[string]string{
		"Csv":  "csv",
		"Tsv":  "tsv",
		"Json": "json",
	}

	// Iterate through each logType and test
	for _, logType := range logTypes {
		formats, resp, err := GetFormats(service, logType)
		if err != nil {
			t.Errorf("Failed to get formats for logType %s: %v", logType, err)
			continue
		}
		if resp.StatusCode != 200 {
			t.Errorf("For logType %s, expected status code 200, got %d", logType, resp.StatusCode)
			continue
		}

		// Validate response for non-empty formats
		for formatName := range expectedFormats {
			var actualValue string
			switch formatName {
			case "Csv":
				actualValue = formats.Csv
			case "Tsv":
				actualValue = formats.Tsv
			case "Json":
				actualValue = formats.Json
			default:
				t.Errorf("Unknown format: %s", formatName)
				continue
			}

			if actualValue == "" {
				t.Errorf("For logType %s and format %s, received an empty response", logType, formatName)
			}
		}
	}
}
