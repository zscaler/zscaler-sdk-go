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

	newUser := &scim_api.SCIMUser{
		Schemas: []string{
			"urn:ietf:params:scim:schemas:core:2.0:User",
			"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User",
		},
		UserName:    "jdoe@acme.com",
		DisplayName: "John Doe",
		EnterpriseExtension: &scim_api.EnterpriseUser{
			Department: "Finance",
		},
	}

	createdUser, _, err := scim_api.CreateUser(ctx, service, newUser)
	if err != nil {
		log.Fatalf("Error creating SCIM user: %v", err)
	}
	log.Printf("User created: ID=%s, Department=%s", createdUser.ID, createdUser.EnterpriseExtension.Department)
}
```