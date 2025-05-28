package inspection_predefined_controls

import (
	"context"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestInspectionPredefinedControls(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	// Corrected this line to include the version
	controls, err := GetAll(context.Background(), service, "OWASP_CRS/3.3.0")
	if err != nil {
		t.Errorf("Error getting predefined controls: %v", err)
		return
	}
	if len(controls) == 0 {
		t.Errorf("No predefined controls found")
		return
	}
	name := controls[0].Name
	t.Log("Getting predefined control by name:" + name)
	// Corrected this line to include the version
	control, _, err := GetByName(context.Background(), service, name, "OWASP_CRS/3.3.0")
	if err != nil {
		t.Errorf("Error getting predefined control by name: %v", err)
		return
	}
	if control.Name != name {
		t.Errorf("predefined control name does not match: expected %s, got %s", name, control.Name)
		return
	}

	// Negative Test: Try to retrieve a control with a non-existent name
	nonExistentName := "ThisControlNameDoesNotExist"
	_, _, err = GetByName(context.Background(), service, nonExistentName, "OWASP_CRS/3.3.0")
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}
}

func TestGetAllByGroup(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	version := "OWASP_CRS/3.3.0"

	// Call GetAll to retrieve all control groups
	allControls, err := GetAll(context.Background(), service, version)
	if err != nil {
		t.Fatalf("Error getting all controls: %v", err)
	}

	if len(allControls) == 0 {
		t.Fatalf("No controls found")
	}

	// Use the first control group for the test
	firstControlGroup := allControls[0].ControlGroup
	t.Logf("Fetching details for control group: %s", firstControlGroup)
	controls, err := GetAllByGroup(context.Background(), service, version, firstControlGroup)
	if err != nil {
		t.Fatalf("Error getting details for control group %s: %v", firstControlGroup, err)
	}

	if len(controls) == 0 {
		t.Errorf("No details found for control group: %s", firstControlGroup)
	}

	// Negative Test: Try to retrieve controls for a non-existent control group
	nonExistentGroup := "ThisGroupDoesNotExist"
	controls, err = GetAllByGroup(context.Background(), service, version, nonExistentGroup)
	if err != nil {
		t.Errorf("Error getting details for non-existent control group %s: %v", nonExistentGroup, err)
		return
	}
	if len(controls) != 0 {
		t.Errorf("Expected no controls for non-existent control group, but got %d", len(controls))
	}
}

func TestGetControlGroup(t *testing.T) {
	t.Run("TestValidControlGroup", func(t *testing.T) {
		service, err := tests.NewOneAPIClient()
		if err != nil {
			t.Fatalf("Error creating client: %v", err)
		}

		// service, err := tests.NewZPAClient()
		// if err != nil {
		// 	t.Fatalf("Error creating client: %v", err)
		// }
		version := "OWASP_CRS/3.3.0"
		groupName := "Protocol Issues"

		controls, err := GetAllByGroup(context.Background(), service, version, groupName)
		if err != nil {
			t.Fatalf("Error getting details for control group %s: %v", groupName, err)
		}

		if len(controls) == 0 {
			t.Fatalf("No details found for control group: %s", groupName)
		}
	})

	t.Run("TestNonExistentControlGroup", func(t *testing.T) {
		service, err := tests.NewOneAPIClient()
		if err != nil {
			t.Fatalf("Error creating client: %v", err)
		}

		// service, err := tests.NewZPAClient()
		// if err != nil {
		// 	t.Fatalf("Error creating client: %v", err)
		// }
		version := "OWASP_CRS/3.3.0"
		nonExistentGroupName := "ThisControlGroupNameDoesNotExist"

		controls, err := GetAllByGroup(context.Background(), service, version, nonExistentGroupName)
		if err != nil {
			t.Errorf("Error getting details for non-existent control group %s: %v", nonExistentGroupName, err)
			return
		}
		if len(controls) != 0 {
			t.Errorf("Expected no controls for non-existent control group, but got %d", len(controls))
		}
	})
}
