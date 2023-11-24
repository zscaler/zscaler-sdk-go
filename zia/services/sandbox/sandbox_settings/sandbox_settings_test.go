package sandbox_settings

import (
	"bufio"
	"os"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestUpdateBaAdvancedSettings(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := &Service{Client: client}

	// Read hashes from a file
	file, err := os.Open("hashes.txt")
	if err != nil {
		t.Fatalf("Error opening hashes file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hashes []string
	for scanner.Scan() {
		hashes = append(hashes, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("Error reading hashes from file: %v", err)
	}

	// Define the desired settings for the update
	desiredSettings := BaAdvancedSettings{
		FileHashesToBeBlocked: hashes,
	}

	updatedSettings, err := service.Update(desiredSettings)
	if err != nil {
		t.Errorf("Error updating BA Advanced Settings: %v", err)
	}
	if updatedSettings == nil {
		t.Error("Expected updated BA Advanced Settings, got nil")
	}
}

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
}

func TestEmptyHashList(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := &Service{Client: client}

	// Define the desired settings for the update
	desiredSettings := BaAdvancedSettings{
		FileHashesToBeBlocked: []string{},
	}

	updatedSettings, err := service.Update(desiredSettings)
	if err != nil {
		t.Errorf("Error updating BA Advanced Settings: %v", err)
	}
	if updatedSettings == nil {
		t.Error("Expected updated BA Advanced Settings, got nil")
	}
}
