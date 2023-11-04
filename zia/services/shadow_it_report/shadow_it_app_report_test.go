package shadow_it_report

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestAccCustomTags(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	// Fetch all custom tags.
	tags, err := service.GetAllCustomTags()
	if err != nil {
		t.Errorf("Error getting all custom tags: %v", err)
		return
	}
	if len(tags) == 0 {
		t.Errorf("No custom tags found")
		return
	}

	// Adjusting to query only the first 10 names.
	maxQuery := 10
	if len(tags) < maxQuery {
		maxQuery = len(tags)
	}
	for i := 0; i < maxQuery; i++ {
		name := tags[i].Name
		t.Logf("Getting tag by name: %s", name)
		tag, err := service.GetCustomTagsByName(name)
		if err != nil {
			t.Errorf("Error getting custom tag by name %s: %v", name, err)
			continue // Use continue to proceed with the next iteration even if this fails.
		}
		if tag.Name != name {
			t.Errorf("Tag name does not match: expected %s, got %s", name, tag.Name)
		}
	}
}
