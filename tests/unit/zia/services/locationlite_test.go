package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/location/locationlite"
)

const locationLitePath = "/zia/api/v1/locations/lite"

func sampleLocationLite() locationlite.LocationLite {
	return locationlite.LocationLite{
		Name:              "tests-location-lite",
		ParentID:          0,
		TZ:                "UNITED_STATES_AMERICA_LOS_ANGELES",
		XFFForwardEnabled: true,
		AUPEnabled:        true,
		CautionEnabled:    true,
		SurrogateIP:       true,
		SurrogateIPEnforcedForKnownBrowsers: true,
		OFWEnabled:        true,
		IPSControl:        true,
	}
}

func TestLocationLite_GetLocationLiteID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	locationID := 100
	loc := sampleLocationLite()
	loc.ID = locationID

	server.On("GET", locationLitePath+"/100", common.SuccessResponse(loc))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := locationlite.GetLocationLiteID(context.Background(), service, locationID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, locationID, result.ID)
	assert.Equal(t, "tests-location-lite", result.Name)
	assert.True(t, result.OFWEnabled)
}

func TestLocationLite_GetLocationLiteByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	locationName := "tests-location-lite"
	server.On("GET", locationLitePath, common.SuccessResponse([]locationlite.LocationLite{
		func() locationlite.LocationLite {
			l := sampleLocationLite()
			l.ID = 100
			return l
		}(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := locationlite.GetLocationLiteByName(context.Background(), service, locationName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, locationName, result.Name)
}

func TestLocationLite_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", locationLitePath, common.SuccessResponse([]locationlite.LocationLite{
		sampleLocationLite(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := locationlite.GetAll(context.Background(), service, nil)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestLocationLite_Structure(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		loc := sampleLocationLite()
		loc.ID = 100

		data, err := json.Marshal(loc)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"surrogateIP":true`)
		assert.Contains(t, string(data), `"ipsControl":true`)
	})
}
