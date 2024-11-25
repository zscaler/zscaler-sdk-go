```go
package main

import (
	"log"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/appservercontroller"
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
