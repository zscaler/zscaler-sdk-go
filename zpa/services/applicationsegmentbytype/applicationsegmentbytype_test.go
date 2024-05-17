package applicationsegmentbytype

import (
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/applicationsegmentinspection"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/applicationsegmentpra"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/bacertificate"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/segmentgroup"
)

func TestCreateApplicationSegmentPRA(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	segmentGroupName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	// create application segment group for testing
	appGroupService := segmentgroup.New(client)
	appGroup := segmentgroup.SegmentGroup{
		Name:        segmentGroupName,
		Description: segmentGroupName,
	}
	createdAppGroup, _, err := appGroupService.Create(&appGroup)
	if err != nil {
		t.Errorf("Error creating application segment group: %v", err)
		return
	}

	service := applicationsegmentpra.New(client)
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
		DomainNames:      []string{"rdp_pra.bd-hashicorp.com"},
		TCPAppPortRange: []common.NetworkPorts{
			{
				From: "3389",
				To:   "3389",
			},
		},
		CommonAppsDto: applicationsegmentpra.CommonAppsDto{
			AppsConfig: []applicationsegmentpra.AppsConfig{
				{
					Name:                name,
					Description:         name,
					Enabled:             true,
					AppTypes:            []string{"SECURE_REMOTE_ACCESS"},
					ApplicationPort:     "3389",
					ApplicationProtocol: "RDP",
					ConnectionSecurity:  "ANY",
					Domain:              "rdp_pra.bd-hashicorp.com",
				},
			},
		},
	}

	// Test resource creation
	createdResource, _, err := service.Create(appSegment)
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
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	// create application segment group for testing
	appGroupService := segmentgroup.New(client)
	appGroup := segmentgroup.SegmentGroup{
		Name:        segmentGroupName,
		Description: segmentGroupName,
	}
	createdAppGroup, _, err := appGroupService.Create(&appGroup)
	if err != nil {
		t.Errorf("Error creating application segment group: %v", err)
		return
	}

	baCertificateService := bacertificate.New(client)
	certificateList, _, err := baCertificateService.GetAll()
	if err != nil {
		t.Errorf("Error getting saml attributes: %v", err)
		return
	}
	if len(certificateList) == 0 {
		t.Error("Expected retrieved saml attributes to be non-empty, but got empty slice")
	}

	service := applicationsegmentinspection.New(client)
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
				From: "8443",
				To:   "8443",
			},
		},
		CommonAppsDto: applicationsegmentinspection.CommonAppsDto{
			AppsConfig: []applicationsegmentinspection.AppsConfig{
				{
					Name:                name,
					Description:         name,
					Enabled:             true,
					AppTypes:            []string{"INSPECT"},
					ApplicationPort:     "8443",
					ApplicationProtocol: "HTTPS",
					Domain:              "server1.bd-hashicorp.com",
					CertificateID:       certificateList[0].ID,
				},
			},
		},
	}

	// Test resource creation
	createdResource, _, err := service.Create(appSegment)
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

	service := New(client)
	expandAll := true
	applicationTypes := []string{"INSPECT", "SECURE_REMOTE_ACCESS"}

	for _, applicationType := range applicationTypes {
		t.Run(applicationType, func(t *testing.T) {
			retrievedByTypeResources, _, err := service.GetByApplicationType("", applicationType, expandAll)
			if err != nil {
				t.Errorf("Error retrieving resource by application type '%s': %v", applicationType, err)
			}
			if len(retrievedByTypeResources) == 0 {
				t.Errorf("Expected to retrieve at least one resource by application type '%s', but got none", applicationType)
			}
		})
	}
}

func cleanupResources(client *zpa.Client) error {
	appSegmentPra := applicationsegmentpra.New(client)
	resources, _, err := appSegmentPra.GetAll()
	if err != nil {
		log.Printf("[ERROR] Failed to get app connector groups: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, err := appSegmentPra.Delete(r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete app connector group with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}

	appSegmentInspection := applicationsegmentinspection.New(client)
	inspectionResources, _, err := appSegmentInspection.GetAll()
	if err != nil {
		log.Printf("[ERROR] Failed to get app connector groups: %v", err)
		return err
	}

	for _, r := range inspectionResources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, err := appSegmentInspection.Delete(r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete app connector group with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}

	segmentGroupservice := segmentgroup.New(client)
	segmentGroupResources, _, err := segmentGroupservice.GetAll()
	if err != nil {
		log.Printf("[ERROR] Failed to get app connector groups: %v", err)
		return err
	}

	for _, r := range segmentGroupResources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, err := segmentGroupservice.Delete(r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete app connector group with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}
