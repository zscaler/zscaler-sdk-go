package praconsole

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/applicationsegmentpra"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/bacertificate"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/privilegedremoteaccess/praportal"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/segmentgroup"
)

func TestPRAConsole(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	updateName := "updated-" + acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// create segment group for testing
	segGroupService := segmentgroup.New(client)
	appGroup := segmentgroup.SegmentGroup{
		Name:        name,
		Description: name,
		Enabled:     true,
	}
	createdSegGroup, _, err := segGroupService.Create(&appGroup)
	if err != nil {
		t.Errorf("Error creating segment group: %v", err)
		return
	}
	defer func() {
		time.Sleep(time.Second * 2) // Sleep for 2 seconds before deletion
		_, _, getErr := segGroupService.Get(createdSegGroup.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := segGroupService.Delete(createdSegGroup.ID)
			if err != nil {
				t.Errorf("Error deleting segment group: %v", err)
			}
		}
	}()

	// create pra application segment for testing
	praSegmentService := applicationsegmentpra.New(client)
	praAppSeg := applicationsegmentpra.AppSegmentPRA{
		Name:            name,
		Description:     name,
		Enabled:         true,
		SegmentGroupID:  createdSegGroup.ID,
		IsCnameEnabled:  true,
		BypassType:      "NEVER",
		IcmpAccessType:  "PING_TRACEROUTING",
		HealthReporting: "ON_ACCESS",
		HealthCheckType: "DEFAULT",
		TCPKeepAlive:    "1",
		DomainNames:     []string{"rdp_pra.example.com", "ssh_pra.example.com"},
		TCPAppPortRange: []common.NetworkPorts{
			{
				From: "3389",
				To:   "3389",
			},
			{
				From: "22",
				To:   "22",
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
					Domain:              "rdp_pra.example.com",
				},
				{
					Name:                name,
					Description:         name,
					Enabled:             true,
					AppTypes:            []string{"SECURE_REMOTE_ACCESS"},
					ApplicationPort:     "22",
					ApplicationProtocol: "SSH",
					Domain:              "ssh_pra.example.com",
				},
			},
		},
	}
	createdpraAppSeg, _, err := praSegmentService.Create(praAppSeg)
	if err != nil {
		t.Errorf("Error creating pra application segment: %v", err)
		return
	}

	// Assuming the praSegmentService.Get correctly returns the payload as described
	retrievedpraAppSeg, _, err := praSegmentService.Get(createdpraAppSeg.ID)
	if err != nil {
		t.Errorf("Error retrieving created pra application segment: %v", err)
		return
	}

	baCertificateService := bacertificate.New(client)
	baCertList, _, err := baCertificateService.GetAll()
	if err != nil {
		t.Errorf("Error getting certificates: %v", err)
		return
	}
	if len(baCertList) == 0 {
		t.Error("Expected retrieved certificates to be non-empty, but got empty slice")
	}
	praPortalService := praportal.New(client)
	// Create multiple PRA Portals and collect their IDs
	var praPortalIDs []string
	for i, cert := range baCertList[:2] { // Assuming you need two PRA Portals and there are at least two certificates
		praPortal, _, err := praPortalService.Create(&praportal.PRAPortal{
			Name:                    name + fmt.Sprintf("_pra_portal_%d", i),
			Description:             name + fmt.Sprintf(" Description %d", i),
			Enabled:                 true,
			Domain:                  name + fmt.Sprintf("_domain_%d.example.com", i),
			UserNotification:        "This is an automated integration test",
			UserNotificationEnabled: true,
			CertificateID:           cert.ID,
		})
		if err != nil {
			t.Errorf("Error creating PRA portal %d for testing PRA console: %v", i, err)
			return
		}
		defer func(portalID string) {
			_, err := praPortalService.Delete(portalID)
			if err != nil {
				t.Logf("Error deleting PRA portal with ID %s: %v", portalID, err)
			}
		}(praPortal.ID)
		praPortalIDs = append(praPortalIDs, praPortal.ID)
	}
	service := New(client)

	var praConsoles []PRAConsole

	// Assuming retrievedpraAppSeg.SRAAppsDto correctly holds a slice of SRAAppsDto objects
	if len(retrievedpraAppSeg.PRAApps) >= 2 {
		// Example to address the indexing issue based on your setup
		for i := 0; i < len(retrievedpraAppSeg.PRAApps); i += 2 {
			// Ensure there's at least one more element for the second console
			if i+1 < len(retrievedpraAppSeg.PRAApps) {
				praConsole1 := PRAConsole{
					Name:        name + "_rdp_pra.example.com",
					Description: name + "_rdp_pra.example.com",
					Enabled:     true,
					PRAApplication: PRAApplication{
						ID: retrievedpraAppSeg.PRAApps[i].ID, // Properly using indexing on the slice
					},
					PRAPortals: []PRAPortals{
						{ID: praPortalIDs[0]},
						{ID: praPortalIDs[1]},
					},
				}

				praConsole2 := PRAConsole{
					Name:        name + "_ssh_pra.example.com",
					Description: name + "_ssh_pra.example.com",
					Enabled:     true,
					PRAApplication: PRAApplication{
						ID: retrievedpraAppSeg.PRAApps[i+1].ID, // Properly using indexing on the slice
					},
					PRAPortals: []PRAPortals{
						{ID: praPortalIDs[0]},
						{ID: praPortalIDs[1]},
					},
				}

				praConsoles = append(praConsoles, praConsole1, praConsole2)
			}
		}
	}

	// Test resource creation
	createdResources, _, err := service.CreatePraBulk(praConsoles)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making POST request: %v", err)
	}
	// Collect IDs of created PRAConsole resources for later deletion
	var createdConsoleIDs []string
	for _, createdResource := range createdResources {
		createdConsoleIDs = append(createdConsoleIDs, createdResource.ID)
	}

	// Retrieve and Update all PRA Consoles
	allPRAConsoles, _, err := service.GetAll()
	if err != nil {
		t.Errorf("Error retrieving PRA Consoles: %v", err)
		return
	}

	for _, console := range allPRAConsoles {
		// Prepare the update - e.g., updating the description for simplicity
		console.Description = updateName
		_, err := service.Update(console.ID, &console)
		if err != nil {
			t.Errorf("Error updating PRA console with ID %s: %v", console.ID, err)
		}
	}

	// Delete PRA Console resources after updates
	for _, consoleID := range createdConsoleIDs {
		_, err := service.Delete(consoleID)
		if err != nil {
			t.Errorf("Error deleting PRA console with ID %s: %v", consoleID, err)
		}
	}

	// Defer the deletion of the praAppSeg resource with a delay
	defer func() {
		time.Sleep(2 * time.Second) // Delay to ensure all deletions have propagated
		_, err := praSegmentService.Delete(createdpraAppSeg.ID)
		if err != nil {
			t.Errorf("Error deleting pra application segment: %v", err)
		}
	}()
}

func TestRetrieveNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, _, err = service.Get("non_existent_id")
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, err = service.Delete("non_existent_id")
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, err = service.Update("non_existent_id", &PRAConsole{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, _, err = service.GetByName("non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
