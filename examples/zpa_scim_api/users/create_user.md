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

	newUser := &scim_api.ScimUser{
		Schemas: []string{
			"urn:ietf:params:scim:schemas:core:2.0:User",
			"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User",
		},
		NickName:    "John Doe",
		UserType:    "ZPA User",
		UserName:    "johndoe@acme.com",
		Active:      true,
		DisplayName: "John Doe",
		Enterprise: scim_api.EnterpriseFields{
			Department: "A000",
		},
		Name: scim_api.Name{
			Formatted:  "John Doe <acme.com>",
			FamilyName: "Doe",
			GivenName:  "John",
		},
		Emails: []scim_api.Email{
			{Value: "johndoe@acme.com"},
		},
	}
	createdUser, _, err := scim_api.CreateUser(ctx, service, newUser)
	if err != nil {
		log.Fatalf("Error creating SCIM user: %v", err)
	}
	userID := createdUser.ID
	log.Printf("User created: ID=%s, UserName=%s", userID, createdUser.UserName)
}
```