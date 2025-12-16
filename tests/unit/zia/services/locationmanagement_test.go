// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/location/locationmanagement"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestLocationManagement_GetLocation_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	locationID := 12345
	path := "/zia/api/v1/locations/12345"

	server.On("GET", path, common.SuccessResponse(locationmanagement.Locations{
		ID:          locationID,
		Name:        "HQ Office",
		Description: "Headquarters location",
		Country:     "US",
		State:       "CA",
		TZ:          "America/Los_Angeles",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := locationmanagement.GetLocation(context.Background(), service, locationID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, locationID, result.ID)
	assert.Equal(t, "HQ Office", result.Name)
	assert.Equal(t, "US", result.Country)
}

func TestLocationManagement_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/locations"

	server.On("GET", path, common.SuccessResponse([]locationmanagement.Locations{
		{ID: 1, Name: "Location 1", Country: "US"},
		{ID: 2, Name: "Location 2", Country: "UK"},
		{ID: 3, Name: "Location 3", Country: "DE"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := locationmanagement.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 3)
}

func TestLocationManagement_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/locations"

	server.On("POST", path, common.SuccessResponse(locationmanagement.Locations{
		ID:          99999,
		Name:        "New Office",
		Description: "New branch office",
		Country:     "US",
		State:       "NY",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newLocation := &locationmanagement.Locations{
		Name:        "New Office",
		Description: "New branch office",
		Country:     "US",
		State:       "NY",
	}

	result, err := locationmanagement.Create(context.Background(), service, newLocation)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
	assert.Equal(t, "New Office", result.Name)
}

func TestLocationManagement_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	locationID := 12345
	path := "/zia/api/v1/locations/12345"

	server.On("PUT", path, common.SuccessResponse(locationmanagement.Locations{
		ID:          locationID,
		Name:        "Updated Office",
		Description: "Updated description",
		Country:     "US",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateLocation := &locationmanagement.Locations{
		ID:          locationID,
		Name:        "Updated Office",
		Description: "Updated description",
	}

	result, _, err := locationmanagement.Update(context.Background(), service, locationID, updateLocation)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Office", result.Name)
}

func TestLocationManagement_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	locationID := 12345
	path := "/zia/api/v1/locations/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = locationmanagement.Delete(context.Background(), service, locationID)

	require.NoError(t, err)
}

func TestLocationManagement_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	locationName := "HQ Office"
	path := "/zia/api/v1/locations"

	server.On("GET", path, common.SuccessResponse([]locationmanagement.Locations{
		{ID: 1, Name: "Other Office", Country: "US"},
		{ID: 2, Name: locationName, Country: "US"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := locationmanagement.GetLocationByName(context.Background(), service, locationName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, locationName, result.Name)
}

func TestLocationManagement_GetSublocations_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	parentID := 12345
	path := "/zia/api/v1/locations/12345/sublocations"

	server.On("GET", path, common.SuccessResponse([]locationmanagement.Locations{
		{ID: 100, Name: "Sublocation 1", ParentID: parentID},
		{ID: 200, Name: "Sublocation 2", ParentID: parentID},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := locationmanagement.GetSublocations(context.Background(), service, parentID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

// Note: GetSubLocation, GetSubLocationBySubID, GetSubLocationByNames, and GetSubLocationByName tests omitted
// due to complex internal calls that require multi-step mocking

func TestLocationManagement_BulkDelete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/locations/bulkDelete"

	server.On("POST", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	ids := []int{1, 2, 3}
	_, err = locationmanagement.BulkDelete(context.Background(), service, ids)

	require.NoError(t, err)
}

func TestLocationManagement_GetLocationSupportedCountries_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/locations/supportedCountries"

	server.On("GET", path, common.SuccessResponse([]string{"US", "CA", "GB", "DE", "FR"}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := locationmanagement.GetLocationSupportedCountries(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 5)
	assert.Contains(t, result, "US")
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestLocations_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Locations JSON marshaling", func(t *testing.T) {
		loc := locationmanagement.Locations{
			ID:                12345,
			Name:              "Headquarters",
			ParentID:          0,
			UpBandwidth:       100000,
			DnBandwidth:       200000,
			Country:           "US",
			State:             "California",
			Language:          "en-US",
			TZ:                "America/Los_Angeles",
			IPAddresses:       []string{"203.0.113.0/24"},
			Ports:             []int{80, 443},
			AuthRequired:      true,
			BasicAuthEnabled:  true,
			DigestAuthEnabled: false,
			KerberosAuth:      false,
			SurrogateIP:       true,
			IdleTimeInMinutes: 60,
			OFWEnabled:        true,
			IPSControl:        true,
			AUPEnabled:        true,
			CautionEnabled:    true,
			IPv6Enabled:       false,
		}

		data, err := json.Marshal(loc)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Headquarters"`)
		assert.Contains(t, string(data), `"country":"US"`)
		assert.Contains(t, string(data), `"ofwEnabled":true`)
	})

	t.Run("Locations JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Branch Office",
			"parentId": 12345,
			"upBandwidth": 50000,
			"dnBandwidth": 100000,
			"country": "CA",
			"state": "Ontario",
			"tz": "America/Toronto",
			"ipAddresses": ["198.51.100.0/24"],
			"ports": [80, 443, 8080],
			"authRequired": true,
			"basicAuthEnabled": true,
			"digestAuthEnabled": true,
			"kerberosAuth": false,
			"iotDiscoveryEnabled": true,
			"surrogateIP": true,
			"idleTimeInMinutes": 30,
			"displayTimeUnit": "MINUTE",
			"surrogateIPEnforcedForKnownBrowsers": true,
			"surrogateRefreshTimeInMinutes": 120,
			"ofwEnabled": true,
			"ipsControl": true,
			"aupEnabled": true,
			"aupTimeoutInDays": 7,
			"profile": "CORPORATE",
			"vpnCredentials": [
				{
					"id": 100,
					"type": "UFQDN",
					"fqdn": "vpn.branch.company.com"
				}
			],
			"staticLocationGroups": [
				{"id": 200, "name": "US Locations"}
			],
			"dynamiclocationGroups": [
				{"id": 300, "name": "Auto-Assigned"}
			]
		}`

		var loc locationmanagement.Locations
		err := json.Unmarshal([]byte(jsonData), &loc)
		require.NoError(t, err)

		assert.Equal(t, 54321, loc.ID)
		assert.Equal(t, 12345, loc.ParentID)
		assert.Equal(t, "CA", loc.Country)
		assert.Len(t, loc.IPAddresses, 1)
		assert.Len(t, loc.Ports, 3)
		assert.Len(t, loc.VPNCredentials, 1)
		assert.True(t, loc.IOTDiscoveryEnabled)
	})

	t.Run("VPNCredentials JSON marshaling", func(t *testing.T) {
		vpn := locationmanagement.VPNCredentials{
			ID:           12345,
			Type:         "UFQDN",
			FQDN:         "vpn.company.com",
			IPAddress:    "203.0.113.1",
			PreSharedKey: "secret-key",
			Comments:     "Primary VPN",
		}

		data, err := json.Marshal(vpn)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"type":"UFQDN"`)
		assert.Contains(t, string(data), `"fqdn":"vpn.company.com"`)
	})

	t.Run("Location JSON marshaling", func(t *testing.T) {
		loc := locationmanagement.Location{
			ID:   12345,
			Name: "HQ Location",
		}

		data, err := json.Marshal(loc)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"HQ Location"`)
	})

	t.Run("ManagedBy JSON marshaling", func(t *testing.T) {
		managed := locationmanagement.ManagedBy{
			ID:   100,
			Name: "SD-WAN Partner",
		}

		data, err := json.Marshal(managed)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":100`)
		assert.Contains(t, string(data), `"name":"SD-WAN Partner"`)
	})
}

func TestLocations_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse locations list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "HQ", "country": "US", "parentId": 0},
			{"id": 2, "name": "Branch 1", "country": "US", "parentId": 1},
			{"id": 3, "name": "Branch 2", "country": "CA", "parentId": 1}
		]`

		var locations []locationmanagement.Locations
		err := json.Unmarshal([]byte(jsonResponse), &locations)
		require.NoError(t, err)

		assert.Len(t, locations, 3)
		assert.Equal(t, 0, locations[0].ParentID)
		assert.Equal(t, 1, locations[1].ParentID)
	})

	t.Run("Parse location with all security features", func(t *testing.T) {
		jsonResponse := `{
			"id": 100,
			"name": "Secure Location",
			"authRequired": true,
			"basicAuthEnabled": true,
			"digestAuthEnabled": true,
			"kerberosAuth": true,
			"surrogateIP": true,
			"surrogateIPEnforcedForKnownBrowsers": true,
			"ofwEnabled": true,
			"ipsControl": true,
			"aupEnabled": true,
			"cautionEnabled": true,
			"sslScanEnabled": true,
			"zappSSLScanEnabled": true,
			"xffForwardEnabled": true,
			"ipv6Enabled": true
		}`

		var loc locationmanagement.Locations
		err := json.Unmarshal([]byte(jsonResponse), &loc)
		require.NoError(t, err)

		assert.True(t, loc.AuthRequired)
		assert.True(t, loc.KerberosAuth)
		assert.True(t, loc.OFWEnabled)
		assert.True(t, loc.IPSControl)
		assert.True(t, loc.AUPEnabled)
		assert.True(t, loc.IPv6Enabled)
	})

	t.Run("Parse supported countries", func(t *testing.T) {
		jsonResponse := `["US", "CA", "GB", "DE", "FR", "JP", "AU"]`

		var countries []string
		err := json.Unmarshal([]byte(jsonResponse), &countries)
		require.NoError(t, err)

		assert.Len(t, countries, 7)
		assert.Contains(t, countries, "US")
	})
}

func TestLocations_Sublocation(t *testing.T) {
	t.Parallel()

	t.Run("Sublocation with scope", func(t *testing.T) {
		jsonData := `{
			"id": 200,
			"name": "Sublocation 1",
			"parentId": 100,
			"subLocScopeEnabled": true,
			"subLocScope": "AWS_TAG",
			"subLocScopeValues": ["env:production", "app:web"],
			"subLocAccIds": ["123456789012"],
			"otherSubLocation": false,
			"other6SubLocation": false
		}`

		var subloc locationmanagement.Locations
		err := json.Unmarshal([]byte(jsonData), &subloc)
		require.NoError(t, err)

		assert.Equal(t, 100, subloc.ParentID)
		assert.True(t, subloc.SubLocScopeEnabled)
		assert.Equal(t, "AWS_TAG", subloc.SubLocScope)
		assert.Len(t, subloc.SubLocScopeValues, 2)
	})

	t.Run("Other sublocation", func(t *testing.T) {
		jsonData := `{
			"id": 300,
			"name": "Other",
			"parentId": 100,
			"otherSubLocation": true,
			"other6SubLocation": false
		}`

		var subloc locationmanagement.Locations
		err := json.Unmarshal([]byte(jsonData), &subloc)
		require.NoError(t, err)

		assert.True(t, subloc.OtherSubLocation)
		assert.False(t, subloc.Other6SubLocation)
	})
}

