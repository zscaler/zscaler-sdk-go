package sandbox_settings

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestGetBaAdvancedSettings(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := &Service{Client: client}

	settings, err := service.Get()
	if err != nil {
		t.Errorf("Error getting BA Advanced Settings: %v", err)
	}
	if settings == nil {
		t.Error("Expected BA Advanced Settings, got nil")
	}
	// You can add more assertions here to validate the contents of `settings`
}

func TestUpdateBaAdvancedSettings(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := &Service{Client: client}

	updatedSettings, err := service.Update()
	if err != nil {
		t.Errorf("Error updating BA Advanced Settings: %v", err)
	}
	if updatedSettings == nil {
		t.Error("Expected updated BA Advanced Settings, got nil")
	}
	// You can add more assertions here to validate the contents of `updatedSettings`
}

func TestGetFileHashCount(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := &Service{Client: client}

	hashCount, err := service.GetFileHashCount()
	if err != nil {
		t.Errorf("Error getting file hash count: %v", err)
	}
	if hashCount == nil {
		t.Error("Expected file hash count, got nil")
	}
	// You can add more assertions here to validate the contents of `hashCount`
}
