package applicationsegment_move

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/microtenants"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/servergroup"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestApplicationSegmentMove(t *testing.T) {
	// Generate base random strings
	baseName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	baseDescription := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Function to create microtenant with retries
	// Function to create microtenant with retries
	createMicrotenantWithRetry := func(name, description string) (*microtenants.MicroTenant, error) {
		microtenant := microtenants.MicroTenant{
			Name:                       name,
			Description:                description,
			Enabled:                    true,
			PrivilegedApprovalsEnabled: true,
			CriteriaAttribute:          "AuthDomain",
		}
		var createdMicrotenant *microtenants.MicroTenant
		var err error
		for i := 0; i < 3; i++ { // Retry up to 3 times
			createdMicrotenant, _, err = microtenants.Create(context.Background(), service, microtenant)
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
	createdMicrotenant, err := createMicrotenantWithRetry(baseName+"-microtenant", baseDescription+"-microtenant")
	if err != nil {
		t.Fatalf("Failed to create microtenant: %v", err)
	}
	defer func() {
		_, err := microtenants.Delete(context.Background(), service, createdMicrotenant.ID)
		if err != nil {
			t.Errorf("Error deleting microtenant: %v", err)
		}
	}()

	microtenantID := createdMicrotenant.ID

	// Step 2: Create resources in the new Microtenant
	appConnGroup, _, err := appconnectorgroup.Create(context.Background(), service, appconnectorgroup.AppConnectorGroup{
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

	serverGroup, _, err := servergroup.Create(context.Background(), service, &servergroup.ServerGroup{
		Name:             baseName + "-microtenant-server",
		Description:      baseDescription + "-microtenant-server",
		Enabled:          true,
		DynamicDiscovery: true,
		MicroTenantID:    microtenantID,
		AppConnectorGroups: []appconnectorgroup.AppConnectorGroup{
			{ID: appConnGroup.ID},
		},
	})
	if err != nil {
		t.Fatalf("Error creating server group: %v", err)
	}

	segGroup := segmentgroup.SegmentGroup{
		Name:          baseName + "-microtenant-seg",
		Description:   baseDescription + "-microtenant-seg",
		MicroTenantID: microtenantID,
	}
	createdSegGroup, _, err := segmentgroup.Create(context.Background(), service, &segGroup)
	if err != nil {
		t.Fatalf("Error creating segment group: %v", err)
	}

	// Step 3: Create resources in the Parent Tenant
	appConnGroupParent, _, err := appconnectorgroup.Create(context.Background(), service, appconnectorgroup.AppConnectorGroup{
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

	serverGroupParent, _, err := servergroup.Create(context.Background(), service, &servergroup.ServerGroup{
		Name:             baseName + "-parent-server",
		Description:      baseDescription + "-parent-server",
		Enabled:          true,
		DynamicDiscovery: true,
		AppConnectorGroups: []appconnectorgroup.AppConnectorGroup{
			{ID: appConnGroupParent.ID},
		},
	})
	if err != nil {
		t.Fatalf("Error creating server group in parent tenant: %v", err)
	}

	segGroupParent := segmentgroup.SegmentGroup{
		Name:        baseName + "-parent-seg",
		Enabled:     true,
		Description: baseDescription + "-parent-seg",
	}
	createdSegGroupParent, _, err := segmentgroup.Create(context.Background(), service, &segGroupParent)
	if err != nil {
		t.Fatalf("Error creating segment group in parent tenant: %v", err)
	}

	// Step 4: Create Application Segment in the Parent Tenant
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
		ServerGroups: []servergroup.ServerGroup{
			{ID: serverGroupParent.ID},
		},
		TCPAppPortRange: []common.NetworkPorts{
			{
				From: "5443",
				To:   "5443",
			},
		},
	}
	createdAppSegment, _, err := applicationsegment.Create(context.Background(), service, appSegment)
	if err != nil {
		t.Fatalf("Error creating application segment: %v", err)
	}

	// Step 5: Move Application Segment to the Microtenant
	moveRequest := AppSegmentMicrotenantMoveRequest{
		ApplicationID:        createdAppSegment.ID,
		TargetSegmentGroupID: createdSegGroup.ID,
		TargetMicrotenantID:  microtenantID,
		TargetServerGroupID:  serverGroup.ID,
	}

	resp, err := AppSegmentMicrotenantMove(context.Background(), service, createdAppSegment.ID, moveRequest)
	if err != nil {
		t.Fatalf("Error moving application segment: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to move application segment, status code: %d", resp.StatusCode)
	}

	// Cleanup: Resources created in Parent Tenant (except those in Microtenant)
	defer func() {
		time.Sleep(time.Second * 2) // Sleep for 2 seconds before deletion
		_, err = servergroup.Delete(context.Background(), service, serverGroupParent.ID)
		if err != nil {
			t.Errorf("Error deleting server group in parent tenant: %v", err)
		}
		_, err = appconnectorgroup.Delete(context.Background(), service, appConnGroupParent.ID)
		if err != nil {
			t.Errorf("Error deleting app connector group in parent tenant: %v", err)
		}
		_, err = segmentgroup.Delete(context.Background(), service, createdSegGroupParent.ID)
		if err != nil {
			t.Errorf("Error deleting segment group in parent tenant: %v", err)
		}
	}()
}
