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

		_, err = scim_api.DeleteGroup(ctx, service, groupID)
		if err != nil {
			log.Fatalf("Error deleting SCIM group: %v", err)
		}
		log.Printf("Group deleted: ID=%s", groupID)
}
```