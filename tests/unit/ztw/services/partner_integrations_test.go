// Package services provides unit tests for ZTW services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/partner_integrations"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/partner_integrations/account_groups"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/partner_integrations/public_cloud_info"
)

func TestPartnerIntegrations_Structure(t *testing.T) {
	t.Parallel()

	t.Run("WorkloadDiscoverySettings JSON marshaling", func(t *testing.T) {
		settings := partner_integrations.WorkloadDiscoverySettings{
			TrustedAccountId: "123456789012",
			TrustedRoleName:  "ZscalerDiscoveryRole",
		}

		data, err := json.Marshal(settings)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"trustedAccountId":"123456789012"`)
		assert.Contains(t, string(data), `"trustedRoleName":"ZscalerDiscoveryRole"`)
	})

	t.Run("DiscoveryPermissions JSON marshaling", func(t *testing.T) {
		permissions := partner_integrations.DiscoveryPermissions{
			DiscoveryRole: "arn:aws:iam::123456789012:role/ZscalerRole",
			ExternalID:    "ext-id-12345",
		}

		data, err := json.Marshal(permissions)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"discoveryRole"`)
		assert.Contains(t, string(data), `"externalId":"ext-id-12345"`)
	})
}

func TestAccountGroups_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AccountGroups JSON marshaling", func(t *testing.T) {
		group := account_groups.AccountGroups{
			ID:          12345,
			Name:        "Production AWS Accounts",
			Description: "All production AWS accounts",
			CloudType:   "AWS",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Production AWS Accounts"`)
		assert.Contains(t, string(data), `"cloudType":"AWS"`)
	})

	t.Run("AccountGroups JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Development Azure Accounts",
			"description": "Azure dev accounts",
			"cloudType": "AZURE",
			"publicCloudAccounts": [
				{"id": 100, "name": "Account-1"},
				{"id": 101, "name": "Account-2"}
			],
			"cloudConnectorGroups": [
				{"id": 200, "name": "CC-Group-1"}
			]
		}`

		var group account_groups.AccountGroups
		err := json.Unmarshal([]byte(jsonData), &group)
		require.NoError(t, err)

		assert.Equal(t, 54321, group.ID)
		assert.Equal(t, "AZURE", group.CloudType)
		assert.Len(t, group.PublicCloudAccounts, 2)
		assert.Len(t, group.CloudConnectorGroups, 1)
	})
}

func TestPublicCloudInfo_Structure(t *testing.T) {
	t.Parallel()

	t.Run("PublicCloudInfo JSON marshaling", func(t *testing.T) {
		info := public_cloud_info.PublicCloudInfo{
			ID:         12345,
			Name:       "AWS-Prod-Account",
			CloudType:  "AWS",
			ExternalID: "ext-aws-123",
			AccountDetails: &public_cloud_info.AccountDetails{
				Name:         "Production Account",
				AwsAccountID: "123456789012",
				AwsRoleName:  "ZscalerRole",
				ExternalID:   "ext-aws-123",
			},
		}

		data, err := json.Marshal(info)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"cloudType":"AWS"`)
		assert.Contains(t, string(data), `"accountDetails"`)
	})

	t.Run("AccountDetails JSON marshaling", func(t *testing.T) {
		details := public_cloud_info.AccountDetails{
			Name:                   "AWS Account",
			AwsAccountID:           "123456789012",
			AwsRoleName:            "ZscalerDiscoveryRole",
			CloudWatchGroupArn:     "arn:aws:logs:us-east-1:123456789012:log-group:zscaler",
			EventBusName:           "zscaler-event-bus",
			ExternalID:             "ext-12345",
			LogInfoType:            "INFO",
			TroubleShootingLogging: true,
			TrustedAccountID:       "987654321098",
			TrustedRole:            "TrustedRole",
		}

		data, err := json.Marshal(details)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"awsAccountId":"123456789012"`)
		assert.Contains(t, string(data), `"troubleShootingLogging":true`)
	})

	t.Run("PublicCloudInfoLite JSON marshaling", func(t *testing.T) {
		lite := public_cloud_info.PublicCloudInfoLite{
			ID:        1,
			Name:      "Lite Account",
			AccountId: "123456789012",
			CloudType: "AWS",
		}

		data, err := json.Marshal(lite)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"accountId":"123456789012"`)
	})
}

func TestPartnerIntegrations_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse account groups list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Group-1", "cloudType": "AWS"},
			{"id": 2, "name": "Group-2", "cloudType": "AZURE"},
			{"id": 3, "name": "Group-3", "cloudType": "GCP"}
		]`

		var groups []account_groups.AccountGroups
		err := json.Unmarshal([]byte(jsonResponse), &groups)
		require.NoError(t, err)

		assert.Len(t, groups, 3)
		assert.Equal(t, "AWS", groups[0].CloudType)
		assert.Equal(t, "AZURE", groups[1].CloudType)
		assert.Equal(t, "GCP", groups[2].CloudType)
	})

	t.Run("Parse public cloud info list", func(t *testing.T) {
		jsonResponse := `[
			{
				"id": 1,
				"name": "AWS-Account-1",
				"cloudType": "AWS",
				"externalId": "ext-1",
				"regionStatus": [
					{"id": 1, "name": "us-east-1", "status": true}
				]
			}
		]`

		var infos []public_cloud_info.PublicCloudInfo
		err := json.Unmarshal([]byte(jsonResponse), &infos)
		require.NoError(t, err)

		assert.Len(t, infos, 1)
		assert.Len(t, infos[0].RegionStatus, 1)
	})
}

