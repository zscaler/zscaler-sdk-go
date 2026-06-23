package adaptive_access

import (
	"context"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

// TestAdaptiveAccessGetAll validates that the list endpoint returns data and
// that GetByName can locate one of the returned profiles.
func TestAdaptiveAccessGetAll(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	profiles, err := GetAll(context.Background(), service)
	if err != nil {
		t.Fatalf("Error retrieving adaptive access profiles: %v", err)
	}

	if len(profiles) == 0 {
		t.Log("No adaptive access profiles returned; skipping name lookup assertions")
		return
	}

	// Validate the basic shape of the first profile.
	first := profiles[0]
	if first.Name == "" {
		t.Errorf("Expected adaptive access profile to have a name, got empty string")
	}

	// Test GetByName using a profile name returned by GetAll.
	byName, err := GetByName(context.Background(), service, first.Name)
	if err != nil {
		t.Fatalf("Error retrieving adaptive access profile by name '%s': %v", first.Name, err)
	}
	if byName == nil {
		t.Fatalf("Expected adaptive access profile by name '%s', got nil", first.Name)
	}
	if !strings.EqualFold(byName.Name, first.Name) {
		t.Errorf("Expected profile name '%s', got '%s'", first.Name, byName.Name)
	}
}

// TestAdaptiveAccessGetProfileRules exercises the /profiles/rules endpoint
// with no filters, then with the optional iamAapIds filter derived from the
// data returned by the unfiltered call.
func TestAdaptiveAccessGetProfileRules(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Unfiltered call.
	rules, err := GetProfileRules(context.Background(), service, nil)
	if err != nil {
		t.Fatalf("Error retrieving adaptive access profile rules: %v", err)
	}

	if len(rules) == 0 {
		t.Log("No adaptive access profile rules returned; skipping filter assertions")
		return
	}

	// Collect a non-empty iamAapId to use as a filter.
	var filterID string
	for _, r := range rules {
		if r.IamAapID != "" {
			filterID = r.IamAapID
			break
		}
	}

	if filterID == "" {
		t.Log("No iamAapId present in returned rules; skipping iamAapIds filter assertions")
		return
	}

	// Filtered call using the iamAapIds optional parameter.
	filtered, err := GetProfileRules(context.Background(), service, &GetFilterOptions{
		IAMAapIDs: []string{filterID},
	})
	if err != nil {
		t.Fatalf("Error retrieving adaptive access profile rules with iamAapIds filter: %v", err)
	}

	// Every returned record should match the requested filter ID.
	for _, r := range filtered {
		if r.IamAapID != "" && r.IamAapID != filterID {
			t.Errorf("Expected filtered rule iamAapId '%s', got '%s'", filterID, r.IamAapID)
		}
	}
}

// TestAdaptiveAccessGetProfileRulesWithOrgID validates that supplying the orgId
// optional parameter does not error against a live tenant.
func TestAdaptiveAccessGetProfileRulesWithOrgID(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	rules, err := GetProfileRules(context.Background(), service, nil)
	if err != nil {
		t.Fatalf("Error retrieving adaptive access profile rules: %v", err)
	}

	if len(rules) == 0 {
		t.Log("No adaptive access profile rules returned; skipping orgId filter assertions")
		return
	}

	// Re-query passing both optional parameters together to exercise the full
	// query-string construction path.
	var filterID string
	for _, r := range rules {
		if r.IamAapID != "" {
			filterID = r.IamAapID
			break
		}
	}

	orgID := 0
	opts := &GetFilterOptions{OrgID: &orgID}
	if filterID != "" {
		opts.IAMAapIDs = []string{filterID}
	}

	if _, err := GetProfileRules(context.Background(), service, opts); err != nil {
		t.Fatalf("Error retrieving adaptive access profile rules with combined filters: %v", err)
	}
}

func TestAdaptiveAccessGetByNameNonExistent(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = GetByName(context.Background(), service, "non_existent_adaptive_access_profile")
	if err == nil {
		t.Error("Expected error retrieving adaptive access profile by non-existent name, but got nil")
	}
}
