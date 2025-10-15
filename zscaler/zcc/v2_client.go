package zcc

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
	rl "github.com/zscaler/zscaler-sdk-go/v3/ratelimiter"
	"github.com/zscaler/zscaler-sdk-go/v3/utils"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
)

type Client struct {
	sync.Mutex
	Config *Configuration
}

func NewClient(config *Configuration) (*Client, error) {
	if config == nil {
		return nil, errors.New("configuration cannot be nil")
	}

	// Enable Debug logging if the Debug flag is set
	if config.Debug {
		_ = os.Setenv("ZSCALER_SDK_LOG", "true")
		_ = os.Setenv("ZSCALER_SDK_VERBOSE", "true")
		config.Logger = logger.GetDefaultLogger("zcc-logger: ")
	}

	// Ensure HTTP clients are properly initialized
	if config.HTTPClient == nil {
		config.HTTPClient = getHTTPClient(config.Logger, nil, config)
	}

	// Authenticate the client using the configuration
	authToken, err := Authenticate(config.Context, config, config.Logger)
	if err != nil {
		config.Logger.Printf("[ERROR] Failed to authenticate client: %v\n", err)
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}

	config.ZCC.Client.AuthToken = authToken

	client := &Client{
		Config: config,
	}

	return client, nil
}

func (cfg *Configuration) SetBackoffConfig(backoffConf *BackoffConfig) {
	cfg.ZCC.Client.RateLimit.BackoffConf = backoffConf
}

// getHTTPClient sets up the retryable HTTP client with backoff and retry policies.
func getHTTPClient(l logger.Logger, rateLimiter *rl.RateLimiter, cfg *Configuration) *http.Client {
	if cfg.HTTPClient == nil {
		if cfg.ZCC.Client.RateLimit.BackoffConf != nil && cfg.ZCC.Client.RateLimit.BackoffConf.Enabled {
			retryableClient := retryablehttp.NewClient()
			retryableClient.RetryWaitMin = cfg.ZCC.Client.RateLimit.RetryWaitMin
			retryableClient.RetryWaitMax = cfg.ZCC.Client.RateLimit.RetryWaitMax
			retryableClient.RetryMax = cfg.ZCC.Client.RateLimit.BackoffConf.MaxNumOfRetries
			retryableClient.Logger = cfg.Logger
			retryableClient.HTTPClient.Transport = logging.NewSubsystemLoggingHTTPTransport("gozscaler", retryableClient.HTTPClient.Transport)
			retryableClient.CheckRetry = checkRetry
			retryableClient.Logger = l
			retryableClient.HTTPClient.Timeout = cfg.ZCC.Client.RequestTimeout

			retryableClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
				if resp != nil {
					if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
						endpoint := ""
						if resp.Request != nil {
							endpoint = resp.Request.URL.Path // Extract endpoint from the request URL
						}
						// Call getRetryAfter with endpoint for ZCC's endpoint-specific retry logic
						retryAfter := getRetryAfter(endpoint, l)
						if retryAfter > 0 {
							return retryAfter
						}
					}
				}
				// Use exponential backoff for all retries
				// API's own rate limiting handles rate limits
				mult := math.Pow(2, float64(attemptNum)) * float64(min)
				sleep := time.Duration(mult)
				if float64(sleep) != mult || sleep > max {
					sleep = max
				}
				return sleep
			}

			// Configure proxy settings from configuration
			proxyFunc := http.ProxyFromEnvironment // Default behavior (uses system/env variables)
			if cfg.ZCC.Client.Proxy.Host != "" {
				// Include username and password if provided
				proxyURLString := fmt.Sprintf("http://%s:%d", cfg.ZCC.Client.Proxy.Host, cfg.ZCC.Client.Proxy.Port)
				if cfg.ZCC.Client.Proxy.Username != "" && cfg.ZCC.Client.Proxy.Password != "" {
					// URL-encode the username and password
					proxyAuth := url.UserPassword(cfg.ZCC.Client.Proxy.Username, cfg.ZCC.Client.Proxy.Password)
					proxyURLString = fmt.Sprintf("http://%s@%s:%d", proxyAuth.String(), cfg.ZCC.Client.Proxy.Host, cfg.ZCC.Client.Proxy.Port)
				}

				proxyURL, err := url.Parse(proxyURLString)
				if err == nil {
					proxyFunc = http.ProxyURL(proxyURL) // Use custom proxy from configuration
				} else {
					cfg.Logger.Printf("[ERROR] Invalid proxy URL: %v", err)
				}
			}

			// Setup transport with custom proxy, if applicable, and check for HTTPS certificate check disabling
			transport := &http.Transport{
				Proxy:               proxyFunc,
				MaxIdleConnsPerHost: maxIdleConnections,
			}

			// Disable HTTPS check if the configuration requests it
			if cfg.ZCC.Testing.DisableHttpsCheck {
				transport.TLSClientConfig = &tls.Config{
					InsecureSkipVerify: true, // This disables HTTPS certificate validation
				}
				cfg.Logger.Printf("[INFO] HTTPS certificate validation is disabled (testing mode).")
			}

			retryableClient.HTTPClient.Transport = transport
			cfg.HTTPClient = retryableClient.StandardClient()
		}
	}
	return cfg.HTTPClient
}

func containsInt(codes []int, code int) bool {
	for _, a := range codes {
		if a == code {
			return true
		}
	}
	return false
}

func getRetryAfter(endpoint string, l logger.Logger) time.Duration {
	// Handle specific rate limits for `/downloadDevices` endpoint
	if strings.Contains(endpoint, "/downloadDevices") {
		l.Printf("[INFO] Reached rate limit for /downloadDevices. Retrying after 24 hours.")
		return 24 * time.Hour
	}

	// Default rate limit for all other endpoints: 100 calls per hour
	l.Printf("[INFO] General rate limit reached for endpoint: %s. Retrying after 1 hour.", endpoint)
	return time.Hour
}

// getRetryOnStatusCodes return a list of http status codes we want to apply retry on.
// Return empty slice to enable retry on all connection & server errors.
// Or return []int{429}  to retry on only TooManyRequests error.
func getRetryOnStatusCodes() []int {
	return []int{http.StatusTooManyRequests}
}

// Used to make http client retry on provided list of response status codes.
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

func Authenticate(ctx context.Context, cfg *Configuration, logger logger.Logger) (*AuthToken, error) {
	cfg.Lock()
	defer cfg.Unlock()

	logger.Printf("[DEBUG] Authenticating client: clientID=%s", cfg.ZCC.Client.ZCCClientID)

	// Check if a valid token is available
	if cfg.ZCC.Client.AuthToken == nil || cfg.ZCC.Client.AuthToken.AccessToken == "" || utils.IsTokenExpired(cfg.ZCC.Client.AuthToken.AccessToken) {
		logger.Printf("[DEBUG] No valid auth token found. Starting authentication process.")

		// Validate client credentials
		if cfg.ZCC.Client.ZCCClientID == "" || cfg.ZCC.Client.ZCCClientSecret == "" {
			logger.Printf("[ERROR] No client credentials provided. Set %s and %s environment variables or use ConfigSetter.", ZCC_CLIENT_ID, ZCC_CLIENT_SECRET)
			return nil, fmt.Errorf("missing client credentials: %s and/or %s", ZCC_CLIENT_ID, ZCC_CLIENT_SECRET)
		}

		// Prepare the authentication request
		authReq := AuthRequest{
			APIKey:    cfg.ZCC.Client.ZCCClientID,
			SecretKey: cfg.ZCC.Client.ZCCClientSecret,
		}
		data, err := json.Marshal(authReq)
		if err != nil {
			logger.Printf("[ERROR] Failed to serialize authentication request: %v", err)
			return nil, fmt.Errorf("failed to serialize authentication request: %w", err)
		}

		// Create the HTTP request
		req, err := http.NewRequestWithContext(ctx, "POST", cfg.BaseURL.String()+"/auth/v1/login", bytes.NewBuffer(data))
		if err != nil {
			logger.Printf("[ERROR] Failed to create authentication request: %v", err)
			return nil, fmt.Errorf("failed to create authentication request: %w", err)
		}
		req.Header.Add("Content-Type", "application/json")
		if cfg.UserAgent != "" {
			req.Header.Add("User-Agent", cfg.UserAgent)
		}

		// Execute the request using the configured HTTP client
		resp, err := cfg.HTTPClient.Do(req)
		if err != nil {
			logger.Printf("[ERROR] Authentication request failed: %v", err)
			return nil, fmt.Errorf("authentication request failed: %w", err)
		}
		defer resp.Body.Close()

		// Read the response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Printf("[ERROR] Failed to read authentication response: %v", err)
			return nil, fmt.Errorf("failed to read authentication response: %w", err)
		}

		// Handle non-2xx responses
		if resp.StatusCode >= 300 {
			logger.Printf("[ERROR] Authentication failed with status %d. Response: %s", resp.StatusCode, respBody)
			return nil, fmt.Errorf("authentication failed: HTTP %d, response: %s", resp.StatusCode, respBody)
		}

		// Parse the authentication response
		var authToken AuthToken
		err = json.Unmarshal(respBody, &authToken)
		if err != nil {
			logger.Printf("[ERROR] Failed to parse authentication response: %v", err)
			return nil, fmt.Errorf("failed to parse authentication response: %w", err)
		}

		// Log successful authentication
		logger.Printf("[DEBUG] Authentication successful. Token received.")

		// Store the token in the configuration
		cfg.ZCC.Client.AuthToken = &authToken
		return &authToken, nil
	}

	// Return the existing valid token
	logger.Printf("[DEBUG] Using existing valid auth token.")
	return cfg.ZCC.Client.AuthToken, nil
}

func (client *Client) NewRequestDo(method, url string, options, body, v interface{}) (*http.Response, error) {
	client.Config.Logger.Printf("[DEBUG] Creating new request: method=%s, url=%s", method, url)
	return client.newRequestDoCustom(method, url, options, body, v, client.Config)
}

func (client *Client) newRequestDoCustom(method, urlStr string, options, body, v interface{}, config *Configuration) (*http.Response, error) {
	client.Config.Logger.Printf("[DEBUG] newRequestDoCustom called with method=%s, urlStr=%s", method, urlStr)

	// Authenticate and handle errors
	if _, err := Authenticate(context.Background(), config, config.Logger); err != nil {
		client.Config.Logger.Printf("[ERROR] Authentication failed: %v", err)
		return nil, err
	}

	// Create a new HTTP request
	req, err := client.newRequest(method, urlStr, options, body, client.Config)
	if err != nil {
		client.Config.Logger.Printf("[ERROR] Failed to create new request: %v", err)
		return nil, err
	}

	// Start tracking request execution time
	start := time.Now()
	reqID := uuid.NewString()

	// Log the request
	logger.LogRequest(client.Config.Logger, req, reqID, nil, true)

	// Execute the request
	resp, err := client.do(req, v, start, reqID)
	if err != nil {
		client.Config.Logger.Printf("[ERROR] Request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Handle unauthorized or forbidden responses
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		client.Config.Logger.Printf("[WARN] Unauthorized or forbidden response. Retrying authentication.")

		// Re-authenticate and handle errors
		if _, err := Authenticate(context.Background(), config, config.Logger); err != nil {
			client.Config.Logger.Printf("[ERROR] Re-authentication failed: %v", err)
			return nil, err
		}

		// Retry the original request
		resp, err = client.do(req, v, start, reqID)
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Request failed after re-authentication: %v", err)
			return nil, err
		}
		resp.Body.Close()
		return resp, nil
	}

	// Return the response
	return resp, err
}

// Generating the Http request.
func (client *Client) newRequest(method, urlPath string, options, body interface{}, config *Configuration) (*http.Request, error) {
	if config.ZCC.Client.AuthToken == nil || config.ZCC.Client.AuthToken.AccessToken == "" {
		client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s\n", ZCC_CLIENT_ID, config.ZCC.Client.ZCCClientID)
		return nil, fmt.Errorf("failed to signin the user %s=%s", ZCC_CLIENT_ID, config.ZCC.Client.ZCCClientID)
	}
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	// Join the path to the base-url
	u := *client.Config.BaseURL
	unescaped, err := url.PathUnescape(urlPath)
	if err != nil {
		return nil, err
	}

	// Set the encoded path data
	u.RawPath = u.Path + urlPath
	u.Path += unescaped

	// Set the query parameters
	if options != nil {
		q, err := query.Values(options)
		if err != nil {
			return nil, err
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("auth-token", config.ZCC.Client.AuthToken.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	if client.Config.UserAgent != "" {
		req.Header.Add("User-Agent", client.Config.UserAgent)
	}

	return req, nil
}

func (client *Client) do(req *http.Request, v interface{}, start time.Time, reqID string) (*http.Response, error) {
	// Initialize the HTTP client using the configuration's method
	httpClient := getHTTPClient(client.Config.Logger, nil, client.Config)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Check for errors in the HTTP response
	if err := errorx.CheckErrorInResponse(resp, err); err != nil {
		return resp, err
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}

	// Log the raw response body for debugging
	client.Config.Logger.Printf("[DEBUG] Raw response body: %s", string(bodyBytes))

	if v != nil {
		// Decode JSON from the raw response body
		if err := json.Unmarshal(bodyBytes, v); err != nil {
			client.Config.Logger.Printf("[ERROR] Failed to parse JSON response: %v", err)
			return resp, err
		}
	}

	// Log the response details
	logger.LogResponse(client.Config.Logger, resp, start, reqID)

	// Unescape HTML entities in the response
	unescapeHTML(v)

	return resp, nil
}

func unescapeHTML(entity interface{}) {
	if entity == nil {
		return
	}
	data, err := json.Marshal(entity)
	if err != nil {
		return
	}
	var mapData map[string]interface{}
	err = json.Unmarshal(data, &mapData)
	if err != nil {
		return
	}
	for _, field := range []string{"name", "description"} {
		if v, ok := mapData[field]; ok && v != nil {
			str, ok := v.(string)
			if ok {
				mapData[field] = html.UnescapeString(html.UnescapeString(str))
			}
		}
	}
	data, err = json.Marshal(mapData)
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, entity)
}
