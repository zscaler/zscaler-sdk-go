package praapproval

import (
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentpra"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"
)

// getCurrentEpochTime returns the current time in epoch format.
func getCurrentEpochTime() int64 {
	return time.Now().Unix()
}

func getStartTime() string {
	// Adjusting the start time to be 5 minutes into the future.
	startTime := getCurrentEpochTime() + 5*60 // Adding 300 seconds (5 minutes)
	return strconv.FormatInt(startTime, 10)
}

func getEndTime() string {
	endTime := getCurrentEpochTime() + (364 * 24 * 60 * 60) // Adding 31536000 seconds (365 days)
	return strconv.FormatInt(endTime, 10)
}

// A sample list of IANA Time Zones.
// Extend this list based on your requirements.
var timeZones = []string{
	"America/New_York",
	"America/Chicago",
	"America/Denver",
	"America/Los_Angeles",
	"America/Vancouver",
	"Europe/London",
	"Europe/Berlin",
	"Asia/Tokyo",
	"Asia/Shanghai",
	"Asia/Kolkata",
	"Australia/Sydney",
}

// randTimeZone selects a random time zone from the timeZones slice.
func randTimeZone() string {
	rand.Seed(time.Now().UnixNano()) // Ensure different output for each program run
	return timeZones[rand.Intn(len(timeZones))]
}

// getRandomTimeZone ensures the randomly selected time zone is valid by trying to load it.
func getRandomTimeZone() (string, error) {
	tz := randTimeZone()
	if _, err := time.LoadLocation(tz); err != nil {
		return "", err
	}
	return tz, nil
}
func TestCredentialController(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	//updateName := "tests-" + acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	appGroup := segmentgroup.SegmentGroup{
		Name:        name,
		Description: name,
		Enabled:     true,
	}
	createdSegGroup, _, err := segmentgroup.Create(service, &appGroup)
	if err != nil {
		t.Errorf("Error creating segment group: %v", err)
		return
	}
	defer func() {
		time.Sleep(time.Second * 2) // Sleep for 2 seconds before deletion
		_, _, getErr := segmentgroup.Get(service, createdSegGroup.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := segmentgroup.Delete(service, createdSegGroup.ID)
			if err != nil {
				t.Errorf("Error deleting segment group: %v", err)
			}
		}
	}()

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
				From: "3390",
				To:   "3390",
			},
			{
				From: "2222",
				To:   "2222",
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
					Domain:              "rdp_pra.example.com",
				},
				{
					Name:                name,
					Description:         name,
					Enabled:             true,
					AppTypes:            []string{"SECURE_REMOTE_ACCESS"},
					ApplicationPort:     "2222",
					ApplicationProtocol: "SSH",
					Domain:              "ssh_pra.example.com",
				},
			},
		},
	}
	createdpraAppSeg, _, err := applicationsegmentpra.Create(service, praAppSeg)
	if err != nil {
		t.Errorf("Error creating pra application segment: %v", err)
		return
	}

	// Assuming the praSegmentService.Get correctly returns the payload as described
	retrievedpraAppSeg, _, err := applicationsegmentpra.Get(service, createdpraAppSeg.ID)
	if err != nil {
		t.Errorf("Error retrieving created pra application segment: %v", err)
		return
	}

	defer func() {
		time.Sleep(time.Second * 2) // Sleep for 2 seconds before deletion
		_, _, getErr := applicationsegmentpra.Get(service, createdpraAppSeg.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := applicationsegmentpra.Delete(service, createdpraAppSeg.ID)
			if err != nil {
				t.Errorf("Error deleting pra application segment: %v", err)
			}
		}
	}()
	// Attempt to get a random but valid time zone
	tz, err := getRandomTimeZone()
	if err != nil {
		t.Fatalf("Failed to load random time zone: %v", err)
	}

	credController := PrivilegedApproval{
		EmailIDs:  []string{"carol.kirk@securitygeek.io"},
		StartTime: getStartTime(), // Dynamically generate valid start time
		EndTime:   getEndTime(),   // Dynamically generate valid end time
		Status:    "ACTIVE",
		WorkingHours: &WorkingHours{
			Days:          []string{"MON", "TUE", "WED", "THU", "FRI"},
			TimeZone:      tz,
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
	createdResource, _, err := Create(service, &credController)
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

	if createdResource.ID == "" {
		t.Fatal("Expected created resource ID to be non-empty")
	}
	// Test resource retrieval
	retrievedResource, _, err := Get(service, createdResource.ID)
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
	_, err = Update(service, createdResource.ID, &credController)
	if err != nil {
		t.Errorf("Error updating resource: %v", err)
		return
	}

	// Retrieve the resource again to verify the update was successful
	updatedResource, _, err := Get(service, createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving updated resource: %v", err)
		return
	}

	// Verify the 'Days' attribute was successfully updated
	if !reflect.DeepEqual(updatedResource.WorkingHours.Days, []string{"MON", "WED", "FRI"}) {
		t.Errorf("Expected updated working days to be 'MON', 'WED', 'FRI', but got '%v'", updatedResource.WorkingHours.Days)
	}

	// Test resources retrieval
	resources, _, err := GetAll(service)
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
	_, err = Delete(service, createdResource.ID)
	if err != nil {
		t.Errorf("Error deleting resource: %v", err)
		return
	}

	// Test resource retrieval after deletion
	_, _, err = Get(service, createdResource.ID)
	if err == nil {
		t.Errorf("Expected error retrieving deleted resource, but got nil")
	}
}

func TestRetrieveNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, _, err = Get(service, "non_existent_id")
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = Update(service, "non_existent_id", &PrivilegedApproval{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = Delete(service, "non_existent_id")
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestDeleteExpiredResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	resp, err := DeleteExpired(service)
	if err != nil {
		t.Errorf("Unexpected error when calling DeleteExpired: %v", err)
	} else if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP status code 200 OK, got: %d", resp.StatusCode)
	} else {
		t.Log("DeleteExpired completed successfully, potentially with no resources to delete.")
	}
}
