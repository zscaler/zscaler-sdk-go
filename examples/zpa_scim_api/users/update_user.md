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
	scimToken := os.Getenv("ZPA_SCIM_TOKEN")
	scimCloud := os.Getenv("ZPA_SCIM_CLOUD")
	IdpID 	  := os.Getenv("ZPA_IDP_ID")

	scimClient, err := zpa.NewScimConfig(
		zpa.WithScimToken(scimToken),
		zpa.WithScimCloud(scimCloud),
		zpa.WithIDPId(IdpID),
	)
	if err != nil {
		log.Fatalf("Failed to create SCIM client: %v", err)
	}

	service := zscaler.NewZPAScimService(scimClient)
	ctx := context.Background()

	updatedUser := &scim_api.ScimUser{
		Schemas: []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
		ExternalID: "enterprise-" + userID,
	}
	if _, err := scim_api.UpdateUser(ctx, service, userID, updatedUser); err != nil {
		log.Fatalf("Error updating SCIM user: %v", err)
	}
	log.Printf("User updated: ID=%s, ExternalID=%s", userID, updatedUser.ExternalID)
}
```