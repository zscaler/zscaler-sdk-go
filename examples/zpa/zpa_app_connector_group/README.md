```go
package main

import (
	"log"
	"os"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"
)

func main() {
	zpa_client_id := os.Getenv("ZPA_CLIENT_ID")
	zpa_client_secret := os.Getenv("ZPA_CLIENT_SECRET")
	zpa_customer_id := os.Getenv("ZPA_CUSTOMER_ID")
	zpa_cloud := os.Getenv("ZPA_CLOUD")
	config, err := zpa.NewConfig(zpa_client_id, zpa_client_secret, zpa_customer_id, zpa_cloud, "userAgent")
	if err != nil {
		log.Printf("[ERROR] creating config failed: %v\n", err)
		return
	}
	zpaClient := zpa.NewClient(config)
	appConnectorGroupService := appconnectorgroup.New(zpaClient)
	app := appconnectorgroup.AppConnectorGroup{
		Name:                   "Example app connector group",
		Description:            "Example  app connector group",
		Enabled:                true,
		CityCountry:            "California, US",
		CountryCode:            "US",
		Latitude:               "37.3382082",
		Longitude:              "-121.8863286",
		Location:               "San Jose, CA, USA",
		UpgradeDay:             "SUNDAY",
		UpgradeTimeInSecs:      "66600",
		OverrideVersionProfile: true,
		VersionProfileID:       "0",
		DNSQueryType:           "IPV4",
	}
	// Create new app connector group
	createdResource, _, err := appConnectorGroupService.Create(app)
	if err != nil {
		log.Printf("[ERROR] creating app connector group failed: %v\n", err)
		return
	}
	// Update app connector group
	createdResource.Description = "New description"
	_, err = appConnectorGroupService.Update(createdResource.ID, createdResource)
	if err != nil {
		log.Printf("[ERROR] updating app connector group failed: %v\n", err)
		return
	}
	// Delete app connector group
	_, err = appConnectorGroupService.Delete(createdResource.ID)
	if err != nil {
		log.Printf("[ERROR] deleting app connector group failed: %v\n", err)
		return
	}
}
```
