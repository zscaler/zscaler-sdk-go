package downloaddevices

import (
	"os"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestDownloadDevices(t *testing.T) {
	client, err := tests.NewZccClient()
	if err != nil {
		t.Fatalf("Failed to create ZCC client: %v", err)
	}
	service := New(client)

	osTypes := "1,2"           // iOS and Android
	registrationTypes := "1,4" // Registered and Unregistered

	file, err := os.CreateTemp("", "devices-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name()) // clean up

	err = service.DownloadDevices(osTypes, registrationTypes, file)
	if err != nil {
		t.Fatalf("Error downloading devices: %v", err)
	}

	t.Logf("Devices information downloaded successfully to %s", file.Name())
}
