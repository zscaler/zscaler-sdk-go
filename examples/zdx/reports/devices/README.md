```go
package main

import (
 "log"
 "os"

 "github.com/zscaler/zscaler-sdk-go/v2/zdx"
 "github.com/zscaler/zscaler-sdk-go/v2/zdx/services/reports/users"
)

func main() {
	apiKey := os.Getenv("ZDX_API_KEY_ID")
	apiSecret := os.Getenv("ZDX_API_SECRET")
	cfg, err := zdx.NewConfig(apiKey, apiSecret, "userAgent")
	if err != nil {
		log.Printf("[ERROR] creating client failed: %v\n", err)
		return
	}
	cli := zdx.NewClient(cfg)
	deviceService := devices.New(cli)
	devices, _, err := deviceService.GetAll(devices.GetDevicesFilters{
		Limit: 1000,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("devices: %v", devices)
}
```
