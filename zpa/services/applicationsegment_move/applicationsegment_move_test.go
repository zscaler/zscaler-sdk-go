package applicationsegment_move

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/authdomain"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/microtenants"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/segmentgroup"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/servergroup"
)

func TestApplicationSegmentMove(t *testing.T) {
	// Generate base random strings
	baseName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	baseDescription := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	rPort := strconv.Itoa(acctest.RandIntRange(1000, 9999))

	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// Step 1: Get available auth domains
	authDomainService := authdomain.New(client)
	authDomainList, _, err := authDomainService.GetAllAuthDomains()
	if err != nil {
		t.Errorf("Error getting auth domains: %v", err)
		return
	}
	if len(authDomainList.AuthDomains) == 0 {
		t.Error("Expected retrieved auth domains to be non-empty, but got empty slice")
		return
	}

	// Function to create microtenant with retries
	createMicrotenantWithRetry := func(name, description string, domains []string) (*microtenants.MicroTenant, error) {
		microtenantService := microtenants.New(client)
		microtenant := microtenants.MicroTenant{
			Name:                       name,
			Description:                description,
			Enabled:                    true,
			PrivilegedApprovalsEnabled: true,
			CriteriaAttribute:          "AuthDomain",
			CriteriaAttributeValues:    domains,
		}
		var createdMicrotenant *microtenants.MicroTenant
		var err error
		for i := 0; i < 3; i++ { // Retry up to 3 times
			createdMicrotenant, _, err = microtenantService.Create(microtenant)
			if err == nil {
				break
			}
			if strings.Contains(err.Error(), "domains.already.exists.in.other.microtenant") {
				t.Logf("Retry %d: Domain already exists in another microtenant, retrying...", i+1)
				time.Sleep(time.Second * 2) // Sleep for 2 seconds before retrying
				continue
			}
			break
		}
		return createdMicrotenant, err
	}

	// Create Microtenant
	createdMicrotenant, err := createMicrotenantWithRetry(baseName+"-microtenant", baseDescription+"-microtenant", []string{authDomainList.AuthDomains[0]})
	if err != nil {
		t.Fatalf("Failed to create microtenant: %v", err)
	}
	defer func() {
		microtenantService := microtenants.New(client)
		_, err := microtenantService.Delete(createdMicrotenant.ID)
		if err != nil {
			t.Errorf("Error deleting microtenant: %v", err)
		}
	}()

	microtenantID := createdMicrotenant.ID

	// Step 2: Create resources in the new Microtenant
	appConnGroupService := appconnectorgroup.New(client)
	appConnGroup, _, err := appConnGroupService.Create(appconnectorgroup.AppConnectorGroup{
		Name:                     baseName + "-microtenant-appconn",
		Description:              baseDescription + "-microtenant-appconn",
		Enabled:                  true,
		CityCountry:              "San Jose, US",
		Latitude:                 "37.33874",
		Longitude:                "-121.8852525",
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
		MicroTenantID:            microtenantID,
	})
	if err != nil {
		t.Fatalf("Error creating app connector group: %v", err)
	}

	serverGroupService := servergroup.New(client)
	serverGroup, _, err := serverGroupService.Create(&servergroup.ServerGroup{
		Name:             baseName + "-microtenant-server",
		Description:      baseDescription + "-microtenant-server",
		Enabled:          true,
		DynamicDiscovery: true,
		MicroTenantID:    microtenantID,
		AppConnectorGroups: []servergroup.AppConnectorGroups{
			{ID: appConnGroup.ID},
		},
	})
	if err != nil {
		t.Fatalf("Error creating server group: %v", err)
	}

	segGroupService := segmentgroup.New(client)
	segGroup := segmentgroup.SegmentGroup{
		Name:          baseName + "-microtenant-seg",
		Description:   baseDescription + "-microtenant-seg",
		MicroTenantID: microtenantID,
	}
	createdSegGroup, _, err := segGroupService.Create(&segGroup)
	if err != nil {
		t.Fatalf("Error creating segment group: %v", err)
	}

	// Step 3: Create resources in the Parent Tenant
	appConnGroupServiceParent := appconnectorgroup.New(client)
	appConnGroupParent, _, err := appConnGroupServiceParent.Create(appconnectorgroup.AppConnectorGroup{
		Name:                     baseName + "-parent-appconn",
		Description:              baseDescription + "-parent-appconn",
		Enabled:                  true,
		CityCountry:              "San Jose, US",
		Latitude:                 "37.33874",
		Longitude:                "-121.8852525",
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
	})
	if err != nil {
		t.Fatalf("Error creating app connector group in parent tenant: %v", err)
	}

	serverGroupServiceParent := servergroup.New(client)
	serverGroupParent, _, err := serverGroupServiceParent.Create(&servergroup.ServerGroup{
		Name:             baseName + "-parent-server",
		Description:      baseDescription + "-parent-server",
		Enabled:          true,
		DynamicDiscovery: true,
		AppConnectorGroups: []servergroup.AppConnectorGroups{
			{ID: appConnGroupParent.ID},
		},
	})
	if err != nil {
		t.Fatalf("Error creating server group in parent tenant: %v", err)
	}

	segGroupServiceParent := segmentgroup.New(client)
	segGroupParent := segmentgroup.SegmentGroup{
		Name:        baseName + "-parent-seg",
		Enabled:     true,
		Description: baseDescription + "-parent-seg",
	}
	createdSegGroupParent, _, err := segGroupServiceParent.Create(&segGroupParent)
	if err != nil {
		t.Fatalf("Error creating segment group in parent tenant: %v", err)
	}

	// Step 4: Create Application Segment in the Parent Tenant
	appSegmentService := applicationsegment.New(client)
	appSegment := applicationsegment.ApplicationSegmentResource{
		Name:                  baseName + "-parent-appseg",
		Description:           baseDescription + "-parent-appseg",
		Enabled:               true,
		SegmentGroupID:        createdSegGroupParent.ID,
		IsCnameEnabled:        true,
		BypassType:            "NEVER",
		IcmpAccessType:        "PING_TRACEROUTING",
		HealthReporting:       "ON_ACCESS",
		HealthCheckType:       "DEFAULT",
		TCPKeepAlive:          "1",
		InspectTrafficWithZia: false,
		MatchStyle:            "EXCLUSIVE",
		DomainNames:           []string{"test.example.com"},
		ServerGroups: []applicationsegment.AppServerGroups{
			{ID: serverGroupParent.ID},
		},
		TCPAppPortRange: []common.NetworkPorts{
			{
				From: rPort,
				To:   rPort,
			},
		},
	}
	createdAppSegment, _, err := appSegmentService.Create(appSegment)
	if err != nil {
		t.Fatalf("Error creating application segment: %v", err)
	}

	// Step 5: Move Application Segment to the Microtenant
	appSegmentMoveService := New(client)
	moveRequest := AppSegmentMicrotenantMove{
		ApplicationID:        createdAppSegment.ID,
		TargetSegmentGroupID: createdSegGroup.ID,
		TargetMicrotenantID:  microtenantID,
		TargetServerGroupID:  serverGroup.ID,
	}

	resp, err := appSegmentMoveService.AppSegmentMicrotenantMove(createdAppSegment.ID, moveRequest)
	if err != nil {
		t.Fatalf("Error moving application segment: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Failed to move application segment, status code: %d", resp.StatusCode)
	}

	// Cleanup: Resources created in Parent Tenant (except those in Microtenant)
	defer func() {
		time.Sleep(time.Second * 2) // Sleep for 2 seconds before deletion
		_, err = serverGroupServiceParent.Delete(serverGroupParent.ID)
		if err != nil {
			t.Errorf("Error deleting server group in parent tenant: %v", err)
		}
		_, err = appConnGroupServiceParent.Delete(appConnGroupParent.ID)
		if err != nil {
			t.Errorf("Error deleting app connector group in parent tenant: %v", err)
		}
		_, err = segGroupServiceParent.Delete(createdSegGroupParent.ID)
		if err != nil {
			t.Errorf("Error deleting segment group in parent tenant: %v", err)
		}
	}()
}
