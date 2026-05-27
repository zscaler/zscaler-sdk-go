package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/adminuserrolemgmt/admins"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/adminuserrolemgmt/roles"
	ziacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/filetypecontrol"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/filteringrules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/networkservices"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/location/locationgroups"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/location/locationlite"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/location/locationmanagement"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/malware_protection"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/pacfiles"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/remote_assistance"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sandbox/sandbox_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sandbox/sandbox_settings"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/security_policy_settings"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/traffic_capture"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/dc_exclusions"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/extranet"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/gretunnels"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/staticips"
	virtualipaddress "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/virtualipaddress"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/vpncredentials"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlfilteringpolicies"
)

func TestAdminUsers_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/adminUsers/99999"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := admins.GetAdminUsers(context.Background(), service, 99999)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestAdminRoles_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/adminRoles/99999"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := roles.Get(context.Background(), service, 99999)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestMalwareProtection_UpdateATPMalwareInspection_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", "/zia/api/v1/cyberThreatProtection/atpMalwareInspection", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = malware_protection.UpdateATPMalwareInspection(context.Background(), service, malware_protection.ATPMalwareInspection{
		InspectInbound: true,
	})
	require.Error(t, err)
}

func TestMalwareProtection_UpdateATPMalwareProtocol_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", "/zia/api/v1/cyberThreatProtection/atpMalwareProtocols", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = malware_protection.UpdateATPMalwareProtocol(context.Background(), service, malware_protection.ATPMalwareProtocols{})
	require.Error(t, err)
}

func TestMalwareProtection_UpdateATPMalwarePolicy_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", "/zia/api/v1/cyberThreatProtection/malwarePolicy", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = malware_protection.UpdateATPMalwarePolicy(context.Background(), service, malware_protection.MalwarePolicy{})
	require.Error(t, err)
}

func TestMalwareProtection_UpdateATPMalwareSettings_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", "/zia/api/v1/cyberThreatProtection/malwareSettings", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = malware_protection.UpdateATPMalwareSettings(context.Background(), service, malware_protection.MalwareSettings{})
	require.Error(t, err)
}

func TestPacFiles_GetPacFileByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/pacFiles", common.SuccessResponse([]pacfiles.PACFileConfig{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := pacfiles.GetPacFileByName(context.Background(), service, "missing.pac")
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestPacFiles_DeletePacFile_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("DELETE", "/zia/api/v1/pacFiles/99999", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = pacfiles.DeletePacFile(context.Background(), service, 99999)
	require.Error(t, err)
}

func TestRemoteAssistance_UpdateRemoteAssistance_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", "/zia/api/v1/remoteAssistance", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = remote_assistance.UpdateRemoteAssistance(context.Background(), service, remote_assistance.RemoteAssistance{})
	require.Error(t, err)
}

func TestSandboxRules_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/sandboxRules", common.SuccessResponse([]sandbox_rules.SandboxRules{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := sandbox_rules.GetByName(context.Background(), service, "missing-rule")
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestSandboxSettings_Update_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", "/zia/api/v1/behavioralAnalysisAdvancedSettings", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = sandbox_settings.Update(context.Background(), service, sandbox_settings.BaAdvancedSettings{})
	require.Error(t, err)
}

func TestSecurityPolicySettings_UpdateListUrls_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", "/zia/api/v1/security/advanced", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = security_policy_settings.UpdateListUrls(context.Background(), service, security_policy_settings.ListUrls{})
	require.Error(t, err)
}

func TestFileTypeControl_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/fileTypeRules/99999", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := filetypecontrol.Get(context.Background(), service, 99999)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestFileTypeControl_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/fileTypeRules", common.SuccessResponse([]filetypecontrol.FileTypeRules{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := filetypecontrol.GetByName(context.Background(), service, "missing-rule")
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestURLFilteringRules_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/urlFilteringRules", common.SuccessResponse([]urlfilteringpolicies.URLFilteringRule{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := urlfilteringpolicies.GetByName(context.Background(), service, "missing-rule")
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestLocationManagement_GetLocation_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/locations/99999", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := locationmanagement.GetLocation(context.Background(), service, 99999)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestLocationManagement_BulkDelete_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", "/zia/api/v1/locations/bulkDelete", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = locationmanagement.BulkDelete(context.Background(), service, []int{1, 2})
	require.Error(t, err)
}

func TestLocationLite_GetLocationLiteByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/locations/lite", common.SuccessResponse([]locationlite.LocationLite{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := locationlite.GetLocationLiteByName(context.Background(), service, "missing-location")
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestNetworkServices_GetByName_WithProtocol_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	protocol := "TCP"
	server.On("GET", "/zia/api/v1/networkServices", common.SuccessResponse([]networkservices.NetworkServices{
		{ID: 1, Name: "HTTPS", Type: "STANDARD"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservices.GetByName(context.Background(), service, "", &protocol, nil)
	require.NoError(t, err)
	assert.Equal(t, "HTTPS", result.Name)
}

func TestNetworkServices_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/networkServices", common.SuccessResponse([]networkservices.NetworkServices{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservices.GetByName(context.Background(), service, "missing-service", nil, nil)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestURLFilteringRules_Get_ISOLATE_CBIProfileFallback_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	getPath := "/zia/api/v1/urlFilteringRules/12345"
	listPath := "/zia/api/v1/urlFilteringRules"

	server.On("GET", getPath, common.SuccessResponse(urlfilteringpolicies.URLFilteringRule{
		ID: ruleID, Name: "Isolate Rule", Action: "ISOLATE", CBIProfileID: 99,
	}))
	server.On("GET", listPath, common.SuccessResponse([]urlfilteringpolicies.URLFilteringRule{
		{
			ID: ruleID, Name: "Isolate Rule", Action: "ISOLATE", CBIProfileID: 99,
			CBIProfile: &ziacommon.CBIProfile{ID: "99", Name: "Default Profile"},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := urlfilteringpolicies.Get(context.Background(), service, ruleID)
	require.NoError(t, err)
	require.NotNil(t, result.CBIProfile)
	assert.Equal(t, "99", result.CBIProfile.ID)
}

func TestFileTypeControl_GetFileTypeCategories_WithFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/fileTypeCategories", common.SuccessResponse([]filetypecontrol.FileTypeCategory{
		{ID: 1, Name: "ALZ"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	exclude := true
	result, err := filetypecontrol.GetFileTypeCategories(context.Background(), service, &filetypecontrol.GetFileTypeCategoriesFilterOptions{
		Enums:                  []string{"FILETYPECATEGORYFORFILETYPECONTROL"},
		ExcludeCustomFileTypes: &exclude,
	})
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestPacFiles_CreatePacFile_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", "/zia/api/v1/pacFiles", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = pacfiles.CreatePacFile(context.Background(), service, &pacfiles.PACFileConfig{Name: "test.pac"})
	require.Error(t, err)
}

func TestPacFiles_ValidatePacFile_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", "/zia/api/v1/pacFiles/validate", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = pacfiles.ValidatePacFile(context.Background(), service, "bad content")
	require.Error(t, err)
}

func TestLocationGroups_GetAll_WithFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	name := "US Locations"
	groupType := "STATIC_GROUP"
	server.On("GET", "/zia/api/v1/locations/groups", common.SuccessResponse([]locationgroups.LocationGroup{
		{ID: 1, Name: name, GroupType: groupType},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := locationgroups.GetAll(context.Background(), service, &locationgroups.GetAllFilterOptions{
		Name:      &name,
		GroupType: &groupType,
	})
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestLocationGroups_GetLocationGroupByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/locations/groups", common.SuccessResponse([]locationgroups.LocationGroup{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := locationgroups.GetLocationGroupByName(context.Background(), service, "missing-group")
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestLocationManagement_GetSubLocationByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/locations", common.SuccessResponse([]locationmanagement.Locations{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := locationmanagement.GetSubLocationByName(context.Background(), service, "missing-sub")
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestFirewallFilteringRules_Create_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", "/zia/api/v1/firewallFilteringRules", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = filteringrules.Create(context.Background(), service, &filteringrules.FirewallFilteringRules{Name: "fail"})
	require.Error(t, err)
}

func TestTrafficCapture_Create_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", "/zia/api/v1/trafficCaptureRules", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = traffic_capture.Create(context.Background(), service, &traffic_capture.TrafficCaptureRules{Name: "fail"})
	require.Error(t, err)
}

func TestTrafficCapture_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/trafficCaptureRules", common.SuccessResponse([]traffic_capture.TrafficCaptureRules{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := traffic_capture.GetByName(context.Background(), service, "missing")
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDCExclusions_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/dcExclusions", common.SuccessResponse([]dc_exclusions.DCExclusions{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dc_exclusions.GetByName(context.Background(), service, "missing-dc")
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestExtranet_GetExtranetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/extranet", common.SuccessResponse([]extranet.Extranet{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := extranet.GetExtranetByName(context.Background(), service, "missing")
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestStaticIPs_GetByIPAddress_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/staticIP", common.SuccessResponse([]staticips.StaticIP{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := staticips.GetByIPAddress(context.Background(), service, "203.0.113.99")
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestVirtualIPAddress_GetZSGREVirtualIPList_InsufficientVIPs_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/vips/recommendedList", common.SuccessResponse([]virtualipaddress.GREVirtualIPList{
		{ID: 1, VirtualIp: "192.0.2.1"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = virtualipaddress.GetZSGREVirtualIPList(context.Background(), service, "203.0.113.1", 2)
	require.Error(t, err)
}

func TestFirewallFilteringRules_GetAll_AllFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/firewallFilteringRules", common.SuccessResponse([]filteringrules.FirewallFilteringRules{
		{ID: 1, Name: "Filtered"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := filteringrules.GetAll(context.Background(), service, &filteringrules.GetAllFilterOptions{
		RuleName: "Filtered", RuleLabel: "Default", RuleLabelId: 1, RuleOrder: "1",
		RuleDescription: "desc", RuleAction: "BLOCK", Location: "HQ", Department: "Eng",
		Group: "Admins", User: "admin", Device: "laptop", DeviceGroup: "corp",
		DeviceTrustLevel: "HIGH_TRUST", SrcIps: "10.0.0.1", DestAddresses: "8.8.8.8",
		SrcIpGroups: "src-grp", DestIpGroups: "dst-grp", NwApplication: "APNS",
		NwServices: "HTTP", DestIpCategories: "OFFICE365",
	})
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestFirewallFilteringRules_GetFirewallFilteringRuleCount_AllFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/firewallFilteringRules/count", common.SuccessResponse(4))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	count, err := filteringrules.GetFirewallFilteringRuleCount(context.Background(), service, &filteringrules.GetAllFilterOptions{
		PredefinedRuleCount: true, RuleName: "Block", RuleLabel: "Default", RuleLabelId: 1,
		RuleOrder: "1", RuleDescription: "desc", RuleAction: "BLOCK", Location: "HQ",
		Department: "Eng", Group: "Admins", User: "admin", Device: "laptop",
		DeviceGroup: "corp", DeviceTrustLevel: "HIGH_TRUST", SrcIps: "10.0.0.1",
		DestAddresses: "8.8.8.8", SrcIpGroups: "src-grp", DestIpGroups: "dst-grp",
		NwApplication: "APNS", NwServices: "HTTP", DestIpCategories: "OFFICE365",
	})
	require.NoError(t, err)
	assert.Equal(t, 4, count)
}

func TestTrafficCapture_Delete_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("DELETE", "/zia/api/v1/trafficCaptureRules/99999", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = traffic_capture.Delete(context.Background(), service, 99999)
	require.Error(t, err)
}

func TestTrafficCapture_GetTrafficCaptureRuleCount_AllFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/trafficCaptureRules/count", common.SuccessResponse(9))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	count, err := traffic_capture.GetTrafficCaptureRuleCount(context.Background(), service, &traffic_capture.TrafficCaptureRulesCountQuery{
		PredefinedRuleCount: true, RuleName: "Capture", RuleDescription: "desc",
		RuleOrder: "1", RuleAction: "CAPTURE", Location: "HQ", Department: "Eng",
		Group: "Admins", User: "admin", DeviceGroup: "corp", Device: "laptop",
		DeviceTrustLevel: "HIGH_TRUST",
	})
	require.NoError(t, err)
	assert.Equal(t, 9, count)
}

func TestVirtualIPAddress_GetVIPRecommendedList_AllOptions_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/vips/recommendedList", common.SuccessResponse([]virtualipaddress.GREVirtualIPList{
		{ID: 1, VirtualIp: "192.0.2.1"},
		{ID: 2, VirtualIp: "192.0.2.2"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := virtualipaddress.GetVIPRecommendedList(context.Background(), service,
		virtualipaddress.WithRoutableIP(true),
		virtualipaddress.WithWithinCountryOnly(false),
		virtualipaddress.WithIncludePrivateServiceEdge(true),
		virtualipaddress.WithIncludeCurrentVips(true),
		virtualipaddress.WithSourceIP("203.0.113.1"),
		virtualipaddress.WithLatitude(37.7749),
		virtualipaddress.WithLongitude(-122.4194),
		virtualipaddress.WithSubcloud("default"),
	)
	require.NoError(t, err)
	assert.Len(t, *result, 2)
}

func TestVirtualIPAddress_GetPairZSGREVirtualIPsWithinCountry_Fallback_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/vips/recommendedList", common.SuccessResponse([]virtualipaddress.GREVirtualIPList{
		{ID: 1, VirtualIp: "192.0.2.1", CountryCode: "CA"},
		{ID: 2, VirtualIp: "192.0.2.2", CountryCode: "US"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := virtualipaddress.GetPairZSGREVirtualIPsWithinCountry(context.Background(), service, "203.0.113.1", "US")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(*result), 2)
}

func TestVPNCredentials_GetVPNByType_WithFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/vpnCredentials", common.SuccessResponse([]vpncredentials.VPNCredentials{
		{ID: 1, Type: "UFQDN", FQDN: "vpn.example.com"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	includeOnly := true
	locationID := 1
	managedBy := 2
	result, err := vpncredentials.GetVPNByType(context.Background(), service, "UFQDN", &includeOnly, &locationID, &managedBy)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestVPNCredentials_GetVPNByType_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/vpnCredentials", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := vpncredentials.GetVPNByType(context.Background(), service, "UFQDN", nil, nil, nil)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestGRETunnels_GetByIPAddress_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/greTunnels", common.SuccessResponse([]gretunnels.GreTunnels{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := gretunnels.GetByIPAddress(context.Background(), service, "203.0.113.99")
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestStaticIPs_Create_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", "/zia/api/v1/staticIP", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = staticips.Create(context.Background(), service, &staticips.StaticIP{IpAddress: "203.0.113.1"})
	require.Error(t, err)
}
