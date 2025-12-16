// Package services provides unit tests for ZTW services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/locationmanagement/location"
)

// =====================================================
// SDK Function Tests
// =====================================================

func TestLocation_GetLocation_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	locationID := 12345
	path := "/ztw/api/v1/location/12345"

	server.On("GET", path, common.SuccessResponse(location.Locations{
		ID:         locationID,
		Name:       "US-East-HQ",
		Country:    "US",
		ECLocation: true,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := location.GetLocation(context.Background(), service, locationID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, locationID, result.ID)
	assert.Equal(t, "US-East-HQ", result.Name)
}

func TestLocation_GetLocationByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	locationName := "US-East-HQ"
	path := "/ztw/api/v1/location"

	server.On("GET", path, common.SuccessResponse([]location.Locations{
		{ID: 1, Name: "Other Location", Country: "UK"},
		{ID: 2, Name: locationName, Country: "US"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := location.GetLocationByName(context.Background(), service, locationName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, locationName, result.Name)
}

func TestLocation_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/location"

	server.On("GET", path, common.SuccessResponse([]location.Locations{
		{ID: 1, Name: "Location 1", Country: "US"},
		{ID: 2, Name: "Location 2", Country: "UK"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := location.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestLocation_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/location"

	server.On("POST", path, common.SuccessResponse(location.Locations{
		ID:         99999,
		Name:       "New Location",
		Country:    "US",
		ECLocation: true,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newLoc := &location.Locations{
		Name:       "New Location",
		Country:    "US",
		ECLocation: true,
	}

	result, err := location.Create(context.Background(), service, newLoc)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestLocation_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	locationID := 12345
	path := "/ztw/api/v1/location/12345"

	server.On("PUT", path, common.SuccessResponse(location.Locations{
		ID:      locationID,
		Name:    "Updated Location",
		Country: "US",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateLoc := &location.Locations{
		ID:      locationID,
		Name:    "Updated Location",
		Country: "US",
	}

	result, _, err := location.Update(context.Background(), service, locationID, updateLoc)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Location", result.Name)
}

func TestLocation_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	locationID := 12345
	path := "/ztw/api/v1/location/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = location.Delete(context.Background(), service, locationID)

	require.NoError(t, err)
}

// =====================================================
// Structure Tests
// =====================================================

func TestLocation_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Locations JSON marshaling", func(t *testing.T) {
		loc := location.Locations{
			ID:                      12345,
			Name:                    "US-East-HQ",
			ParentID:                0,
			EnforceBandwidthControl: true,
			UpBandwidth:             100000,
			DnBandwidth:             200000,
			Country:                 "US",
			State:                   "Virginia",
			TZ:                      "America/New_York",
			AuthRequired:            true,
			SSLScanEnabled:          true,
			XFFForwardEnabled:       true,
			ECLocation:              true,
			SurrogateIP:             true,
			IdleTimeInMinutes:       30,
			OFWEnabled:              true,
			IPSControl:              true,
			AUPEnabled:              true,
			IPv6Enabled:             true,
		}

		data, err := json.Marshal(loc)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"US-East-HQ"`)
		assert.Contains(t, string(data), `"country":"US"`)
		assert.Contains(t, string(data), `"ecLocation":true`)
	})

	t.Run("Locations JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "EU-West-Branch",
			"parentId": 12345,
			"enforceBandwidthControl": true,
			"upBandwidth": 50000,
			"dnBandwidth": 100000,
			"country": "DE",
			"state": "Bavaria",
			"tz": "Europe/Berlin",
			"authRequired": true,
			"sslScanEnabled": true,
			"xffForwardEnabled": true,
			"ecLocation": true,
			"surrogateIP": true,
			"idleTimeInMinutes": 60,
			"ofwEnabled": true,
			"ipsControl": true,
			"aupEnabled": false,
			"ipv6Enabled": false,
			"kerberosAuth": true,
			"digestAuthEnabled": false
		}`

		var loc location.Locations
		err := json.Unmarshal([]byte(jsonData), &loc)
		require.NoError(t, err)

		assert.Equal(t, 54321, loc.ID)
		assert.Equal(t, "EU-West-Branch", loc.Name)
		assert.Equal(t, 12345, loc.ParentID)
		assert.Equal(t, "DE", loc.Country)
		assert.Equal(t, "Bavaria", loc.State)
		assert.True(t, loc.ECLocation)
		assert.True(t, loc.KerberosAuth)
	})

	t.Run("VPNCredentials JSON marshaling", func(t *testing.T) {
		cred := location.VPNCredentials{
			ID:           1,
			Type:         "UFQDN",
			FQDN:         "vpn.company.com",
			IPAddress:    "192.168.1.1",
			PreSharedKey: "secret-key",
			Comments:     "Primary VPN connection",
		}

		data, err := json.Marshal(cred)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"type":"UFQDN"`)
		assert.Contains(t, string(data), `"fqdn":"vpn.company.com"`)
	})

	t.Run("VPCInfo JSON marshaling", func(t *testing.T) {
		vpc := location.VPCInfo{
			CloudProvider: "AWS",
			CloudMeta: location.CloudMeta{
				ID:   1,
				Name: "vpc-12345",
			},
		}

		data, err := json.Marshal(vpc)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"cloudProvider":"AWS"`)
		assert.Contains(t, string(data), `"name":"vpc-12345"`)
	})
}

func TestLocation_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse locations list response", func(t *testing.T) {
		jsonResponse := `[
			{
				"id": 1,
				"name": "Headquarters",
				"country": "US",
				"ecLocation": true
			},
			{
				"id": 2,
				"name": "Branch-EU",
				"parentId": 1,
				"country": "DE",
				"ecLocation": true
			},
			{
				"id": 3,
				"name": "Branch-APAC",
				"parentId": 1,
				"country": "JP",
				"ecLocation": true
			}
		]`

		var locations []location.Locations
		err := json.Unmarshal([]byte(jsonResponse), &locations)
		require.NoError(t, err)

		assert.Len(t, locations, 3)
		assert.Equal(t, "Headquarters", locations[0].Name)
		assert.Equal(t, 0, locations[0].ParentID)
		assert.Equal(t, 1, locations[1].ParentID)
	})

	t.Run("Parse location with VPN credentials", func(t *testing.T) {
		jsonResponse := `{
			"id": 100,
			"name": "Secure Branch",
			"country": "US",
			"vpnCredentials": [
				{
					"id": 1,
					"type": "UFQDN",
					"fqdn": "branch.vpn.company.com",
					"ipAddress": "10.0.0.1"
				},
				{
					"id": 2,
					"type": "IP",
					"ipAddress": "10.0.0.2"
				}
			]
		}`

		var loc location.Locations
		err := json.Unmarshal([]byte(jsonResponse), &loc)
		require.NoError(t, err)

		assert.Equal(t, "Secure Branch", loc.Name)
		assert.Len(t, loc.VPNCredentials, 2)
		assert.Equal(t, "UFQDN", loc.VPNCredentials[0].Type)
		assert.Equal(t, "IP", loc.VPNCredentials[1].Type)
	})
}

