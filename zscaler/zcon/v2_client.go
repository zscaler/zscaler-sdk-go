package zcon

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

	// Validate ZIA Cloud
	if config.ZCON.Client.ZCONCloud == "" {
		logger.Printf("[ERROR] Missing ZIA cloud configuration.")
		return nil, errors.New("ZIACloud configuration is missing")
	}

	// Construct the base URL
	baseURL := config.BaseURL.String()

	// Validate authentication credentials
	if config.ZCON.Client.ZCONUsername == "" || config.ZCON.Client.ZCONPassword == "" || config.ZCON.Client.ZCONApiKey == "" {
		logger.Printf("[ERROR] Missing required ZIA credentials (username, password, or API key).")
		return nil, errors.New("missing required ZIA credentials")
	}

	// Initialize rate limiter
	rateLimiter := rl.NewRateLimiter(
		int(config.ZCON.Client.RateLimit.MaxRetries),
		int(config.ZCON.Client.RateLimit.RetryWaitMin.Seconds()),
		int(config.ZCON.Client.RateLimit.RetryWaitMax.Seconds()),
		int(config.ZCON.Client.RequestTimeout.Seconds()),
	)

	// Initialize HTTP client
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = getHTTPClient(logger, rateLimiter, config)
	}

	// Initialize the Client instance
	cli := &Client{
		userName:         config.ZCON.Client.ZCONUsername,
		password:         config.ZCON.Client.ZCONPassword,
		apiKey:           config.ZCON.Client.ZCONApiKey,
		cloud:            config.ZCON.Client.ZCONCloud,
		HTTPClient:       httpClient,
		URL:              baseURL,
		Logger:           logger,
		UserAgent:        config.UserAgent,
		cacheEnabled:     config.ZCON.Client.Cache.Enabled,
		cacheTtl:         config.ZCON.Client.Cache.DefaultTtl,
		cacheCleanwindow: config.ZCON.Client.Cache.DefaultTti,
		cacheMaxSizeMB:   int(config.ZCON.Client.Cache.DefaultCacheMaxSizeMB),
		rateLimiter:      rateLimiter,
		stopTicker:       make(chan bool),
		sessionTimeout:   JSessionIDTimeout * time.Minute,
		sessionRefreshed: time.Time{},
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

	//logger.Printf("[DEBUG] ZIA client successfully initialized with base URL: %s", baseURL)
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
func MakeAuthRequestZCON(credentials *Credentials, baseURL string, client *http.Client, userAgent string) (*Session, error) {
	if credentials == nil {
		return nil, fmt.Errorf("empty credentials")
	}

	data, err := json.Marshal(credentials)
	if err != nil {
		return nil, err
	}

	// Use baseURL directly and append only the endpoint
	authURL := fmt.Sprintf("%s%s", baseURL, zconAPIAuthURL)
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
	session, err := MakeAuthRequestZCON(&credentialData, c.URL, c.HTTPClient, c.UserAgent)
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

	// Initialize or refresh session if necessary
	if c.session == nil || now.After(c.sessionRefreshed.Add(c.sessionTimeout-jSessionTimeoutOffset)) {
		c.Logger.Printf("[INFO] Session is invalid or expired. Refreshing session...")
		if !c.refreshing {
			c.refreshing = true
			defer func() { c.refreshing = false }()

			// Refresh session
			timeStamp := getCurrentTimestampMilisecond()
			obfuscatedKey, err := obfuscateAPIKey(c.apiKey, timeStamp)
			if err != nil {
				return err
			}

			credentials := &Credentials{
				Username:  c.userName,
				Password:  c.password,
				APIKey:    obfuscatedKey,
				TimeStamp: timeStamp,
			}
			session, err := MakeAuthRequestZCON(credentials, c.URL, c.HTTPClient, c.UserAgent)
			if err != nil {
				c.Logger.Printf("[ERROR] Failed to refresh session: %v", err)
				return err
			}

			// Update session and timeout
			c.session = session
			c.sessionRefreshed = time.Now()
			if c.session.PasswordExpiryTime > 0 {
				c.sessionTimeout = time.Duration(c.session.PasswordExpiryTime) * time.Second
			} else {
				c.sessionTimeout = JSessionIDTimeout * time.Minute
			}
		} else {
			c.Logger.Printf("[INFO] Another goroutine is refreshing the session. Waiting...")
			// Wait for ongoing refresh to complete
			for c.refreshing {
				time.Sleep(100 * time.Millisecond)
			}
		}
	} else {
		c.Logger.Printf("[INFO] Session is valid, no refresh needed.")
	}

	// Ensure the JSESSIONID is set in the HTTP client cookies
	url, err := url.Parse(c.URL)
	if err != nil {
		c.Logger.Printf("[ERROR] Failed to parse URL: %v", err)
		return err
	}
	if c.HTTPClient.Jar == nil {
		c.HTTPClient.Jar, err = cookiejar.New(nil)
		if err != nil {
			c.Logger.Printf("[ERROR] Failed to create HTTP cookie jar: %v", err)
			return err
		}
	}
	c.HTTPClient.Jar.SetCookies(url, []*http.Cookie{
		{Name: cookieName, Value: c.session.JSessionID},
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
	retryableClient.RetryWaitMin = cfg.ZCON.Client.RateLimit.RetryWaitMin
	retryableClient.RetryWaitMax = cfg.ZCON.Client.RateLimit.RetryWaitMax

	if cfg.ZCON.Client.RateLimit.MaxRetries == 0 {
		// Set RetryMax to a very large number to simulate indefinite retries within the timeout duration.
		retryableClient.RetryMax = math.MaxInt32
	} else {
		retryableClient.RetryMax = int(cfg.ZCON.Client.RateLimit.MaxRetries)
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
	if cfg.ZCON.Client.RequestTimeout == 0 {
		retryableClient.HTTPClient.Timeout = time.Second * 60 // Default to 60 seconds if not specified.
	} else {
		retryableClient.HTTPClient.Timeout = cfg.ZCON.Client.RequestTimeout
	}

	// Configure proxy settings from configuration
	proxyFunc := http.ProxyFromEnvironment // Default behavior (uses system/env variables)
	if cfg.ZCON.Client.Proxy.Host != "" {
		// Include username and password if provided
		proxyURLString := fmt.Sprintf("http://%s:%d", cfg.ZCON.Client.Proxy.Host, cfg.ZCON.Client.Proxy.Port)
		if cfg.ZCON.Client.Proxy.Username != "" && cfg.ZCON.Client.Proxy.Password != "" {
			// URL-encode the username and password
			proxyAuth := url.UserPassword(cfg.ZCON.Client.Proxy.Username, cfg.ZCON.Client.Proxy.Password)
			proxyURLString = fmt.Sprintf("http://%s@%s:%d", proxyAuth.String(), cfg.ZCON.Client.Proxy.Host, cfg.ZCON.Client.Proxy.Port)
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
	if cfg.ZCON.Testing.DisableHttpsCheck {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: false, // This disables HTTPS certificate validation
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

func (c *Client) Logout(ctx context.Context) error {
	_, err := c.Request(ctx, zconAPIAuthURL, "DELETE", nil, "application/json") // Pass context as the first argument
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
func (c *Client) GenericRequest(ctx context.Context, baseUrl, endpoint, method string, body io.Reader, urlParams url.Values, contentType string) ([]byte, error) {
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
	req, err = http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	if c.UserAgent != "" {
		req.Header.Add("User-Agent", c.UserAgent)
	}

	err = c.checkSession()
	if err != nil {
		return nil, err
	}
	jSessionID := c.session.JSessionID
	req.Header.Set("JSessionID", jSessionID)

	otherHeaders := map[string]string{}

	reqID := uuid.New().String()
	start := time.Now()
	logger.LogRequest(c.Logger, req, reqID, otherHeaders, true)

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
func (c *Client) Request(ctx context.Context, endpoint, method string, data []byte, contentType string) ([]byte, error) {
	return c.GenericRequest(ctx, c.URL, endpoint, method, bytes.NewReader(data), nil, contentType)
}

func (client *Client) WithFreshCache() {
	client.freshCache = true
}

// Create sends an HTTP POST request.
func (c *Client) Create(ctx context.Context, endpoint string, o interface{}) (interface{}, error) {
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

	resp, err := c.Request(ctx, endpoint, "POST", data, "application/json")
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

func (c *Client) CreateWithSlicePayload(ctx context.Context, endpoint string, slice interface{}) ([]byte, error) {
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

	resp, err := c.Request(ctx, endpoint, "POST", data, "application/json")
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

func (c *Client) UpdateWithSlicePayload(ctx context.Context, endpoint string, slice interface{}) ([]byte, error) {
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

	resp, err := c.Request(ctx, endpoint, "PUT", data, "application/json")
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CreateWithRawPayload sends an HTTP POST request with a raw string payload.
func (c *Client) CreateWithRawPayload(ctx context.Context, endpoint string, payload string) ([]byte, error) {
	if payload == "" {
		return nil, errors.New("tried to create with an empty string payload")
	}

	// Convert the string payload to []byte
	data := []byte(payload)

	// Send the raw string as a POST request
	resp, err := c.Request(ctx, endpoint, "POST", data, "application/json")
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
func (c *Client) Read(ctx context.Context, endpoint string, o interface{}) error {
	contentType := c.GetContentType()
	resp, err := c.Request(ctx, endpoint, "GET", nil, contentType)
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
func (c *Client) UpdateWithPut(ctx context.Context, endpoint string, o interface{}) (interface{}, error) {
	return c.updateGeneric(ctx, endpoint, o, "PUT", "application/json")
}

// Update ...
func (c *Client) Update(ctx context.Context, endpoint string, o interface{}) (interface{}, error) {
	return c.updateGeneric(ctx, endpoint, o, "PATCH", "application/merge-patch+json")
}

// Update ...
func (c *Client) updateGeneric(ctx context.Context, endpoint string, o interface{}, method, contentType string) (interface{}, error) {
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

	resp, err := c.Request(ctx, endpoint, method, data, contentType)
	if err != nil {
		return nil, err
	}

	responseObject := reflect.New(t).Interface()
	err = json.Unmarshal(resp, &responseObject)
	return responseObject, err
}

// Delete ...
func (c *Client) Delete(ctx context.Context, endpoint string) error {
	_, err := c.Request(ctx, endpoint, "DELETE", nil, "application/json")
	if err != nil {
		return err
	}
	return nil
}

// BulkDelete sends an HTTP POST request for bulk deletion and expects a 204 No Content response.
func (c *Client) BulkDelete(ctx context.Context, endpoint string, payload interface{}) (*http.Response, error) {
	if payload == nil {
		return nil, errors.New("tried to delete with a nil payload, expected a struct")
	}

	// Marshal the payload into JSON
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Send the POST request
	resp, err := c.Request(ctx, endpoint, "POST", data, "application/json")
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
