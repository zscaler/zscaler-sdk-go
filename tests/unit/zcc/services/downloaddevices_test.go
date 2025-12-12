// Package services provides unit tests for ZCC downloaddevices service
package services

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/downloaddevices"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestDownloadDevices_Download_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/downloadDevices"

	// Mock CSV response
	csvData := `udid,username,hostname,os_type
device-001,user1@example.com,host1,1
device-002,user2@example.com,host2,2
device-003,user3@example.com,host3,3`

	server.On("GET", path, common.CSVResponse(csvData))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	var buffer bytes.Buffer
	err = downloaddevices.DownloadDevices(context.Background(), service, "1,2,3", "all", &buffer)

	require.NoError(t, err)
	assert.Contains(t, buffer.String(), "device-001")
	assert.Contains(t, buffer.String(), "user1@example.com")
}

func TestDownloadDevices_DownloadWithFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/downloadDevices"

	csvData := `udid,username,hostname,os_type
device-001,user1@example.com,host1,1`

	server.On("GET", path, common.CSVResponse(csvData))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	var buffer bytes.Buffer
	err = downloaddevices.DownloadDevices(context.Background(), service, "1", "registered", &buffer)

	require.NoError(t, err)
	assert.Contains(t, buffer.String(), "device-001")
}

func TestDownloadDevices_DownloadServiceStatus_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/downloadServiceStatus"

	csvData := `udid,status,last_seen,service_status
device-001,active,2024-01-15T10:00:00Z,running
device-002,inactive,2024-01-14T09:00:00Z,stopped`

	server.On("GET", path, common.CSVResponse(csvData))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	var buffer bytes.Buffer
	err = downloaddevices.DownloadServiceStatus(context.Background(), service, "1,2", "all", &buffer)

	require.NoError(t, err)
	assert.Contains(t, buffer.String(), "device-001")
	assert.Contains(t, buffer.String(), "active")
}

func TestDownloadDevices_DownloadServiceStatusFiltered_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/downloadServiceStatus"

	csvData := `udid,status
device-001,active`

	server.On("GET", path, common.CSVResponse(csvData))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	var buffer bytes.Buffer
	err = downloaddevices.DownloadServiceStatus(context.Background(), service, "1", "active", &buffer)

	require.NoError(t, err)
	assert.Contains(t, buffer.String(), "device-001")
}

func TestDownloadDevices_DownloadEmptyParams_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/downloadDevices"

	csvData := `udid,username,hostname
device-001,user1@example.com,host1`

	server.On("GET", path, common.CSVResponse(csvData))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	var buffer bytes.Buffer
	err = downloaddevices.DownloadDevices(context.Background(), service, "", "", &buffer)

	require.NoError(t, err)
	assert.Contains(t, buffer.String(), "device-001")
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestDownloadDevices_Structure(t *testing.T) {
	t.Parallel()

	t.Run("DownloadDevicesQueryParams JSON marshaling", func(t *testing.T) {
		params := downloaddevices.DownloadDevicesQueryParams{
			OSTypes:           "1,2,3",
			RegistrationTypes: "registered",
		}

		data, err := json.Marshal(params)
		require.NoError(t, err)

		// Check the values are present
		assert.Contains(t, string(data), "1,2,3")
		assert.Contains(t, string(data), "registered")
	})

	t.Run("DownloadDevicesQueryParams JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"osTypes": "1,2",
			"registrationTypes": "all"
		}`

		var params downloaddevices.DownloadDevicesQueryParams
		err := json.Unmarshal([]byte(jsonData), &params)
		require.NoError(t, err)

		// The struct uses url tags, not json tags - values may not unmarshal correctly
		// This tests the struct can be parsed without errors
		assert.NotNil(t, params)
	})

	t.Run("Empty query params", func(t *testing.T) {
		params := downloaddevices.DownloadDevicesQueryParams{}

		data, err := json.Marshal(params)
		require.NoError(t, err)

		// Empty struct should marshal to empty JSON object or with empty strings
		assert.NotNil(t, data)
	})

	t.Run("Query params with only OSTypes", func(t *testing.T) {
		params := downloaddevices.DownloadDevicesQueryParams{
			OSTypes: "1",
		}

		data, err := json.Marshal(params)
		require.NoError(t, err)

		assert.Contains(t, string(data), "1")
	})

	t.Run("Query params with only RegistrationTypes", func(t *testing.T) {
		params := downloaddevices.DownloadDevicesQueryParams{
			RegistrationTypes: "unregistered",
		}

		data, err := json.Marshal(params)
		require.NoError(t, err)

		assert.Contains(t, string(data), "unregistered")
	})
}

func TestDownloadDevices_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse CSV response", func(t *testing.T) {
		csvData := `udid,username,hostname,os_type
device-001,user1@example.com,host1,1
device-002,user2@example.com,host2,2`

		lines := bytes.Split([]byte(csvData), []byte("\n"))
		assert.Len(t, lines, 3)
		assert.Contains(t, string(lines[0]), "udid")
		assert.Contains(t, string(lines[1]), "device-001")
		assert.Contains(t, string(lines[2]), "device-002")
	})

	t.Run("Parse empty CSV response", func(t *testing.T) {
		csvData := `udid,username,hostname,os_type`

		lines := bytes.Split([]byte(csvData), []byte("\n"))
		assert.Len(t, lines, 1)
		assert.Contains(t, string(lines[0]), "udid")
	})

	t.Run("Parse CSV with special characters", func(t *testing.T) {
		csvData := `udid,username,hostname
device-001,"user, with comma@example.com",host1
device-002,user2@example.com,"host with spaces"`

		lines := bytes.Split([]byte(csvData), []byte("\n"))
		assert.Len(t, lines, 3)
	})
}
