package locationgroups

import (
	"fmt"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestLocationGroups(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Test resources retrieval
	resources, err := GetAll(service)
	if err != nil {
		t.Errorf("Error retrieving resources: %v", err)
	}
	if len(resources) == 0 {
		return
	}
	name := resources[0].Name
	resourceByName, err := GetLocationGroupByName(service, name)
	if err != nil {
		t.Errorf("Error retrieving resource by name: %v", err)
	}

	_, err = GetLocationGroup(service, resourceByName.ID)
	if err != nil {
		t.Errorf("expected resource to exist: %v", err)
	}

	// Test GetGroupType with both DYNAMIC_GROUP and STATIC_GROUP
	groupTypes := []string{"DYNAMIC_GROUP", "STATIC_GROUP"}
	for _, gType := range groupTypes {
		t.Run(fmt.Sprintf("GroupType-%s", gType), func(t *testing.T) {
			resourceByType, err := GetGroupType(service, gType)
			if err != nil {
				t.Errorf("Error retrieving resource by type '%s': %v", gType, err)
			} else {
				_, err = GetLocationGroup(service, resourceByType.ID)
				if err != nil {
					t.Errorf("expected resource to exist for group type '%s': %v", gType, err)
				}
			}
		})
	}
}
