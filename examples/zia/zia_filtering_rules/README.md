```go
package main

import (
	"log"
	"os"

	"github.com/zscaler/zscaler-sdk-go/v2/zia"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/filteringrules"
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
	firewallRuleService := filteringrules.New(cli)
	rules := filteringrules.FirewallFilteringRules{
		Name:              "Example100",
		Description:       "Example100",
		Action:            "ALLOW",
		Order:             1,
		EnableFullLogging: true,
	}
	// Create new static
	createFirewallRules, err := firewallRuleService.Create(&rules)
	if err != nil {
		log.Printf("[ERROR] creating firewall rule failed: %v\n", err)
		return
	}
	// Update static
	createFirewallRules.Description = "New comment"
	_, err = firewallRuleService.Update(createFirewallRules.ID, createFirewallRules)
	if err != nil {
		log.Printf("[ERROR] updating firewall rule comment failed: %v\n", err)
		return
	}
	// Delete static
	_, err = firewallRuleService.Delete(createFirewallRules.ID)
	if err != nil {
		log.Printf("[ERROR] deleting firewall rule failed: %v\n", err)
		return
	}
}
```
