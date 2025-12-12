// Package services provides unit tests for ZTW services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
)

func TestZTWCommon_Structure(t *testing.T) {
	t.Parallel()

	t.Run("IDName JSON marshaling", func(t *testing.T) {
		item := common.IDName{
			ID:   123,
			Name: "Test Item",
		}

		data, err := json.Marshal(item)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":123`)
		assert.Contains(t, string(data), `"name":"Test Item"`)
	})

	t.Run("IDNameExtensions JSON marshaling", func(t *testing.T) {
		item := common.IDNameExtensions{
			ID:   456,
			Name: "Extended Item",
			Extensions: map[string]interface{}{
				"customField": "value",
				"priority":    1,
			},
		}

		data, err := json.Marshal(item)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":456`)
		assert.Contains(t, string(data), `"extensions"`)
	})

	t.Run("CommonIDNameExternalID JSON marshaling", func(t *testing.T) {
		item := common.CommonIDNameExternalID{
			ID:              789,
			Name:            "External Item",
			IsNameL10nTag:   false,
			Deleted:         false,
			ExternalID:      "ext-789",
			AssociationTime: 1699000000,
		}

		data, err := json.Marshal(item)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":789`)
		assert.Contains(t, string(data), `"externalId":"ext-789"`)
		assert.Contains(t, string(data), `"associationTime":1699000000`)
	})

	t.Run("ECVMs JSON marshaling", func(t *testing.T) {
		vm := common.ECVMs{
			ID:                1001,
			Name:              "EC-VM-1",
			Status:            []string{"RUNNING", "HEALTHY"},
			OperationalStatus: "ACTIVE",
			FormFactor:        "LARGE",
			CityGeoId:         12345,
			NATIP:             "10.0.0.1",
			ZiaGateway:        "zia-gw-1.zscaler.net",
			ZpaBroker:         "zpa-broker-1.zscaler.net",
			BuildVersion:      "22.1.0",
			LastUpgradeTime:   1699000000,
		}

		data, err := json.Marshal(vm)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":1001`)
		assert.Contains(t, string(data), `"operationalStatus":"ACTIVE"`)
		assert.Contains(t, string(data), `"buildVersion":"22.1.0"`)
	})

	t.Run("ManagementNw JSON marshaling", func(t *testing.T) {
		mgmt := common.ManagementNw{
			ID:             1,
			IPStart:        "10.0.0.1",
			IPEnd:          "10.0.0.254",
			Netmask:        "255.255.255.0",
			DefaultGateway: "10.0.0.1",
			NWType:         "MANAGEMENT",
			DNS: &common.DNS{
				ID:      1,
				IPs:     []string{"8.8.8.8", "8.8.4.4"},
				DNSType: "PUBLIC",
			},
		}

		data, err := json.Marshal(mgmt)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"ipStart":"10.0.0.1"`)
		assert.Contains(t, string(data), `"netmask":"255.255.255.0"`)
		assert.Contains(t, string(data), `"dns"`)
	})

	t.Run("ECInstances JSON marshaling", func(t *testing.T) {
		instance := common.ECInstances{
			ID:             100,
			ECInstanceType: "PRIMARY",
			ServiceIPs: &common.CommonIPs{
				IPStart: "10.1.0.1",
				IPEnd:   "10.1.0.10",
			},
			OutGwIp: "10.1.0.254",
			NatIP:   "203.0.113.1",
			DNSIP:   []string{"10.1.0.2", "10.1.0.3"},
		}

		data, err := json.Marshal(instance)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"ecInstanceType":"PRIMARY"`)
		assert.Contains(t, string(data), `"outGwIp":"10.1.0.254"`)
	})

	t.Run("ZPAAppSegments JSON marshaling", func(t *testing.T) {
		segment := common.ZPAAppSegments{
			ID:         12345,
			Name:       "Internal App",
			ExternalID: "zpa-app-123",
		}

		data, err := json.Marshal(segment)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"externalId":"zpa-app-123"`)
	})

	t.Run("RegionStatus JSON marshaling", func(t *testing.T) {
		region := common.RegionStatus{
			ID:        1,
			Name:      "US-East",
			CloudType: "AWS",
			Status:    true,
		}

		data, err := json.Marshal(region)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"name":"US-East"`)
		assert.Contains(t, string(data), `"cloudType":"AWS"`)
		assert.Contains(t, string(data), `"status":true`)
	})

	t.Run("SupportedRegions JSON marshaling", func(t *testing.T) {
		region := common.SupportedRegions{
			ID:        1,
			Name:      "West-Europe",
			CloudType: "AZURE",
		}

		data, err := json.Marshal(region)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"name":"West-Europe"`)
		assert.Contains(t, string(data), `"cloudType":"AZURE"`)
	})
}

func TestZTWCommon_SortParams(t *testing.T) {
	t.Parallel()

	t.Run("GetSortParams with both params", func(t *testing.T) {
		result := common.GetSortParams(common.NameSortField, common.ASCSortOrder)
		assert.Contains(t, result, "sortBy=name")
		assert.Contains(t, result, "sortOrder=asc")
	})

	t.Run("GetSortParams with only sortBy", func(t *testing.T) {
		result := common.GetSortParams(common.IDSortField, "")
		assert.Equal(t, "sortBy=id", result)
	})

	t.Run("GetSortParams with only sortOrder", func(t *testing.T) {
		result := common.GetSortParams("", common.DESCSortOrder)
		assert.Equal(t, "sortOrder=desc", result)
	})

	t.Run("GetPageSize returns correct value", func(t *testing.T) {
		assert.Equal(t, 1000, common.GetPageSize())
	})
}

