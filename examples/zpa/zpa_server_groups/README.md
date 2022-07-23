```go
package main

import (
	"log"

	"github.com/zscaler/zscaler-sdk-go/zpa"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/servergroup"
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
	config, err := zpa.NewConfig("clientID", "clientSecret", "customerID", "baseURL", "userAgent")
	if err != nil {
		log.Printf("[ERROR] creating config failed: %v\n", err)
		return
	}
	zpaClient := zpa.NewClient(config)
	zpaServerGroupsService := servergroup.New(zpaClient)
	appConnectorGroupService := appconnectorgroup.New(zpaClient)
	appConnectorGroup := appconnectorgroup.AppConnectorGroup{
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
	createdAppConnectorGroup, _, err := appConnectorGroupService.Create(appConnectorGroup)
	if err != nil {
		log.Printf("[ERROR] creating app connector group failed: %v\n", err)
		return
	}
	app := servergroup.ServerGroup{
		Name:             "Example server group",
		Description:      "Example server group",
		Enabled:          true,
		DynamicDiscovery: true,
		AppConnectorGroups: []servergroup.AppConnectorGroups{
			{
				ID: createdAppConnectorGroup.ID,
			},
		},
	}
	// Create new server group
	createdResource, _, err := zpaServerGroupsService.Create(&app)
	if err != nil {
		log.Printf("[ERROR] creating server group failed: %v\n", err)
		return
	}
	// Update server group
	createdResource.Description = "New description"
	_, err = zpaServerGroupsService.Update(createdResource.ID, createdResource)
	if err != nil {
		log.Printf("[ERROR] updating server group failed: %v\n", err)
		return
	}
	// Delete server group
	_, err = zpaServerGroupsService.Delete(createdResource.ID)
	if err != nil {
		log.Printf("[ERROR] deleting server group failed: %v\n", err)
		return
	}

	// Delete app connector group
	_, err = appConnectorGroupService.Delete(createdAppConnectorGroup.ID)
	if err != nil {
		log.Printf("[ERROR] deleting app connector group failed: %v\n", err)
		return
	}
}

```