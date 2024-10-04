package applicationsegmentbytype

/*
import (

	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/applicationsegmentinspection"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/applicationsegmentpra"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/bacertificate"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/browseraccess"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/segmentgroup"

)

	func TestCreateApplicationSegmentPRA(t *testing.T) {
		name_pra := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha) + "-pra"
		name_inspection := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha) + "-inspection"
		name_baApp := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha) + "-baApp"
		segmentGroupName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
		client, err := tests.NewZpaClient()
		if err != nil {
			t.Errorf("Error creating client: %v", err)
			return
		}

		service := services.New(client)

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

		appSegmentPra := applicationsegmentpra.AppSegmentPRA{
			Name:            name_pra,
			Description:     name_pra,
			Enabled:         true,
			SegmentGroupID:  createdAppGroup.ID,
			IsCnameEnabled:  true,
			BypassType:      "NEVER",
			IcmpAccessType:  "PING_TRACEROUTING",
			HealthReporting: "ON_ACCESS",
			HealthCheckType: "DEFAULT",
			TCPKeepAlive:    "1",
			DomainNames:     []string{"rdp_pra1.bd-hashicorp.com"},
			TCPAppPortRange: []common.NetworkPorts{
				{
					From: "3390",
					To:   "3390",
				},
			},
			CommonAppsDto: applicationsegmentpra.CommonAppsDto{
				AppsConfig: []applicationsegmentpra.AppsConfig{
					{
						Name:                name_pra,
						Description:         name_pra,
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
		createdPraResource, _, err := applicationsegmentpra.Create(service, appSegmentPra)
		// Check if the request was successful
		if err != nil {
			t.Errorf("Error making POST request: %v", err)
		}

		if createdPraResource.ID == "" {
			t.Error("Expected created resource ID to be non-empty, but got ''")
		}
		if createdPraResource.Name != name_pra {
			t.Errorf("Expected created resource name '%s', but got '%s'", name_pra, createdPraResource.Name)
		}

		// Create Application Segment Inspection
		certificateList, _, err := bacertificate.GetAll(service)
		if err != nil {
			t.Errorf("Error getting saml attributes: %v", err)
			return
		}
		if len(certificateList) == 0 {
			t.Error("Expected retrieved saml attributes to be non-empty, but got empty slice")
		}

		appSegmentInspection := applicationsegmentinspection.AppSegmentInspection{
			Name:            name_inspection,
			Description:     name_inspection,
			Enabled:         true,
			SegmentGroupID:  createdAppGroup.ID,
			IsCnameEnabled:  true,
			BypassType:      "NEVER",
			ICMPAccessType:  "PING_TRACEROUTING",
			HealthReporting: "ON_ACCESS",
			HealthCheckType: "DEFAULT",
			TCPKeepAlive:    "1",
			DomainNames:     []string{"server1.bd-hashicorp.com"},
			TCPAppPortRange: []common.NetworkPorts{
				{
					From: "8444",
					To:   "8444",
				},
			},
			CommonAppsDto: applicationsegmentinspection.CommonAppsDto{
				AppsConfig: []applicationsegmentinspection.AppsConfig{
					{
						Name:                name_inspection,
						Description:         name_inspection,
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
		createdInspectionResource, _, err := applicationsegmentinspection.Create(service, appSegmentInspection)
		// Check if the request was successful
		if err != nil {
			t.Errorf("Error making POST request: %v", err)
		}

		if createdInspectionResource.ID == "" {
			t.Error("Expected created resource ID to be non-empty, but got ''")
		}
		if createdInspectionResource.Name != name_inspection {
			t.Errorf("Expected created resource name '%s', but got '%s'", name_inspection, createdInspectionResource.Name)
		}

		BaAppSegment := browseraccess.BrowserAccess{
			Name:            name_baApp,
			Description:     name_baApp,
			Enabled:         true,
			SegmentGroupID:  createdAppGroup.ID,
			IsCnameEnabled:  true,
			BypassType:      "NEVER",
			ICMPAccessType:  "PING_TRACEROUTING",
			HealthReporting: "ON_ACCESS",
			HealthCheckType: "DEFAULT",
			TCPKeepAlive:    "1",
			DomainNames:     []string{name_baApp + ".bd-hashicorp"},
			ClientlessApps: []browseraccess.ClientlessApps{
				{
					Name:                name_baApp + ".bd-hashicorp",
					Description:         name_baApp + ".bd-hashicorp",
					Enabled:             true,
					TrustUntrustedCert:  true,
					Domain:              name_baApp + ".bd-hashicorp",
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
		createdBaResource, _, err := browseraccess.Create(service, BaAppSegment)
		// Check if the request was successful
		if err != nil {
			t.Errorf("Error making POST request: %v", err)
		}

		if createdBaResource.ID == "" {
			t.Error("Expected created resource ID to be non-empty, but got ''")
		}
		if createdBaResource.Name != name_baApp {
			t.Errorf("Expected created resource name '%s', but got '%s'", name_baApp, createdBaResource.Name)
		}
	}

/*

	func TestAppSegmentInspectionInspection(t *testing.T) {
		name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
		segmentGroupName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
		client, err := tests.NewZpaClient()
		if err != nil {
			t.Errorf("Error creating client: %v", err)
			return
		}

		service := services.New(client)

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

		appSegmentInspection := applicationsegmentinspection.AppSegmentInspection{
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
		createdResource, _, err := applicationsegmentinspection.Create(service, appSegmentInspection)
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
		client, err := tests.NewZpaClient()
		if err != nil {
			t.Errorf("Error creating client: %v", err)
			return
		}

		service := services.New(client)

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

		BaAppSegment := browseraccess.BrowserAccess{
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
		createdBaResource, _, err := browseraccess.Create(service, BaAppSegment)
		// Check if the request was successful
		if err != nil {
			t.Errorf("Error making POST request: %v", err)
		}

		if createdBaResource.ID == "" {
			t.Error("Expected created resource ID to be non-empty, but got ''")
		}
		if createdBaResource.Name != name {
			t.Errorf("Expected created resource name '%s', but got '%s'", name, createdBaResource.Name)
		}
	}

func TestGetByApplicationType(t *testing.T) {
	client, err := tests.NewZpaClient()
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

	service := services.New(client)
	expandAll := true
	applicationTypes := []string{"BROWSER_ACCESS", "INSPECT", "SECURE_REMOTE_ACCESS"}

	// Test valid application types with and without appName
	for _, applicationType := range applicationTypes {
		t.Run("Without appName "+applicationType, func(t *testing.T) {
			retrievedByTypeResources, _, err := GetByApplicationType(service, "", applicationType, expandAll)
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
			allResources, _, err := GetByApplicationType(service, "", applicationType, expandAll)
			if err != nil {
				t.Errorf("Error retrieving resources by application type '%s': %v", applicationType, err)
			}

			var filteredResources []AppSegmentBaseAppDto
			for _, resource := range allResources {
				if strings.Contains(resource.Name, "tests-") {
					filteredResources = append(filteredResources, resource)
				}
			}

			if len(filteredResources) == 0 {
				t.Logf("No resources found for application type '%s' containing 'tests-'", applicationType)
			} else {
				t.Logf("Retrieved %d resources for application type '%s' containing 'tests-'", len(filteredResources), applicationType)
			}
		})

	}

	// Test invalid application type
	t.Run("Invalid applicationType", func(t *testing.T) {
		invalidApplicationType := "INVALID_TYPE"
		_, _, err := GetByApplicationType(service, "", invalidApplicationType, expandAll)
		if err == nil {
			t.Errorf("Expected error for invalid application type '%s', but got nil", invalidApplicationType)
		}
	})
}

func cleanupResources(client *zpa.Client) error {
	service := services.New(client)

	// First delete all Application Segment PRA resources
	praResources, _, err := applicationsegmentpra.GetAll(service)
	if err != nil {
		log.Printf("[ERROR] Failed to get application segment PRA: %v", err)
		return err
	}

	for _, r := range praResources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		_, err := applicationsegmentpra.Delete(service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete application segment PRA with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}

	// Then delete all Application Segment Inspection resources
	inspectionResources, _, err := applicationsegmentinspection.GetAll(service)
	if err != nil {
		log.Printf("[ERROR] Failed to get application segment inspection: %v", err)
		return err
	}

	for _, r := range inspectionResources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		_, err := applicationsegmentinspection.Delete(service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete application segment inspection with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}

	// Finally, delete the Segment Groups
	segmentGroupResources, _, err := segmentgroup.GetAll(service)
	if err != nil {
		log.Printf("[ERROR] Failed to get segment groups: %v", err)
		return err
	}

	for _, r := range segmentGroupResources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		_, err := segmentgroup.Delete(service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete segment group with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}
*/
