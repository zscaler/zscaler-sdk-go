package zdx

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
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
		config.Logger = logger.GetDefaultLogger("zdx-logger: ")
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

	config.ZDX.Client.AuthToken = authToken

	client := &Client{
		Config: config,
	}

	return client, nil
}

func (cfg *Configuration) SetBackoffConfig(backoffConf *BackoffConfig) {
	cfg.ZDX.Client.RateLimit.BackoffConf = backoffConf
}

// getHTTPClient sets up the retryable HTTP client with backoff and retry policies.
func getHTTPClient(l logger.Logger, rateLimiter *rl.RateLimiter, cfg *Configuration) *http.Client {
	retryableClient := retryablehttp.NewClient()

	// Set retry settings
	retryableClient.RetryWaitMin = cfg.ZDX.Client.RateLimit.RetryWaitMin
	retryableClient.RetryWaitMax = cfg.ZDX.Client.RateLimit.RetryWaitMax

	retryableClient.RetryMax = int(cfg.ZDX.Client.RateLimit.MaxRetries)
	if retryableClient.RetryMax == 0 {
		retryableClient.RetryMax = math.MaxInt32
	}

	// Backoff logic with rate limit headers
	retryableClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		if resp != nil {
			if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
				retryAfter := getRetryAfter(resp, l)
				if retryAfter > 0 {
					return retryAfter
				}
			}

			if resp.Request != nil {
				wait, delay := rateLimiter.Wait(resp.Request.Method)
				if wait {
					return delay
				}
				return 0
			}
		}

		// Default exponential backoff
		multiplier := math.Pow(2, float64(attemptNum)) * float64(min)
		sleep := time.Duration(multiplier)
		if float64(sleep) != multiplier || sleep > max {
			sleep = max
		}
		return sleep
	}

	retryableClient.CheckRetry = checkRetry
	retryableClient.Logger = l

	// Set request timeout
	retryableClient.HTTPClient.Timeout = cfg.ZDX.Client.RequestTimeout
	if retryableClient.HTTPClient.Timeout == 0 {
		retryableClient.HTTPClient.Timeout = time.Second * 60
	}

	// Configure proxy settings
	proxyFunc := http.ProxyFromEnvironment
	if cfg.ZDX.Client.Proxy.Host != "" {
		proxyURLString := fmt.Sprintf("http://%s:%d", cfg.ZDX.Client.Proxy.Host, cfg.ZDX.Client.Proxy.Port)
		if cfg.ZDX.Client.Proxy.Username != "" && cfg.ZDX.Client.Proxy.Password != "" {
			proxyAuth := url.UserPassword(cfg.ZDX.Client.Proxy.Username, cfg.ZDX.Client.Proxy.Password)
			proxyURLString = fmt.Sprintf("http://%s@%s:%d", proxyAuth.String(), cfg.ZDX.Client.Proxy.Host, cfg.ZDX.Client.Proxy.Port)
		}

		proxyURL, err := url.Parse(proxyURLString)
		if err == nil {
			proxyFunc = http.ProxyURL(proxyURL)
		} else {
			l.Printf("[ERROR] Invalid proxy URL: %v", err)
		}
	}

	transport := &http.Transport{
		Proxy:               proxyFunc,
		MaxIdleConnsPerHost: maxIdleConnections,
	}

	if cfg.ZDX.Testing.DisableHttpsCheck {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: false}
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
	// Extract rate limit headers
	rateLimitRemaining := resp.Header.Get("RateLimit-Remaining")
	rateLimitReset := resp.Header.Get("RateLimit-Reset")

	l.Printf("[INFO] RateLimit-Remaining: %s", rateLimitRemaining)
	l.Printf("[INFO] RateLimit-Reset: %s", rateLimitReset)

	// Parse RateLimit-Remaining
	remaining, err := strconv.Atoi(rateLimitRemaining)
	if err != nil {
		l.Printf("[WARN] Failed to parse RateLimit-Remaining header: %v", err)
		remaining = 1 // Assume remaining is 1 to prevent immediate retries
	}

	// If remaining requests are 0, calculate sleep time from RateLimit-Reset
	if remaining == 0 && rateLimitReset != "" {
		resetTime, err := strconv.ParseInt(rateLimitReset, 10, 64)
		if err != nil {
			l.Printf("[WARN] Failed to parse RateLimit-Reset header: %v", err)
		} else {
			currentTime := time.Now().Unix()
			sleepTime := resetTime - currentTime
			if sleepTime > 0 {
				l.Printf("[INFO] Rate limit reached. Retrying in %d seconds.", sleepTime)
				return time.Duration(sleepTime) * time.Second
			}
		}
	}

	// Fallback to exponential backoff if headers are insufficient
	l.Printf("[INFO] Falling back to default retry delay.")
	return time.Second * 5
}

// getRetryOnStatusCodes return a list of http status codes we want to apply retry on.
// return empty slice to enable retry on all connection & server errors.
// or return []int{429}  to retry on only TooManyRequests error
func getRetryOnStatusCodes() []int {
	return []int{http.StatusTooManyRequests}
}

// Used to make http client retry on provided list of response status codes
func checkRetry(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	if resp != nil && containsInt(getRetryOnStatusCodes(), resp.StatusCode) {
		return true, nil
	}
	return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
}

type ApiErr struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Authenticate(ctx context.Context, cfg *Configuration, logger logger.Logger) (*AuthToken, error) {
	cfg.Lock()
	defer cfg.Unlock()

	if cfg.ZDX.Client.AuthToken == nil || cfg.ZDX.Client.AuthToken.AccessToken == "" || utils.IsTokenExpired(cfg.ZDX.Client.AuthToken.AccessToken) {
		if cfg.ZDX.Client.ZDXAPIKeyID == "" || cfg.ZDX.Client.ZDXAPISecret == "" {
			logger.Printf("[ERROR] No client credentials provided. Set %s and %s environment variables or use ConfigSetter.", ZDX_API_KEY_ID, ZDX_API_SECRET)
			return nil, fmt.Errorf("missing client credentials: %s and/or %s", ZDX_API_KEY_ID, ZDX_API_SECRET)
		}

		maskedAPIKeyID := maskAPIKeyID(cfg.ZDX.Client.ZDXAPIKeyID)
		currTimestamp := time.Now().Unix()
		authReq := AuthRequest{
			Timestamp:    currTimestamp,
			APIKeyID:     cfg.ZDX.Client.ZDXAPIKeyID,
			APIKeySecret: generateHash(cfg.ZDX.Client.ZDXAPISecret, currTimestamp),
		}

		data, _ := json.Marshal(authReq)
		url := cfg.BaseURL.String() + "/v1/oauth/token"

		attempts := 0
		maxAttempts := 5
		backoffFactor := 1
		additionalDelay := 5 * time.Second // Optional delay for retries

		for attempts < maxAttempts {
			req, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
			if err != nil {
				logger.Printf("[ERROR] Failed to create request for user %s=%s, err: %v\n", ZDX_API_KEY_ID, maskedAPIKeyID, err)
				return nil, fmt.Errorf("[ERROR] Failed to create request for user %s=%s, err: %v", ZDX_API_KEY_ID, maskedAPIKeyID, err)
			}

			req.Header.Add("Content-Type", contentTypeJSON)
			if cfg.UserAgent != "" {
				req.Header.Add("User-Agent", cfg.UserAgent)
			}

			resp, err := cfg.HTTPClient.Do(req)
			if err != nil {
				logger.Printf("[ERROR] Failed to sign in the user %s=%s, err: %v\n", ZDX_API_KEY_ID, maskedAPIKeyID, err)
				return nil, fmt.Errorf("[ERROR] Failed to sign in the user %s=%s, err: %v", ZDX_API_KEY_ID, maskedAPIKeyID, err)
			}
			defer resp.Body.Close()

			// Read the response body
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Printf("[ERROR] Failed to read response for user %s=%s, err: %v\n", ZDX_API_KEY_ID, maskedAPIKeyID, err)
				return nil, fmt.Errorf("[ERROR] Failed to read response for user %s=%s, err: %v", ZDX_API_KEY_ID, maskedAPIKeyID, err)
			}

			// Handle rate-limited responses
			if resp.StatusCode == http.StatusTooManyRequests {
				rateLimitReset := resp.Header.Get("RateLimit-Reset")
				var sleepTime time.Duration
				if rateLimitReset != "" {
					resetTime, err := strconv.ParseInt(rateLimitReset, 10, 64)
					if err == nil {
						sleepTime = time.Until(time.Unix(resetTime, 0))
					}
				}

				if sleepTime <= 0 {
					// Fallback to exponential backoff
					sleepTime = time.Duration(backoffFactor) * time.Second
					backoffFactor *= 2
				}

				logger.Printf("[WARN] Rate limit exceeded. Retrying in %s (attempt %d/%d)", sleepTime, attempts+1, maxAttempts)
				time.Sleep(sleepTime + additionalDelay)
				attempts++
				continue
			}

			// Handle non-2xx responses
			if resp.StatusCode >= 300 {
				logger.Printf("[ERROR] Failed to sign in the user %s=%s, got HTTP status: %d, response body: %s, url: %s\n", ZDX_API_KEY_ID, maskedAPIKeyID, resp.StatusCode, respBody, url)
				return nil, fmt.Errorf("[ERROR] Failed to sign in the user %s=%s, got HTTP status: %d, response body: %s, url: %s", ZDX_API_KEY_ID, maskedAPIKeyID, resp.StatusCode, respBody, url)
			}

			// Parse the successful authentication response
			var authToken AuthToken
			err = json.Unmarshal(respBody, &authToken)
			if err != nil {
				logger.Printf("[ERROR] Failed to parse response for user %s=%s, err: %v\n", ZDX_API_KEY_ID, maskedAPIKeyID, err)
				return nil, fmt.Errorf("[ERROR] Failed to parse response for user %s=%s, err: %v", ZDX_API_KEY_ID, maskedAPIKeyID, err)
			}

			// Store the token in the configuration
			cfg.ZDX.Client.AuthToken = &authToken
			return &authToken, nil
		}

		// Exceeded maximum attempts
		logger.Printf("[ERROR] Rate limit retries exceeded for user %s=%s\n", ZDX_API_KEY_ID, maskedAPIKeyID)
		return nil, fmt.Errorf("[ERROR] Rate limit retries exceeded for user %s=%s", ZDX_API_KEY_ID, maskedAPIKeyID)
	}

	// Return the existing valid token
	return cfg.ZDX.Client.AuthToken, nil
}

func (client *Client) NewRequestDo(ctx context.Context, method, urlStr string, options, body, v interface{}) (*http.Response, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil: ensure the client is properly initialized")
	}
	return client.newRequestDoCustom(ctx, method, urlStr, options, body, v, client.Config)
}

func (client *Client) newRequestDoCustom(ctx context.Context, method, urlStr string, options, body, v interface{}, config *Configuration) (*http.Response, error) {
	// Authenticate and log errors
	if _, err := Authenticate(ctx, config, config.Logger); err != nil {
		client.Config.Logger.Printf("[ERROR] Authentication failed: %v", err)
		return nil, err
	}

	// Create the request
	req, err := client.newRequest(method, urlStr, options, body, client.Config)
	if err != nil {
		client.Config.Logger.Printf("[ERROR] Failed to create request: %v", err)
		return nil, err
	}

	req = req.WithContext(ctx)

	reqID := uuid.NewString()
	start := time.Now()
	logger.LogRequest(client.Config.Logger, req, reqID, nil, true)

	// Perform the request
	resp, err := client.do(req, v, start, reqID)
	if err != nil {
		client.Config.Logger.Printf("[ERROR] Request failed: %v", err)
		return resp, err
	}

	// Safeguard against nil response
	if resp == nil {
		client.Config.Logger.Printf("[ERROR] Received nil response from API.")
		return nil, fmt.Errorf("received nil response from API")
	}

	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	// Handle unauthorized or forbidden responses
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		client.Config.Logger.Printf("[WARN] Unauthorized or forbidden response. Retrying authentication.")

		// Re-authenticate and log errors
		if _, err := Authenticate(ctx, config, config.Logger); err != nil {
			client.Config.Logger.Printf("[ERROR] Re-authentication failed: %v", err)
			return nil, err
		}

		// Retry the original request
		resp, err = client.do(req, v, start, reqID)
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Request failed after re-authentication: %v", err)
			return nil, err
		}

		// Ensure the response is not nil before returning
		if resp == nil {
			client.Config.Logger.Printf("[ERROR] Nil response received after re-authentication.")
			return nil, fmt.Errorf("nil response received after re-authentication")
		}
	}

	return resp, err
}

func generateHash(apiSecret string, currTimestamp int64) string {
	currTimestampStr := strconv.FormatInt(currTimestamp, 10)
	hash := sha256.New()
	hash.Write([]byte(apiSecret + ":" + currTimestampStr))
	return hex.EncodeToString(hash.Sum(nil))
}

func maskAPIKeyID(apiKeyID string) string {
	if len(apiKeyID) <= 4 {
		return "****"
	}
	return apiKeyID[:2] + strings.Repeat("*", len(apiKeyID)-4) + apiKeyID[len(apiKeyID)-2:]
}

// Generating the Http request
func (client *Client) newRequest(method, urlPath string, options, body interface{}, cfg *Configuration) (*http.Request, error) {
	if cfg.ZDX.Client.AuthToken == nil || cfg.ZDX.Client.AuthToken.AccessToken == "" {
		maskedAPIKeyID := maskAPIKeyID(cfg.ZDX.Client.ZDXAPIKeyID)
		client.Config.Logger.Printf("[ERROR] Failed to sign in the user %s=%s\n", ZDX_API_KEY_ID, maskedAPIKeyID)
		return nil, fmt.Errorf("failed to sign in the user %s=%s", ZDX_API_KEY_ID, maskedAPIKeyID)
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
	u.Path = u.Path + unescaped

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

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.ZDX.Client.AuthToken.AccessToken))
	req.Header.Add("Content-Type", contentTypeJSON)

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
	// Read and log the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(respBody)) // Reset the response body

	logger.LogResponse(client.Config.Logger, resp, start, reqID)
	logger.WriteLog(client.Config.Logger, "Response Body: %s", string(respBody)) // Log the response body separately

	if err := errorx.CheckErrorInResponse(resp, err); err != nil {
		return resp, err
	}

	if v != nil {
		// Reset the response body again for unmarshalling
		resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
		if err := decodeJSON(resp, v); err != nil {
			return resp, err
		}
	}
	unescapeHTML(v)
	return resp, nil
}

func decodeJSON(res *http.Response, v interface{}) error {
	return json.NewDecoder(res.Body).Decode(&v)
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
	json.Unmarshal(data, entity)
}
