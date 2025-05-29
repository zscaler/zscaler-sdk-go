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

	userByID, _, err := scim_api.GetUser(ctx, service, "83b86518-77f3-4059-9aa5-be999c5d1fc9")
	if err != nil {
		log.Fatalf("Error retrieving SCIM user by ID: %v", err)
	}
	log.Printf("User fetched by ID: ID=%s, UserName=%s", userByID.ID, userByID.UserName)
}
```