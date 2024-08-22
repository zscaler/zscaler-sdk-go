package cmd

import (
	"fmt"
	"os"

	"github.com/briandowns/spinner"
	"github.com/mattn/go-isatty"
	"github.com/spf13/viper"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx"
	"go.uber.org/zap"
)

type cliState struct {
	Profile    string
	Account    string
	KeyID      string
	Secret     string
	Token      string
	OrgLevel   bool
	CfgVersion int

	// Zscaler API Client (assuming you have a client package)
	ZDXApi     *zdx.Client
	Log        *zap.SugaredLogger
	JsonOutput bool
	YamlOutput bool
	CsvOutput  bool

	spinner         *spinner.Spinner
	nonInteractive  bool
	profileDetails  map[string]interface{}
	tokenCache      string
	componentParser componentArgParser
}

// NewDefaultState creates a new cliState with some defaults
func NewDefaultState() *cliState {
	c := &cliState{
		Profile:        "default",
		CfgVersion:     1,
		Log:            zap.NewExample().Sugar(),
		JsonOutput:     false,
		YamlOutput:     false,
		CsvOutput:      false,
		nonInteractive: !isatty.IsTerminal(os.Stdout.Fd()),
	}

	return c
}

// SetProfile sets the provided profile into the cliState and loads the entire
// state of the Zscaler CLI by calling 'LoadState()'
func (c *cliState) SetProfile(profile string) error {
	if profile == "" {
		return fmt.Errorf("specify a profile")
	}

	c.Profile = profile
	c.Log.Debugw("custom profile", "profile", profile)
	return c.LoadState()
}

// LoadState loads the state of the CLI by reading the configuration and environment variables
func (c *cliState) LoadState() error {
	c.profileDetails = viper.GetStringMap(c.Profile)
	if len(c.profileDetails) == 0 && c.Profile != "default" {
		return fmt.Errorf("the profile '%s' could not be found. Try running 'zscaler configure --profile %s'", c.Profile, c.Profile)
	}

	c.Token = c.extractValueString("api_token")
	c.KeyID = c.extractValueString("api_key")
	c.Secret = c.extractValueString("api_secret")
	c.Account = c.extractValueString("account")
	version := c.extractValueInt("version")
	if version > 1 {
		c.CfgVersion = version
	}

	c.Log.Debugw("state loaded",
		"profile", c.Profile,
		"account", c.Account,
		"api_token", c.Token,
		"api_key", c.KeyID,
		"api_secret", c.Secret,
		"config_version", c.CfgVersion,
	)

	return nil
}

// NewClient creates and stores a new Zscaler API client
func (c *cliState) NewClient() error {
	// Implement the logic to initialize your API client here
	// Example:
	apiOpts := []zdx.Option{
		zdx.WithApiKeys(c.KeyID, c.Secret),
	}

	client, err := zdx.NewClient(c.Account, apiOpts...)
	if err != nil {
		return fmt.Errorf("unable to generate API client: %w", err)
	}

	c.ZDXApi = client
	return nil
}

func (c *cliState) extractValueString(key string) string {
	if val, ok := c.profileDetails[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
		c.Log.Warnw("config value type mismatch",
			"expected_type", "string",
			"actual_type", fmt.Sprintf("%T", val),
			"file", viper.ConfigFileUsed(),
			"profile", c.Profile,
			"key", key,
			"value", val,
		)
		return ""
	}
	c.Log.Warnw("unable to find key from config",
		"file", viper.ConfigFileUsed(),
		"profile", c.Profile,
		"key", key,
	)
	return ""
}

func (c *cliState) extractValueInt(key string) int {
	if val, ok := c.profileDetails[key]; ok {
		if i, ok := val.(int); ok {
			return i
		}
		c.Log.Warnw("config value type mismatch",
			"expected_type", "int",
			"actual_type", fmt.Sprintf("%T", val),
			"file", viper.ConfigFileUsed(),
			"profile", c.Profile,
			"key", key,
			"value", val,
		)
		return 0
	}
	c.Log.Warnw("unable to find key from config",
		"file", viper.ConfigFileUsed(),
		"profile", c.Profile,
		"key", key,
	)
	return 0
}
