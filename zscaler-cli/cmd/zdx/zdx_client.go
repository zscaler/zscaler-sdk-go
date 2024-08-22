package zdx

import (
	"fmt"

	"github.com/zscaler/zscaler-sdk-go/v2/zdx"
)

// InitZDXClient initializes the ZDX API client using the configuration from cliState
func InitZDXClient(state *cmd.cliState) (*zdx.Client, error) {
	// Create a new Config using the API Key and Secret from the cliState
	config, err := zdx.NewConfig(state.KeyID, state.Secret, "zscaler-cli")
	if err != nil {
		return nil, fmt.Errorf("error creating ZDX config: %w", err)
	}

	// Optionally, set the base URL if it's different from the default
	if state.Account != "" {
		config.BaseURL, err = zdx.NewBaseURL(state.Account)
		if err != nil {
			return nil, fmt.Errorf("invalid ZDX API URL: %w", err)
		}
	}

	// Instantiate the ZDX API client
	client := zdx.NewClient(config)
	return client, nil
}
