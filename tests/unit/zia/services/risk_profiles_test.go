// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	ziacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloudapplications/risk_profiles"
)

const riskProfilesPath = "/zia/api/v1/riskProfiles"

// sampleRiskProfile mirrors the integration test payload in risk_profiles_test.go.
func sampleRiskProfile(name string) risk_profiles.RiskProfiles {
	return risk_profiles.RiskProfiles{
		ProfileName:               name,
		Status:                    "SANCTIONED",
		RiskIndex:                 []int{1, 2, 3, 4, 5},
		PasswordStrength:          "GOOD",
		PoorItemsOfService:        "YES",
		AdminAuditLogs:            "YES",
		DataBreach:                "YES",
		SourceIpRestrictions:      "YES",
		FileSharing:               "YES",
		MfaSupport:                "YES",
		SslPinned:                 "YES",
		Certifications:            []string{"AICPA", "CCPA", "CISP"},
		DataEncryptionInTransit:   []string{"SSLV2", "SSLV3", "TLSV1_0", "TLSV1_1", "TLSV1_2", "TLSV1_3", "UN_KNOWN"},
		HttpSecurityHeaders:       "YES",
		Evasive:                   "YES",
		DnsCaaPolicy:              "YES",
		SslCertValidity:           "YES",
		WeakCipherSupport:         "YES",
		Vulnerability:             "YES",
		VulnerableToHeartBleed:    "YES",
		SslCertKeySize:            "BITS_2048",
		VulnerableToPoodle:        "YES",
		SupportForWaf:             "YES",
		VulnerabilityDisclosure:   "YES",
		DomainKeysIdentifiedMail:  "YES",
		MalwareScanningForContent: "YES",
		DomainBasedMessageAuth:    "YES",
		SenderPolicyFramework:     "YES",
		RemoteScreenSharing:       "YES",
		VulnerableToLogJam:        "YES",
		ProfileType:               "CLOUD_APPLICATIONS",
	}
}

// =====================================================
// SDK Function Tests
// =====================================================

func TestRiskProfiles_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileID := 1001
	path := "/zia/api/v1/riskProfiles/1001"

	profile := sampleRiskProfile("tests-risk-profile")
	profile.ID = profileID

	server.On("GET", path, common.SuccessResponse(profile))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := risk_profiles.Get(context.Background(), service, profileID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, profileID, result.ID)
	assert.Equal(t, "SANCTIONED", result.Status)
	assert.Equal(t, "GOOD", result.PasswordStrength)
	assert.Equal(t, "BITS_2048", result.SslCertKeySize)
	assert.Equal(t, "CLOUD_APPLICATIONS", result.ProfileType)
	assert.Len(t, result.RiskIndex, 5)
	assert.Len(t, result.Certifications, 3)
}

func TestRiskProfiles_Get_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/riskProfiles/9999", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := risk_profiles.Get(context.Background(), service, 9999)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestRiskProfiles_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileName := "tests-risk-profile"
	server.On("GET", riskProfilesPath, common.SuccessResponse([]risk_profiles.RiskProfiles{
		{ID: 1, ProfileName: "other-profile", Status: "SANCTIONED"},
		func() risk_profiles.RiskProfiles {
			p := sampleRiskProfile(profileName)
			p.ID = 2
			return p
		}(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := risk_profiles.GetByName(context.Background(), service, profileName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, profileName, result.ProfileName)
	assert.Equal(t, 2, result.ID)
}

func TestRiskProfiles_GetByName_CaseInsensitive_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", riskProfilesPath, common.SuccessResponse([]risk_profiles.RiskProfiles{
		{ID: 5, ProfileName: "Tests-Risk-Profile", Status: "SANCTIONED", ProfileType: "CLOUD_APPLICATIONS"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := risk_profiles.GetByName(context.Background(), service, "TESTS-RISK-PROFILE")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Tests-Risk-Profile", result.ProfileName)
}

func TestRiskProfiles_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", riskProfilesPath, common.SuccessResponse([]risk_profiles.RiskProfiles{
		{ID: 1, ProfileName: "existing-profile"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := risk_profiles.GetByName(context.Background(), service, "non_existent_name")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no risk profiles found with name: non_existent_name")
}

func TestRiskProfiles_GetByName_APIError_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", riskProfilesPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := risk_profiles.GetByName(context.Background(), service, "any")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestRiskProfiles_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", riskProfilesPath, common.SuccessResponse(func() risk_profiles.RiskProfiles {
		p := sampleRiskProfile("tests-new-profile")
		p.ID = 88888
		return p
	}()))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newProfile := sampleRiskProfile("tests-new-profile")

	result, _, err := risk_profiles.Create(context.Background(), service, &newProfile)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 88888, result.ID)
	assert.Equal(t, "tests-new-profile", result.ProfileName)
	assert.Equal(t, "SANCTIONED", result.Status)
}

func TestRiskProfiles_Create_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", riskProfilesPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newProfile := sampleRiskProfile("fail-profile")

	result, _, err := risk_profiles.Create(context.Background(), service, &newProfile)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestRiskProfiles_Create_UnexpectedResponseType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", riskProfilesPath, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newProfile := sampleRiskProfile("no-body-profile")

	result, _, err := risk_profiles.Create(context.Background(), service, &newProfile)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "object returned from api was not a risk profile pointer")
}

func TestRiskProfiles_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileID := 1001
	path := "/zia/api/v1/riskProfiles/1001"

	server.On("PUT", path, common.SuccessResponse(func() risk_profiles.RiskProfiles {
		p := sampleRiskProfile("tests-updated-profile")
		p.ID = profileID
		return p
	}()))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateProfile := sampleRiskProfile("tests-updated-profile")
	updateProfile.ID = profileID

	result, _, err := risk_profiles.Update(context.Background(), service, profileID, &updateProfile)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "tests-updated-profile", result.ProfileName)
}

func TestRiskProfiles_Update_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", "/zia/api/v1/riskProfiles/1001", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateProfile := sampleRiskProfile("tests-updated-profile")

	result, _, err := risk_profiles.Update(context.Background(), service, 1001, &updateProfile)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestRiskProfiles_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("DELETE", "/zia/api/v1/riskProfiles/1001", common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = risk_profiles.Delete(context.Background(), service, 1001)

	require.NoError(t, err)
}

func TestRiskProfiles_Delete_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("DELETE", "/zia/api/v1/riskProfiles/1001", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = risk_profiles.Delete(context.Background(), service, 1001)

	require.Error(t, err)
}

func TestRiskProfiles_GetAllLite_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", riskProfilesPath+"/lite", common.SuccessResponse([]risk_profiles.RiskProfiles{
		{ID: 1, ProfileName: "lite-profile-1", Status: "SANCTIONED"},
		{ID: 2, ProfileName: "lite-profile-2", Status: "SANCTIONED"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := risk_profiles.GetAllLite(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestRiskProfiles_GetAllLite_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", riskProfilesPath+"/lite", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := risk_profiles.GetAllLite(context.Background(), service)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestRiskProfiles_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", riskProfilesPath, common.SuccessResponse([]risk_profiles.RiskProfiles{
		func() risk_profiles.RiskProfiles {
			p := sampleRiskProfile("profile-a")
			p.ID = 10
			return p
		}(),
		func() risk_profiles.RiskProfiles {
			p := sampleRiskProfile("profile-b")
			p.ID = 11
			return p
		}(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := risk_profiles.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "SANCTIONED", result[0].Status)
}

func TestRiskProfiles_GetAll_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", riskProfilesPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := risk_profiles.GetAll(context.Background(), service)

	require.Error(t, err)
	assert.Nil(t, result)
}

// =====================================================
// Structure Tests
// =====================================================

func TestRiskProfiles_Structure(t *testing.T) {
	t.Parallel()

	t.Run("RiskProfiles JSON marshaling", func(t *testing.T) {
		profile := sampleRiskProfile("tests-structure")
		profile.ID = 12345
		profile.ModifiedBy = &ziacommon.IDNameExtensions{ID: 1, Name: "admin@example.com"}
		profile.CustomTags = []ziacommon.IDNameExternalID{
			{ID: 10, Name: "tag-1", ExternalID: "ext-1"},
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"profileName":"tests-structure"`)
		assert.Contains(t, string(data), `"status":"SANCTIONED"`)
		assert.Contains(t, string(data), `"passwordStrength":"GOOD"`)
		assert.Contains(t, string(data), `"sslCertKeySize":"BITS_2048"`)
		assert.Contains(t, string(data), `"profileType":"CLOUD_APPLICATIONS"`)
		assert.Contains(t, string(data), `"certifications"`)
		assert.Contains(t, string(data), `"dataEncryptionInTransit"`)
	})

	t.Run("RiskProfiles JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 200,
			"profileName": "tests-unmarshal",
			"status": "SANCTIONED",
			"profileType": "CLOUD_APPLICATIONS",
			"riskIndex": [1, 2, 3, 4, 5],
			"passwordStrength": "GOOD",
			"poorItemsOfService": "YES",
			"certifications": ["AICPA", "CCPA", "CISP"],
			"dataEncryptionInTransit": ["TLSV1_2", "TLSV1_3"],
			"sslCertKeySize": "BITS_2048",
			"modifiedBy": {"id": 5, "name": "ops@example.com"}
		}`

		var profile risk_profiles.RiskProfiles
		err := json.Unmarshal([]byte(jsonData), &profile)
		require.NoError(t, err)

		assert.Equal(t, 200, profile.ID)
		assert.Equal(t, "tests-unmarshal", profile.ProfileName)
		assert.Equal(t, "GOOD", profile.PasswordStrength)
		require.NotNil(t, profile.ModifiedBy)
		assert.Equal(t, "ops@example.com", profile.ModifiedBy.Name)
		assert.Len(t, profile.RiskIndex, 5)
	})
}

func TestRiskProfiles_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse risk profiles list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "profileName": "profile-1", "status": "SANCTIONED", "profileType": "CLOUD_APPLICATIONS"},
			{"id": 2, "profileName": "profile-2", "status": "SANCTIONED", "profileType": "CLOUD_APPLICATIONS"}
		]`

		var profiles []risk_profiles.RiskProfiles
		err := json.Unmarshal([]byte(jsonResponse), &profiles)
		require.NoError(t, err)

		assert.Len(t, profiles, 2)
		assert.Equal(t, "CLOUD_APPLICATIONS", profiles[0].ProfileType)
	})
}
