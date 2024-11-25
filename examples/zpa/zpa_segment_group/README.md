```go
package main

import (
	"log"
	"os"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"
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
	segmentGroupService := segmentgroup.New(zpaClient)
	app := segmentgroup.SegmentGroup{
		Name:                "Example application server ",
		Description:         "Example application server ",
		Enabled:             true,
		PolicyMigrated:      true,
		TcpKeepAliveEnabled: "1",
	}
	// Create new segment group
	createdResource, _, err := segmentGroupService.Create(&app)
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
