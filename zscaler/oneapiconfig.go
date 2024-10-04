package zscaler

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zscaler/zscaler-sdk-go/v3/cache"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
	rl "github.com/zscaler/zscaler-sdk-go/v3/ratelimiter"
	"github.com/zscaler/zscaler-sdk-go/v3/utils"
)

const (
	maxIdleConnections  int = 40
	requestTimeout      int = 60
	contentTypeJSON         = "application/json"
	MaxNumOfRetries         = 100
	RetryWaitMaxSeconds     = 20
	RetryWaitMinSeconds     = 5
	loggerPrefix            = "oneapi-logger: "
)

// Client defines the ZIA client structure.
type Client struct {
	sync.Mutex
	oauth2Credentials *Configuration
	stopTicker        chan bool
}

// NewOneAPIClient creates a new client using OAuth2 authentication for any service.
func NewOneAPIClient(config *Configuration) (*Service, error) {
	cli := &Client{
		oauth2Credentials: config,
		stopTicker:        make(chan bool),
	}

	// Start token renewal ticker
	cli.startTokenRenewalTicker()

	// Return the service directly
	return NewService(cli), nil
}

// startTokenRenewalTicker starts a ticker to renew the token before it expires.
func (c *Client) startTokenRenewalTicker() {
	tokenExpiry := time.Now()
	if c.oauth2Credentials.Zscaler.Client.AuthToken != nil {
		tokenExpiry = c.oauth2Credentials.Zscaler.Client.AuthToken.Expiry
	}
	renewalInterval := time.Until(tokenExpiry) - (time.Minute * 1) // Renew 1 minute before expiration

	if renewalInterval > 0 {
		ticker := time.NewTicker(renewalInterval)
		go func() {
			for {
				select {
				case <-ticker.C:
					// Refresh the token
					authToken, err := Authenticate(c.oauth2Credentials.Context, c.oauth2Credentials, c.oauth2Credentials.Logger)
					if err != nil {
						c.oauth2Credentials.Logger.Printf("[ERROR] Failed to renew OAuth2 token: %v", err)
					} else {
						c.oauth2Credentials.Zscaler.Client.AuthToken = authToken
						c.oauth2Credentials.Logger.Printf("[INFO] OAuth2 token successfully renewed")
						// Reset the ticker for the next renewal
						renewalInterval = time.Until(authToken.Expiry) - (time.Minute * 1)
						ticker.Reset(renewalInterval)
					}
				case <-c.stopTicker:
					ticker.Stop()
					return
				}
			}
		}()
	}
}

func (client *Client) GetLogger() logger.Logger {
	return client.oauth2Credentials.Logger
}

// getHTTPClient sets up the retryable HTTP client with backoff and retry policies.
func getHTTPClient(l logger.Logger, rateLimiter *rl.RateLimiter, cfg *Configuration) *http.Client {
	retryableClient := retryablehttp.NewClient()
	retryableClient.RetryWaitMin = cfg.Zscaler.Client.RateLimit.RetryWaitMin
	retryableClient.RetryWaitMax = cfg.Zscaler.Client.RateLimit.RetryWaitMax
	retryableClient.RetryMax = int(cfg.Zscaler.Client.RateLimit.MaxRetries)

	// Configure backoff and retry policies
	retryableClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		if resp != nil {
			if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
				retryAfter := getRetryAfter(resp, l)
				if retryAfter > 0 {
					return retryAfter
				}
			}
			if resp.Request != nil {
				wait, d := rateLimiter.Wait(resp.Request.Method)
				if wait {
					return d
				}
				return 0
			}
		}
		mult := math.Pow(2, float64(attemptNum)) * float64(min)
		sleep := time.Duration(mult)
		if float64(sleep) != mult || sleep > max {
			sleep = max
		}
		return sleep
	}
	retryableClient.CheckRetry = checkRetry
	retryableClient.Logger = l
	retryableClient.HTTPClient.Timeout = time.Duration(requestTimeout) * time.Second

	// Configure proxy settings from configuration
	proxyFunc := http.ProxyFromEnvironment // Default behavior (uses system/env variables)
	if cfg.Zscaler.Client.Proxy.Host != "" {
		// Include username and password if provided
		proxyURLString := fmt.Sprintf("http://%s:%d", cfg.Zscaler.Client.Proxy.Host, cfg.Zscaler.Client.Proxy.Port)
		if cfg.Zscaler.Client.Proxy.Username != "" && cfg.Zscaler.Client.Proxy.Password != "" {
			// URL-encode the username and password
			proxyAuth := url.UserPassword(cfg.Zscaler.Client.Proxy.Username, cfg.Zscaler.Client.Proxy.Password)
			proxyURLString = fmt.Sprintf("http://%s@%s:%d", proxyAuth.String(), cfg.Zscaler.Client.Proxy.Host, cfg.Zscaler.Client.Proxy.Port)
		}

		proxyURL, err := url.Parse(proxyURLString)
		if err == nil {
			proxyFunc = http.ProxyURL(proxyURL) // Use custom proxy from configuration
		} else {
			l.Printf("[ERROR] Invalid proxy URL: %v", err)
		}
	}

	// Setup transport with custom proxy, if applicable, and check for HTTPS certificate check disabling
	transport := &http.Transport{
		Proxy:               proxyFunc,
		MaxIdleConnsPerHost: maxIdleConnections,
	}

	// Disable HTTPS check if the configuration requests it
	if cfg.Zscaler.Testing.DisableHttpsCheck {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true, // This disables HTTPS certificate validation
		}
		l.Printf("[INFO] HTTPS certificate validation is disabled (testing mode).")
	}

	retryableClient.HTTPClient.Transport = transport
	return retryableClient.StandardClient()
}

// getRetryAfter checks for the Retry-After header or response body to determine retry wait time.
func getRetryAfter(resp *http.Response, l logger.Logger) time.Duration {
	retryAfterHeader := resp.Header.Get("Retry-After")
	// x-ratelimit-reset: The time (in seconds) remaining in the current window after which the rate limit resets.
	ratelimitReset := resp.Header.Get("x-ratelimit-reset")
	// Date header from the response
	dateHeader := resp.Header.Get("Date")

	if retryAfterHeader != "" {
		// Try to parse the Retry-After value as an integer (seconds)
		if sleep, err := strconv.ParseInt(retryAfterHeader, 10, 64); err == nil {
			l.Printf("[INFO] got Retry-After from header: %s\n", retryAfterHeader)
			return time.Second * time.Duration(sleep+1) // Add 1 second padding
		} else {
			// Fallback: try parsing it as a duration (like "13s" from ZPA)
			dur, err := time.ParseDuration(retryAfterHeader)
			if err == nil {
				l.Printf("[INFO] got Retry-After duration from header: %s\n", retryAfterHeader)
				return dur + time.Second // Add 1 second padding
			}
			l.Printf("[INFO] error parsing Retry-After header: %v\n", err)
		}
	} else if ratelimitReset != "" {
		// Handle x-ratelimit-reset and Date headers
		resetSeconds, err := strconv.Atoi(ratelimitReset)
		if err == nil {
			// Use the Date header for more accurate timing
			requestDate, err := time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", dateHeader)
			if err != nil {
				// Log the error but proceed without using the Date header if parsing fails
				l.Printf("[INFO] error parsing Date header, defaulting to x-ratelimit-reset: %v\n", err)
				return time.Duration(resetSeconds+1) * time.Second // Fallback to just x-ratelimit-reset + 1 second padding
			}

			// Calculate the remaining time until rate limit reset using the Date header
			currentUnix := requestDate.Unix()
			waitTime := int64(resetSeconds) - currentUnix + 1 // Adding 1 second padding
			l.Printf("[INFO] x-ratelimit-reset header used, wait time calculated using Date header: %d seconds\n", waitTime)
			return time.Duration(waitTime) * time.Second
		}
		l.Printf("[INFO] error parsing x-ratelimit-reset header: %v\n", err)
	}

	// Fallback to default wait time if no Retry-After or x-ratelimit-reset headers exist
	return time.Second * time.Duration(RetryWaitMinSeconds)
}

// getRetryOnStatusCodes return a list of http status codes we want to apply retry on.
// Return empty slice to enable retry on all connection & server errors.
// Or return []int{429}  to retry on only TooManyRequests error.
func getRetryOnStatusCodes() []int {
	return []int{http.StatusTooManyRequests}
}

// checkRetry defines the retry logic based on status codes or response body errors.
func checkRetry(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	if resp != nil && containsInt(getRetryOnStatusCodes(), resp.StatusCode) {
		return true, nil
	}
	return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
}

func (c *Client) buildRequest(ctx context.Context, method, endpoint string, body io.Reader, urlParams url.Values, contentType string) (*http.Request, error) {

	if contentType == "" {
		contentType = contentTypeJSON
	}

	// Initialize urlParams if it's nil to prevent panic when calling urlParams.Set()
	if urlParams == nil {
		urlParams = make(url.Values)
	}

	isSandboxRequest := strings.Contains(endpoint, "/zscsb")
	isZPARequest := strings.Contains(endpoint, "/zpa")
	isZCCRequest := strings.Contains(endpoint, "/zcc")

	// Build the full URL for Sandbox, ZPA, ZCC, or OAuth2-based requests
	fullURL := ""
	baseUrl := ""

	if isSandboxRequest {
		baseUrl = c.GetSandboxURL()
	} else {
		baseUrl = GetAPIBaseURL(c.oauth2Credentials.Zscaler.Client.Cloud)
	}
	if isSandboxRequest {
		fullURL = fmt.Sprintf("%s%s", c.GetSandboxURL(), endpoint)
		urlParams.Set("api_token", c.GetSandboxToken()) // Append Sandbox token
	} else if isZPARequest {
		// Only append customerId to query parameters if it's not already in the URL path
		if !strings.Contains(endpoint, fmt.Sprintf("/customers/%s", c.oauth2Credentials.Zscaler.Client.CustomerID)) && c.oauth2Credentials.Zscaler.Client.CustomerID != "" {
			urlParams.Set("customerId", c.oauth2Credentials.Zscaler.Client.CustomerID)
		}
		fullURL = fmt.Sprintf("%s%s", baseUrl, endpoint)
	} else if isZCCRequest {
		fullURL = fmt.Sprintf("%s%s", baseUrl, endpoint)
	} else {
		fullURL = fmt.Sprintf("%s%s", baseUrl, endpoint)
	}

	// Add URL parameters to the endpoint
	params := ""
	if urlParams != nil {
		params = urlParams.Encode()
	}
	if strings.Contains(endpoint, "?") && params != "" {
		fullURL += "&" + params
	} else if params != "" {
		fullURL += "?" + params
	}

	// Create the HTTP request with context
	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	if c.oauth2Credentials.UserAgent != "" {
		req.Header.Add("User-Agent", c.oauth2Credentials.UserAgent)
	}

	// For non-sandbox requests, handle OAuth2 authentication
	if !isSandboxRequest {
		err = c.authenticate()
		if err != nil {
			return nil, err
		}
		// Extract token from context if available
		if token, ok := ctx.Value(ContextAccessToken).(string); ok && token != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		} else {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.oauth2Credentials.Zscaler.Client.AuthToken.AccessToken))
		}
	}

	return req, nil
}

func (c *Client) ExecuteRequest(ctx context.Context, method, endpoint string, body io.Reader, urlParams url.Values, contentType string) ([]byte, *http.Request, error) {
	req, err := c.buildRequest(ctx, method, endpoint, body, urlParams, contentType)
	if err != nil {
		return nil, nil, err
	}

	isSandboxRequest := strings.Contains(endpoint, "/zscsb")

	// Create cache key using the actual request
	key := cache.CreateCacheKey(req)
	if c.oauth2Credentials.Zscaler.Client.Cache.Enabled && !isSandboxRequest {
		if method != http.MethodGet {
			c.oauth2Credentials.CacheManager.Delete(key)
			c.oauth2Credentials.CacheManager.ClearAllKeysWithPrefix(strings.Split(key, "?")[0])
		}
		resp := c.oauth2Credentials.CacheManager.Get(key)
		inCache := resp != nil
		if inCache {
			respData, err := io.ReadAll(resp.Body)
			if err == nil {
				resp.Body = io.NopCloser(bytes.NewBuffer(respData))
			}
			c.oauth2Credentials.Logger.Printf("[INFO] served from cache, key:%s\n", key)
			return respData, req, nil
		}
	}

	// Execute the request with retries
	var resp *http.Response
	for retry := 1; retry <= 5; retry++ {
		start := time.Now()
		reqID := uuid.New().String()
		logger.LogRequest(c.oauth2Credentials.Logger, req, reqID, nil, !isSandboxRequest)
		httpClient := c.getServiceHTTPClient(endpoint)
		resp, err = httpClient.Do(req)
		logger.LogResponse(c.oauth2Credentials.Logger, resp, start, reqID)
		if err != nil {
			return nil, nil, err
		}
		if !isSandboxRequest && (resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden) {
			err = c.authenticate()
			if err != nil {
				return nil, nil, err
			}
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.oauth2Credentials.Zscaler.Client.AuthToken.AccessToken))
		} else if resp.StatusCode > 299 {
			resp.Body.Close()
			return nil, nil, checkErrorInResponse(resp, fmt.Errorf("api responded with code: %d", resp.StatusCode))
		} else if resp.StatusCode < 300 {
			break
		}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	// Cache logic for successful GET requests
	if !isSandboxRequest && c.oauth2Credentials.Zscaler.Client.Cache.Enabled && method == http.MethodGet {
		resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		c.oauth2Credentials.Logger.Printf("[INFO] saving to cache, key:%s\n", key)
		c.oauth2Credentials.CacheManager.Set(key, cache.CopyResponse(resp))
	}
	_ = tryDrainBody(resp.Body)
	return bodyBytes, req, nil
}

func tryDrainBody(body io.ReadCloser) error {
	defer body.Close()
	_, err := io.Copy(io.Discard, io.LimitReader(body, 4096))
	return err
}

// GetSandboxURL retrieves the sandbox URL for the ZIA service.
func (c *Client) GetSandboxURL() string {
	return "https://csbapi." + c.oauth2Credentials.Zscaler.Client.SandboxCloud + ".net"
}

// GetSandboxToken retrieves the sandbox token from the configuration or environment.
func (c *Client) GetSandboxToken() string {
	// Check if oauth2Credentials or the relevant fields are nil
	if c.oauth2Credentials == nil || c.oauth2Credentials.Zscaler.Client.SandboxToken == "" {
		// Fallback to environment variable if not set in the configuration
		return os.Getenv("ZSCALER_SANDBOX_TOKEN")
	}
	// Return the token from the configuration
	return c.oauth2Credentials.Zscaler.Client.SandboxToken
}

func (c *Client) authValid() bool {
	return c.oauth2Credentials.Zscaler.Client.AuthToken != nil && c.oauth2Credentials.Zscaler.Client.AuthToken.AccessToken != "" && !utils.IsTokenExpired(c.oauth2Credentials.Zscaler.Client.AuthToken.AccessToken)
}

// Unified authentication function to refresh OAuth2 tokens
func (c *Client) authenticate() error {
	c.Lock()
	defer c.Unlock()

	// Check if the AuthToken is nil, empty, or expired
	if !c.authValid() {
		// Pass the context from the Configuration along with the other arguments
		authToken, err := Authenticate(c.oauth2Credentials.Context, c.oauth2Credentials, c.oauth2Credentials.Logger)
		if err != nil {
			return err
		}
		c.oauth2Credentials.Zscaler.Client.AuthToken = authToken
		return nil
	}
	return nil
}

func containsInt(codes []int, code int) bool {
	for _, a := range codes {
		if a == code {
			return true
		}
	}
	return false
}
