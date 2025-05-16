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

		newGroup := &scim_api.SCIMGroup{
			Schemas:     []string{"urn:ietf:params:scim:schemas:core:2.0:Group"},
			DisplayName: groupName,
			Members: []scim_api.SCIMGroupMember{},
		}

		createdGroup, _, err := scim_api.CreateGroup(ctx, service, newGroup)
		if err != nil {
			log.Fatalf("Error creating SCIM group: %v", err)
		}
		groupID := createdGroup.ID
		log.Printf("Group created: ID=%s, DisplayName=%s", groupID, createdGroup.DisplayName)

}
```