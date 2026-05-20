package notification_template

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestNotificationTemplate(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updatedName := "updated-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)

	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Mirrors the attached notification_template.json payload, with a unique
	// name and isDefaultTemplate left false so the create does not collide
	// with the tenant's existing default template.
	template := NotificationTemplate{
		Name:                name,
		IsDefaultTemplate:   false,
		EnableClient:        false,
		EnableZia:           false,
		EnableAppUpdates:    false,
		EnableServiceStatus: false,
		DurationInSeconds:   5,
		EnablePersistent:    false,
		EnableDoNotDisturb:  false,
		ZIANotificationTemplate: ZIANotificationTemplate{
			EnableZiaFirewall:      false,
			EnableZiaFirewallPopup: false,
			EnableZiaDNS:           false,
			EnableZiaDNSPopup:      false,
			EnableZiaIPS:           false,
			EnableZiaIPSPopup:      false,
			EnableZiaPersistent:    false,
		},
		ZPANotificationTemplate: ZPANotificationTemplate{
			EnableDevicePostureFailure: false,
			EnableZpaReauth:            true,
			ZpaReauthIntervalInMinutes: 5,
			DelayPostureFailureSeconds: 0,
		},
	}

	createdResource, _, err := Create(context.Background(), service, &template)
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

	if createdResource.ID == 0 {
		t.Fatal("Expected created resource ID to be non-zero, but got 0")
	}
	if !strings.EqualFold(createdResource.Name, name) {
		t.Errorf("Expected created notification template name '%s', but got '%s'", name, createdResource.Name)
	}
	if createdResource.DurationInSeconds != 5 {
		t.Errorf("Expected created durationInSeconds '5', but got '%d'", createdResource.DurationInSeconds)
	}
	if !createdResource.ZPANotificationTemplate.EnableZpaReauth {
		t.Error("Expected created enableZpaReauth to be true")
	}
	if createdResource.ZPANotificationTemplate.ZpaReauthIntervalInMinutes != 5 {
		t.Errorf("Expected created zpaReauthIntervalInMinutes '5', but got '%d'", createdResource.ZPANotificationTemplate.ZpaReauthIntervalInMinutes)
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
		t.Errorf("Expected retrieved notification template name '%s', but got '%s'", name, retrievedResource.Name)
	}

	// Test resource update (PUT) — rename, enable ZIA + nested ZIA toggles, bump duration.
	retrievedResource.Name = updatedName
	retrievedResource.EnableZia = true
	retrievedResource.DurationInSeconds = 10
	retrievedResource.ZIANotificationTemplate.EnableZiaFirewall = true
	retrievedResource.ZIANotificationTemplate.EnableZiaFirewallPopup = true
	retrievedResource.ZIANotificationTemplate.EnableZiaIPS = true

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
		t.Errorf("Expected updated notification template name '%s', but got '%s'", updatedName, updatedResource.Name)
	}
	if updatedResource.DurationInSeconds != 10 {
		t.Errorf("Expected updated durationInSeconds '10', but got '%d'", updatedResource.DurationInSeconds)
	}
	if !updatedResource.EnableZia {
		t.Error("Expected updated enableZia to be true")
	}
	if !updatedResource.ZIANotificationTemplate.EnableZiaFirewall {
		t.Error("Expected updated enableZiaFirewall to be true")
	}
	if !updatedResource.ZIANotificationTemplate.EnableZiaIPS {
		t.Error("Expected updated enableZiaIPS to be true")
	}

	// Test resource partial update (PATCH).
	// NOTE: PATCH on this endpoint is not a JSON merge — the API treats the
	// body as a full replace. Read the current state, mutate only the field we
	// care about, and send the merged object. EnableZpaReauth must stay true
	// or the API normalizes ZpaReauthIntervalInMinutes back to its default (5).
	// The PATCH response body also omits "id", so assertions are made against
	// a follow-up GET rather than the PATCH return value.
	patchedInterval := 15
	patchSource, err := Get(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource before partial update: %v", err)
	}
	patchSource.ZPANotificationTemplate.EnableZpaReauth = true
	patchSource.ZPANotificationTemplate.ZpaReauthIntervalInMinutes = patchedInterval

	if _, _, err := PartialUpdate(context.Background(), service, createdResource.ID, patchSource); err != nil {
		t.Fatalf("Error partially updating resource: %v", err)
	}

	patchedFromGet, err := Get(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving partially updated resource: %v", err)
	}
	if patchedFromGet.ZPANotificationTemplate.ZpaReauthIntervalInMinutes != patchedInterval {
		t.Errorf("Expected partially updated zpaReauthIntervalInMinutes '%d', but got '%d'", patchedInterval, patchedFromGet.ZPANotificationTemplate.ZpaReauthIntervalInMinutes)
	}
	if !strings.EqualFold(patchedFromGet.Name, updatedName) {
		t.Errorf("Expected PATCH to preserve name '%s', but got '%s'", updatedName, patchedFromGet.Name)
	}
	if !patchedFromGet.ZIANotificationTemplate.EnableZiaFirewall {
		t.Error("Expected PATCH to preserve enableZiaFirewall=true")
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
