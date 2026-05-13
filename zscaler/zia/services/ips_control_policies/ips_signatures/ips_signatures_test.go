package ips_signatures

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	maxRetries    = 3
	retryInterval = 2 * time.Second

	maxConflictRetries    = 5
	conflictRetryInterval = 1 * time.Second

	// advancedSecurityCategoryID is the threat category ID returned by the
	// API for ADVANCED_SECURITY. It is stable across tenants and used here
	// to keep the test self-contained without needing a lookup endpoint.
	advancedSecurityCategoryID = 64
)

// retryOnConflict retries an operation while the API returns
// EDIT_LOCK_NOT_AVAILABLE, the ZIA global edit lock error.
func retryOnConflict(operation func() error) error {
	var lastErr error
	for i := 0; i < maxConflictRetries; i++ {
		lastErr = operation()
		if lastErr == nil {
			return nil
		}

		if strings.Contains(lastErr.Error(), `"code":"EDIT_LOCK_NOT_AVAILABLE"`) {
			log.Printf("Conflict error detected, retrying in %v... (Attempt %d/%d)", conflictRetryInterval, i+1, maxConflictRetries)
			time.Sleep(conflictRetryInterval)
			continue
		}

		return lastErr
	}
	return lastErr
}

// tryRetrieveResource retrieves an IPS signature rule with a small retry loop
// to absorb propagation delay between Create and the first readable Get.
func tryRetrieveResource(ctx context.Context, service *zscaler.Service, id int) (*IPSSignatureRules, error) {
	var resource *IPSSignatureRules
	var err error

	for i := 0; i < maxRetries; i++ {
		resource, err = Get(ctx, service, id)
		if err == nil && resource != nil && resource.ID == id {
			return resource, nil
		}
		log.Printf("Attempt %d: Error retrieving IPS signature rule, retrying in %v...", i+1, retryInterval)
		time.Sleep(retryInterval)
	}
	return nil, err
}

// buildRuleText returns a valid Suricata rule text accepted by the Zscaler
// validator (content >= 5 chars, http_uri modifier, single line).
func buildRuleText(sid int) string {
	return fmt.Sprintf(
		`alert http any any -> any any (msg:"Test HTTP rule sid %d"; content:"/admin"; http_uri; nocase; sid:%d; rev:1;)`,
		sid, sid,
	)
}

func TestIPSSignatureRules(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(30, acctest.CharSetAlpha)
	updateDescription := acctest.RandStringFromCharSet(30, acctest.CharSetAlpha)

	sid := 1000000 + int(time.Now().UnixNano()%900000)
	ruleText := buildRuleText(sid)

	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	ctx := context.Background()

	signature := IPSSignatureRules{
		Name:        name,
		Description: description,
		RuleText:    ruleText,
		Enabled:     false,
		Category: &IPSSignatureCategory{
			ID: advancedSecurityCategoryID,
		},
	}

	var createdResource *IPSSignatureRules
	err = retryOnConflict(func() error {
		createdResource, _, err = Create(ctx, service, &signature)
		return err
	})
	if err != nil {
		t.Fatalf("Error creating IPS signature rule: %v", err)
	}

	// Ensure cleanup runs even if a sub-test fails later
	t.Cleanup(func() {
		cleanupErr := retryOnConflict(func() error {
			_, delErr := Delete(context.Background(), service, createdResource.ID)
			return delErr
		})
		if cleanupErr != nil {
			t.Logf("[CLEANUP] failed to delete IPS signature rule %d: %v", createdResource.ID, cleanupErr)
		}
	})

	if createdResource.ID == 0 {
		t.Fatal("Expected created resource ID to be non-zero, but got 0")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.Name)
	}
	if createdResource.RuleText != ruleText {
		t.Errorf("Expected created resource ruleText to match input, but it differed")
	}
	if createdResource.Category == nil || createdResource.Category.ID != advancedSecurityCategoryID {
		t.Errorf("Expected category ID %d, but got %+v", advancedSecurityCategoryID, createdResource.Category)
	}

	t.Run("Get", func(t *testing.T) {
		retrieved, err := tryRetrieveResource(ctx, service, createdResource.ID)
		if err != nil {
			t.Fatalf("Error retrieving resource: %v", err)
		}
		if retrieved.ID != createdResource.ID {
			t.Errorf("Expected retrieved ID '%d', but got '%d'", createdResource.ID, retrieved.ID)
		}
		if retrieved.Name != name {
			t.Errorf("Expected retrieved name '%s', but got '%s'", name, retrieved.Name)
		}
		if retrieved.Category == nil || retrieved.Category.Name == "" {
			t.Errorf("Expected category to be populated on Get, got: %+v", retrieved.Category)
		}
	})

	t.Run("Update", func(t *testing.T) {
		retrieved, err := tryRetrieveResource(ctx, service, createdResource.ID)
		if err != nil {
			t.Fatalf("Error retrieving resource for update: %v", err)
		}
		retrieved.Description = updateDescription
		err = retryOnConflict(func() error {
			_, _, err = Update(ctx, service, createdResource.ID, retrieved)
			return err
		})
		if err != nil {
			t.Fatalf("Error updating resource: %v", err)
		}

		updated, err := Get(ctx, service, createdResource.ID)
		if err != nil {
			t.Fatalf("Error retrieving updated resource: %v", err)
		}
		if updated.Description != updateDescription {
			t.Errorf("Expected updated description '%s', but got '%s'", updateDescription, updated.Description)
		}
	})

	t.Run("GetByName", func(t *testing.T) {
		retrieved, err := GetByName(ctx, service, name)
		if err != nil {
			t.Fatalf("Error retrieving resource by name: %v", err)
		}
		if retrieved.ID != createdResource.ID {
			t.Errorf("Expected ID '%d', but got '%d'", createdResource.ID, retrieved.ID)
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		resources, err := GetAll(ctx, service)
		if err != nil {
			t.Fatalf("Error retrieving resources: %v", err)
		}
		if len(resources) == 0 {
			t.Fatal("Expected retrieved resources to be non-empty, but got empty slice")
		}
		found := false
		for _, r := range resources {
			if r.ID == createdResource.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected retrieved resources to contain created resource '%d', but it didn't", createdResource.ID)
		}
	})

	t.Run("ValidateRuleText_Valid", func(t *testing.T) {
		validSID := sid + 1
		validation, err := ValidateIPSSignatureRuleText(ctx, service, buildRuleText(validSID))
		if err != nil {
			t.Fatalf("Error validating well-formed rule text: %v", err)
		}
		if validation == nil {
			t.Fatal("Expected validation result, got nil")
		}
		// The API uses an empty/OK-style status for successful validation. A
		// non-empty ErrMsg indicates a real failure that should fail the test.
		if validation.ErrMsg != "" {
			t.Errorf("Expected no validation error for well-formed rule, got: %s (param: %q, suggestion: %q)",
				validation.ErrMsg, validation.ErrParameter, validation.ErrSuggestion)
		}
	})

	t.Run("ValidateRuleText_Invalid", func(t *testing.T) {
		// content shorter than 5 characters is rejected by the Zscaler validator.
		// The API returns HTTP 400 INVALID_INPUT_ARGUMENT (with the diagnostic
		// in the standard error envelope), which the SDK surfaces as a Go error.
		invalid := fmt.Sprintf(
			`alert http any any -> any any (msg:"Test HTTP rule"; content:"abc"; sid:%d; rev:1;)`,
			sid+2,
		)
		validation, err := ValidateIPSSignatureRuleText(ctx, service, invalid)
		if err == nil {
			t.Errorf("Expected validation failure for invalid rule, but got clean result: %+v", validation)
			return
		}
		if !strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") &&
			!strings.Contains(err.Error(), "Minimum length") {
			t.Errorf("Expected an INVALID_INPUT_ARGUMENT / pattern-length error, got: %v", err)
		}
		t.Logf("Got expected validation error: %v", err)
	})

}

func TestRetrieveNonExistentIPSSignatureRule(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	_, err = Get(context.Background(), service, 0)
	if err == nil {
		t.Error("Expected error retrieving non-existent IPS signature rule, but got nil")
	}
}

func TestDeleteNonExistentIPSSignatureRule(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	_, err = Delete(context.Background(), service, 0)
	if err == nil {
		t.Error("Expected error deleting non-existent IPS signature rule, but got nil")
	}
}

func TestUpdateNonExistentIPSSignatureRule(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	_, _, err = Update(context.Background(), service, 0, &IPSSignatureRules{})
	if err == nil {
		t.Error("Expected error updating non-existent IPS signature rule, but got nil")
	}
}

func TestGetByNameNonExistentIPSSignatureRule(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	_, err = GetByName(context.Background(), service, "non_existent_signature_name_xyz")
	if err == nil {
		t.Error("Expected error retrieving by non-existent name, but got nil")
	}
}

func TestValidateIPSSignatureRuleText_EmptyInput(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	_, err = ValidateIPSSignatureRuleText(context.Background(), service, "")
	if err == nil {
		t.Error("Expected error validating empty rule text, but got nil")
	}
}
