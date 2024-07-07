package inspection_profile

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/inspectioncontrol/inspection_custom_controls"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/inspectioncontrol/inspection_predefined_controls"
)

func TestInspectionProfile(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	var resourcesToDelete []string

	service := services.New(client)

	// create Inspection Custom Control for testing
	customControl := inspection_custom_controls.InspectionCustomControl{
		Name:          name,
		Description:   name,
		Action:        "PASS",
		DefaultAction: "PASS",
		ParanoiaLevel: "1",
		Severity:      "CRITICAL",
		Type:          "RESPONSE",
		Rules: []inspection_custom_controls.Rules{
			{
				Names: []string{name, name, name},
				Type:  "RESPONSE_HEADERS",
				Conditions: []inspection_custom_controls.Conditions{
					{
						LHS: "SIZE",
						OP:  "GE",
						RHS: "1000",
					},
				},
			},
			{
				Type: "RESPONSE_BODY",
				Conditions: []inspection_custom_controls.Conditions{
					{
						LHS: "SIZE",
						OP:  "GE",
						RHS: "1000",
					},
				},
			},
		},
	}

	createdCustomControl, _, err := inspection_custom_controls.Create(service, customControl)
	if err != nil || createdCustomControl == nil || createdCustomControl.ID == "" {
		t.Fatalf("Error creating inspection custom control or ID is empty")
	}
	fmt.Printf("Created custom inspection control with ID: %s\n", createdCustomControl.ID)

	resourcesToDelete = append(resourcesToDelete, createdCustomControl.ID) // Add to our list

	defer func() {
		existingControl, _, errCheck := inspection_custom_controls.Get(service, createdCustomControl.ID)
		if errCheck == nil && existingControl != nil {
			_, errDelete := inspection_custom_controls.Delete(service, createdCustomControl.ID)
			if errDelete != nil {
				t.Errorf("Error deleting inspection custom control: %v", errDelete)
			}
		}
	}()

	predefinedControlsByGroup, err := inspection_predefined_controls.GetAllByGroup(service, "OWASP_CRS/3.3.0", "Preprocessors")
	if err != nil {
		t.Errorf("Error getting predefined controls by group: %v", err)
		return
	}

	controlByName, _, err := inspection_predefined_controls.GetByName(service, "Multipart request body failed strict validation", "OWASP_CRS/3.3.0")
	if err != nil {
		t.Errorf("Error getting predefined control by name: %v", err)
		return
	}

	predefinedControls := make([]common.CustomCommonControls, len(predefinedControlsByGroup)+1)
	for i, control := range predefinedControlsByGroup {
		predefinedControls[i] = common.CustomCommonControls{
			ID:          control.ID,
			Action:      "BLOCK",
			ActionValue: control.ActionValue,
		}
	}
	predefinedControls[len(predefinedControlsByGroup)] = common.CustomCommonControls{
		ID:     controlByName.ID,
		Action: "BLOCK",
	}

	profile := InspectionProfile{
		Name:                      name,
		Description:               name,
		ParanoiaLevel:             "1",
		PredefinedControlsVersion: "OWASP_CRS/3.3.0",
		IncarnationNumber:         "6",
		ControlInfoResource: []ControlInfoResource{
			{
				ControlType: "CUSTOM",
			},
		},
		CustomControls: []InspectionCustomControl{
			{
				ID:     createdCustomControl.ID,
				Action: "BLOCK",
			},
		},
		PredefinedControls:   predefinedControls,
		GlobalControlActions: []string{"PREDEFINED:PASS", "CUSTOM:NONE", "OVERRIDE_ACTION:COMMON"},
		CommonGlobalOverrideActionsConfig: map[string]interface{}{
			"PREDEF_CNTRL_GLOBAL_ACTION": interface{}("PASS"),
			"IS_OVERRIDE_ACTION_COMMON":  interface{}("TRUE"),
		},
	}

	createdResource, _, err := Create(service, profile)
	if err != nil || createdResource == nil {
		t.Fatalf("Error making POST request: %v or createdResource is nil", err)
	}
	resourcesToDelete = append(resourcesToDelete, createdResource.ID) // Add to our list

	retrievedResourceAfterCreation, _, err := Get(service, createdResource.ID)
	if err != nil || retrievedResourceAfterCreation == nil {
		t.Fatalf("Failed to verify the creation of the resource: %v", err)
	}

	defer func() {
		for _, resourceID := range resourcesToDelete {
			existingResource, _, errCheck := Get(service, resourceID)
			if errCheck == nil && existingResource != nil {
				_, errDelete := Delete(service, resourceID)
				if errDelete != nil {
					t.Errorf("Error deleting inspection profile with ID %s: %v", resourceID, errDelete)
				}
			}

			existingControl, _, errCheckControl := inspection_custom_controls.Get(service, resourceID)
			if errCheckControl == nil && existingControl != nil {
				_, errDeleteControl := inspection_custom_controls.Delete(service, resourceID)
				if errDeleteControl != nil {
					t.Errorf("Error deleting inspection custom control with ID %s: %v", resourceID, errDeleteControl)
				}
			}
		}
	}()

	if createdResource.ID == "" {
		t.Error("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.Name)
	}
	// Test resource retrieval
	retrievedResource, _, err := Get(service, createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
		return
	}
	// Add this check to ensure that retrievedResource is not nil
	if retrievedResource == nil {
		t.Error("Retrieved resource is nil.")
		return
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != name {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", name, retrievedResource.Name)
	}
	// Test resource update
	retrievedResource.Name = updateName
	_, err = Update(service, createdResource.ID, retrievedResource)
	if err != nil {
		t.Errorf("Error updating resource: %v", err)
	}
	updatedResource, _, err := Get(service, createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%s', but got '%s'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}

	// Test resource retrieval by name
	retrievedResource, _, err = GetByName(service, updateName)
	if err != nil {
		t.Errorf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != updateName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, retrievedResource.Name)
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
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	_, _, err = Get(service, "non_existent_id")
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	_, err = Delete(service, "non_existent_id")
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	_, err = Update(service, "non_existent_id", &InspectionProfile{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	_, _, err = GetByName(service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
