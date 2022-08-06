```go
package main

import (
	"log"

	"github.com/zscaler/zscaler-sdk-go/zia"
	"github.com/zscaler/zscaler-sdk-go/zia/services/firewallpolicies/ipsourcegroups"
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

	ipSourceGroupService := ipsourcegroups.New(cli)
	ipSourceGroups := ipsourcegroups.IPSourceGroups{
		Name:              "Example IP Source Groups",
		Description:       "Example IP Source Groups",
        IPAddresses:       []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"},
	}
	// Create new ip source group
	createIPSourceGroup, err := ipSourceGroupService.Create(&ipSourceGroups)
	if err != nil {
		log.Printf("[ERROR] creating ip source group failed: %v\n", err)
		return
	}
	// Update ip source group
	createIPSourceGroup.Description = "New description"
	_, _, err = ipSourceGroupService.Update(createIPSourceGroup.ID, createIPSourceGroup)
	if err != nil {
		log.Printf("[ERROR] updating ip source group failed: %v\n", err)
		return
	}
	// Delete ip source group
	_, err = ipSourceGroupService.Delete(createIPSourceGroup.ID)
	if err != nil {
		log.Printf("[ERROR] deleting ip source group failed: %v\n", err)
		return
	}
}

```
