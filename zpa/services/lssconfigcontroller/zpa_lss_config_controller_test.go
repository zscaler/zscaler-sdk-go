package lssconfigcontroller

/*
import (
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appconnectorgroup"
)


// clean all resources
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	cleanResources() // clean up at the beginning
}

func teardown() {
	cleanResources() // clean up at the end
}

func shouldClean() bool {
	val, present := os.LookupEnv("ZSCALER_SDK_TEST_SWEEP")
	return !present || (present && (val == "" || val == "true")) // simplified for clarity
}

func cleanResources() {
	if !shouldClean() {
		return
	}

	client, err := tests.NewZpaClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	service := New(client)
	resources, _, _ := service.GetAll()
	for _, r := range resources {
		if !strings.HasPrefix(r.LSSConfig.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.LSSConfig.Name)
		_, _ = service.Delete(r.ID)
	}
}

func TestLSSConfigController(t *testing.T) {
	ipAddress, _ := acctest.RandIpAddress("192.168.0.0/24")
	rPort := strconv.Itoa(acctest.RandIntRange(1000, 9999))
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	appConnGroupName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	appConnectorGroupService := appconnectorgroup.New(client)
	appGroup := appconnectorgroup.AppConnectorGroup{
		Name:                     appConnGroupName,
		Description:              appConnGroupName,
		Enabled:                  true,
		CityCountry:              "San Jose, US",
		Latitude:                 "37.3382082",
		Longitude:                "-121.8863286",
		Location:                 "San Jose, CA, USA",
		UpgradeDay:               "SUNDAY",
		UpgradeTimeInSecs:        "66600",
		OverrideVersionProfile:   true,
		VersionProfileName:       "Default",
		VersionProfileID:         "0",
		DNSQueryType:             "IPV4_IPV6",
		PRAEnabled:               false,
		WAFDisabled:              true,
		TCPQuickAckApp:           true,
		TCPQuickAckAssistant:     true,
		TCPQuickAckReadAssistant: true,
	}

	createdAppConnGroup, _, err := appConnectorGroupService.Create(appGroup)
	if err != nil || createdAppConnGroup == nil || createdAppConnGroup.ID == "" {
		t.Fatalf("Error creating application connector group or ID is empty")
		return
	}

	defer func() {
		if createdAppConnGroup != nil && createdAppConnGroup.ID != "" {
			existingGroup, _, errCheck := appConnectorGroupService.Get(createdAppConnGroup.ID)
			if errCheck == nil && existingGroup != nil {
				_, errDelete := appConnectorGroupService.Delete(createdAppConnGroup.ID)
				if errDelete != nil {
					t.Errorf("Error deleting application connector group: %v", errDelete)
				}
			}
		}
	}()

	// lssConfigFormatService := New(client)
	// lssLogType, _, err := lssConfigFormatService.GetFormats()
	// if err != nil {
	// 	t.Errorf("Error getting LSS Log Type: %v", err)
	// 	return
	// }

	service := New(client)
	lssConfig := &LSSResource{
		LSSConfig: &LSSConfig{
			Name:          name,
			Description:   name,
			Enabled:       true,
			LSSHost:       ipAddress,
			LSSPort:       rPort,
			Format:        "json",
			SourceLogType: "zpn_trans_log",
			UseTLS:        true,
		},
		ConnectorGroups: []ConnectorGroups{
			{
				ID: createdAppConnGroup.ID,
			},
		},
		PolicyRuleResource: &PolicyRuleResource{
			Name:        name,
			Description: name,
			Action:      "ALLOW",
			Conditions: []PolicyRuleResourceConditions{
				{
					Negated:  false,
					Operator: "OR",
					Operands: &[]PolicyRuleResourceOperands{
						{
							ObjectType: "CLIENT_TYPE",
							Values:     []string{"zpn_client_type_exporter"},
						},
						{
							ObjectType: "CLIENT_TYPE",
							Values:     []string{"zpn_client_type_ip_anchoring"},
						},
						{
							ObjectType: "CLIENT_TYPE",
							Values:     []string{"zpn_client_type_zapp"},
						},
						{
							ObjectType: "CLIENT_TYPE",
							Values:     []string{"zpn_client_type_edge_connector"},
						},
						{
							ObjectType: "CLIENT_TYPE",
							Values:     []string{"zpn_client_type_machine_tunnel"},
						},
						{
							ObjectType: "CLIENT_TYPE",
							Values:     []string{"zpn_client_type_browser_isolation"},
						},
						{
							ObjectType: "CLIENT_TYPE",
							Values:     []string{"zpn_client_type_slogger"},
						},
					},
				},
			},
		},
	}

	createdResource, _, err := service.Create(lssConfig)
	if err != nil || createdResource == nil {
		t.Fatalf("Error making POST request: %v or createdResource is nil", err)
	}

	// Fetch the resource again to get full details
	createdResource, _, err = service.Get(createdResource.ID)
	if err != nil {
		t.Fatalf("Error fetching the created resource: %v", err)
	}
	if createdResource.ID == "" {
		t.Error("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.LSSConfig.Name != name {
		t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.LSSConfig.Name)
	}

	// Test resource retrieval
	retrievedResource, _, err := service.Get(createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.LSSConfig.Name != name {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", name, createdResource.LSSConfig.Name)
	}
	// Test resource update
	retrievedResource.LSSConfig.Name = updateName
	_, err = service.Update(createdResource.ID, retrievedResource)
	if err != nil {
		t.Errorf("Error updating resource: %v", err)
	}
	updatedResource, _, err := service.Get(createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%s', but got '%s'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.LSSConfig.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.LSSConfig.Name)
	}
	// Test resource retrieval by name
	retrievedResource, _, err = service.GetByName(updateName)
	if err != nil {
		t.Errorf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.LSSConfig.Name != updateName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, createdResource.LSSConfig.Name)
	}
	// Test resources retrieval
	resources, _, err := service.GetAll()
	if err != nil {
		t.Errorf("Error retrieving groups: %v", err)
	}
	if len(resources) == 0 {
		t.Error("Expected retrieved resources to be non-empty, but got empty slice")
	}
	// check if the created resource is in the list
	found := false
	for _, resource := range resources {
		if resource.ID == createdResource.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected retrieved groups to contain created resource '%s', but it didn't", createdResource.ID)
	}
	// Test resource removal
	_, err = service.Delete(createdResource.ID)
	if err != nil {
		t.Errorf("Error deleting resource: %v", err)
		return
	}

	// Test resource retrieval after deletion
	_, _, err = service.Get(createdResource.ID)
	if err == nil {
		t.Errorf("Expected error retrieving deleted resource, but got nil")
	}
}
*/
