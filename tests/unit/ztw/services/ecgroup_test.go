// Package services provides unit tests for ZTW services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/ecgroup"
)

func TestEcGroup_Structure(t *testing.T) {
	t.Parallel()

	t.Run("EcGroup JSON marshaling", func(t *testing.T) {
		group := ecgroup.EcGroup{
			ID:                    12345,
			Name:                  "AWS-US-East-EC-Group",
			Description:           "Cloud Connector group for US East region",
			DeployType:            "AWS",
			Status:                []string{"ACTIVE", "PROVISIONED"},
			Platform:              "AWS",
			AWSAvailabilityZone:   "us-east-1a",
			MaxEcCount:            4,
			TunnelMode:            "GRE",
			Location: &common.CommonIDNameExternalID{
				ID:   100,
				Name: "US-East",
			},
			ProvTemplate: &common.CommonIDNameExternalID{
				ID:   200,
				Name: "Standard-Template",
			},
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"AWS-US-East-EC-Group"`)
		assert.Contains(t, string(data), `"deployType":"AWS"`)
		assert.Contains(t, string(data), `"platform":"AWS"`)
	})

	t.Run("EcGroup JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Azure-WestEurope-EC-Group",
			"desc": "Cloud Connector group for West Europe",
			"deployType": "AZURE",
			"status": ["ACTIVE"],
			"platform": "AZURE",
			"azureAvailabilityZone": "westeurope-1",
			"maxEcCount": 2,
			"tunnelMode": "IPSEC",
			"location": {
				"id": 300,
				"name": "West-Europe",
				"externalId": "ext-300"
			},
			"provTemplate": {
				"id": 400,
				"name": "Azure-Template"
			},
			"ecVMs": [
				{
					"id": 1001,
					"name": "EC-VM-1",
					"status": ["RUNNING"],
					"operationalStatus": "ACTIVE"
				}
			]
		}`

		var group ecgroup.EcGroup
		err := json.Unmarshal([]byte(jsonData), &group)
		require.NoError(t, err)

		assert.Equal(t, 54321, group.ID)
		assert.Equal(t, "Azure-WestEurope-EC-Group", group.Name)
		assert.Equal(t, "AZURE", group.DeployType)
		assert.Equal(t, "westeurope-1", group.AzureAvailabilityZone)
		assert.NotNil(t, group.Location)
		assert.Equal(t, "West-Europe", group.Location.Name)
		assert.Len(t, group.ECVMs, 1)
	})
}

func TestEcGroup_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse ec groups list response", func(t *testing.T) {
		jsonResponse := `[
			{
				"id": 1,
				"name": "Group-1",
				"deployType": "AWS",
				"status": ["ACTIVE"],
				"platform": "AWS"
			},
			{
				"id": 2,
				"name": "Group-2",
				"deployType": "AZURE",
				"status": ["PROVISIONING"],
				"platform": "AZURE"
			},
			{
				"id": 3,
				"name": "Group-3",
				"deployType": "GCP",
				"status": ["INACTIVE"],
				"platform": "GCP"
			}
		]`

		var groups []ecgroup.EcGroup
		err := json.Unmarshal([]byte(jsonResponse), &groups)
		require.NoError(t, err)

		assert.Len(t, groups, 3)
		assert.Equal(t, "AWS", groups[0].Platform)
		assert.Equal(t, "AZURE", groups[1].Platform)
		assert.Equal(t, "GCP", groups[2].Platform)
	})
}

