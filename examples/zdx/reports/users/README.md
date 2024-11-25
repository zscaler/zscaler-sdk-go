```go
package main

import (
 "log"
 "os"

 "github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zdx"
 "github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zdx/services/reports/users"
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
	userService := users.New(cli)
	users, _, err := userService.GetAll(users.GetUsersFilters{
		Limit: 1000,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("users: %v", users)
}
```
