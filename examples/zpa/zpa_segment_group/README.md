```go
package main

import (
	"log"

	"github.com/zscaler/zscaler-sdk-go/zpa"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/segmentgroup"
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
	segmentGroupService := segmentgroup.New(zpaClient)
	app := segmentgroup.SegmentGroup{
		Name:                "Example application server ",
		Description:         "Example application server ",
		Enabled:             true,
		PolicyMigrated:      true,
		TcpKeepAliveEnabled  "1"
	}
	// Create new segment group
	createdResource, _, err := segmentGroupService.Create(app)
	if err != nil {
		log.Printf("[ERROR] creating segment group failed: %v\n", err)
		return
	}
	// Update segment group
	createdResource.Description = "New description"
	_, err = segmentGroupService.Update(createdResource.ID, createdResource)
	if err != nil {
		log.Printf("[ERROR] updating segment group failed: %v\n", err)
		return
	}
	// Delete segment group
	_, err = segmentGroupService.Delete(createdResource.ID)
	if err != nil {
		log.Printf("[ERROR] deleting segment group failed: %v\n", err)
		return
	}
}

```