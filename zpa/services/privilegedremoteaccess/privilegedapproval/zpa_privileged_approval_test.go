package privilegedapproval

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
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
		// Check each email ID for the prefix "tests-"
		for _, emailID := range r.EmailIDs {
			if strings.HasPrefix(emailID, "tests-") {
				// If any email ID matches, delete the resource and log the action
				log.Printf("Deleting resource with ID: %s, EmailIDs: %v", r.ID, r.EmailIDs)
				_, _ = service.Delete(r.ID)
				break // Break the inner loop after deletion to avoid multiple attempts
			}
		}
	}
}

func TestCredentialController(t *testing.T) {
	// name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	// updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
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

	credController := PrivilegedApproval{
		EmailIDs:  []string{"wguilherme@securitygeek.io"},
		StartTime: "1708497551",
		EndTime:   "1710999551",
		Status:    "ACTIVE",
		WorkingHours: &WorkingHours{
			Days:          []string{"MON", "TUE", "WED", "THU", "FRI"},
			TimeZone:      "Asia/Calcutta",
			StartTime:     "09:00",
			EndTime:       "17:00",
			StartTimeCron: "0 0 17 ? * MON,TUE,WED,THU,FRI",
			EndTimeCron:   "0 0 1 ? * TUE,WED,THU,FRI,SAT",
		},
		Applications: []Applications{
			{ID: "144124980601389314"},
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

	// Retrieving by EmailID (adjustment for requirement 1)
	// retrievedResourceByEmail, _, err := service.GetByEmailID("wguilherme@securitygeek.io")
	// if err != nil {
	// 	t.Fatalf("Error retrieving resource by email: %v", err)
	// }
	// if retrievedResourceByEmail.ID != createdResource.ID {
	// 	t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResourceByEmail.ID)
	// }

	// Updating the "Days" attribute (adjustment for requirement 2)
	// retrievedResourceByEmail.WorkingHours.Days = []string{"MON", "TUE", "WED", "THU"} // Updated days
	// _, err = service.Update(createdResource.ID, retrievedResourceByEmail)
	// if err != nil {
	// 	t.Fatalf("Error updating resource: %v", err)
	// }

	// Verifying the update
	// updatedResource, _, err := service.Get(createdResource.ID)
	// if err != nil {
	// 	t.Fatalf("Error retrieving updated resource: %v", err)
	// }
	// if len(updatedResource.WorkingHours.Days) != 4 {
	// 	t.Errorf("Expected 4 days in updated resource, but got %d", len(updatedResource.WorkingHours.Days))
	// }

	// Continue with the rest of your test...
}
