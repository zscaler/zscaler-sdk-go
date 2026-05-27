package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/administrator_controller"
	appconnectorschedule "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorschedule"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentbrowseraccess"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentinspection"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentpra"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/bacertificate"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/branch_connector"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/client_settings"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloud_connector"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloud_connector_group"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/customerversionprofile"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/extranet_resource"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/location_controller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud_controller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud_group"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praapproval"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praconsole"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/pracredential"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praportal"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/provisioningkey"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/scimattributeheader"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/scimgroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgecontroller"
	serviceedgeschedule "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgeschedule"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/aup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_controller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_link"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/workload_tag_group"
)

func TestBACertificate_Get_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "certificate", "missing")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := bacertificate.Get(context.Background(), api.Service, "missing")
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestAdministratorController_Get_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "admin", "missing")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := administrator_controller.Get(context.Background(), api.Service, "missing")
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestPrivateCloudGroup_Get_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "privateCloudGroup", "missing")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := private_cloud_group.Get(context.Background(), api.Service, "missing")
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestPrivateCloudController_Get_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "privateCloudController", "missing")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := private_cloud_controller.Get(context.Background(), api.Service, "missing")
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestServiceEdgeController_Get_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "serviceEdge", "missing")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := serviceedgecontroller.Get(context.Background(), api.Service, "missing")
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestCloudConnectorGroup_Get_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "cloudConnectorGroup", "missing")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := cloud_connector_group.Get(context.Background(), api.Service, "missing")
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestBranchConnector_GetAll_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "branchConnector")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := branch_connector.GetAll(context.Background(), api.Service)
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestCloudConnector_GetAll_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "cloudConnector")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := cloud_connector.GetAll(context.Background(), api.Service)
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestExtranetResource_GetExtranetResourcePartner_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "extranetResource", "partner")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := extranet_resource.GetExtranetResourcePartner(context.Background(), api.Service)
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestWorkloadTagGroup_GetWorkloadTagGroup_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "workloadTagGroup", "summary")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := workload_tag_group.GetWorkloadTagGroup(context.Background(), api.Service)
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestApplicationSegmentInspection_Get_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "application", "missing")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := applicationsegmentinspection.Get(context.Background(), api.Service, "missing")
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestAppConnectorSchedule_GetSchedule_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "connectorSchedule")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := appconnectorschedule.GetSchedule(context.Background(), api.Service)
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestServiceEdgeSchedule_GetSchedule_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "serviceEdgeSchedule")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := serviceedgeschedule.GetSchedule(context.Background(), api.Service)
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestApplicationSegmentBrowserAccess_Get_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "application", "missing")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := applicationsegmentbrowseraccess.Get(context.Background(), api.Service, "missing")
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestApplicationSegmentPRA_Get_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "application", "missing")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := applicationsegmentpra.Get(context.Background(), api.Service, "missing")
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestPRAPortal_Get_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "praPortal", "missing")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := praportal.Get(context.Background(), api.Service, "missing")
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestPRAConsole_Get_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "praConsole", "missing")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := praconsole.Get(context.Background(), api.Service, "missing")
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestPRAApproval_Get_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "approval", "missing")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := praapproval.Get(context.Background(), api.Service, "missing")
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestPRACredential_Get_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "credential", "missing")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := pracredential.Get(context.Background(), api.Service, "missing")
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestUserPortalController_Get_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "userPortal", "missing")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := portal_controller.Get(context.Background(), api.Service, "missing")
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestUserPortalLink_Get_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "userPortalLink", "missing")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := portal_link.Get(context.Background(), api.Service, "missing")
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestUserPortalAUP_Get_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "userportal", "aup", "missing")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := aup.Get(context.Background(), api.Service, "missing")
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestClientSettings_GetAllClientSettings_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "clientSetting", "all")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := client_settings.GetAllClientSettings(context.Background(), api.Service)
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestCustomerVersionProfile_GetAll_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "visible", "versionProfiles")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := customerversionprofile.GetAll(context.Background(), api.Service)
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestLocationController_GetLocationSummary_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "location", "summary")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := location_controller.GetLocationSummary(context.Background(), api.Service)
	require.Error(t, err)
	assert.Nil(t, got)
}


func TestScimGroup_Get_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	path := common.ZPAUserConfigPath(api.CustomerID, "scimgroup", "missing")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := scimgroup.Get(context.Background(), api.Service, "missing")
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestScimAttributeHeader_Get_NotFound_SDK(t *testing.T) {
	t.Parallel()
	api := common.NewZPATest(t)
	idpID := "idp-1"
	path := common.ZPAPath(api.CustomerID, "idp", idpID, "scimattribute", "missing")
	api.On("GET", path, common.NotFoundResponse())
	got, _, err := scimattributeheader.Get(context.Background(), api.Service, idpID, "missing")
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestZPAErrorPaths_GetByName_NotFound(t *testing.T) {
	t.Parallel()

	t.Run("applicationsegment", func(t *testing.T) {
		api := common.NewZPATest(t)
		path := common.ZPAPath(api.CustomerID, "application")
		api.On("GET", path, common.SuccessResponse(common.ZPAList([]applicationsegment.ApplicationSegmentResource{})))
		got, _, err := applicationsegment.GetByName(context.Background(), api.Service, "missing-app")
		require.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("provisioningkey", func(t *testing.T) {
		api := common.NewZPATest(t)
		keyType := "CONNECTOR_GRP"
		path := common.ZPAPath(api.CustomerID, "associationType", keyType, "provisioningKey")
		api.On("GET", path, common.SuccessResponse(common.ZPAList([]provisioningkey.ProvisioningKey{})))
		got, _, err := provisioningkey.GetByName(context.Background(), api.Service, keyType, "missing-key")
		require.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("scimgroup", func(t *testing.T) {
		api := common.NewZPATest(t)
		idpID := "idp-123"
		path := common.ZPAUserConfigPath(api.CustomerID, "scimgroup", "idpId", idpID)
		api.On("GET", path, common.SuccessResponse(common.ZPAList([]scimgroup.ScimGroup{})))
		got, _, err := scimgroup.GetByName(context.Background(), api.Service, "missing-group", idpID)
		require.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("policysetcontroller", func(t *testing.T) {
		api := common.NewZPATest(t)
		policyType := "ACCESS_POLICY"
		path := common.ZPAPath(api.CustomerID, "policySet", "rules", "policyType", policyType)
		api.On("GET", path, common.SuccessResponse(common.ZPAList([]policysetcontroller.PolicyRule{})))
		got, _, err := policysetcontroller.GetByNameAndType(context.Background(), api.Service, policyType, "missing-rule")
		require.Error(t, err)
		assert.Nil(t, got)
	})
}
