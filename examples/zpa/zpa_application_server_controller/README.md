```go
package main

import (
	"log"

	"github.com/zscaler/zscaler-sdk-go/zpa"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/appservercontroller"
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
	appServerControllerService := appservercontroller.New(zpaClient)
	app := appservercontroller.ApplicationServer{
		Name:                "Example application server ",
		Description:         "Example application server ",
		Enabled:             true,
		Address:             "192.168.1.1"
	}
	// Create new application server
	createdResource, _, err := appServerControllerService.Create(app)
	if err != nil {
		log.Printf("[ERROR] creating application server failed: %v\n", err)
		return
	}
	// Update application server
	createdResource.Description = "New description"
	_, err = appServerControllerService.Update(createdResource.ID, createdResource)
	if err != nil {
		log.Printf("[ERROR] updating application server  failed: %v\n", err)
		return
	}
	// Delete application server
	_, err = appServerControllerService.Delete(createdResource.ID)
	if err != nil {
		log.Printf("[ERROR] deleting application server failed: %v\n", err)
		return
	}
}

```