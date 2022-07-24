```go
package main

import (
	"log"

	"github.com/zscaler/zscaler-sdk-go/zia"
	"github.com/zscaler/zscaler-sdk-go/zia/services/trafficforwarding/staticips"
)

func main() {
	cli, err := zia.NewClient("username@acme.com", "pwd", "apiKey", "zscalerthree", "userAgent")
	if err != nil {
		log.Printf("[ERROR] creating client failed: %v\n", err)
		return
	}
	staticIpsService := staticips.New(cli)
	staticIps := staticips.StaticIP{
		IpAddress:   "3.1.1.1",
		RoutableIP:  true,
		Comment:     "Example static ip",
		GeoOverride: true,
		Latitude:    -36.848461,
		Longitude:   174.763336,
	}
	// Create new static
	createStaticIps, _, err := staticIpsService.Create(&staticIps)
	if err != nil {
		log.Printf("[ERROR] creating static failed: %v\n", err)
		return
	}
	// Update static
	createStaticIps.Comment = "New comment"
	_, _, err = staticIpsService.Update(createStaticIps.ID, createStaticIps)
	if err != nil {
		log.Printf("[ERROR] updating static failed: %v\n", err)
		return
	}
	// Delete static
	_, err = staticIpsService.Delete(createStaticIps.ID)
	if err != nil {
		log.Printf("[ERROR] deleting static failed: %v\n", err)
		return
	}
}

```