package inspection_predefined_controls

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestInspectionPredefinedControls(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	// Corrected this line to include the version
	controls, err := service.GetAll("OWASP_CRS/3.3.0")
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
	control, _, err := service.GetByName(name, "OWASP_CRS/3.3.0")
	if err != nil {
		t.Errorf("Error getting predefined control by name: %v", err)
		return
	}
	if control.Name != name {
		t.Errorf("predefined control name does not match: expected %s, got %s", name, control.Name)
		return
	}
}

func TestGetAllByGroup(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := New(client)

	version := "OWASP_CRS/3.3.0"

	// Call GetAll to retrieve all control groups
	allControls, err := service.GetAll(version)
	if err != nil {
		t.Fatalf("Error getting all controls: %v", err)
	}

	if len(allControls) == 0 {
		t.Fatalf("No controls found")
	}

	// Use a map to store unique controlGroup names
	controlGroups := make(map[string]bool)
	for _, control := range allControls {
		controlGroups[control.ControlGroup] = true
	}

	// Now call GetAllByGroup for each unique controlGroup
	for group := range controlGroups {
		t.Logf("Fetching details for control group: %s", group)
		controls, err := service.GetAllByGroup(version, group)
		if err != nil {
			t.Fatalf("Error getting details for control group %s: %v", group, err)
		}

		if len(controls) == 0 {
			t.Errorf("No details found for control group: %s", group)
		}
	}
}

func TestGetControlGroup(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := New(client)

	version := "OWASP_CRS/3.3.0"
	groupName := "Protocol Issues"

	controls, err := service.GetAllByGroup(version, groupName)
	if err != nil {
		t.Fatalf("Error getting details for control group %s: %v", groupName, err)
	}

	if len(controls) == 0 {
		t.Fatalf("No details found for control group: %s", groupName)
	}
}
