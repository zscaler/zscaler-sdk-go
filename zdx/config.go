package zdx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/zscaler/zscaler-sdk-go/v2/logger"
)

const (
	defaultBaseURL        = "https://api.zdxcloud.net"
	defaultTimeout        = 240 * time.Second
	loggerPrefix          = "zdx-logger: "
	ZDX_API_KEY_ID        = "ZDX_API_KEY_ID"
	ZDX_API_SECRET        = "ZDX_API_SECRET"
	configPath     string = ".zdx/credentials.json"
)

var defaultBackoffConf = &BackoffConfig{
	Enabled:             true,
	MaxNumOfRetries:     100,
	RetryWaitMaxSeconds: 20,
	RetryWaitMinSeconds: 5,
}

type BackoffConfig struct {
	Enabled             bool // Set to true to enable backoff and retry mechanism
	RetryWaitMinSeconds int  // Minimum time to wait
	RetryWaitMaxSeconds int  // Maximum time to wait
	MaxNumOfRetries     int  // Maximum number of retries
}

type AuthToken struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"token"`
	ExpiresIn   int    `json:"expires_in"`
}

type CredentialsConfig struct {
	APIKeyID  string `json:"zdx_api_key_id"`
	APISecret string `json:"zdx_api_secret"`
}

// Config contains all the configuration data for the API client
type Config struct {
	BaseURL    *url.URL
	httpClient *http.Client
	// The logger writer interface to write logging messages to. Defaults to standard out.
	Logger logger.Logger
	// Credentials for basic authentication.
	APIKeyID, APISecret string
	// Backoff config
	BackoffConf *BackoffConfig
	AuthToken   *AuthToken
	sync.Mutex
	UserAgent string
}

func NewConfig(apiKeyID, apiSecret, userAgent string) (*Config, error) {
	var logger logger.Logger = logger.GetDefaultLogger(loggerPrefix)
	if apiKeyID == "" || apiSecret == "" {
		apiKeyID = os.Getenv(ZDX_API_KEY_ID)
		apiSecret = os.Getenv(ZDX_API_SECRET)
	}
	// last resort to configuration file:
	if apiKeyID == "" || apiSecret == "" {
		creds, err := loadCredentialsFromConfig(logger)
		if err != nil || creds == nil {
			return nil, err
		}
		apiKeyID = creds.APIKeyID
		apiSecret = creds.APISecret
	}
	rawUrl := defaultBaseURL

	baseURL, err := url.Parse(rawUrl)
	if err != nil {
		logger.Printf("[ERROR] error occurred while configuring the client: %v", err)
	}
	return &Config{
		BaseURL:     baseURL,
		Logger:      logger,
		httpClient:  nil,
		APIKeyID:    apiKeyID,
		APISecret:   apiSecret,
		BackoffConf: defaultBackoffConf,
		UserAgent:   userAgent,
	}, err
}

func (c *Config) SetBackoffConfig(backoffConf BackoffConfig) {
	c.BackoffConf = &backoffConf
}

// loadCredentialsFromConfig Returns the credentials found in a config file
func loadCredentialsFromConfig(logger logger.Logger) (*CredentialsConfig, error) {
	usr, _ := user.Current()
	dir := usr.HomeDir
	path := filepath.Join(dir, configPath)
	logger.Printf("[INFO]Loading configuration file at:%s", path)
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New("Could not open credentials file, needs to contain one json object with keys: zdx_api_key_id, zdx_api_secret, and zdx_cloud. " + err.Error())
	}
	configBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var config CredentialsConfig
	err = json.Unmarshal(configBytes, &config)
	if err != nil || config.APIKeyID == "" || config.APISecret == "" {
		return nil, fmt.Errorf("could not parse credentials file, needs to contain one json object with keys: zdx_api_key_id and zdx_api_secret. error: %v", err)
	}
	return &config, nil
}

func (c *Config) GetHTTPClient() *http.Client {
	if c.httpClient == nil {
		if c.BackoffConf != nil && c.BackoffConf.Enabled {
			retryableClient := retryablehttp.NewClient()
			retryableClient.Logger = c.Logger
			retryableClient.RetryWaitMin = time.Second * time.Duration(c.BackoffConf.RetryWaitMinSeconds)
			retryableClient.RetryWaitMax = time.Second * time.Duration(c.BackoffConf.RetryWaitMaxSeconds)
			retryableClient.RetryMax = c.BackoffConf.MaxNumOfRetries
			retryableClient.HTTPClient.Transport = logging.NewSubsystemLoggingHTTPTransport("gozscaler", retryableClient.HTTPClient.Transport)
			retryableClient.CheckRetry = checkRetry
			retryableClient.HTTPClient.Timeout = defaultTimeout
			c.httpClient = retryableClient.StandardClient()
		} else {
			c.httpClient = &http.Client{
				Timeout: defaultTimeout,
			}
		}
	}
	return c.httpClient
}

func containsInt(codes []int, code int) bool {
	for _, a := range codes {
		if a == code {
			return true
		}
	}
	return false
}

// getRetryOnStatusCodes return a list of http status codes we want to apply retry on.
// return empty slice to enable retry on all connection & server errors.
// or return []int{429}  to retry on only TooManyRequests error
func getRetryOnStatusCodes() []int {
	return []int{http.StatusTooManyRequests}
}

// Used to make http client retry on provided list of response status codes
func checkRetry(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	if resp != nil && containsInt(getRetryOnStatusCodes(), resp.StatusCode) {
		return true, nil
	}
	return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
}
