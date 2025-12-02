package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/vcr"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa"
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

func NewOneAPIClient() (*zscaler.Service, error) {
	// Fetch credentials directly from environment variables
	clientID := os.Getenv("ZSCALER_CLIENT_ID")
	clientSecret := os.Getenv("ZSCALER_CLIENT_SECRET")
	vanityDomain := os.Getenv("ZSCALER_VANITY_DOMAIN")
	zscalerCloud := os.Getenv("ZSCALER_CLOUD")         // Optional
	sandboxToken := os.Getenv("ZSCALER_SANDBOX_TOKEN") // Optional
	sandboxCloud := os.Getenv("ZSCALER_SANDBOX_CLOUD") // Optional

	// Collect missing keys
	var missing []string
	if clientID == "" {
		missing = append(missing, "ZSCALER_CLIENT_ID")
	}
	if clientSecret == "" {
		missing = append(missing, "ZSCALER_CLIENT_SECRET")
	}
	if vanityDomain == "" {
		missing = append(missing, "ZSCALER_VANITY_DOMAIN")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	// Build the configuration
	config, err := zscaler.NewConfiguration(
		zscaler.WithClientID(clientID),
		zscaler.WithClientSecret(clientSecret),
		zscaler.WithVanityDomain(vanityDomain),
		zscaler.WithZscalerCloud(zscalerCloud),
		zscaler.WithSandboxToken(sandboxToken),
		zscaler.WithSandboxCloud(sandboxCloud),
		zscaler.WithTestingDisableHttpsCheck(false),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating configuration: %v", err)
	}

	client, err := zscaler.NewOneAPIClient(config)
	if err != nil {
		return nil, fmt.Errorf("error creating OneAPI client: %v", err)
	}

	return client, nil
}

// ============================================================================
// VCR Testing Support
// ============================================================================

// VCRTestClient wraps a client with optional VCR recorder
type VCRTestClient struct {
	Service    *zscaler.Service
	vcrClient  *vcr.VCRClient
	isVCRMode  bool
}

// Stop stops the VCR recorder if in VCR mode
func (c *VCRTestClient) Stop() {
	if c.vcrClient != nil {
		c.vcrClient.Stop()
	}
}

// IsVCRMode returns true if MOCK_TESTS is set (either recording or playback)
func IsVCRMode() bool {
	mockTests := os.Getenv("MOCK_TESTS")
	return mockTests == "true" || mockTests == "false"
}

// IsVCRPlayback returns true if in VCR playback mode (MOCK_TESTS=true)
func IsVCRPlayback() bool {
	return os.Getenv("MOCK_TESTS") == "true"
}

// IsVCRRecording returns true if in VCR recording mode (MOCK_TESTS=false)
func IsVCRRecording() bool {
	return os.Getenv("MOCK_TESTS") == "false"
}

// NewVCRTestClient creates a client that supports both VCR and real API modes
// - MOCK_TESTS=true: VCR playback mode (uses recorded cassettes)
// - MOCK_TESTS=false: VCR recording mode (records to cassettes)
// - MOCK_TESTS not set: Real API mode (no VCR)
func NewVCRTestClient(t *testing.T, cassetteName string, service string) (*VCRTestClient, error) {
	mockTests := os.Getenv("MOCK_TESTS")
	
	// If MOCK_TESTS is not set, use regular client (no VCR)
	if mockTests == "" {
		client, err := NewOneAPIClient()
		if err != nil {
			return nil, err
		}
		return &VCRTestClient{
			Service:   client,
			isVCRMode: false,
		}, nil
	}
	
	// Use VCR client
	vcrClient, err := vcr.NewVCRClient(t, cassetteName, service)
	if err != nil {
		return nil, err
	}
	
	return &VCRTestClient{
		Service:   vcrClient.Service,
		vcrClient: vcrClient,
		isVCRMode: true,
	}, nil
}

// Name generation for tests
var (
	testNameCounter int
	testNameMutex   sync.Mutex
)

// ResetTestNameCounter resets the deterministic name counter
// Call this at the beginning of each test when using VCR
func ResetTestNameCounter() {
	testNameMutex.Lock()
	defer testNameMutex.Unlock()
	testNameCounter = 0
}

// GetTestNameCounter returns the next counter value (for names that can't have separators)
func GetTestNameCounter() int {
	testNameMutex.Lock()
	defer testNameMutex.Unlock()
	testNameCounter++
	return testNameCounter
}

// GetTestName returns a name for test resources.
// - VCR mode: deterministic counter-based names for cassette consistency
// - Non-VCR mode: random suffix to avoid duplicates between runs
// Format: prefix-suffix (e.g., tests-sslins-0001 or tests-sslins-abcdefghij)
func GetTestName(prefix string) string {
	testNameMutex.Lock()
	defer testNameMutex.Unlock()
	testNameCounter++
	
	if IsVCRMode() {
		// VCR mode: deterministic names for consistent cassette matching
		return fmt.Sprintf("%s-%04d", prefix, testNameCounter)
	}
	// Non-VCR mode: random suffix to avoid duplicates
	return fmt.Sprintf("%s-%s", prefix, randomString(10))
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// GetTestIP returns a deterministic IP address in VCR mode
// In non-VCR mode, returns a random IP from the given CIDR
func GetTestIP(cidr string) string {
	if IsVCRMode() {
		testNameMutex.Lock()
		defer testNameMutex.Unlock()
		testNameCounter++
		// Generate a deterministic IP based on counter
		return fmt.Sprintf("192.168.%d.%d", (testNameCounter/256)%256, testNameCounter%256)
	}
	// For non-VCR, use a simple random IP
	return fmt.Sprintf("192.168.%d.%d", rand.Intn(256), rand.Intn(256))
}

// GetTestLatitude returns a deterministic latitude for VCR
func GetTestLatitude() float64 {
	if IsVCRMode() {
		return 37.3382
	}
	return 37.0 + rand.Float64()*10
}

// GetTestLongitude returns a deterministic longitude for VCR
func GetTestLongitude() float64 {
	if IsVCRMode() {
		return -121.8863
	}
	return -122.0 + rand.Float64()*10
}

// NewOneAPIClient instantiates a new OneAPI client for testing
// func NewOneAPIClient() (*zscaler.Service, error) {
// 	// Fetch credentials directly from environment variables
// 	clientID := os.Getenv("ZSCALER_CLIENT_ID")
// 	clientSecret := os.Getenv("ZSCALER_CLIENT_SECRET")
// 	vanityDomain := os.Getenv("ZSCALER_VANITY_DOMAIN")
// 	zscalerCloud := os.Getenv("ZSCALER_CLOUD")         // Optional, set this if needed
// 	sandboxToken := os.Getenv("ZSCALER_SANDBOX_TOKEN") // Optional, set this if needed
// 	sandboxCloud := os.Getenv("ZSCALER_SANDBOX_CLOUD") // Optional, set this if needed

// 	// Ensure required environment variables are set
// 	if clientID == "" || clientSecret == "" || vanityDomain == "" {
// 		return nil, fmt.Errorf("required environment variables (ZSCALER_CLIENT_ID, ZSCALER_CLIENT_SECRET, ZSCALER_VANITY_DOMAIN) are not set")
// 	}

// 	// Build the configuration using the environment variables
// 	config, err := zscaler.NewConfiguration(
// 		zscaler.WithClientID(clientID),
// 		zscaler.WithClientSecret(clientSecret),
// 		zscaler.WithVanityDomain(vanityDomain),
// 		zscaler.WithZscalerCloud(zscalerCloud),
// 		zscaler.WithSandboxToken(sandboxToken),
// 		zscaler.WithSandboxCloud(sandboxCloud),
// 		// zscaler.WithDebug(true),
// 		zscaler.WithTestingDisableHttpsCheck(false),
// 		// zscaler.WithUserAgentExtra("zscaler-sdk-go"),
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("error creating configuration: %v", err)
// 	}

// 	// Instantiate the OneAPI client and pass the service name (e.g., "zia")
// 	client, err := zscaler.NewOneAPIClient(config)
// 	if err != nil {
// 		return nil, fmt.Errorf("error creating OneAPI client: %v", err)
// 	}

// 	return client, nil
// }

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

func NewZPAClient() (*zscaler.Service, error) {
	// Fetch credentials from environment variables
	client_id := os.Getenv("ZPA_CLIENT_ID")
	client_secret := os.Getenv("ZPA_CLIENT_SECRET")
	customer_id := os.Getenv("ZPA_CUSTOMER_ID")
	cloud := os.Getenv("ZPA_CLOUD")

	if client_id == "" || client_secret == "" || customer_id == "" || cloud == "" {
		return nil, fmt.Errorf("missing ZPA credentials: ensure ZPA_CLIENT_ID, ZPA_CLIENT_SECRET, ZPA_CUSTOMER_ID and ZPA_CLOUD environment variables are set")
	}

	// Create a new ZPA configuration
	zpaCfg, err := zpa.NewConfiguration(
		zpa.WithZPAClientID(client_id),
		zpa.WithZPAClientSecret(client_secret),
		zpa.WithZPACustomerID(customer_id),
		zpa.WithZPACloud(cloud),
		zpa.WithDebug(false),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ZPA configuration: %w", err)
	}

	// ✅ Return the legacy client (type *zscaler.Service)
	service, err := zscaler.NewLegacyZpaClient(zpaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create ZPA legacy client: %w", err)
	}

	return service, nil
}
