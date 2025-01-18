```go
package main

import (
	"log"

	"github.com/zscaler/zscaler-sdk-go/v3/zia"
	"github.com/zscaler/zscaler-sdk-go/v3/zia/services/firewallpolicies/networkapplications"
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

	networkApplicationGroupsService := networkapplications.New(cli)
	networkApplicationGroupsService := networkapplications.NetworkApplicationGroups{
		Name:           "Example Network Application Groups",
		Description:    "Example Network Application Groups",
        NetworkApplications:      []string{"YAMMER", "OFFICE365", "SKYPE_FOR_BUSINESS", "OUTLOOK","SHAREPOINT",
											"SHAREPOINT_ADMIN",
											"SHAREPOINT_BLOG",
											"SHAREPOINT_CALENDAR",
											"SHAREPOINT_DOCUMENT",
											"SHAREPOINT_ONLINE",
											"ONEDRIVE"
			},
	}
	// Create new ip destination group
	createNetworkApplicationGroups, err := networkApplicationGroupsService.Create(&NetworkApplicationGroups)
	if err != nil {
		log.Printf("[ERROR] creating ip destination group failed: %v\n", err)
		return
	}
	// Update ip destination group
	createNetworkApplicationGroups.Description = "New description"
	_, _, err = networkApplicationGroupsService.Update(networkApplicationGroupsService.ID, createNetworkApplicationGroups)
	if err != nil {
		log.Printf("[ERROR] updating ip destination group failed: %v\n", err)
		return
	}
	// Delete ip destination group
	_, err = networkApplicationGroupsService.Delete(createNetworkApplicationGroups.ID)
	if err != nil {
		log.Printf("[ERROR] deleting ip destination group failed: %v\n", err)
		return
	}
}

```
