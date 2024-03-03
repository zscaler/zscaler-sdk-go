package praapproval

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/applicationsegmentpra"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/segmentgroup"
)

func TestCredentialController(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	//updateName := "tests-" + acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	service := New(client)
	// Adjusting StartTime and EndTime
	// startTime, err := time.Parse("15:04", "10:00") // Assuming you want 10:00 as start time
	// if err != nil {
	// 	t.Fatalf("Failed to parse start time: %v", err)
	// }
	// endTime, err := time.Parse("15:04", "17:00") // Assuming you want 17:00 as end time
	// if err != nil {
	// 	t.Fatalf("Failed to parse end time: %v", err)
	// }

	// Adjusting TimeZone
	// loc, err := time.LoadLocation("Asia/Calcutta")
	// if err != nil {
	// 	t.Fatalf("Failed to load location: %v", err)
	// }

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

	defer func() {
		time.Sleep(time.Second * 2) // Sleep for 2 seconds before deletion
		_, _, getErr := praSegmentService.Get(createdpraAppSeg.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := praSegmentService.Delete(createdpraAppSeg.ID)
			if err != nil {
				t.Errorf("Error deleting pra application segment: %v", err)
			}
		}
	}()
	credController := PrivilegedApproval{
		EmailIDs:  []string{"wxiiqedzjo@bd-hashicorp.com"},
		StartTime: "1709596800",
		EndTime:   "1741132800",
		Status:    "ACTIVE",
		WorkingHours: &WorkingHours{
			Days:          []string{"MON", "TUE", "WED", "THU", "FRI"},
			TimeZone:      "America/Vancouver",
			StartTime:     "09:00",
			EndTime:       "17:00",
			StartTimeCron: "0 0 17 ? * MON,TUE,WED,THU,FRI",
			EndTimeCron:   "0 0 1 ? * TUE,WED,THU,FRI,SAT",
		},
		Applications: []Applications{
			{ID: retrievedpraAppSeg.ID},
		},
	}

	// Test resource creation
	createdResource, _, err := service.Create(&credController)
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

	if createdResource.ID == "" {
		t.Fatal("Expected created resource ID to be non-empty")
	}
	// Test resource retrieval
	retrievedResource, _, err := service.Get(createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}

	// Update the 'Days' attribute within 'WorkingHours' to only "MON", "WED", "FRI"
	credController.WorkingHours.Days = []string{"MON", "WED", "FRI"}

	// Update the StartTimeCron and EndTimeCron to reflect the new working days
	credController.WorkingHours.StartTimeCron = "0 0 17 ? * MON,WED,FRI"
	credController.WorkingHours.EndTimeCron = "0 0 1 ? * TUE,THU,SAT"

	// Call the Update function with the modified 'credController' struct
	_, err = service.Update(createdResource.ID, &credController)
	if err != nil {
		t.Errorf("Error updating resource: %v", err)
		return
	}

	// Retrieve the resource again to verify the update was successful
	updatedResource, _, err := service.Get(createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving updated resource: %v", err)
		return
	}

	// Verify the 'Days' attribute was successfully updated
	if !reflect.DeepEqual(updatedResource.WorkingHours.Days, []string{"MON", "WED", "FRI"}) {
		t.Errorf("Expected updated working days to be 'MON', 'WED', 'FRI', but got '%v'", updatedResource.WorkingHours.Days)
	}

	// Test resources retrieval
	resources, _, err := service.GetAll()
	if err != nil {
		t.Errorf("Error retrieving resources: %v", err)
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
		t.Errorf("Expected retrieved resources to contain created resource '%s', but it didn't", createdResource.ID)
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

func TestUpdateNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, err = service.Update("non_existent_id", &PrivilegedApproval{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
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

func TestDeleteExpiredResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	resp, err := service.DeleteExpired()
	if err != nil {
		t.Errorf("Unexpected error when calling DeleteExpired: %v", err)
	} else if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP status code 200 OK, got: %d", resp.StatusCode)
	} else {
		t.Log("DeleteExpired completed successfully, potentially with no resources to delete.")
	}
}
