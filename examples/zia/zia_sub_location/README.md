```go
package main

import (
    "log"
    "os"

    "github.com/zscaler/zscaler-sdk-go/v3/zia"
    "github.com/zscaler/zscaler-sdk-go/v3/zia/services/locationmanagement"
)

func main() {
    username := os.Getenv("ZIA_USERNAME")
    password := os.Getenv("ZIA_PASSWORD")
    apiKey := os.Getenv("ZIA_API_KEY")
    ziaCloud := os.Getenv("ZIA_CLOUD")

    cli, err := zia.NewClient(username, password, apiKey, ziaCloud, "userAgent")
    if err != nil {
        log.Printf("[ERROR] creating client failed: %v\n", err)
        return
    }

    locationmanagementService := locationmanagement.New(cli)

    locations, err := locationmanagementService.GetSublocations(12008136)
    if err == nil {
        log.Printf("GetSublocations -> locations: %v\n", locations)
    } else {
        log.Printf("GetSublocations -> error: %v", err)
    }

    location, err := locationmanagementService.GetSubLocation(12008136, 12008137)
    if err == nil {
        log.Printf("GetSubLocation -> location: %v\n", location)
    } else {
        log.Printf("GetSubLocation -> error: %v", err)
    }

    location, err = locationmanagementService.GetSubLocationByName("Guest Wi-Fi - Branch01")
    if err == nil {
        log.Printf("GetSubLocationByName -> location: %v\n", location)
    } else {
        log.Printf("GetSubLocationByName -> error: %v", err)
    }

    location, err = locationmanagementService.GetSubLocationByNames("BR - Sao Paulo - Branch01", "Guest Wi-Fi - Branch01")
    if err == nil {
        log.Printf("GetSubLocationByNames -> location: %v\n", location)
    } else {
        log.Printf("GetSubLocationByNames -> error: %v", err)
    }

    location, err = locationmanagementService.GetSubLocationBySubID(12008137)
    if err == nil {
        log.Printf("GetSubLocationBySubID -> location: %v\n", location)
    } else {
        log.Printf("GetSubLocationBySubID -> error: %v", err)
    }

    location, err = locationmanagementService.GetLocationOrSublocationByID(12008137)
    if err == nil {
        log.Printf("GetLocationOrSublocationByID -> location: %v\n", location)
    } else {
        log.Printf("GetLocationOrSublocationByID -> error: %v", err)
    }

    location, err = locationmanagementService.GetLocationOrSublocationByName("Guest Wi-Fi - Branch01")
    if err == nil {
        log.Printf("GetLocationOrSublocationByName -> location: %v\n", location)
    } else {
        log.Printf("GetLocationOrSublocationByName -> error: %v", err)
    }
}
```
