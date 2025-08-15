# Get User Groups

This example demonstrates how to retrieve groups associated with a specific user using the ZIdentity API.

## Example Code

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zidentity"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zidentity/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zidentity/services/users"
)

func main() {
	// Initialize the ZIdentity service
	service := &zscaler.Service{
		Client: &zscaler.Client{
			BaseURL: "https://your-zscaler-instance.com",
			// Add your authentication configuration here
		},
	}

	// Create zidentity service
	zidentityService := zidentity.New(service)

	// User ID to get groups for
	userID := "your-user-id"

	// Create pagination parameters (optional)
	queryParams := common.NewPaginationQueryParams(50) // Get 50 records per page
	queryParams.WithOffset(0)                          // Start from the beginning

	// Get all groups for the user
	groups, err := users.GetGroupsByUser(context.Background(), service, userID, &queryParams)
	if err != nil {
		log.Fatalf("Failed to get user groups: %v", err)
	}

	fmt.Printf("Found %d groups for user %s:\n", len(groups), userID)
	for _, group := range groups {
		fmt.Printf("- Group: %s (ID: %s, Description: %s)\n", 
			group.Name, group.ID, group.Description)
	}

	// Alternative: Get a single page of results
	pageResponse, err := users.GetGroupsByUserPage(context.Background(), service, userID, &queryParams)
	if err != nil {
		log.Fatalf("Failed to get user groups page: %v", err)
	}

	fmt.Printf("\nPage Info:\n")
	fmt.Printf("- Total Results: %d\n", pageResponse.ResultsTotal)
	fmt.Printf("- Page Offset: %d\n", pageResponse.PageOffset)
	fmt.Printf("- Page Size: %d\n", pageResponse.PageSize)
	fmt.Printf("- Records in this page: %d\n", len(pageResponse.Records))
	
	if pageResponse.NextLink != "" {
		fmt.Printf("- Has next page: Yes\n")
	}
}
```

## API Endpoint

This example uses the following API endpoint:
- **GET** `/admin/api/v1/users/{id}/groups`

## Parameters

- `id` (path, required): The user ID of the individual whose groups details are to be retrieved
- `offset` (query, optional): The starting point for pagination, with the number of records that can be skipped before fetching results (minimum: 0)
- `limit` (query, optional): The maximum number of records to return per request (minimum: 0, maximum: 1000)

## Response Structure

The API returns a paginated response with the following structure:

```json
{
  "results_total": 0,
  "pageOffset": 0,
  "pageSize": 0,
  "next_link": "string",
  "prev_link": "string",
  "records": [
    {
      "name": "string",
      "description": "string",
      "id": "string",
      "source": "UI",
      "idp": {
        "id": "string",
        "name": "string",
        "displayName": "string"
      },
      "isDynamicGroup": true,
      "adminEntitlementEnabled": true,
      "serviceEntitlementEnabled": true
    }
  ]
}
```

## Notes

- The function now supports pagination parameters (`offset` and `limit`)
- Use `GetGroupsByUser()` to retrieve all groups across all pages
- Use `GetGroupsByUserPage()` to retrieve a single page of results
- The response includes pagination metadata for handling large result sets 