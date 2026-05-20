package zia_posture

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestZIAPosture(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updatedName := "updated-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)

	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Mirrors the attached zia_posture.json payload (CrowdStrike criteria
	// referenced by id / udid), with a unique name per run.
	posture := ZIAPosture{
		Name:     name,
		Platform: 3,
		HighTrustCriteria: HighTrustCriteria{
			Cs: []TrustCriteriaSet{
				{
					Cn: []TrustCriterion{
						{ID: "9911", Name: "CrowdStrike_ZPA_ZTA_40", UDID: "6e36dd2f-ce19-47b3-8f26-f1e4e8f6313e"},
						{ID: "9913", Name: "CrowdStrike_ZPA_ZTA_80", UDID: "fc73ffb2-3ad7-49d5-9bff-10480589d188"},
						{ID: "9915", Name: "CrowdStrike_ZPA_Pre-ZTA", UDID: "cfab2ee9-9bf4-4482-9dcc-dadf7311c49b"},
					},
				},
			},
		},
		MediumTrustCriteria: MediumTrustCriteria{
			Cs: []TrustCriteriaSet{
				{
					Cn: []TrustCriterion{
						{ID: "9911", Name: "CrowdStrike_ZPA_ZTA_40", UDID: "6e36dd2f-ce19-47b3-8f26-f1e4e8f6313e"},
						{ID: "9913", Name: "CrowdStrike_ZPA_ZTA_80", UDID: "fc73ffb2-3ad7-49d5-9bff-10480589d188"},
					},
				},
			},
		},
		LowTrustCriteria: LowTrustCriteria{
			Cs: []TrustCriteriaSet{
				{
					Cn: []TrustCriterion{
						{ID: "9911", Name: "CrowdStrike_ZPA_ZTA_40", UDID: "6e36dd2f-ce19-47b3-8f26-f1e4e8f6313e"},
						{ID: "9913", Name: "CrowdStrike_ZPA_ZTA_80", UDID: "fc73ffb2-3ad7-49d5-9bff-10480589d188"},
					},
				},
			},
		},
	}

	createdResource, _, err := Create(context.Background(), service, &posture)
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

	if createdResource.ID == 0 {
		t.Fatal("Expected created resource ID to be non-zero, but got 0")
	}
	if !strings.EqualFold(createdResource.Name, name) {
		t.Errorf("Expected created zia posture name '%s', but got '%s'", name, createdResource.Name)
	}
	if createdResource.Platform != 3 {
		t.Errorf("Expected created platform '3', but got '%d'", createdResource.Platform)
	}
	if len(createdResource.HighTrustCriteria.Cs) == 0 || len(createdResource.HighTrustCriteria.Cs[0].Cn) != 3 {
		t.Errorf("Expected created highTrustCriteria to contain 3 criteria, but got: %+v", createdResource.HighTrustCriteria)
	}
	if len(createdResource.MediumTrustCriteria.Cs) == 0 || len(createdResource.MediumTrustCriteria.Cs[0].Cn) != 2 {
		t.Errorf("Expected created mediumTrustCriteria to contain 2 criteria, but got: %+v", createdResource.MediumTrustCriteria)
	}
	if len(createdResource.LowTrustCriteria.Cs) == 0 || len(createdResource.LowTrustCriteria.Cs[0].Cn) != 2 {
		t.Errorf("Expected created lowTrustCriteria to contain 2 criteria, but got: %+v", createdResource.LowTrustCriteria)
	}

	// Test resource retrieval
	retrievedResource, err := Get(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if !strings.EqualFold(retrievedResource.Name, name) {
		t.Errorf("Expected retrieved zia posture name '%s', but got '%s'", name, retrievedResource.Name)
	}

	// Test resource update (PUT) — rename and narrow highTrustCriteria down to 2 entries.
	retrievedResource.Name = updatedName
	if len(retrievedResource.HighTrustCriteria.Cs) > 0 {
		retrievedResource.HighTrustCriteria.Cs[0].Cn = retrievedResource.HighTrustCriteria.Cs[0].Cn[:2]
	}

	if _, _, err := Update(context.Background(), service, createdResource.ID, retrievedResource); err != nil {
		t.Fatalf("Error updating resource: %v", err)
	}

	updatedResource, err := Get(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving updated resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%d', but got '%d'", createdResource.ID, updatedResource.ID)
	}
	if !strings.EqualFold(updatedResource.Name, updatedName) {
		t.Errorf("Expected updated zia posture name '%s', but got '%s'", updatedName, updatedResource.Name)
	}
	if len(updatedResource.HighTrustCriteria.Cs) == 0 || len(updatedResource.HighTrustCriteria.Cs[0].Cn) != 2 {
		t.Errorf("Expected updated highTrustCriteria to contain 2 criteria, but got: %+v", updatedResource.HighTrustCriteria)
	}

	// Test resource partial update (PATCH).
	// NOTE: PATCH on this endpoint is not a JSON merge — the API treats the
	// body as a full replace. Read current state, mutate only the field we
	// care about, then send the merged object. Assertions are made against a
	// follow-up GET because the PATCH response body may omit "id".
	patchSource, err := Get(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource before partial update: %v", err)
	}
	if len(patchSource.LowTrustCriteria.Cs) > 0 {
		patchSource.LowTrustCriteria.Cs[0].Cn = patchSource.LowTrustCriteria.Cs[0].Cn[:1]
	}

	if _, _, err := PartialUpdate(context.Background(), service, createdResource.ID, patchSource); err != nil {
		t.Fatalf("Error partially updating resource: %v", err)
	}

	patchedFromGet, err := Get(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving partially updated resource: %v", err)
	}
	if len(patchedFromGet.LowTrustCriteria.Cs) == 0 || len(patchedFromGet.LowTrustCriteria.Cs[0].Cn) != 1 {
		t.Errorf("Expected partially updated lowTrustCriteria to contain 1 criterion, but got: %+v", patchedFromGet.LowTrustCriteria)
	}
	if !strings.EqualFold(patchedFromGet.Name, updatedName) {
		t.Errorf("Expected PATCH to preserve name '%s', but got '%s'", updatedName, patchedFromGet.Name)
	}
	if len(patchedFromGet.HighTrustCriteria.Cs) == 0 || len(patchedFromGet.HighTrustCriteria.Cs[0].Cn) != 2 {
		t.Errorf("Expected PATCH to preserve highTrustCriteria with 2 criteria, but got: %+v", patchedFromGet.HighTrustCriteria)
	}

	// Test resource retrieval by name (using the updated name)
	retrievedByName, err := GetByName(context.Background(), service, updatedName)
	if err != nil {
		t.Fatalf("Error retrieving resource by name: %v", err)
	}
	if retrievedByName.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedByName.ID)
	}

	// Test resources retrieval (v2 paginated list)
	resources, err := GetAll(context.Background(), service, nil)
	if err != nil {
		t.Fatalf("Error retrieving resources: %v", err)
	}
	if len(resources) == 0 {
		t.Fatal("Expected retrieved resources to be non-empty, but got empty slice")
	}

	found := false
	for _, resource := range resources {
		if resource.ID == createdResource.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected retrieved resources to contain created resource '%d', but it didn't", createdResource.ID)
	}

	// Test resource removal
	if _, err := Delete(context.Background(), service, createdResource.ID); err != nil {
		t.Fatalf("Error deleting resource: %v", err)
	}

	_, err = Get(context.Background(), service, createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = GetByName(context.Background(), service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
