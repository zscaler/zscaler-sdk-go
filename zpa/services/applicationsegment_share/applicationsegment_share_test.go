package applicationsegment_share

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/authdomain"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/microtenants"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/segmentgroup"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/servergroup"
)

func TestApplicationSegmentShare(t *testing.T) {
	// Generate base random strings
	baseName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	baseDescription := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// Step 1: Get available auth domains
	service := services.New(client)

	authDomainList, _, err := authdomain.GetAllAuthDomains(service)
	if err != nil {
		t.Errorf("Error getting auth domains: %v", err)
		return
	}
	if len(authDomainList.AuthDomains) < 2 {
		t.Error("Expected at least two retrieved auth domains, but got less")
		return
	}

	// Function to create microtenant with retries
	createMicrotenantWithRetry := func(name, description string, domains []string) (*microtenants.MicroTenant, error) {
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
			createdMicrotenant, _, err = microtenants.Create(service, microtenant)
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

	// Create Microtenant A
	microtenantAID, err := createMicrotenantWithRetry(baseName+"-microtenantA", baseDescription+"-microtenantA", []string{authDomainList.AuthDomains[0]})
	if err != nil {
		t.Fatalf("Failed to create microtenant A: %v", err)
	}
	defer func() {
		_, err := microtenants.Delete(service, microtenantAID.ID)
		if err != nil {
			t.Errorf("Error deleting microtenant A: %v", err)
		}
	}()

	// Create Microtenant B
	microtenantBID, err := createMicrotenantWithRetry(baseName+"-microtenantB", baseDescription+"-microtenantB", []string{authDomainList.AuthDomains[1]})
	if err != nil {
		t.Fatalf("Failed to create microtenant B: %v", err)
	}
	defer func() {
		_, err := microtenants.Delete(service, microtenantBID.ID)
		if err != nil {
			t.Errorf("Error deleting microtenant B: %v", err)
		}
	}()

	appConnGroupA, _, err := appconnectorgroup.Create(service, appconnectorgroup.AppConnectorGroup{
		Name:                     baseName + "-microtenantA-appconn",
		Description:              baseDescription + "-microtenantA-appconn",
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
		MicroTenantID:            microtenantAID.ID,
	})
	if err != nil {
		t.Fatalf("Error creating app connector group A: %v", err)
	}

	serverGroupA, _, err := servergroup.Create(service, &servergroup.ServerGroup{
		Name:             baseName + "-microtenantA-server",
		Description:      baseDescription + "-microtenantA-server",
		Enabled:          true,
		DynamicDiscovery: true,
		MicroTenantID:    microtenantAID.ID,
		AppConnectorGroups: []servergroup.AppConnectorGroups{
			{ID: appConnGroupA.ID},
		},
	})
	if err != nil {
		t.Fatalf("Error creating server group A: %v", err)
	}

	segGroupA, _, err := segmentgroup.Create(service, &segmentgroup.SegmentGroup{
		Name:          baseName + "-microtenantA-seg",
		Description:   baseDescription + "-microtenantA-seg",
		Enabled:       true,
		MicroTenantID: microtenantAID.ID,
	})
	if err != nil {
		t.Fatalf("Error creating segment group A: %v", err)
	}

	// Step 3: Create a single Application Segment in Microtenant A
	appSegment := applicationsegment.ApplicationSegmentResource{
		Name:                  baseName + "-appseg",
		Description:           baseDescription + "-appseg",
		Enabled:               true,
		SegmentGroupID:        segGroupA.ID,
		IsCnameEnabled:        true,
		BypassType:            "NEVER",
		IcmpAccessType:        "PING_TRACEROUTING",
		HealthReporting:       "ON_ACCESS",
		HealthCheckType:       "DEFAULT",
		TCPKeepAlive:          "1",
		InspectTrafficWithZia: false,
		MatchStyle:            "EXCLUSIVE",
		MicroTenantID:         microtenantAID.ID,
		DomainNames:           []string{"test.example.com"},
		ServerGroups: []applicationsegment.AppServerGroups{
			{ID: serverGroupA.ID},
		},
		TCPAppPortRange: []common.NetworkPorts{
			{
				From: "6443",
				To:   "6443",
			},
		},
	}
	createdAppSegment, _, err := applicationsegment.Create(service, appSegment)
	if err != nil {
		t.Fatalf("Error creating application segment: %v", err)
	}

	// Step 4: Share Application Segment from Microtenant A to Microtenant B
	shareRequest := AppSegmentSharedToMicrotenant{
		ApplicationID:       createdAppSegment.ID,
		ShareToMicrotenants: []string{microtenantBID.ID},
		MicroTenantID:       microtenantAID.ID,
	}

	resp, err := AppSegmentMicrotenantShare(service, createdAppSegment.ID, shareRequest)
	if err != nil {
		t.Fatalf("Error sharing application segment: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Failed to share application segment, status code: %d", resp.StatusCode)
	}
}
