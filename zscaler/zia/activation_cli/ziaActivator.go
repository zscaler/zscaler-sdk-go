package main

import (
	"log"
	"os"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/activation"
)

func getEnvVarOrFail(k string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	log.Fatalf("[ERROR] Couldn't find environment variable %s\n", k)
	return ""
}

func main() {
	log.Printf("[INFO] Initializing ZIA client\n")

	// Retrieve the necessary environment variables
	clientID := getEnvVarOrFail("ZSCALER_CLIENT_ID")
	clientSecret := getEnvVarOrFail("ZSCALER_CLIENT_SECRET")
	vanityDomain := getEnvVarOrFail("ZSCALER_VANITY_DOMAIN")
	zscalerCloud := os.Getenv("ZSCALER_CLOUD") // Optional: might not be required, set as empty if not

	// Create a configuration using the environment variables
	config, err := zscaler.NewConfiguration(
		zscaler.WithClientID(clientID),
		zscaler.WithClientSecret(clientSecret),
		zscaler.WithVanityDomain(vanityDomain),
		zscaler.WithZscalerCloud(zscalerCloud), // Optional, can be an empty string if not set
		zscaler.WithUserAgentExtra("zscaler-sdk-go"),
	)
	if err != nil {
		log.Fatalf("[ERROR] Failed creating ZIA configuration: %v\n", err)
	}

	// Instantiate the ZIA client with the service name "zia"
	cli, err := zscaler.NewOneAPIClient(config)
	if err != nil {
		log.Fatalf("[ERROR] Failed Initializing ZIA client: %v\n", err)
	}

	// Call the activation API using the instantiated client
	resp, err := activation.CreateActivation(cli, activation.Activation{
		Status: "active",
	})
	if err != nil {
		log.Printf("[ERROR] Activation Failed: %v\n", err)
	} else {
		log.Printf("[INFO] Activation succeeded: %#v\n", resp)
	}

	os.Exit(0)
}
