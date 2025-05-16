```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/scim_api"
)

func main() {
	scimToken := os.Getenv("ZIA_SCIM_API_TOKEN")
	scimCloud := os.Getenv("ZIA_SCIM_CLOUD")
	tenantID := os.Getenv("ZIA_SCIM_TENANT_ID")

	scimClient, err := zia.NewScimConfig(
		zia.WithScimToken(scimToken),
		zia.WithScimCloud(scimCloud),
		zia.WithTenantID(tenantID),
	)
	if err != nil {
		log.Fatalf("Failed to create SCIM client: %v", err)
	}

	service := zscaler.NewZIAScimService(scimClient)
	ctx := context.Background()

		groupByName, _, err := scim_api.GetGroupByName(ctx, service, updatedGroup.DisplayName)
		if err != nil {
			log.Fatalf("Error retrieving SCIM group by name: %v", err)
		}
		log.Printf("Group fetched by name: ID=%s, DisplayName=%s", groupByName.ID, groupByName.DisplayName)
}
```