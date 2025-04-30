package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw"
)

const (
	charSetAlphaUpper  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charSetAlphaLower  = "abcdefghijklmnopqrstuvwxyz"
	charSetNumeric     = "0123456789"
	charSetSpecialChar = "!@#$%^&*"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func TestPassword(length int) string {
	if length < 8 {
		length = 8
	} else if length > 100 {
		length = 100
	}

	result := make([]byte, length)
	result[0] = charSetAlphaLower[rand.Intn(len(charSetAlphaLower))]
	result[1] = charSetAlphaUpper[rand.Intn(len(charSetAlphaUpper))]
	result[2] = charSetNumeric[rand.Intn(len(charSetNumeric))]
	result[3] = charSetSpecialChar[rand.Intn(len(charSetSpecialChar))]

	charSetAll := charSetAlphaLower + charSetAlphaUpper + charSetNumeric + charSetSpecialChar
	for i := 4; i < length; i++ {
		result[i] = charSetAll[rand.Intn(len(charSetAll))]
	}
	// Shuffle the result to avoid predictable patterns (lower, upper, numeric, special)
	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	return string(result)
}

// NewOneAPIClient instantiates a new OneAPI client for testing
func NewOneAPIClient() (*zscaler.Service, error) {
	// Fetch credentials directly from environment variables
	clientID := os.Getenv("ZSCALER_CLIENT_ID")
	clientSecret := os.Getenv("ZSCALER_CLIENT_SECRET")
	vanityDomain := os.Getenv("ZSCALER_VANITY_DOMAIN")
	zscalerCloud := os.Getenv("ZSCALER_CLOUD")         // Optional, set this if needed
	sandboxToken := os.Getenv("ZSCALER_SANDBOX_TOKEN") // Optional, set this if needed
	sandboxCloud := os.Getenv("ZSCALER_SANDBOX_CLOUD") // Optional, set this if needed

	// Ensure required environment variables are set
	if clientID == "" || clientSecret == "" || vanityDomain == "" {
		return nil, fmt.Errorf("required environment variables (ZSCALER_CLIENT_ID, ZSCALER_CLIENT_SECRET, ZSCALER_VANITY_DOMAIN) are not set")
	}

	// Build the configuration using the environment variables
	config, err := zscaler.NewConfiguration(
		zscaler.WithClientID(clientID),
		zscaler.WithClientSecret(clientSecret),
		zscaler.WithVanityDomain(vanityDomain),
		zscaler.WithZscalerCloud(zscalerCloud),
		zscaler.WithSandboxToken(sandboxToken),
		zscaler.WithSandboxCloud(sandboxCloud),
		// zscaler.WithDebug(true),
		zscaler.WithTestingDisableHttpsCheck(false),
		// zscaler.WithUserAgentExtra("zscaler-sdk-go"),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating configuration: %v", err)
	}

	// Instantiate the OneAPI client and pass the service name (e.g., "zia")
	client, err := zscaler.NewOneAPIClient(config)
	if err != nil {
		return nil, fmt.Errorf("error creating OneAPI client: %v", err)
	}

	return client, nil
}

func NewZTWClient() (*ztw.Client, error) {
	// Fetch credentials from environment variables
	username := os.Getenv("ZTW_USERNAME")
	password := os.Getenv("ZTW_PASSWORD")
	apiKey := os.Getenv("ZTW_API_KEY")
	cloud := os.Getenv("ZTW_CLOUD")

	if username == "" || password == "" || apiKey == "" || cloud == "" {
		return nil, fmt.Errorf("missing ZTW credentials: ensure ZTW_USERNAME, ZTW_PASSWORD, ZTW_API_KEY and ZTW_CLOUD environment variables are set")
	}

	// Create a new ZTW configuration
	ztwCfg, err := ztw.NewConfiguration(
		ztw.WithZtwUsername(username),
		ztw.WithZtwPassword(password),
		ztw.WithZtwAPIKey(apiKey),
		ztw.WithZtwCloud(cloud),
		ztw.WithDebug(false),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ZTW configuration: %w", err)
	}

	// Initialize the ZTW client
	ztwClient, err := ztw.NewClient(ztwCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create ZTW client: %w", err)
	}

	return ztwClient, nil
}

// ParseJSONRequest parses the JSON request body from the given HTTP request.
func ParseJSONRequest(t *testing.T, r *http.Request, v interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.Unmarshal(body, v)
}

// WriteJSONResponse writes the JSON response with the given status code and data to the HTTP response writer.
func WriteJSONResponse(t *testing.T, w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		encoder := json.NewEncoder(w)
		return encoder.Encode(data)
	}

	return nil
}

func NewZdxClient() (*zdx.Client, error) {
	// Fetch credentials from environment variables
	clientID := os.Getenv("ZDX_API_KEY_ID")
	clientSecret := os.Getenv("ZDX_API_SECRET")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("missing ZDX credentials: ensure ZDX_API_KEY_ID and ZDX_API_SECRET environment variables are set")
	}

	// Create a new ZDX configuration
	zdxCfg, err := zdx.NewConfiguration(
		zdx.WithZDXAPIKeyID(clientID),
		zdx.WithZDXAPISecret(clientSecret),
		zdx.WithDebug(false),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ZDX configuration: %w", err)
	}

	// Initialize the ZDX client
	zdxClient, err := zdx.NewClient(zdxCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create ZDX client: %w", err)
	}

	return zdxClient, nil
}
