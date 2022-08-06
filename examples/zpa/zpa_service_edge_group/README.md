```go
package main

import (
	"log"
	"os"

	"github.com/zscaler/zscaler-sdk-go/zpa"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/serviceedgegroup"
)

func main() {
	/*
		If you set one of the value of the parameters to empty string, the client will fallback to:
		 - The env variables: ZPA_CLIENT_ID, ZPA_CLIENT_SECRET, ZPA_CUSTOMER_ID, ZPA_CLOUD
		 - Or if the env vars are not set, the client will try to use the config file which should be placed at  $HOME/.zpa/credentials.json on Linux and OS X, or "%USERPROFILE%\.zpa/credentials.json" on windows
		 	with the following format:
			{
				"zpa_client_id": "",
				"zpa_client_secret": "",
				"zpa_customer_id": "",
				"zpa_cloud": ""
			}
	*/
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
	serviceEdgeGroupService := serviceedgegroup.New(zpaClient)
	app := serviceedgegroup.ServiceEdgeGroup{
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
		VersionProfileID:       "2",
	}
	// Create new service edge group
	createdResource, _, err := serviceEdgeGroupService.Create(app)
	if err != nil {
		log.Printf("[ERROR] creating service edge group failed: %v\n", err)
		return
	}
	// Update service edge group
	createdResource.Description = "New description"
	_, err = serviceEdgeGroupService.Update(createdResource.ID, createdResource)
	if err != nil {
		log.Printf("[ERROR] updating service edge group failed: %v\n", err)
		return
	}
	// Delete service edge group
	_, err = serviceEdgeGroupService.Delete(createdResource.ID)
	if err != nil {
		log.Printf("[ERROR] deleting service edge group failed: %v\n", err)
		return
	}
}
```