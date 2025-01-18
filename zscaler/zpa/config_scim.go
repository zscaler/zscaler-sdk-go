package zpa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
	rl "github.com/zscaler/zscaler-sdk-go/v3/ratelimiter"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
)

const (
	ZPA_SCIM_TOKEN     = "ZPA_SCIM_TOKEN"
	ZPA_IDP_ID         = "ZPA_IDP_ID"
	ZPA_SCIM_CLOUD     = "ZPA_SCIM_CLOUD"
	scimdefaultBaseURL = "https://scim1.private.zscaler.com/scim/1/"
	scimbetaBaseURL    = "https://scim1.zpabeta.ne/scim/1/"
	scimzpaTwoBaseUrl  = "https://scim1.zpatwo.net/scim/1/"
	scimgovBaseURL     = "https://scim1.zpagov.net/scim/1/"
	scimgovUsBaseURL   = "https://scim1.zpagov.us/scim/1/"
	scimpreviewBaseUrl = "https://scim1.zpapreview.net/scim/1/"
	defaultTimeout     = 240 * time.Second
)

var defaultBackoffConf = &BackoffConfig{
	Enabled:             true,
	MaxNumOfRetries:     100,
	RetryWaitMaxSeconds: 10,
	RetryWaitMinSeconds: 2,
}

type BackoffConfig struct {
	Enabled             bool // Set to true to enable backoff and retry mechanism
	RetryWaitMinSeconds int  // Minimum time to wait
	RetryWaitMaxSeconds int  // Maximum time to wait
	MaxNumOfRetries     int  // Maximum number of retries
}

type ScimClient struct {
	ScimConfig *ScimConfig
}

type ScimConfig struct {
	BaseURL     *url.URL
	httpClient  *http.Client
	AuthToken   string
	IDPId       string
	Logger      logger.Logger
	rateLimiter *rl.RateLimiter
	BackoffConf *BackoffConfig
	UserAgent   string
}

type ScimConfigSetter func(*ScimConfig)

// WithScimToken sets the SCIM token in the configuration.
func WithScimToken(scimToken string) ScimConfigSetter {
	return func(c *ScimConfig) {
		c.AuthToken = scimToken
	}
}

// WithIDPId sets the IDP ID in the configuration.
func WithIDPId(idpId string) ScimConfigSetter {
	return func(c *ScimConfig) {
		c.IDPId = idpId
	}
}

// WithScimBaseURL sets the SCIM BaseURL in the configuration.
func WithScimCloud(baseURL string) ScimConfigSetter {
	return func(c *ScimConfig) {
		parsedURL, err := url.Parse(baseURL)
		if err == nil {
			c.BaseURL = parsedURL
		} else {
			c.Logger.Printf("[ERROR] Invalid base URL: %v", err)
		}
	}
}

// WithScimUserAgent sets the User-Agent in the configuration.
func WithScimUserAgent(userAgent string) ScimConfigSetter {
	return func(c *ScimConfig) {
		c.UserAgent = userAgent
	}
}

// WithScimTimeout sets the HTTP client timeout in the configuration.
func WithScimTimeout(timeout time.Duration) ScimConfigSetter {
	return func(c *ScimConfig) {
		if c.httpClient != nil {
			c.httpClient.Timeout = timeout
		} else {
			c.httpClient = &http.Client{Timeout: timeout}
		}
	}
}

// WithScimRateLimiter sets the rate limiter in the configuration.
func WithScimRateLimiter(rateLimiter *rl.RateLimiter) ScimConfigSetter {
	return func(c *ScimConfig) {
		c.rateLimiter = rateLimiter
	}
}

// NewScimConfig initializes a configuration specifically for SCIM-based API endpoints
func NewScimConfig(setters ...ScimConfigSetter) (*ScimClient, error) {
	var logger logger.Logger = logger.GetDefaultLogger(loggerPrefix)

	// Default configuration values
	scimConfig := &ScimConfig{
		BaseURL:     nil,
		AuthToken:   os.Getenv(ZPA_SCIM_TOKEN),
		IDPId:       os.Getenv(ZPA_IDP_ID),
		Logger:      logger,
		httpClient:  &http.Client{Timeout: defaultTimeout},
		BackoffConf: defaultBackoffConf,
		UserAgent:   fmt.Sprintf("zscaler-sdk-go/%s golang/%s %s/%s", VERSION, runtime.Version(), runtime.GOOS, runtime.GOARCH),
		rateLimiter: rl.NewRateLimiter(20, 10, 10, 10),
	}

	// Apply setters to customize configuration
	for _, setter := range setters {
		setter(scimConfig)
	}

	// Validate required configuration fields
	if scimConfig.AuthToken == "" || scimConfig.IDPId == "" {
		return nil, fmt.Errorf("scim token and idp id are required for SCIM-based configuration")
	}

	// Set the base URL based on the SCIM cloud environment
	baseURL := os.Getenv(ZPA_SCIM_CLOUD)
	if baseURL == "" {
		baseURL = "PRODUCTION" // Default to production
	}

	switch strings.ToUpper(baseURL) {
	case "BETA":
		scimConfig.BaseURL, _ = url.Parse(scimbetaBaseURL)
	case "ZPATWO":
		scimConfig.BaseURL, _ = url.Parse(scimzpaTwoBaseUrl)
	case "GOV":
		scimConfig.BaseURL, _ = url.Parse(scimgovBaseURL)
	case "GOVUS":
		scimConfig.BaseURL, _ = url.Parse(scimgovUsBaseURL)
	case "PREVIEW":
		scimConfig.BaseURL, _ = url.Parse(scimpreviewBaseUrl)
	case "PRODUCTION", "":
		scimConfig.BaseURL, _ = url.Parse(scimdefaultBaseURL)
	default:
		return nil, fmt.Errorf("invalid SCIM cloud: %s", baseURL)
	}

	// Return the SCIM client
	return &ScimClient{ScimConfig: scimConfig}, nil
}

// DoRequest performs an HTTP request specifically for SCIM endpoints with enhanced logging
func (c *ScimClient) DoRequest(ctx context.Context, method, endpoint string, payload interface{}, target interface{}) (*http.Response, error) {
	fullURL := fmt.Sprintf("%s%s", c.ScimConfig.BaseURL.String(), endpoint)
	reqID := uuid.NewString() // Generate a unique request ID
	start := time.Now()

	// Marshal payload if provided
	var reqBody io.Reader
	if payload != nil {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			c.ScimConfig.Logger.Printf("[ERROR] Failed to marshal payload: %v", err)
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonPayload)
	}

	// Create the HTTP request with context
	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		c.ScimConfig.Logger.Printf("[ERROR] Failed to create request: %v", err)
		return nil, err
	}

	// Add headers, including the Authorization token
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.ScimConfig.AuthToken))
	if c.ScimConfig.UserAgent != "" {
		req.Header.Add("User-Agent", c.ScimConfig.UserAgent)
	}

	// Log the request
	logger.LogRequest(c.ScimConfig.Logger, req, reqID, nil, true)

	// Send the HTTP request
	resp, err := c.ScimConfig.httpClient.Do(req)
	if err != nil {
		c.ScimConfig.Logger.Printf("[ERROR] Error occurred during request: %v", err)
		return nil, err
	}

	// Read and log the response data
	respData, err := io.ReadAll(resp.Body)
	if err == nil {
		resp.Body = io.NopCloser(bytes.NewBuffer(respData))
	}

	// Check for errors in the response status or body
	if err := errorx.CheckErrorInResponse(resp, err); err != nil {
		return resp, err
	}

	// Decode JSON into target if provided and no error
	if target != nil {
		if err := decodeJSON(respData, target); err != nil {
			c.ScimConfig.Logger.Printf("[ERROR] Failed to decode response: %v", err)
			return resp, err
		}
	}

	// Log the response details with the same reqID
	logger.LogResponse(c.ScimConfig.Logger, resp, start, reqID)

	// Optional: Unescape HTML in target, if necessary
	unescapeHTML(target)

	c.ScimConfig.Logger.Printf("[DEBUG] Successfully completed request to %s", fullURL)
	return resp, nil
}
