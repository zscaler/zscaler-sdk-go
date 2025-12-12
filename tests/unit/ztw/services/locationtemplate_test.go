// Package services provides unit tests for ZTW services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/locationmanagement/locationtemplate"
)

func TestLocationTemplate_Structure(t *testing.T) {
	t.Parallel()

	t.Run("LocationTemplate JSON marshaling", func(t *testing.T) {
		template := locationtemplate.LocationTemplate{
			ID:          12345,
			Name:        "Branch-Template",
			Description: "Standard branch location template",
			Editable:    true,
			LastModTime: 1699000000,
			LocationTemplateDetails: &locationtemplate.LocationTemplateDetails{
				TemplatePrefix:                      "BRANCH",
				XFFForwardEnabled:                   true,
				AuthRequired:                        true,
				CautionEnabled:                      true,
				AupEnabled:                          true,
				AupTimeoutInDays:                    30,
				OFWEnabled:                          true,
				IPSControl:                          true,
				EnforceBandwidthControl:             true,
				UpBandwidth:                         100,
				DnBandwidth:                         200,
				SurrogateIP:                         true,
				IdleTimeInMinutes:                   60,
				SurrogateIPEnforcedForKnownBrowsers: true,
			},
		}

		data, err := json.Marshal(template)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Branch-Template"`)
		assert.Contains(t, string(data), `"template"`)
		assert.Contains(t, string(data), `"ofwEnabled":true`)
	})

	t.Run("LocationTemplate JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "HQ-Template",
			"desc": "Headquarters location template",
			"editable": false,
			"lastModTime": 1699500000,
			"template": {
				"templatePrefix": "HQ",
				"xffForwardEnabled": true,
				"authRequired": true,
				"cautionEnabled": false,
				"aupEnabled": true,
				"aupTimeoutInDays": 7,
				"ofwEnabled": true,
				"ipsControl": true,
				"enforceBandwidthControl": true,
				"upBandwidth": 500,
				"dnBandwidth": 1000,
				"displayTimeUnit": "MINUTE",
				"idleTimeInMinutes": 30,
				"surrogateIP": true,
				"surrogateIPEnforcedForKnownBrowsers": true,
				"surrogateRefreshTimeUnit": "HOUR",
				"surrogateRefreshTimeInMinutes": 120
			}
		}`

		var template locationtemplate.LocationTemplate
		err := json.Unmarshal([]byte(jsonData), &template)
		require.NoError(t, err)

		assert.Equal(t, 54321, template.ID)
		assert.Equal(t, "HQ-Template", template.Name)
		assert.False(t, template.Editable)
		assert.NotNil(t, template.LocationTemplateDetails)
		assert.Equal(t, "HQ", template.LocationTemplateDetails.TemplatePrefix)
		assert.Equal(t, 500, template.LocationTemplateDetails.UpBandwidth)
		assert.Equal(t, 1000, template.LocationTemplateDetails.DnBandwidth)
		assert.True(t, template.LocationTemplateDetails.SurrogateIP)
	})

	t.Run("LocationTemplateDetails JSON marshaling", func(t *testing.T) {
		details := locationtemplate.LocationTemplateDetails{
			TemplatePrefix:          "TEST",
			XFFForwardEnabled:       true,
			AuthRequired:            true,
			OFWEnabled:              true,
			IPSControl:              true,
			EnforceBandwidthControl: false,
			DisplayTimeUnit:         "HOUR",
			IdleTimeInMinutes:       120,
		}

		data, err := json.Marshal(details)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"templatePrefix":"TEST"`)
		assert.Contains(t, string(data), `"xffForwardEnabled":true`)
		assert.Contains(t, string(data), `"displayTimeUnit":"HOUR"`)
	})
}

func TestLocationTemplate_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse location templates list", func(t *testing.T) {
		jsonResponse := `[
			{
				"id": 1,
				"name": "Template-1",
				"editable": true,
				"template": {
					"templatePrefix": "T1",
					"ofwEnabled": true
				}
			},
			{
				"id": 2,
				"name": "Template-2",
				"editable": false,
				"template": {
					"templatePrefix": "T2",
					"ofwEnabled": false
				}
			}
		]`

		var templates []locationtemplate.LocationTemplate
		err := json.Unmarshal([]byte(jsonResponse), &templates)
		require.NoError(t, err)

		assert.Len(t, templates, 2)
		assert.True(t, templates[0].Editable)
		assert.False(t, templates[1].Editable)
		assert.Equal(t, "T1", templates[0].LocationTemplateDetails.TemplatePrefix)
		assert.Equal(t, "T2", templates[1].LocationTemplateDetails.TemplatePrefix)
	})

	t.Run("Parse template with all security features", func(t *testing.T) {
		jsonResponse := `{
			"id": 100,
			"name": "Secure-Template",
			"template": {
				"templatePrefix": "SEC",
				"authRequired": true,
				"cautionEnabled": true,
				"aupEnabled": true,
				"aupTimeoutInDays": 1,
				"ofwEnabled": true,
				"ipsControl": true,
				"surrogateIP": true,
				"surrogateIPEnforcedForKnownBrowsers": true
			}
		}`

		var template locationtemplate.LocationTemplate
		err := json.Unmarshal([]byte(jsonResponse), &template)
		require.NoError(t, err)

		details := template.LocationTemplateDetails
		assert.True(t, details.AuthRequired)
		assert.True(t, details.CautionEnabled)
		assert.True(t, details.AupEnabled)
		assert.True(t, details.OFWEnabled)
		assert.True(t, details.IPSControl)
		assert.True(t, details.SurrogateIP)
	})
}

