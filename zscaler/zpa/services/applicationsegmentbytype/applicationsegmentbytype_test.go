package applicationsegmentbytype

/*
import (
	"log"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentinspection"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentpra"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/bacertificate"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/browseraccess"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestCreateApplicationSegmentPRA(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	segmentGroupName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	appGroup := segmentgroup.SegmentGroup{
		Name:        segmentGroupName,
		Description: segmentGroupName,
		Enabled:     true,
	}
	createdAppGroup, _, err := segmentgroup.Create(service, &appGroup)
	if err != nil {
		t.Errorf("Error creating application segment group: %v", err)
		return
	}

	appSegment := applicationsegmentpra.AppSegmentPRA{
		Name:             name,
		Description:      name,
		Enabled:          true,
		SegmentGroupID:   createdAppGroup.ID,
		SegmentGroupName: createdAppGroup.Name,
		IsCnameEnabled:   true,
		BypassType:       "NEVER",
		IcmpAccessType:   "PING_TRACEROUTING",
		HealthReporting:  "ON_ACCESS",
		HealthCheckType:  "DEFAULT",
		TCPKeepAlive:     "1",
		DomainNames:      []string{"rdp_pra1.bd-hashicorp.com"},
		TCPAppPortRange: []common.NetworkPorts{
			{
				From: "3390",
				To:   "3390",
			},
		},
		CommonAppsDto: applicationsegmentpra.CommonAppsDto{
			AppsConfig: []applicationsegmentpra.AppsConfig{
				{
					Name:                name,
					Description:         name,
					Enabled:             true,
					AppTypes:            []string{"SECURE_REMOTE_ACCESS"},
					ApplicationPort:     "3390",
					ApplicationProtocol: "RDP",
					ConnectionSecurity:  "ANY",
					Domain:              "rdp_pra1.bd-hashicorp.com",
				},
			},
		},
	}

	// Test resource creation
	createdResource, _, err := applicationsegmentpra.Create(service, appSegment)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making POST request: %v", err)
	}

	if createdResource.ID == "" {
		t.Error("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.Name)
	}
}

func TestAppSegmentInspectionInspection(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	segmentGroupName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	appGroup := segmentgroup.SegmentGroup{
		Name:        segmentGroupName,
		Description: segmentGroupName,
		Enabled:     true,
	}
	createdAppGroup, _, err := segmentgroup.Create(service, &appGroup)
	if err != nil {
		t.Errorf("Error creating application segment group: %v", err)
		return
	}

	certificateList, _, err := bacertificate.GetAll(service)
	if err != nil {
		t.Errorf("Error getting saml attributes: %v", err)
		return
	}
	if len(certificateList) == 0 {
		t.Error("Expected retrieved saml attributes to be non-empty, but got empty slice")
	}

	appSegment := applicationsegmentinspection.AppSegmentInspection{
		Name:             name,
		Description:      name,
		Enabled:          true,
		SegmentGroupID:   createdAppGroup.ID,
		SegmentGroupName: createdAppGroup.Name,
		IsCnameEnabled:   true,
		BypassType:       "NEVER",
		ICMPAccessType:   "PING_TRACEROUTING",
		HealthReporting:  "ON_ACCESS",
		HealthCheckType:  "DEFAULT",
		TCPKeepAlive:     "1",
		DomainNames:      []string{"server1.bd-hashicorp.com"},
		TCPAppPortRange: []common.NetworkPorts{
			{
				From: "8444",
				To:   "8444",
			},
		},
		CommonAppsDto: applicationsegmentinspection.CommonAppsDto{
			AppsConfig: []applicationsegmentinspection.AppsConfig{
				{
					Name:                name,
					Description:         name,
					Enabled:             true,
					AppTypes:            []string{"INSPECT"},
					ApplicationPort:     "8444",
					ApplicationProtocol: "HTTPS",
					Domain:              "server1.bd-hashicorp.com",
					CertificateID:       certificateList[0].ID,
				},
			},
		},
	}

	// Test resource creation
	createdResource, _, err := applicationsegmentinspection.Create(service, appSegment)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making POST request: %v", err)
	}

	if createdResource.ID == "" {
		t.Error("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.Name)
	}
}

func TestBaApplicationSegment(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	segmentGroupName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	appGroup := segmentgroup.SegmentGroup{
		Name:        segmentGroupName,
		Description: segmentGroupName,
		Enabled:     true,
	}
	createdAppGroup, _, err := segmentgroup.Create(service, &appGroup)
	if err != nil {
		t.Errorf("Error creating application segment group: %v", err)
		return
	}

	certificateList, _, err := bacertificate.GetAll(service)
	if err != nil {
		t.Errorf("Error getting certificates: %v", err)
		return
	}
	if len(certificateList) == 0 {
		t.Error("Expected retrieved certificates to be non-empty, but got empty slice")
		return
	}

	appSegment := browseraccess.BrowserAccess{
		Name:             name,
		Description:      name,
		Enabled:          true,
		SegmentGroupID:   createdAppGroup.ID,
		SegmentGroupName: createdAppGroup.Name,
		IsCnameEnabled:   true,
		BypassType:       "NEVER",
		ICMPAccessType:   "PING_TRACEROUTING",
		HealthReporting:  "ON_ACCESS",
		HealthCheckType:  "DEFAULT",
		TCPKeepAlive:     "1",
		DomainNames:      []string{name + ".bd-hashicorp"},
		ClientlessApps: []browseraccess.ClientlessApps{
			{
				Name:                name + ".bd-hashicorp",
				Description:         name + ".bd-hashicorp",
				Enabled:             true,
				TrustUntrustedCert:  true,
				Domain:              name + ".bd-hashicorp",
				ApplicationProtocol: "HTTPS",
				ApplicationPort:     "7443",
				CertificateID:       certificateList[0].ID,
			},
		},
		TCPAppPortRange: []common.NetworkPorts{
			{
				From: "7443",
				To:   "7443",
			},
		},
	}
	// Test resource creation
	createdResource, _, err := browseraccess.Create(service, appSegment)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making POST request: %v", err)
	}

	if createdResource.ID == "" {
		t.Error("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.Name)
	}
}

func TestGetByApplicationType(t *testing.T) {
	client, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	defer func() {
		err := cleanupResources(client)
		if err != nil {
			t.Errorf("Error during cleanup: %v", err)
		}
	}()

	expandAll := true
	applicationTypes := []string{"BROWSER_ACCESS", "INSPECT", "SECURE_REMOTE_ACCESS"}

	// Test valid application types with and without appName
	for _, applicationType := range applicationTypes {
		t.Run("Without appName "+applicationType, func(t *testing.T) {
			retrievedByTypeResources, _, err := GetByApplicationType("", applicationType, expandAll)
			if err != nil {
				t.Errorf("Error retrieving resource by application type '%s': %v", applicationType, err)
			}
			if len(retrievedByTypeResources) == 0 {
				t.Logf("No resources found for application type '%s'", applicationType)
			} else {
				t.Logf("Retrieved %d resources for application type '%s'", len(retrievedByTypeResources), applicationType)
			}
		})

		t.Run("With appName "+applicationType, func(t *testing.T) {
			appName := "example-app"
			retrievedByTypeResources, _, err := GetByApplicationType(appName, applicationType, expandAll)
			if err != nil {
				t.Errorf("Error retrieving resource by application type '%s' with appName '%s': %v", applicationType, appName, err)
			}
			if len(retrievedByTypeResources) == 0 {
				t.Logf("No resources found for application type '%s' with appName '%s'", applicationType, appName)
			} else {
				t.Logf("Retrieved %d resources for application type '%s' with appName '%s'", len(retrievedByTypeResources), applicationType, appName)
			}
		})
	}

	// Test invalid application type
	t.Run("Invalid applicationType", func(t *testing.T) {
		invalidApplicationType := "INVALID_TYPE"
		_, _, err := GetByApplicationType("", invalidApplicationType, expandAll)
		if err == nil {
			t.Errorf("Expected error for invalid application type '%s', but got nil", invalidApplicationType)
		}
	})
}

func cleanupResources(client *zscaler.Client) error {

	resources, _, err := applicationsegmentpra.GetAll()
	if err != nil {
		log.Printf("[ERROR] Failed to get application segment pra: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		// log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, err := applicationsegmentpra.Delete(r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete application segment pra with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}

	inspectionResources, _, err := applicationsegmentinspection.GetAll()
	if err != nil {
		log.Printf("[ERROR] Failed to get application segment inspection: %v", err)
		return err
	}

	for _, r := range inspectionResources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		// log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, err := applicationsegmentinspection.Delete(service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete application segment inspection with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}

	segmentGroupResources, _, err := segmentgroup.GetAll(service)
	if err != nil {
		log.Printf("[ERROR] Failed to get segment groups: %v", err)
		return err
	}

	for _, r := range segmentGroupResources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		// log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, err := segmentgroup.Delete(service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete segment groups with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}
*/
