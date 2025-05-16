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

	foundByName, _, err := scim_api.GetUserByName(ctx, service, "John Doe")
	if err != nil {
		log.Fatalf("Error fetching SCIM user by name: %v", err)
	}
	log.Printf("User fetched by DisplayName: ID=%s, Department=%s", foundByName.ID, foundByName.EnterpriseExtension.Department)

}
```