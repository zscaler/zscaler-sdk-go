package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	client "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcon"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcon/services"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcon/services/activation"
)

func getEnvVarOrFail(k string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	log.Fatalf("[ERROR] Couldn't find environment variable %s\n", k)
	return ""
}

func main() {
	log.Printf("[INFO] Initializing ZCON client\n")
	username := getEnvVarOrFail("ZCON_USERNAME")
	password := getEnvVarOrFail("ZCON_PASSWORD")
	apiKey := getEnvVarOrFail("ZCON_API_KEY")
	zconCloud := getEnvVarOrFail("ZCON_CLOUD")
	userAgent := fmt.Sprintf("(%s %s) cli/zconActivator", runtime.GOOS, runtime.GOARCH)

	cli, err := client.NewClient(username, password, apiKey, zconCloud, userAgent)
	if err != nil {
		log.Fatalf("[ERROR] Failed Initializing zcon client: %v\n", err)
	}
	service := services.New(cli)

	resp, err := activation.ForceActivationStatus(service, activation.ECAdminActivation{
		OrgEditStatus:         "org_edit_status",
		OrgLastActivateStatus: "org_last_activate_status",
		// AdminStatusMap:        "admin_status_map",
		AdminActivateStatus: "admin_activate_status",
	})
	if err != nil {
		log.Printf("[ERROR] Activation Failed: %v\n", err)
	} else {
		log.Printf("[INFO] Activation succeded: %#v\n", resp)
	}
	log.Printf("[INFO] Destroying session: %#v\n", resp)
	cli.Logout()
	os.Exit(0)
}
