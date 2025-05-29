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

	updatedUser := &scim_api.SCIMUser{
		Schemas: []string{
			"urn:ietf:params:scim:schemas:core:2.0:User",
			"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User",
		},
		UserName:    createdUser.UserName,
		DisplayName: createdUser.DisplayName,
		EnterpriseExtension: &scim_api.EnterpriseUser{
			Department: "Finance",
		},
	}

	_, err = scim_api.UpdateUser(ctx, service, createdUser.ID, updatedUser)
	if err != nil {
		log.Fatalf("Error updating SCIM user: %v", err)
	}
	log.Printf("User updated: ID=%s, Department changed to Finance", createdUser.ID)
}
```