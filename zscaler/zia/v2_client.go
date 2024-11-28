package zia

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zscaler/zscaler-sdk-go/v3/cache"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
	rl "github.com/zscaler/zscaler-sdk-go/v3/ratelimiter"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
)

// NewClient Returns a Client from credentials passed as parameters.
func NewClient(config *Configuration) (*Client, error) {

	if config == nil {
		return nil, errors.New("configuration cannot be nil")
	}

	// Enable Debug logging if the Debug flag is set
	if config.Debug {
		_ = os.Setenv("ZSCALER_SDK_LOG", "true")
		_ = os.Setenv("ZSCALER_SDK_VERBOSE", "true")
		config.Logger = logger.GetDefaultLogger(loggerPrefix)
	}

	logger := logger.GetDefaultLogger(loggerPrefix)
	// logger.Printf("[DEBUG] Initializing client with provided configuration.")

	// Validate ZIA Cloud
	if config.ZIA.Client.ZIACloud == "" {
		logger.Printf("[ERROR] Missing ZIA cloud configuration. Ensure WithZiaCloud is set.")
		return nil, errors.New("ZIACloud configuration is missing")
	}

	// Construct the base URL based on the ZIA cloud
	baseURL := config.BaseURL.String()

	// Validate authentication credentials
	if config.ZIA.Client.ZIAUsername == "" || config.ZIA.Client.ZIAPassword == "" || config.ZIA.Client.ZIAApiKey == "" {
		logger.Printf("[ERROR] Missing required ZIA credentials (username, password, or API key).")
		return nil, errors.New("missing required ZIA credentials")
	}

	// Perform authentication using the provided credentials
	timeStamp := getCurrentTimestampMilisecond()
	obfuscatedKey, err := obfuscateAPIKey(config.ZIA.Client.ZIAApiKey, timeStamp)
	if err != nil {
		logger.Printf("[ERROR] Failed to obfuscate API key: %v", err)
		return nil, fmt.Errorf("failed to obfuscate API key: %w", err)
	}

	credentials := &Credentials{
		Username:  config.ZIA.Client.ZIAUsername,
		Password:  config.ZIA.Client.ZIAPassword,
		APIKey:    obfuscatedKey,
		TimeStamp: timeStamp,
	}

	// Initialize rate limiter
	rateLimiter := rl.NewRateLimiter(
		int(config.ZIA.Client.RateLimit.MaxRetries),
		int(config.ZIA.Client.RateLimit.RetryWaitMin.Seconds()),
		int(config.ZIA.Client.RateLimit.RetryWaitMax.Seconds()),
		int(config.ZIA.Client.RequestTimeout.Seconds()),
	)

	// Initialize HTTP client
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = getHTTPClient(logger, rateLimiter, config)
	}

	// Perform authentication request
	session, err := MakeAuthRequestZIA(credentials, baseURL, httpClient, config.UserAgent)
	if err != nil {
		logger.Printf("[ERROR] Authentication failed: %v", err)
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Initialize the Client instance
	cli := &Client{
		userName:         config.ZIA.Client.ZIAUsername,
		password:         config.ZIA.Client.ZIAPassword,
		apiKey:           config.ZIA.Client.ZIAApiKey,
		cloud:            config.ZIA.Client.ZIACloud,
		HTTPClient:       config.HTTPClient,
		URL:              baseURL,
		Logger:           logger,
		UserAgent:        config.UserAgent,
		cacheEnabled:     config.ZIA.Client.Cache.Enabled,
		cacheTtl:         config.ZIA.Client.Cache.DefaultTtl,
		cacheCleanwindow: config.ZIA.Client.Cache.DefaultTti,
		cacheMaxSizeMB:   int(config.ZIA.Client.Cache.DefaultCacheMaxSizeMB),
		rateLimiter:      rateLimiter,
		stopTicker:       make(chan bool),
		sessionTimeout:   JSessionIDTimeout * time.Minute, // Default session timeout
		session:          session,
		sessionRefreshed: time.Now(),
	}

	// Initialize the cache
	cche, err := cache.NewCache(cli.cacheTtl, cli.cacheCleanwindow, cli.cacheMaxSizeMB)
	if err != nil {
		logger.Printf("[WARN] Failed to initialize cache, using NopCache: %v", err)
		cche = cache.NewNopCache()
	}
	cli.cache = cche

	// Start the session refresh ticker
	cli.startSessionTicker()

	// logger.Printf("[DEBUG] ZIA client successfully initialized with base URL: %s", baseURL)
	return cli, nil
}

func obfuscateAPIKey(apiKey, timeStamp string) (string, error) {
	// check min required size
	if len(timeStamp) < 6 || len(apiKey) < 12 {
		return "", errors.New("time stamp or api key doesn't have required sizes")
	}

	seed := apiKey

	high := timeStamp[len(timeStamp)-6:]
	highInt, _ := strconv.Atoi(high)
	low := fmt.Sprintf("%06d", highInt>>1)
	key := ""

	for i := 0; i < len(high); i++ {
		index, _ := strconv.Atoi((string)(high[i]))
		key += (string)(seed[index])
	}
	for i := 0; i < len(low); i++ {
		index, _ := strconv.Atoi((string)(low[i]))
		key += (string)(seed[index+2])
	}

	return key, nil
}

// MakeAuthRequestZIA authenticates using the provided credentials and returns the session or an error.
func MakeAuthRequestZIA(credentials *Credentials, url string, client *http.Client, userAgent string) (*Session, error) {
	if credentials == nil {
		return nil, fmt.Errorf("empty credentials")
	}

	data, err := json.Marshal(credentials)
	if err != nil {
		return nil, err
	}

	// Add `/api/v1` only for the authentication request
	authURL := fmt.Sprintf("%s/%s%s", url, ziaAPIVersion, ziaAPIAuthURL)
	req, err := http.NewRequest("POST", authURL, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentTypeJSON)
	if userAgent != "" {
		req.Header.Add("User-Agent", userAgent)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body for use in error messages
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var session Session
		if err = json.Unmarshal(body, &session); err != nil {
			return nil, fmt.Errorf("error unmarshalling response: %v", err)
		}
		session.JSessionID, err = extractJSessionIDFromHeaders(resp.Header)
		if err != nil {
			return nil, err
		}
		return &session, nil
	case http.StatusBadRequest:
		return nil, fmt.Errorf("HTTP 400 Bad Request: %s", string(body))
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("HTTP 401 Unauthorized: %s", string(body))
	case http.StatusForbidden:
		return nil, fmt.Errorf("HTTP 403 Forbidden: %s", string(body))
	default:
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
}

func extractJSessionIDFromHeaders(header http.Header) (string, error) {
	sessionIdStr := header.Get("Set-Cookie")
	if sessionIdStr == "" {
		return "", fmt.Errorf("no Set-Cookie header received")
	}
	regex := regexp.MustCompile("JSESSIONID=(.*?);")
	// look for the first match we find
	result := regex.FindStringSubmatch(sessionIdStr)
	if len(result) < 2 {
		return "", fmt.Errorf("couldn't find JSESSIONID in header value")
	}
	return result[1], nil
}

func getCurrentTimestampMilisecond() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
}

// RefreshSession .. The caller should require lock.
func (c *Client) refreshSession() error {
	timeStamp := getCurrentTimestampMilisecond()
	obfuscatedKey, err := obfuscateAPIKey(c.apiKey, timeStamp)
	if err != nil {
		return err
	}
	credentialData := Credentials{
		Username:  c.userName,
		Password:  c.password,
		APIKey:    obfuscatedKey,
		TimeStamp: timeStamp,
	}
	session, err := MakeAuthRequestZIA(&credentialData, c.URL, c.HTTPClient, c.UserAgent)
	if err != nil {
		c.Logger.Printf("[ERROR] Failed to make auth request: %v\n", err)
		return err
	}
	c.session = session
	c.sessionRefreshed = time.Now()
	if c.session.PasswordExpiryTime == -1 {
		c.Logger.Printf("[INFO] PasswordExpiryTime is -1, setting sessionTimeout to 30 minutes")
		c.sessionTimeout = 30 * time.Minute
	} else {
		//c.Logger.Printf("[INFO] Setting session timeout based on PasswordExpiryTime: %v seconds", c.session.PasswordExpiryTime)
		c.sessionTimeout = time.Duration(c.session.PasswordExpiryTime) * time.Second
	}
	return nil
}

// checkSession checks if the session is valid and refreshes it if necessary.
func (c *Client) checkSession() error {
	c.Lock()
	defer c.Unlock()

	now := time.Now()
	if c.session == nil {
		c.Logger.Printf("[INFO] No session found, refreshing session")
		err := c.refreshSession()
		if err != nil {
			c.Logger.Printf("[ERROR] Failed to get session id: %v\n", err)
			return err
		}
	} else {
		c.Logger.Printf("[INFO] Current time: %v\nSession Refreshed: %v\nSession Timeout: %v\n",
			now.Format("2006-01-02 15:04:05 MST"),
			c.sessionRefreshed.Format("2006-01-02 15:04:05 MST"),
			c.sessionTimeout)

		if c.session.PasswordExpiryTime > 0 && now.After(c.sessionRefreshed.Add(c.sessionTimeout-jSessionTimeoutOffset)) {
			c.Logger.Printf("[INFO] Session timeout reached, refreshing session")
			if !c.refreshing {
				c.refreshing = true
				c.Unlock()
				err := c.refreshSession()
				c.Lock()
				c.refreshing = false
				if err != nil {
					c.Logger.Printf("[ERROR] Failed to refresh session id: %v\n", err)
					return err
				}
			} else {
				c.Logger.Printf("[INFO] Another refresh is in progress, waiting for it to complete")
			}
		} else {
			c.Logger.Printf("[INFO] Session is still valid, no need to refresh")
		}
	}

	url, err := url.Parse(c.URL)
	if err != nil {
		c.Logger.Printf("[ERROR] Failed to parse url %s: %v\n", c.URL, err)
		return err
	}

	if c.HTTPClient.Jar == nil {
		c.HTTPClient.Jar, err = cookiejar.New(nil)
		if err != nil {
			c.Logger.Printf("[ERROR] Failed to create new http cookie jar %v\n", err)
			return err
		}
	}

	c.HTTPClient.Jar.SetCookies(url, []*http.Cookie{
		{
			Name:  cookieName,
			Value: c.session.JSessionID,
		},
	})

	return nil
}

func (c *Client) GetContentType() string {
	return contentTypeJSON
}

// getHTTPClient sets up the retryable HTTP client with backoff and retry policies.
func getHTTPClient(l logger.Logger, rateLimiter *rl.RateLimiter, cfg *Configuration) *http.Client {
	retryableClient := retryablehttp.NewClient()

	// Set the retry settings, allowing user to override defaults.
	// Defaults are set by the config, which is initially read from constants but can be overridden.
	retryableClient.RetryWaitMin = cfg.ZIA.Client.RateLimit.RetryWaitMin
	retryableClient.RetryWaitMax = cfg.ZIA.Client.RateLimit.RetryWaitMax

	if cfg.ZIA.Client.RateLimit.MaxRetries == 0 {
		// Set RetryMax to a very large number to simulate indefinite retries within the timeout duration.
		retryableClient.RetryMax = math.MaxInt32
	} else {
		retryableClient.RetryMax = int(cfg.ZIA.Client.RateLimit.MaxRetries)
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

	// Set the request timeout, allowing user-defined override.
	if cfg.ZIA.Client.RequestTimeout == 0 {
		retryableClient.HTTPClient.Timeout = time.Second * 60 // Default to 60 seconds if not specified.
	} else {
		retryableClient.HTTPClient.Timeout = cfg.ZIA.Client.RequestTimeout
	}

	// Configure proxy settings from configuration
	proxyFunc := http.ProxyFromEnvironment // Default behavior (uses system/env variables)
	if cfg.ZIA.Client.Proxy.Host != "" {
		// Include username and password if provided
		proxyURLString := fmt.Sprintf("http://%s:%d", cfg.ZIA.Client.Proxy.Host, cfg.ZIA.Client.Proxy.Port)
		if cfg.ZIA.Client.Proxy.Username != "" && cfg.ZIA.Client.Proxy.Password != "" {
			// URL-encode the username and password
			proxyAuth := url.UserPassword(cfg.ZIA.Client.Proxy.Username, cfg.ZIA.Client.Proxy.Password)
			proxyURLString = fmt.Sprintf("http://%s@%s:%d", proxyAuth.String(), cfg.ZIA.Client.Proxy.Host, cfg.ZIA.Client.Proxy.Port)
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
	if cfg.ZIA.Testing.DisableHttpsCheck {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true, // This disables HTTPS certificate validation
		}
		l.Printf("[INFO] HTTPS certificate validation is disabled (testing mode).")
	}

	retryableClient.HTTPClient.Transport = transport
	return retryableClient.StandardClient()
}

func containsInt(codes []int, code int) bool {
	for _, a := range codes {
		if a == code {
			return true
		}
	}
	return false
}

func getRetryAfter(resp *http.Response, l logger.Logger) time.Duration {
	if s := resp.Header.Get("Retry-After"); s != "" {
		if sleep, err := strconv.ParseInt(s, 10, 64); err == nil {
			l.Printf("[INFO] got Retry-After from header:%s\n", s)
			return time.Second * time.Duration(sleep)
		} else {
			dur, err := time.ParseDuration(s)
			if err == nil {
				return dur
			}
			l.Printf("[INFO] error getting Retry-After from header:%s\n", err)
		}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Printf("[INFO] error getting Retry-After from body:%s\n", err)
		return 0
	}
	data := map[string]string{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		l.Printf("[INFO] error getting Retry-After from body:%s\n", err)
		return 0
	}
	if retryAfterStr, ok := data["Retry-After"]; ok && retryAfterStr != "" {
		l.Printf("[INFO] got Retry-After from body:%s\n", retryAfterStr)
		secondsStr := strings.Split(retryAfterStr, " ")[0]
		seconds, err := strconv.Atoi(secondsStr)
		if err != nil {
			l.Printf("[INFO] error getting Retry-After from body:%s\n", err)
			return 0
		}
		return time.Duration(seconds) * time.Second
	}
	return 0
}

// getRetryOnStatusCodes return a list of http status codes we want to apply retry on.
// Return empty slice to enable retry on all connection & server errors.
// Or return []int{429}  to retry on only TooManyRequests error.
func getRetryOnStatusCodes() []int {
	return []int{http.StatusTooManyRequests}
}

type ApiErr struct {
	Code    string `json:"code"`
	Message string `json:"message"`
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

	if resp != nil && (resp.StatusCode == http.StatusPreconditionFailed || resp.StatusCode == http.StatusConflict || resp.StatusCode == http.StatusUnauthorized) {
		apiRespErr := ApiErr{}
		data, err := io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(data))
		if err == nil {
			err = json.Unmarshal(data, &apiRespErr)
			if err == nil {
				if apiRespErr.Code == "UNEXPECTED_ERROR" && apiRespErr.Message == "Failed during enter Org barrier" ||
					apiRespErr.Code == "EDIT_LOCK_NOT_AVAILABLE" || apiRespErr.Message == "Resource Access Blocked" ||
					apiRespErr.Code == "UNEXPECTED_ERROR" && apiRespErr.Message == "Request processing failed, possibly because an expected precondition was not met" {
					return true, nil
				}
			}
		}
	}
	return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
}

func (c *Client) Logout() error {
	_, err := c.Request(ziaAPIAuthURL, "DELETE", nil, "application/json")
	if err != nil {
		return err
	}
	return nil
}

// startSessionTicker starts a ticker to refresh the session periodically
func (c *Client) startSessionTicker() {
	if c.sessionTimeout > 0 {
		c.sessionTicker = time.NewTicker(c.sessionTimeout - jSessionTimeoutOffset)
		go func() {
			for {
				select {
				case <-c.sessionTicker.C:
					c.Lock()
					if !c.refreshing {
						c.refreshing = true
						c.Unlock()
						c.refreshSession()
						c.Lock()
						c.refreshing = false
					}
					c.Unlock()
				case <-c.stopTicker:
					c.sessionTicker.Stop()
					return
				}
			}
		}()
	} else {
		c.Logger.Printf("[ERROR] Invalid session timeout value: %v\n", c.sessionTimeout)
	}
}

func (c *Client) do(req *http.Request, start time.Time, reqID string) (*http.Response, error) {
	key := cache.CreateCacheKey(req)
	if c.cacheEnabled {
		if req.Method != http.MethodGet {
			c.cache.Delete(key)
			c.cache.ClearAllKeysWithPrefix(strings.Split(key, "?")[0])
		}
		resp := c.cache.Get(key)
		inCache := resp != nil
		if c.freshCache {
			c.cache.Delete(key)
			inCache = false
			c.freshCache = false
		}
		if inCache {
			c.Logger.Printf("[INFO] served from cache, key:%s\n", key)
			return resp, nil
		}
	}

	// Ensure the session is valid before making the request
	err := c.checkSession()
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	logger.LogResponse(c.Logger, resp, start, reqID)
	if err != nil {
		return resp, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// Fallback check for SESSION_NOT_VALID
	if resp.StatusCode == http.StatusUnauthorized || strings.Contains(string(body), "SESSION_NOT_VALID") {
		// Refresh session and retry
		err := c.refreshSession()
		if err != nil {
			return nil, err
		}
		req.Header.Set("JSessionID", c.session.JSessionID)
		resp, err = c.HTTPClient.Do(req)
		logger.LogResponse(c.Logger, resp, start, reqID)
		if err != nil {
			return resp, err
		}
	}

	if c.cacheEnabled && resp.StatusCode >= 200 && resp.StatusCode <= 299 && req.Method == http.MethodGet {
		c.Logger.Printf("[INFO] saving to cache, key:%s\n", key)
		c.cache.Set(key, cache.CopyResponse(resp))
	}

	return resp, nil
}

// Request ... // Needs to review this function.
func (c *Client) GenericRequest(baseUrl, endpoint, method string, body io.Reader, urlParams url.Values, contentType string) ([]byte, error) {
	if contentType == "" {
		contentType = contentTypeJSON
	}

	var req *http.Request
	var resp *http.Response
	var err error
	params := ""
	if urlParams != nil {
		params = urlParams.Encode()
	}
	if strings.Contains(endpoint, "?") && params != "" {
		endpoint += "&" + params
	} else if params != "" {
		endpoint += "?" + params
	}
	fullURL := fmt.Sprintf("%s%s", baseUrl, endpoint)
	req, err = http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	if c.UserAgent != "" {
		req.Header.Add("User-Agent", c.UserAgent)
	}

	reqID := uuid.New().String()
	start := time.Now()
	for retry := 1; retry <= 5; retry++ {
		resp, err = c.do(req, start, reqID)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode <= 299 {
			defer resp.Body.Close()
			break
		}

		resp.Body.Close()
		if resp.StatusCode > 299 && resp.StatusCode != http.StatusUnauthorized {
			return nil, errorx.CheckErrorInResponse(resp, fmt.Errorf("api responded with code: %d", resp.StatusCode))
		}
	}

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bodyResp, nil
}

// Request ... // Needs to review this function.
func (c *Client) Request(endpoint, method string, data []byte, contentType string) ([]byte, error) {
	return c.GenericRequest(c.URL, endpoint, method, bytes.NewReader(data), nil, contentType)
}

func (client *Client) WithFreshCache() {
	client.freshCache = true
}

// Create sends an HTTP POST request.
func (c *Client) Create(endpoint string, o interface{}) (interface{}, error) {
	if o == nil {
		return nil, errors.New("tried to create with a nil payload not a Struct")
	}
	t := reflect.TypeOf(o)
	if t.Kind() != reflect.Struct {
		return nil, errors.New("tried to create with a " + t.Kind().String() + " not a Struct")
	}
	data, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(endpoint, "POST", data, "application/json")
	if err != nil {
		return nil, err
	}
	if len(resp) > 0 {
		// Check if the response is an array of strings
		var stringArrayResponse []string
		if json.Unmarshal(resp, &stringArrayResponse) == nil {
			return stringArrayResponse, nil
		}

		// Otherwise, handle as usual
		responseObject := reflect.New(t).Interface()
		err = json.Unmarshal(resp, &responseObject)
		if err != nil {
			return nil, err
		}
		id := reflect.Indirect(reflect.ValueOf(responseObject)).FieldByName("ID")

		c.Logger.Printf("Created Object with ID %v", id)
		return responseObject, nil
	} else {
		// in case of 204 no content
		return nil, nil
	}
}

func (c *Client) CreateWithSlicePayload(endpoint string, slice interface{}) ([]byte, error) {
	if slice == nil {
		return nil, errors.New("tried to create with a nil payload not a Slice")
	}

	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, errors.New("tried to create with a " + v.Kind().String() + " not a Slice")
	}

	data, err := json.Marshal(slice)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(endpoint, "POST", data, "application/json")
	if err != nil {
		return nil, err
	}
	if len(resp) > 0 {
		return resp, nil
	} else {
		// in case of 204 no content
		return nil, nil
	}
}

func (c *Client) UpdateWithSlicePayload(endpoint string, slice interface{}) ([]byte, error) {
	if slice == nil {
		return nil, errors.New("tried to update with a nil payload not a Slice")
	}

	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, errors.New("tried to update with a " + v.Kind().String() + " not a Slice")
	}

	data, err := json.Marshal(slice)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(endpoint, "PUT", data, "application/json")
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CreateWithRawPayload sends an HTTP POST request with a raw string payload.
func (c *Client) CreateWithRawPayload(endpoint string, payload string) ([]byte, error) {
	if payload == "" {
		return nil, errors.New("tried to create with an empty string payload")
	}

	// Convert the string payload to []byte
	data := []byte(payload)

	// Send the raw string as a POST request
	resp, err := c.Request(endpoint, "POST", data, "application/json")
	if err != nil {
		return nil, err
	}

	// Handle the response
	if len(resp) > 0 {
		return resp, nil
	} else {
		// in case of 204 no content
		return nil, nil
	}
}

// Read ...
func (c *Client) Read(endpoint string, o interface{}) error {
	contentType := c.GetContentType()
	resp, err := c.Request(endpoint, "GET", nil, contentType)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, o)
	if err != nil {
		return err
	}

	return nil
}

// Update ...
func (c *Client) UpdateWithPut(endpoint string, o interface{}) (interface{}, error) {
	return c.updateGeneric(endpoint, o, "PUT", "application/json")
}

// Update ...
func (c *Client) Update(endpoint string, o interface{}) (interface{}, error) {
	return c.updateGeneric(endpoint, o, "PATCH", "application/merge-patch+json")
}

// Update ...
func (c *Client) updateGeneric(endpoint string, o interface{}, method, contentType string) (interface{}, error) {
	if o == nil {
		return nil, errors.New("tried to update with a nil payload not a Struct")
	}
	t := reflect.TypeOf(o)
	if t.Kind() != reflect.Struct {
		return nil, errors.New("tried to update with a " + t.Kind().String() + " not a Struct")
	}
	data, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(endpoint, method, data, contentType)
	if err != nil {
		return nil, err
	}

	responseObject := reflect.New(t).Interface()
	err = json.Unmarshal(resp, &responseObject)
	return responseObject, err
}

// Delete ...
func (c *Client) Delete(endpoint string) error {
	_, err := c.Request(endpoint, "DELETE", nil, "application/json")
	if err != nil {
		return err
	}
	return nil
}

// BulkDelete sends an HTTP POST request for bulk deletion and expects a 204 No Content response.
func (c *Client) BulkDelete(endpoint string, payload interface{}) (*http.Response, error) {
	if payload == nil {
		return nil, errors.New("tried to delete with a nil payload, expected a struct")
	}

	// Marshal the payload into JSON
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Send the POST request
	resp, err := c.Request(endpoint, "POST", data, "application/json")
	if err != nil {
		return nil, err
	}

	// Check the status code (204 No Content expected)
	if len(resp) == 0 {
		c.Logger.Printf("[DEBUG] Bulk delete successful with 204 No Content")
		return &http.Response{StatusCode: 204}, nil
	}

	// If the response is not empty, this might indicate an error or unexpected behavior
	return &http.Response{StatusCode: 200}, fmt.Errorf("unexpected response: %s", string(resp))
}
