package inspection_predefined_controls

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/tests"
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
