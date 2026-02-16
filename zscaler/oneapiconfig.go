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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
)

const (
	maxIdleConnections  int = 40
	requestTimeout      int = 60
	contentTypeJSON         = "application/json"
	MaxNumOfRetries         = 100
	RetryWaitMaxSeconds     = 10
	RetryWaitMinSeconds     = 2
	loggerPrefix            = "oneapi-logger: "
)

// inFlightRequest represents a request that's currently being processed
type inFlightRequest struct {
	wg     sync.WaitGroup
	result []byte
	resp   *http.Response
	req    *http.Request
	err    error
	once   sync.Once
}

// Client defines the ZIA client structure.
type Client struct {
	sync.Mutex
	oauth2Credentials *Configuration
	stopTicker        chan bool
	inFlightRequests  sync.Map // Map[string]*inFlightRequest - tracks in-flight GET requests for deduplication
}

// NewOneAPIClient creates a new client using OAuth2 authentication for any service.
func NewOneAPIClient(config *Configuration) (*Service, error) {
	cli := &Client{
		oauth2Credentials: config,
		stopTicker:        make(chan bool),
	}

	if !config.UseLegacyClient {
		// Authenticate and start token renewal
		if err := cli.authenticate(); err != nil {
			return nil, fmt.Errorf("initial authentication failed: %w", err)
		}
		cli.startTokenRenewalTicker()
	}

	return NewService(cli, nil), nil
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
					// Refresh the token with proper locking to prevent race conditions
					c.Lock()
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
					c.Unlock()
				case <-c.stopTicker:
					ticker.Stop()
					return
				}
			}
		}()
	}
}

// Close stops the token renewal ticker and cleans up resources.
func (c *Client) Close() {
	c.Lock()
	defer c.Unlock()

	if c.stopTicker != nil {
		close(c.stopTicker)
		c.stopTicker = nil
	}
}

func (client *Client) GetLogger() logger.Logger {
	return client.oauth2Credentials.Logger
}

// getHTTPClient sets up the retryable HTTP client with backoff and retry policies.
func getHTTPClient(l logger.Logger, rateLimiter *rl.RateLimiter, cfg *Configuration) *http.Client {
	retryableClient := retryablehttp.NewClient()

	// Set the retry settings, allowing user to override defaults.
	// Defaults are set by the config, which is initially read from constants but can be overridden.
	retryableClient.RetryWaitMin = cfg.Zscaler.Client.RateLimit.RetryWaitMin
	retryableClient.RetryWaitMax = cfg.Zscaler.Client.RateLimit.RetryWaitMax

	// Set RetryMax from configuration - this controls retries at the HTTP client level
	// Both retryablehttp AND ExecuteRequest will handle retries for maximum reliability
	if cfg.Zscaler.Client.RateLimit.MaxRetries > 0 {
		retryableClient.RetryMax = int(cfg.Zscaler.Client.RateLimit.MaxRetries)
	} else {
		retryableClient.RetryMax = MaxNumOfRetries // Default to 100
	}

	// Configure backoff and retry policies for retryablehttp
	retryableClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		if resp != nil {
			if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusUnauthorized {
				retryAfter := getRetryAfter(resp, cfg)
				if retryAfter > 0 {
					return retryAfter
				}
			}
		}
		// Use exponential backoff for all retries
		// Rate limiting is handled by API's 429 errors + ExecuteRequest logic
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
	if cfg.Zscaler.Client.RequestTimeout == 0 {
		retryableClient.HTTPClient.Timeout = time.Second * 60 // Default to 60 seconds if not specified.
	} else {
		retryableClient.HTTPClient.Timeout = cfg.Zscaler.Client.RequestTimeout
	}

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

	// Wrap the transport with rate limiting
	if rateLimiter != nil {
		rateLimitedTransport := &rl.RateLimitTransport{
			Base:            transport,
			Limiter:         rateLimiter,
			Logger:          l,
			AdditionalDelay: 0,
		}
		retryableClient.HTTPClient.Transport = rateLimitedTransport
	} else {
		retryableClient.HTTPClient.Transport = transport
	}

	return retryableClient.StandardClient()
}

func getRetryAfter(resp *http.Response, cfg *Configuration) time.Duration {
	l := cfg.Logger
	retryAfterHeader := resp.Header.Get("Retry-After")
	if retryAfterHeader == "" {
		retryAfterHeader = resp.Header.Get("retry-after")
	}
	ratelimitReset := resp.Header.Get("X-Ratelimit-Reset")
	ratelimitRemaining := resp.Header.Get("X-Ratelimit-Remaining")
	ratelimitLimit := resp.Header.Get("X-Ratelimit-Limit")

	// Use config-defined threshold, or fallback to 2
	threshold := cfg.Zscaler.Client.RateLimit.RetryRemainingThreshold
	if threshold == 0 {
		threshold = 2
	}

	// Log everything if debugging or if we're approaching limits
	if cfg.Debug || ratelimitRemaining != "" {
		l.Printf("[DEBUG] Rate limit headers: Limit=%s, Remaining=%s, Reset=%s, Retry-After=%s, Status=%d",
			ratelimitLimit, ratelimitRemaining, ratelimitReset, retryAfterHeader, resp.StatusCode)
	}

	// Check if we're hitting hourly limits (indicated by very long retry times)
	if retryAfterHeader != "" {
		if sleep, err := strconv.ParseInt(retryAfterHeader, 10, 64); err == nil {
			duration := time.Second * time.Duration(sleep)
			if duration > 5*time.Minute {
				l.Printf("[WARN] Long Retry-After detected: %v (%ds). This likely indicates hourly rate limit reached.", duration, sleep)
			} else {
				l.Printf("[INFO] Retry-After: %v (%ds)", duration, sleep)
			}
			// Trust the API's Retry-After value - it's authoritative
			// Only add buffer for very short delays to account for processing time
			if duration < time.Second {
				return duration + 500*time.Millisecond // Add 500ms buffer for sub-second delays
			}
			return duration
		}
		if dur, err := time.ParseDuration(retryAfterHeader); err == nil {
			if dur > 5*time.Minute {
				l.Printf("[WARN] Long Retry-After detected: %v. This likely indicates hourly rate limit reached.", dur)
			} else {
				l.Printf("[INFO] Retry-After (duration): %v", dur)
			}
			// Trust the API's Retry-After value
			if dur < time.Second {
				return dur + 500*time.Millisecond
			}
			return dur
		}
		l.Printf("[WARN] Could not parse Retry-After header: %s", retryAfterHeader)
	}

	// Check remaining requests before hitting limit
	if ratelimitRemaining != "" {
		if remaining, err := strconv.Atoi(ratelimitRemaining); err == nil && remaining < int(threshold) {
			l.Printf("[WARN] Approaching rate limit: %d requests remaining (threshold: %d)", remaining, threshold)
			if ratelimitReset != "" {
				if resetSecs, err := strconv.Atoi(ratelimitReset); err == nil {
					l.Printf("[INFO] Rate limit resets in %ds, waiting that duration", resetSecs)
					return time.Duration(resetSecs) * time.Second
				}
			}
			l.Printf("[INFO] Approaching rate limit with no reset header; using fallback delay: %ds", RetryWaitMinSeconds)
			return time.Second * time.Duration(RetryWaitMinSeconds)
		}
	}

	// Reset-based retry when no Retry-After header
	if ratelimitReset != "" {
		if resetSecs, err := strconv.Atoi(ratelimitReset); err == nil {
			l.Printf("[INFO] Using X-Ratelimit-Reset: %ds", resetSecs)
			return time.Duration(resetSecs) * time.Second
		}
	}

	// Final fallback - use configured retry wait
	l.Printf("[INFO] No rate limit headers found; using fallback wait: %ds", RetryWaitMinSeconds)
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
	// Retry 409 Conflict and 412 Precondition Failed when body indicates edit lock / org barrier
	// (e.g. "Failed during enter Org barrier") - same logic as ZIA/ZTW legacy clients
	if resp != nil && (resp.StatusCode == http.StatusPreconditionFailed || resp.StatusCode == http.StatusConflict) {
		if errorx.IsEditLockError(resp) {
			return true, nil
		}
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
	isZIdentityRequest := strings.Contains(endpoint, "/admin/api/v1")

	// Build the full URL for Sandbox, ZPA, ZCC, ZIdentity, or OAuth2-based requests
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
	} else if isZIdentityRequest {
		// For zidentity endpoints, construct URL using vanity domain directly
		// Format: https://{vanity_domain}-admin.zslogin{cloud}.net/admin/api/v1/...
		vanityDomain := c.oauth2Credentials.Zscaler.Client.VanityDomain
		cloud := c.oauth2Credentials.Zscaler.Client.Cloud

		var zidentityBaseURL string
		if cloud == "" || strings.EqualFold(cloud, "PRODUCTION") {
			zidentityBaseURL = fmt.Sprintf("https://%s-admin.zslogin.net", vanityDomain)
		} else {
			zidentityBaseURL = fmt.Sprintf("https://%s-admin.zslogin%s.net", vanityDomain, strings.ToLower(cloud))
		}
		fullURL = fmt.Sprintf("%s%s", zidentityBaseURL, endpoint)
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

	// Add x-partner-id header if partnerId is provided in config
	if c.oauth2Credentials.Zscaler.Client.PartnerID != "" {
		req.Header.Set("x-partner-id", c.oauth2Credentials.Zscaler.Client.PartnerID)
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
			if c.oauth2Credentials.Debug {
				c.oauth2Credentials.Logger.Printf("[DEBUG] Using Authorization header from context: Bearer %s...", token[:min(len(token), 20)])
			}
		} else {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.oauth2Credentials.Zscaler.Client.AuthToken.AccessToken))
			if c.oauth2Credentials.Debug {
				token := c.oauth2Credentials.Zscaler.Client.AuthToken.AccessToken
				c.oauth2Credentials.Logger.Printf("[DEBUG] Using Authorization header from AuthToken: Bearer %s...", token[:min(len(token), 20)])
			}
		}
	}

	return req, nil
}

func (c *Client) ExecuteRequest(ctx context.Context, method, endpoint string, body io.Reader, urlParams url.Values, contentType string) ([]byte, *http.Response, *http.Request, error) {
	// Buffer the request body so we can retry with SESSION_NOT_VALID errors
	var requestBodyBytes []byte
	if body != nil {
		var err error
		requestBodyBytes, err = io.ReadAll(body)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to read request body: %w", err)
		}
	}

	req, err := c.buildRequest(ctx, method, endpoint, bytes.NewReader(requestBodyBytes), urlParams, contentType)
	if err != nil {
		return nil, nil, nil, err
	}

	isSandboxRequest := strings.Contains(endpoint, "/zscsb")
	overallStartTime := time.Now()
	totalWaitTime := time.Duration(0) // Track time spent waiting for rate limits

	// Create cache key using the actual request
	key := cache.CreateCacheKey(req)

	// For GET requests, check cache first and handle request deduplication
	var ifr *inFlightRequest
	if method == http.MethodGet && c.oauth2Credentials.Zscaler.Client.Cache.Enabled && !isSandboxRequest {
		// Check cache first
		cachedResp := c.oauth2Credentials.CacheManager.Get(key)
		if cachedResp != nil {
			respData, err := io.ReadAll(cachedResp.Body)
			if err == nil {
				cachedResp.Body = io.NopCloser(bytes.NewBuffer(respData))
			}
			c.oauth2Credentials.Logger.Printf("[INFO] served from cache, key:%s\n", key)
			return respData, cachedResp, req, nil
		}

		// Check if there's an in-flight request for the same URL
		if inFlight, exists := c.inFlightRequests.Load(key); exists {
			if existingIfr, ok := inFlight.(*inFlightRequest); ok {
				c.oauth2Credentials.Logger.Printf("[INFO] waiting for in-flight request to complete, key:%s\n", key)
				existingIfr.wg.Wait() // Wait for the in-flight request to complete
				// Once the request completes, check cache again
				cachedResp := c.oauth2Credentials.CacheManager.Get(key)
				if cachedResp != nil {
					respData, err := io.ReadAll(cachedResp.Body)
					if err == nil {
						cachedResp.Body = io.NopCloser(bytes.NewBuffer(respData))
					}
					c.oauth2Credentials.Logger.Printf("[INFO] served from cache after waiting for in-flight request, key:%s\n", key)
					return respData, cachedResp, req, nil
				}
				// If still not in cache, there was an error - check the in-flight request error
				if existingIfr.err != nil {
					return nil, existingIfr.resp, existingIfr.req, existingIfr.err
				}
				// If no error but not cached, something went wrong - proceed with normal request
			}
		}

		// Create new in-flight request and mark it as in-flight
		ifr = &inFlightRequest{}
		ifr.wg.Add(1)
		c.inFlightRequests.Store(key, ifr)

		// Always clean up the in-flight request when done (success or error)
		defer func() {
			if ifr != nil {
				ifr.once.Do(func() {
					ifr.wg.Done()
					c.inFlightRequests.Delete(key)
				})
			}
		}()
	} else if method != http.MethodGet && c.oauth2Credentials.Zscaler.Client.Cache.Enabled && !isSandboxRequest {
		// For non-GET requests, invalidate cache
		c.oauth2Credentials.CacheManager.Delete(key)
		c.oauth2Credentials.CacheManager.ClearAllKeysWithPrefix(strings.Split(key, "?")[0])
	}

	var resp *http.Response
	sessionNotValidRetryCount := 0
	maxSessionNotValidRetries := int(c.oauth2Credentials.Zscaler.Client.RateLimit.MaxSessionNotValidRetries)

	for retry := 1; ; retry++ { // Infinite loop for retries if MaxRetries=0
		// Check MaxRetries if non-zero
		if c.oauth2Credentials.Zscaler.Client.RateLimit.MaxRetries > 0 && retry > int(c.oauth2Credentials.Zscaler.Client.RateLimit.MaxRetries) {
			return nil, resp, nil, fmt.Errorf("max retries exceeded")
		}

		// Check RequestTimeout - exclude rate limiting wait times from the calculation
		// This ensures that rate limiting delays don't count against the request timeout
		elapsedTime := time.Since(overallStartTime) - totalWaitTime
		if c.oauth2Credentials.Zscaler.Client.RequestTimeout > 0 && elapsedTime >= c.oauth2Credentials.Zscaler.Client.RequestTimeout {
			c.oauth2Credentials.Logger.Printf("[ERROR] Request timeout exceeded: elapsed=%v, waited=%v, total=%v, timeout=%v",
				elapsedTime, totalWaitTime, time.Since(overallStartTime), c.oauth2Credentials.Zscaler.Client.RequestTimeout)
			return nil, resp, nil, fmt.Errorf("request timeout exceeded after %v (excluding %v of rate limit waits)",
				elapsedTime, totalWaitTime)
		}

		start := time.Now()
		reqID := uuid.New().String()
		logger.LogRequest(c.oauth2Credentials.Logger, req, reqID, nil, !isSandboxRequest)
		httpClient := c.getServiceHTTPClient(endpoint)
		resp, err = httpClient.Do(req)
		logger.LogResponse(c.oauth2Credentials.Logger, resp, start, reqID)
		if err != nil {
			return nil, resp, nil, err
		}

		// ✅ Check for session invalidation errors in 401 responses
		if resp.StatusCode == http.StatusUnauthorized {
			bodyCopy, readErr := io.ReadAll(resp.Body)
			if readErr == nil {
				resp.Body = io.NopCloser(bytes.NewReader(bodyCopy)) // rewind

				// Use centralized error checking logic
				if errorx.IsSessionInvalidError(resp) {
					sessionNotValidRetryCount++
					if sessionNotValidRetryCount > maxSessionNotValidRetries {
						return nil, resp, req, fmt.Errorf("max SESSION_NOT_VALID retries exceeded (%d), possible authentication issue", maxSessionNotValidRetries)
					}

					c.oauth2Credentials.Logger.Printf("[WARN] Session invalidation detected (attempt %d, session retry %d/%d), refreshing token and retrying...", retry, sessionNotValidRetryCount, maxSessionNotValidRetries)

					// Enhanced debugging for session invalidation analysis
					if c.oauth2Credentials.Debug {
						tok := c.oauth2Credentials.Zscaler.Client.AuthToken
						c.oauth2Credentials.Logger.Printf("[DEBUG] Session invalidation analysis:")
						c.oauth2Credentials.Logger.Printf("[DEBUG]   - Token exists: %v", tok != nil)
						c.oauth2Credentials.Logger.Printf("[DEBUG]   - Token expiry: %s", tok.Expiry.Format(time.RFC3339))
						c.oauth2Credentials.Logger.Printf("[DEBUG]   - Current time: %s", time.Now().Format(time.RFC3339))
						c.oauth2Credentials.Logger.Printf("[DEBUG]   - Time until expiry: %.2f seconds", time.Until(tok.Expiry).Seconds())
						c.oauth2Credentials.Logger.Printf("[DEBUG]   - Request URL: %s", req.URL.String())
						c.oauth2Credentials.Logger.Printf("[DEBUG]   - Request method: %s", req.Method)
						c.oauth2Credentials.Logger.Printf("[DEBUG]   - Authorization header present: %v", req.Header.Get("Authorization") != "")
						if req.Header.Get("Authorization") != "" {
							authHeader := req.Header.Get("Authorization")
							c.oauth2Credentials.Logger.Printf("[DEBUG]   - Authorization header present: %v", authHeader != "")
						}
					}

					// Force token refresh regardless of client-side validation
					// Session invalidation means the server considers the token invalid
					c.Lock() // Prevent concurrent token refresh
					authToken, err := Authenticate(c.oauth2Credentials.Context, c.oauth2Credentials, c.oauth2Credentials.Logger)
					if err != nil {
						c.Unlock()
						return nil, resp, req, fmt.Errorf("token refresh failed after session invalidation: %w", err)
					}
					c.oauth2Credentials.Zscaler.Client.AuthToken = authToken
					c.Unlock()
					c.oauth2Credentials.Logger.Printf("[INFO] Token refreshed successfully, retrying request...")

					// Add a small delay before retrying to avoid overwhelming the server
					beforeWait := time.Now()
					time.Sleep(time.Second * 2)
					totalWaitTime += time.Since(beforeWait)

					// Rebuild request with fresh body reader
					req, err = c.buildRequest(ctx, method, endpoint, bytes.NewReader(requestBodyBytes), urlParams, contentType)
					if err != nil {
						return nil, nil, nil, err
					}
					continue
				}
			}
			resp.Body = io.NopCloser(bytes.NewReader(bodyCopy)) // rewind even if not retrying
		}

		// Reset session retry counter on successful requests or other errors
		if resp.StatusCode != http.StatusUnauthorized {
			sessionNotValidRetryCount = 0
		}

		// Handle rate-limiting (429 or 503)
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			retryAfter := getRetryAfter(resp, c.oauth2Credentials)
			if retryAfter > 0 {
				// If Retry-After is very long (> 5 minutes), fail fast with clear message
				if retryAfter > 5*time.Minute {
					c.oauth2Credentials.Logger.Printf("[ERROR] Rate limit exceeded with very long Retry-After: %v. This typically indicates hourly rate limits have been reached.", retryAfter)
					return nil, resp, nil, fmt.Errorf("rate limit exceeded: API requires waiting %v before retrying. This typically means hourly rate limits (GET: 1000/hr, POST/PUT: 1000/hr, DELETE: 400/hr) have been reached. Please wait before retrying", retryAfter)
				}

				// Check if context is still valid before sleeping
				select {
				case <-ctx.Done():
					return nil, resp, nil, fmt.Errorf("context cancelled while waiting for rate limit: %w", ctx.Err())
				default:
					c.oauth2Credentials.Logger.Printf("[INFO] Rate limit hit, waiting %v before retry (attempt %d)", retryAfter, retry)
					// Track this wait time so it doesn't count against request timeout
					beforeWait := time.Now()
					time.Sleep(retryAfter)
					totalWaitTime += time.Since(beforeWait)
					continue
				}
			}
			// If no Retry-After header, fall through to exponential backoff
		}

		// Handle success
		if resp.StatusCode < 300 {
			break
		}

		// Handle 409 Conflict and 412 Precondition Failed (EDIT_LOCK_NOT_AVAILABLE, org barrier) - these are retryable
		// The API may return either 409 or 412 for "Failed during enter Org barrier" under high concurrency
		if resp.StatusCode == http.StatusConflict || resp.StatusCode == http.StatusPreconditionFailed {
			bodyCopy, readErr := io.ReadAll(resp.Body)
			if readErr == nil {
				resp.Body = io.NopCloser(bytes.NewReader(bodyCopy)) // rewind

				// Check if this is an EDIT_LOCK_NOT_AVAILABLE / org barrier error
				if errorx.IsEditLockError(resp) {
					backoffDelay := time.Duration(math.Pow(2, float64(retry-1))) * c.oauth2Credentials.Zscaler.Client.RateLimit.RetryWaitMin
					maxBackoff := c.oauth2Credentials.Zscaler.Client.RateLimit.RetryWaitMax
					if backoffDelay > maxBackoff {
						backoffDelay = maxBackoff
					}
					c.oauth2Credentials.Logger.Printf("[WARN] Edit lock/org barrier conflict detected (status %d, attempt %d), waiting %v before retry", resp.StatusCode, retry, backoffDelay)

					// Track this wait time so it doesn't count against request timeout
					beforeWait := time.Now()
					time.Sleep(backoffDelay)
					totalWaitTime += time.Since(beforeWait)

					// Rebuild request with fresh body reader
					req, err = c.buildRequest(ctx, method, endpoint, bytes.NewReader(requestBodyBytes), urlParams, contentType)
					if err != nil {
						return nil, nil, nil, err
					}
					continue
				}
			}
			resp.Body = io.NopCloser(bytes.NewReader(bodyCopy)) // rewind even if not retrying
		}

		// Handle retryable errors with exponential backoff
		if resp.StatusCode >= 500 {
			// 5xx server errors are retryable
			backoffDelay := time.Duration(math.Pow(2, float64(retry-1))) * c.oauth2Credentials.Zscaler.Client.RateLimit.RetryWaitMin
			maxBackoff := c.oauth2Credentials.Zscaler.Client.RateLimit.RetryWaitMax
			if backoffDelay > maxBackoff {
				backoffDelay = maxBackoff
			}
			c.oauth2Credentials.Logger.Printf("[INFO] Server error (status %d), retry %d: waiting %v before retry", resp.StatusCode, retry, backoffDelay)
			// Track this wait time so it doesn't count against request timeout
			beforeWait := time.Now()
			time.Sleep(backoffDelay)
			totalWaitTime += time.Since(beforeWait)
			continue
		}

		// Handle non-retryable client errors (4xx except 401, 409, 429 which are handled above)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return nil, resp, nil, errorx.CheckErrorInResponse(resp, fmt.Errorf("client error"))
		}

		// Handle other non-success status codes
		if resp.StatusCode >= 300 {
			return nil, resp, nil, errorx.CheckErrorInResponse(resp, fmt.Errorf("API error"))
		}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, nil, err
	}

	// Log timing information if there were significant waits
	totalElapsed := time.Since(overallStartTime)
	actualRequestTime := totalElapsed - totalWaitTime
	if totalWaitTime > 0 && c.oauth2Credentials.Debug {
		c.oauth2Credentials.Logger.Printf("[DEBUG] Request completed: total=%v, waited=%v, actual=%v",
			totalElapsed, totalWaitTime, actualRequestTime)
	}

	// Cache logic for successful GET requests
	if !isSandboxRequest && c.oauth2Credentials.Zscaler.Client.Cache.Enabled && method == http.MethodGet {
		resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		c.oauth2Credentials.Logger.Printf("[INFO] saving to cache, key:%s\n", key)
		c.oauth2Credentials.CacheManager.Set(key, cache.CopyResponse(resp))
	}
	_ = tryDrainBody(resp.Body)
	return bodyBytes, resp, req, nil
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
	tok := c.oauth2Credentials.Zscaler.Client.AuthToken

	if tok == nil || tok.AccessToken == "" || tok.Expiry.IsZero() {
		if c.oauth2Credentials.Logger != nil {
			c.oauth2Credentials.Logger.Printf("[DEBUG] authValid: token is nil or expiry not set")
		}
		return false
	}

	expiresIn := time.Until(tok.Expiry).Seconds()
	valid := time.Now().Before(tok.Expiry.Add(-30 * time.Second))

	if c.oauth2Credentials.Logger != nil {
		c.oauth2Credentials.Logger.Printf("[DEBUG] authValid: token exists=%v, expires_in=%.2f, valid=%v",
			true,
			expiresIn,
			valid,
		)
	}

	return valid
}

// Unified authentication function to refresh OAuth2 tokens
func (c *Client) authenticate() error {
	if c.oauth2Credentials.UseLegacyClient {
		return nil // skip authentication for legacy client
	}
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

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
