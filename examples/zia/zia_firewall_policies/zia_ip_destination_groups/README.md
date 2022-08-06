```go
package main

import (
	"log"

	"github.com/zscaler/zscaler-sdk-go/zia"
	"github.com/zscaler/zscaler-sdk-go/zia/services/firewallpolicies/ipdestinationgroups"
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

	ipDestinationService := ipdestinationgroups.New(cli)
	ipDestinationGroups := ipdestinationgroups.IPDestinationGroups{
		Name:           "Example IP Destination Groups",
		Description:    "Example IP Destination Groups",
		Type:           "DSTN_FQDN",
        Addresses:      []string{"test1.acme.com", "test2.acme.com", "test3.acme.com" },
	}
	// Create new ip destination group
	createIPDestinationGroup, err := ipDestinationService.Create(&ipDestinationGroups)
	if err != nil {
		log.Printf("[ERROR] creating ip destination group failed: %v\n", err)
		return
	}
	// Update ip destination group
	createIPDestinationGroup.Description = "New description"
	_, _, err = ipDestinationService.Update(createIPDestinationGroup.ID, createIPDestinationGroup)
	if err != nil {
		log.Printf("[ERROR] updating ip destination group failed: %v\n", err)
		return
	}
	// Delete ip destination group
	_, err = ipDestinationService.Delete(createIPDestinationGroup.ID)
	if err != nil {
		log.Printf("[ERROR] deleting ip destination group failed: %v\n", err)
		return
	}
}

```
