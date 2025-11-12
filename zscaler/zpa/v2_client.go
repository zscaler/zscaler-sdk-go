package zpa

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
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zscaler/zscaler-sdk-go/v3/cache"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
	rl "github.com/zscaler/zscaler-sdk-go/v3/ratelimiter"
	"github.com/zscaler/zscaler-sdk-go/v3/utils"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
)

type Client struct {
	sync.Mutex
	cache  cache.Cache
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
		config.Logger = logger.GetDefaultLogger("zpa-logger: ")
	}

	// Ensure HTTP clients are properly initialized
	if config.HTTPClient == nil {
		// Initialize rate limiter if not already set
		if config.RateLimiter == nil {
			config.RateLimiter = rl.NewRateLimiter(20, 10, 10, 10)
		}
		config.HTTPClient = getHTTPClient(config.Logger, config.RateLimiter, config)
	}

	// Initialize cache if enabled
	if config.ZPA.Client.Cache.Enabled {
		if config.CacheManager == nil {
			config.CacheManager = newCache(config)
		}
	}

	// Authenticate the client using the configuration
	authToken, err := Authenticate(config.Context, config, config.Logger)
	if err != nil {
		config.Logger.Printf("[ERROR] Failed to authenticate client: %v\n", err)
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}

	config.ZPA.Client.AuthToken = authToken

	client := &Client{
		Config: config,
		cache:  config.CacheManager,
	}

	return client, nil
}

func (client *Client) GetLogger() logger.Logger {
	return client.Config.Logger
}

// getHTTPClient sets up the retryable HTTP client with backoff and retry policies.
func getHTTPClient(l logger.Logger, rateLimiter *rl.RateLimiter, cfg *Configuration) *http.Client {
	retryableClient := retryablehttp.NewClient()

	// Set the retry settings, allowing user to override defaults.
	// Defaults are set by the config, which is initially read from constants but can be overridden.
	retryableClient.RetryWaitMin = cfg.ZPA.Client.RateLimit.RetryWaitMin
	retryableClient.RetryWaitMax = cfg.ZPA.Client.RateLimit.RetryWaitMax

	if cfg.ZPA.Client.RateLimit.MaxRetries == 0 {
		retryableClient.RetryMax = math.MaxInt32
	} else {
		retryableClient.RetryMax = int(cfg.ZPA.Client.RateLimit.MaxRetries)
	}

	// Configure backoff and retry policies
	retryableClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		if resp != nil {
			if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
				retryAfter := getRetryAfter(resp, l)
				if retryAfter > 0 {
					return retryAfter
				}
			}
		}
		// Use exponential backoff for all retries
		// The API's own rate limiting (429 + Retry-After) handles rate limits
		// This prevents unnecessary delays from proactive rate limiting
		mult := math.Pow(2, float64(attemptNum)) * float64(min)
		sleep := time.Duration(mult)
		if float64(sleep) != mult || sleep > max {
			sleep = max
		}
		return sleep
	}
	retryableClient.CheckRetry = checkRetry

	retryableClient.Logger = l

	// Set the request timeout, allowing user-defined override.
	if cfg.ZPA.Client.RequestTimeout == 0 {
		retryableClient.HTTPClient.Timeout = time.Second * 240 // Default to 240 seconds to match old SDK.
	} else {
		retryableClient.HTTPClient.Timeout = cfg.ZPA.Client.RequestTimeout
	}

	// Configure proxy settings from configuration
	proxyFunc := http.ProxyFromEnvironment // Default behavior (uses system/env variables)
	if cfg.ZPA.Client.Proxy.Host != "" {
		// Include username and password if provided
		proxyURLString := fmt.Sprintf("http://%s:%d", cfg.ZPA.Client.Proxy.Host, cfg.ZPA.Client.Proxy.Port)
		if cfg.ZPA.Client.Proxy.Username != "" && cfg.ZPA.Client.Proxy.Password != "" {
			// URL-encode the username and password
			proxyAuth := url.UserPassword(cfg.ZPA.Client.Proxy.Username, cfg.ZPA.Client.Proxy.Password)
			proxyURLString = fmt.Sprintf("http://%s@%s:%d", proxyAuth.String(), cfg.ZPA.Client.Proxy.Host, cfg.ZPA.Client.Proxy.Port)
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
	if cfg.ZPA.Testing.DisableHttpsCheck {
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
	retryAfterHeader := resp.Header.Get("retry-after")

	if retryAfterHeader != "" {
		// Try to parse the Retry-After value as an integer (seconds)
		if sleep, err := strconv.ParseInt(retryAfterHeader, 10, 64); err == nil {
			l.Printf("[INFO] got retry-after from header: %s\n", retryAfterHeader)
			return time.Second * time.Duration(sleep+1) // Add 1 second padding
		} else {
			// Fallback: try parsing it as a duration (like "13s" from ZPA)
			dur, err := time.ParseDuration(retryAfterHeader)
			if err == nil {
				l.Printf("[INFO] got retry-after duration from header: %s\n", retryAfterHeader)
				return dur + time.Second // Add 1 second padding
			}
			l.Printf("[INFO] error parsing retry-after header: %v\n", err)
		}
	}

	// Fallback to default wait time if no Retry-After or x-ratelimit-reset headers exist
	return time.Second * time.Duration(RetryWaitMinSeconds)
}

// getRetryOnStatusCodes return a list of http status codes we want to apply retry on.
// Return empty slice to enable retry on all connection & server errors.
// Or return []int{429}  to retry on only TooManyRequests error.
func getRetryOnStatusCodes() []int {
	return []int{http.StatusTooManyRequests, http.StatusConflict}
}

// checkRetry defines the retry logic based on status codes or response body errors.
func checkRetry(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	if resp != nil && containsInt(getRetryOnStatusCodes(), resp.StatusCode) {
		if resp.StatusCode == http.StatusConflict {
			respMap := map[string]string{}
			data, err := io.ReadAll(resp.Body)
			resp.Body = io.NopCloser(bytes.NewBuffer(data))
			if err == nil {
				_ = json.Unmarshal(data, &respMap)
				if errorID, ok := respMap["id"]; ok && (errorID == "api.concurrent.access.error") {
					return true, nil
				}
			}
		}
		return true, nil
	}
	if resp != nil && resp.StatusCode == http.StatusBadRequest {
		respMap := map[string]string{}
		data, err := io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(data))
		if err == nil {
			_ = json.Unmarshal(data, &respMap)
			if errorID, ok := respMap["id"]; ok && (errorID == "non.restricted.entity.authorization.failed" || errorID == "bad.request") {
				return true, nil
			}
		}
		// Implemented to handle upstream restrictions on simultaneous requests when dealing with CRUD operations, related to ZPA Access policy rule order
		// ET-53585: https://jira.corp.ZPA.com/browse/ET-53585
		// ET-48860: https://confluence.corp.ZPA.com/display/ET/ET-48860+incorrect+rules+order
		if err == nil {
			_ = json.Unmarshal(data, &respMap)
			if errorID, ok := respMap["id"]; ok && (errorID == "db.simultaneous.request" || errorID == "bad.request") {
				return true, nil
			}
		}

		// ET-66174: https://jira.corp.ZPA.com/browse/ET-66174
		// DOC-51102: https://jira.corp.ZPA.com/browse/DOC-51102
		if err == nil {
			_ = json.Unmarshal(data, &respMap)
			if errorID, ok := respMap["id"]; ok && (errorID == "api.concurrent.access.error" || errorID == "bad.request") {
				return true, nil
			}
		}
	}
	return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
}

func containsInt(codes []int, code int) bool {
	for _, a := range codes {
		if a == code {
			return true
		}
	}
	return false
}

func Authenticate(ctx context.Context, cfg *Configuration, logger logger.Logger) (*AuthToken, error) {
	cfg.Lock()
	defer cfg.Unlock()

	// Reuse the token if it's still valid
	if cfg.ZPA.Client.AuthToken != nil && cfg.ZPA.Client.AuthToken.AccessToken != "" {
		if !utils.IsTokenExpired(cfg.ZPA.Client.AuthToken.AccessToken) {
			// logger.Printf("[DEBUG] Reusing existing valid authentication token")
			return cfg.ZPA.Client.AuthToken, nil
		}
		// logger.Printf("[DEBUG] Authentication token expired or nearing expiry, refreshing...")
	} else {
		// logger.Printf("[DEBUG] Authentication token not present, initiating new request...")
	}

	clientID := cfg.ZPA.Client.ZPAClientID
	clientSecret := cfg.ZPA.Client.ZPAClientSecret
	baseURL := cfg.BaseURL.String()

	if clientID == "" || clientSecret == "" {
		logger.Printf("[ERROR] Missing client credentials. Ensure ZPA_CLIENT_ID and ZPA_CLIENT_SECRET are set.")
		return nil, errors.New("missing client credentials")
	}

	authURL := fmt.Sprintf("%s/signin", baseURL)
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	// Use getHTTPClient to handle retries, rate-limiting, etc.
	httpClient := getHTTPClient(logger, cfg.RateLimiter, cfg)

	req, err := http.NewRequestWithContext(ctx, "POST", authURL, strings.NewReader(data.Encode()))
	if err != nil {
		logger.Printf("[ERROR] Failed to create authentication request: %v", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if cfg.UserAgent != "" {
		req.Header.Add("User-Agent", cfg.UserAgent)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Printf("[ERROR] Failed to authenticate: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		logger.Printf("[ERROR] Authentication failed with status: %d, response: %s", resp.StatusCode, string(respBody))
		return nil, fmt.Errorf("authentication failed with status: %d, response: %s", resp.StatusCode, string(respBody))
	}

	var token AuthToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		logger.Printf("[ERROR] Failed to decode authentication response: %v", err)
		return nil, err
	}

	// Store the new authentication token
	cfg.ZPA.Client.AuthToken = &token
	// logger.Printf("[DEBUG] Authentication successful. New token acquired.")

	return &token, nil
}

func (client *Client) NewRequestDo(method, url string, options, body, v interface{}) (*http.Response, error) {
	req, err := client.getRequest(method, url, options, body)
	if err != nil {
		return nil, err
	}
	key := cache.CreateCacheKey(req)

	// Check if caching is enabled in the configuration
	if client.Config.ZPA.Client.Cache.Enabled {
		if req.Method != http.MethodGet {
			// Remove resource from cache for non-GET requests
			client.cache.Delete(key)

			// Clear all cache entries with the same URL prefix to handle query param differences
			client.cache.ClearAllKeysWithPrefix(strings.Split(key, "?")[0])
		}

		// Check if response is in cache
		resp := client.cache.Get(key)
		inCache := resp != nil

		// Handle fresh cache logic if the option is enabled
		if inCache && client.Config.ZPA.Client.Cache.DefaultTtl > 0 {
			client.Config.Logger.Printf("[INFO] Cache entry is valid, key:%s\n", key)
		} else {
			client.cache.Delete(key) // Delete stale cache entries
			inCache = false
		}

		if inCache {
			if v != nil {
				respData, err := io.ReadAll(resp.Body)
				if err == nil {
					resp.Body = io.NopCloser(bytes.NewBuffer(respData))
				}
				if err := decodeJSON(respData, v); err != nil {
					return resp, err
				}
			}
			unescapeHTML(v)
			client.Config.Logger.Printf("[INFO] served from cache, key:%s\n", key)
			return resp, nil
		}
	}

	// Make the actual request if not in cache
	resp, err := client.newRequestDoCustom(method, url, options, body, v)
	if err != nil {
		return resp, err
	}

	// Save the response to cache if caching is enabled and the response is cacheable
	if client.Config.ZPA.Client.Cache.Enabled && resp.StatusCode >= 200 && resp.StatusCode <= 299 && req.Method == http.MethodGet && v != nil && reflect.TypeOf(v).Kind() != reflect.Slice {
		d, err := json.Marshal(v)
		if err == nil {
			resp.Body = io.NopCloser(bytes.NewReader(d))
			client.Config.Logger.Printf("[INFO] saving to cache, key:%s\n", key)
			client.cache.Set(key, cache.CopyResponse(resp))
		} else {
			client.Config.Logger.Printf("[ERROR] saving to cache error:%s, key:%s\n", err, key)
		}
	}
	return resp, nil
}

// UnmarshalJSON ensures proper unmarshaling of the AuthToken struct.
func (t *AuthToken) UnmarshalJSON(data []byte) error {
	// Create an alias to avoid infinite recursion
	type Alias AuthToken
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	// Unmarshal into the auxiliary struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Convert RawExpiresIn to int
	if t.RawExpiresIn != "" {
		expiresIn, err := strconv.Atoi(t.RawExpiresIn)
		if err != nil {
			return fmt.Errorf("invalid expires_in value: %v", err)
		}
		t.ExpiresIn = expiresIn
		t.Expiry = time.Now().Add(time.Duration(expiresIn) * time.Second)
	}

	return nil
}

func (client *Client) getRequest(method, urlPath string, options, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	// Parse the URL path to separate the path from any query string
	parsedPath, err := url.Parse(urlPath)
	if err != nil {
		return nil, err
	}

	// Join the parsed path with the base URL
	u := client.Config.BaseURL.ResolveReference(parsedPath)

	// Handle query parameters from options and any additional logic
	if options == nil {
		options = struct{}{}
	}
	q, err := query.Values(options)
	if err != nil {
		return nil, err
	}
	// Here, injectMicrotenantID or any similar function should ensure
	// it's not duplicating query parameters that may already be present in urlPath
	q = common.InjectMicrotentantID(body, q, "")

	// Merge query params from urlPath and options. Avoid overwriting any existing params.
	for key, values := range parsedPath.Query() {
		for _, value := range values {
			q.Add(key, value)
		}
	}

	// Encode the final query, which by default uses '+' for spaces
	encodedQuery := q.Encode()

	// **Here** is the single place we convert '+' to '%20':
	encodedQuery = strings.ReplaceAll(encodedQuery, "+", "%20")

	u.RawQuery = encodedQuery

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (client *Client) newRequestDoCustom(method, urlStr string, options, body, v interface{}) (*http.Response, error) {
	// Authenticate the client using the global Authenticate function
	_, err := Authenticate(client.Config.Context, client.Config, client.Config.Logger)
	if err != nil {
		return nil, err
	}

	// Create a new HTTP request
	req, err := client.newRequest(method, urlStr, options, body)
	if err != nil {
		return nil, err
	}

	reqID := uuid.NewString()
	start := time.Now()

	// Log the request
	logger.LogRequest(client.Config.Logger, req, reqID, nil, true)

	// Execute the HTTP request
	resp, err := client.do(req, v, start, reqID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle unauthorized or forbidden responses by re-authenticating
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		_, err := Authenticate(client.Config.Context, client.Config, client.Config.Logger)
		if err != nil {
			return nil, err
		}

		// Retry the request after re-authentication
		resp, err := client.do(req, v, start, reqID)
		if err != nil {
			return nil, err
		}
		resp.Body.Close()
		return resp, nil
	}

	return resp, err
}

// Generating the Http request
// Generating the HTTP request
func (client *Client) newRequest(method, urlPath string, options, body interface{}) (*http.Request, error) {
	if client.Config.ZPA.Client.AuthToken == nil || client.Config.ZPA.Client.AuthToken.AccessToken == "" {
		client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s\n", ZPA_CLIENT_ID, client.Config.ZPA.Client.ZPAClientID)
		return nil, fmt.Errorf("failed to signin the user %s=%s", ZPA_CLIENT_ID, client.Config.ZPA.Client.ZPAClientID)
	}

	req, err := client.getRequest(method, urlPath, options, body)
	if err != nil {
		return nil, err
	}

	// Add Authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.Config.ZPA.Client.AuthToken.AccessToken))
	req.Header.Add("Content-Type", "application/json")

	// Add User-Agent header if specified
	if client.Config.UserAgent != "" {
		req.Header.Add("User-Agent", client.Config.UserAgent)
	}

	// Add x-partner-id header if partnerId is provided in config
	if client.Config.ZPA.Client.PartnerID != "" {
		req.Header.Set("x-partner-id", client.Config.ZPA.Client.PartnerID)
	}

	return req, nil
}

func (client *Client) do(req *http.Request, v interface{}, start time.Time, reqID string) (*http.Response, error) {
	// Use the appropriate HTTP client
	httpClient := client.Config.HTTPClient
	if httpClient == nil {
		httpClient = client.Config.HTTPClient // Fallback to default HTTP client
	}

	// Execute the HTTP request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Read and buffer the response body
	respData, err := io.ReadAll(resp.Body)
	if err == nil {
		resp.Body = io.NopCloser(bytes.NewBuffer(respData))
	}

	// Check for errors in the response
	if err := errorx.CheckErrorInResponse(resp, err); err != nil {
		return resp, err
	}

	// If 204 No Content, skip unmarshalling
	if resp.StatusCode == http.StatusNoContent || len(respData) == 0 {
		logger.LogResponse(client.Config.Logger, resp, start, reqID)
		return resp, nil
	}

	// Decode the response into the provided variable if applicable
	if v != nil {
		if err := decodeJSON(respData, v); err != nil {
			return resp, err
		}
	}

	// Log the response details
	logger.LogResponse(client.Config.Logger, resp, start, reqID)

	// Unescape any HTML content in the response
	unescapeHTML(v)

	return resp, nil
}

func decodeJSON(respData []byte, v interface{}) error {
	return json.NewDecoder(bytes.NewBuffer(respData)).Decode(&v)
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
