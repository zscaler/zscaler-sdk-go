package zia

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/zscaler/zscaler-sdk-go/v2/cache"
	"github.com/zscaler/zscaler-sdk-go/v2/logger"
	rl "github.com/zscaler/zscaler-sdk-go/v2/ratelimiter"
	"github.com/zscaler/zscaler-sdk-go/v2/utils"
	"github.com/zscaler/zscaler-sdk-go/v2/zidentity"
)

const (
	maxIdleConnections    int = 40
	requestTimeout        int = 60
	JSessionIDTimeout         = 30 // minutes.
	jSessionTimeoutOffset     = 5 * time.Minute
	contentTypeJSON           = "application/json"
	cookieName                = "JSESSIONID"
	MaxNumOfRetries           = 100
	RetryWaitMaxSeconds       = 20
	RetryWaitMinSeconds       = 5
	// API types.
	ziaAPIVersion = "api/v1"
	ziaAPIAuthURL = "/authenticatedSession"
	loggerPrefix  = "zia-logger: "
)

// Client ...
type Client struct {
	sync.Mutex
	userName          string
	password          string
	cloud             string
	apiKey            string
	session           *Session
	sessionRefreshed  time.Time     // Also indicates last usage
	sessionTimeout    time.Duration // in minutes
	URL               string
	HTTPClient        *http.Client
	Logger            logger.Logger
	UserAgent         string
	freshCache        bool
	cacheEnabled      bool
	cache             cache.Cache
	cacheTtl          time.Duration
	cacheCleanwindow  time.Duration
	cacheMaxSizeMB    int
	rateLimiter       *rl.RateLimiter
	sessionTicker     *time.Ticker
	stopTicker        chan bool
	refreshing        bool
	useOneAPI         bool
	oauth2Credentials *zidentity.Credentials
}

// Session ...
type Session struct {
	AuthType           string `json:"authType"`
	ObfuscateAPIKey    bool   `json:"obfuscateApiKey"`
	PasswordExpiryTime int    `json:"passwordExpiryTime"`
	PasswordExpiryDays int    `json:"passwordExpiryDays"`
	Source             string `json:"source"`
	JSessionID         string `json:"jSessionID,omitempty"`
}

// Credentials ...
type Credentials struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	APIKey    string `json:"apiKey"`
	TimeStamp string `json:"timestamp"`
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

// NewOneAPIClient Returns a Client from credentials passed as parameters.
func NewOneAPIClient(clientID, clientSecret, ziaCloud, userAgent, oauth2ProviderUrl string) (*Client, error) {
	logger := logger.GetDefaultLogger(loggerPrefix)
	rateLimiter := rl.NewRateLimiter(2, 1, 1, 1)
	httpClient := getHTTPClient(logger, rateLimiter)

	if clientID == "" || clientSecret == "" {
		clientID = os.Getenv(zidentity.ZIDENTITY_CLIENT_ID)
		clientSecret = os.Getenv(zidentity.ZIDENTITY_CLIENT_SECRET)
	}

	if ziaCloud == "" {
		ziaCloud = os.Getenv("ZIA_CLOUD")
	}

	if oauth2ProviderUrl == "" {
		oauth2ProviderUrl = os.Getenv(zidentity.ZIDENTITY_OAUTH2_PROVIDER_URL)
	}

	var url string
	if strings.EqualFold(ziaCloud, "PRODUCTION") {
		url = "https://api.zsapi.net/zia/" + ziaAPIVersion
	} else {
		url = fmt.Sprintf("https://api.%s.zsapi.net/zia/%s", strings.ToLower(ziaCloud), ziaAPIVersion)
	}
	/*
		TODO: handle this case
			if ziaCloud == "zspreview" {
				url = fmt.Sprintf("https://admin.%s.net/%s", ziaCloud, ziaAPIVersion)
			}
	*/

	cacheDisabled, _ := strconv.ParseBool(os.Getenv("ZSCALER_SDK_CACHE_DISABLED"))
	cli := &Client{
		cloud:            ziaCloud,
		HTTPClient:       httpClient,
		URL:              url,
		Logger:           logger,
		UserAgent:        userAgent,
		cacheEnabled:     !cacheDisabled,
		cacheTtl:         time.Minute * 10,
		cacheCleanwindow: time.Minute * 8,
		cacheMaxSizeMB:   0,
		rateLimiter:      rateLimiter,
		stopTicker:       make(chan bool),
		sessionTimeout:   30 * time.Minute, // Initialize with a default session timeout
		useOneAPI:        true,
		oauth2Credentials: &zidentity.Credentials{
			ClientID:          clientID,
			ClientSecret:      clientSecret,
			Oauth2ProviderUrl: oauth2ProviderUrl,
		},
	}

	cche, err := cache.NewCache(cli.cacheTtl, cli.cacheCleanwindow, cli.cacheMaxSizeMB)
	if err != nil {
		cche = cache.NewNopCache()
	}
	cli.cache = cche

	// Start the session refresh ticker
	cli.startSessionTicker()

	return cli, nil
}

// NewClient Returns a Client from credentials passed as parameters.
func NewClient(username, password, apiKey, ziaCloud, userAgent string) (*Client, error) {
	logger := logger.GetDefaultLogger(loggerPrefix)
	rateLimiter := rl.NewRateLimiter(2, 1, 1, 1)
	httpClient := getHTTPClient(logger, rateLimiter)
	url := fmt.Sprintf("https://zsapi.%s.net/%s", ziaCloud, ziaAPIVersion)
	if ziaCloud == "zspreview" {
		url = fmt.Sprintf("https://admin.%s.net/%s", ziaCloud, ziaAPIVersion)
	}
	cacheDisabled, _ := strconv.ParseBool(os.Getenv("ZSCALER_SDK_CACHE_DISABLED"))
	cli := &Client{
		userName:         username,
		password:         password,
		apiKey:           apiKey,
		cloud:            ziaCloud,
		HTTPClient:       httpClient,
		URL:              url,
		Logger:           logger,
		UserAgent:        userAgent,
		cacheEnabled:     !cacheDisabled,
		cacheTtl:         time.Minute * 10,
		cacheCleanwindow: time.Minute * 8,
		cacheMaxSizeMB:   0,
		rateLimiter:      rateLimiter,
		stopTicker:       make(chan bool),
		sessionTimeout:   30 * time.Minute, // Initialize with a default session timeout
	}

	cche, err := cache.NewCache(cli.cacheTtl, cli.cacheCleanwindow, cli.cacheMaxSizeMB)
	if err != nil {
		cche = cache.NewNopCache()
	}
	cli.cache = cche

	// Start the session refresh ticker
	cli.startSessionTicker()

	return cli, nil
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

	req, err := http.NewRequest("POST", url+ziaAPIAuthURL, bytes.NewReader(data))
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

func (c *Client) WithCache(cache bool) {
	c.cacheEnabled = cache
}

func (c *Client) WithCacheTtl(i time.Duration) {
	c.cacheTtl = i
	c.Lock()
	c.cache.Close()
	cche, err := cache.NewCache(i, c.cacheCleanwindow, c.cacheMaxSizeMB)
	if err != nil {
		cche = cache.NewNopCache()
	}
	c.cache = cche
	c.Unlock()
}

func (c *Client) WithCacheCleanWindow(i time.Duration) {
	c.cacheCleanwindow = i
	c.Lock()
	c.cache.Close()
	cche, err := cache.NewCache(c.cacheTtl, i, c.cacheMaxSizeMB)
	if err != nil {
		cche = cache.NewNopCache()
	}
	c.cache = cche
	c.Unlock()
}

// checkSession checks if the session is valid and refreshes it if necessary.
func (c *Client) checkSession() error {
	c.Lock()
	defer c.Unlock()
	if c.useOneAPI {
		if c.oauth2Credentials != nil && (c.oauth2Credentials.AuthToken == nil || c.oauth2Credentials.AuthToken.AccessToken == "" || utils.IsTokenExpired(c.oauth2Credentials.AuthToken.AccessToken)) {
			a, err := zidentity.Authenticate(
				c.oauth2Credentials.ClientID,
				c.oauth2Credentials.ClientSecret,
				c.oauth2Credentials.Oauth2ProviderUrl,
				c.UserAgent,
				c.HTTPClient,
			)
			if err != nil {
				return err
			}
			c.oauth2Credentials.AuthToken = a
			return nil
		}
		return nil
	}

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

func getHTTPClient(l logger.Logger, rateLimiter *rl.RateLimiter) *http.Client {
	retryableClient := retryablehttp.NewClient()
	retryableClient.RetryWaitMin = time.Second * time.Duration(RetryWaitMinSeconds)
	retryableClient.RetryWaitMax = time.Second * time.Duration(RetryWaitMaxSeconds)
	retryableClient.RetryMax = MaxNumOfRetries

	// Set up the cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		l.Printf("[ERROR] failed to create cookie jar: %v", err)
		// Handle the error, possibly by continuing without a cookie jar
		// or you can choose to halt the execution if the cookie jar is critical
	}

	// Configure the underlying HTTP client
	retryableClient.HTTPClient = &http.Client{
		Jar: jar, // Set the cookie jar
		// ... other configurations ...
	}

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
				} else {
					return 0
				}
			}
		}
		// default to exp backoff
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
	retryableClient.HTTPClient.Transport = &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		MaxIdleConnsPerHost: maxIdleConnections,
	}

	retryableClient.HTTPClient = &http.Client{
		Timeout: time.Duration(requestTimeout) * time.Second,
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			MaxIdleConnsPerHost: maxIdleConnections,
		},
		Jar: jar, // Set the cookie jar
	}
	retryableClient.HTTPClient.Transport = logging.NewSubsystemLoggingHTTPTransport("gozscaler", retryableClient.HTTPClient.Transport)

	retryableClient.CheckRetry = checkRetry
	retryableClient.Logger = l

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

func (c *Client) GetSandboxURL() string {
	return "https://csbapi." + c.cloud + ".net"
}

func (c *Client) GetSandboxToken() string {
	return os.Getenv("ZIA_SANDBOX_TOKEN")
}

// func (c *Client) startSessionTicker() {
// 	c.Lock()
// 	defer c.Unlock()

// 	if c.sessionTicker != nil {
// 		c.stopTicker <- true
// 		c.sessionTicker.Stop()
// 	}

// 	tickerInterval := c.sessionTimeout - 1*time.Minute
// 	c.sessionTicker = time.NewTicker(tickerInterval)
// 	go func() {
// 		for {
// 			select {
// 			case <-c.sessionTicker.C:
// 				err := c.refreshSession()
// 				if err != nil {
// 					c.Logger.Printf("[ERROR] Failed to refresh session: %v\n", err)
// 				}
// 			case <-c.stopTicker:
// 				return
// 			}
// 		}
// 	}()
// }

// startSessionTicker starts a ticker to refresh the session periodically
func (c *Client) startSessionTicker() {
	if c.useOneAPI {
		return
	}

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
