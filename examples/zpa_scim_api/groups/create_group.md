```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/scim_api"
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

		newGroup := &scim_api.ScimGroup{
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