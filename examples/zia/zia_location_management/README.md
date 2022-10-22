```go
package main

import (
	"log"

	"github.com/zscaler/zscaler-sdk-go/zia"
	"github.com/zscaler/zscaler-sdk-go/zia/services/locationmanagement"
)

func main() {
	username := os.Getenv("ZIA_USERNAME")
	password := os.Getenv("ZIA_PASSWORD")
	apiKey := os.Getenv("ZIA_API_KEY")
	ziaCloud := os.Getenv("ZIA_CLOUD")

	cli, err := zia.NewClient(username, password, apiKey, ziaCloud, "userAgent")
	if err != nil {
		log.Printf("[ERROR] creating client failed: %v\n", err)
		return
	}

	locationManagementService := locationmanagement.New(cli)
	locationManagement := locationmanagement.Locations{
		Name:              "Example location managment",
		Description:       "Example location managment",
		Country:           "UNITED_STATES",
		TZ:                "UNITED_STATES_AMERICA_LOS_ANGELES",
		AuthRequired:      true,
		IdleTimeInMinutes: 720,
		DisplayTimeUnit:   "HOUR",
		SurrogateIP:       true,
		XFFForwardEnabled: true,
		OFWEnabled:        true,
		IPSControl:        true,
		IPAddresses:       []string{"1.1.1.1"},
	}
	// Create new location managment
	createLocationManagement, err := locationManagementService.Create(&locationManagement)
	if err != nil {
		log.Printf("[ERROR] creating location managment failed: %v\n", err)
		return
	}
	// Update location managment
	createLocationManagement.Description = "New description"
	_, _, err = locationManagementService.Update(createLocationManagement.ID, createLocationManagement)
	if err != nil {
		log.Printf("[ERROR] updating location managment failed: %v\n", err)
		return
	}
	// Delete location managment
	_, err = locationManagementService.Delete(createLocationManagement.ID)
	if err != nil {
		log.Printf("[ERROR] deleting location managment failed: %v\n", err)
		return
	}
}

```
