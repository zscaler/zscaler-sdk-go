package zscaler

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
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"github.com/zscaler/zscaler-sdk-go/v3/cache"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
	rl "github.com/zscaler/zscaler-sdk-go/v3/ratelimiter"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa"
	ztw "github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw"
	"gopkg.in/yaml.v3"
)

type contextKey string

func (c contextKey) String() string {
	return "zscaler " + string(c)
}

var (
	// ContextAccessToken takes a string OAuth2 access token as authentication for the request.
	ContextAccessToken = contextKey("access_token")
)

const (
	VERSION               = "3.8.11"
	ZSCALER_CLIENT_ID     = "ZSCALER_CLIENT_ID"
	ZSCALER_CLIENT_SECRET = "ZSCALER_CLIENT_SECRET"
	ZSCALER_VANITY_DOMAIN = "ZSCALER_VANITY_DOMAIN"
	ZSCALER_PRIVATE_KEY   = "ZSCALER_PRIVATE_KEY"
	ZSCALER_CLOUD         = "ZSCALER_CLOUD"
)

// AuthToken represents the OAuth2 authentication token and its expiration time.
type AuthToken struct {
	TokenType   string      `json:"token_type"`
	AccessToken string      `json:"access_token"`
	ExpiresIn   json.Number `json:"expires_in"` // <- FIXED
	Expiry      time.Time
}

// Legacy struct holds and instance of each legacy API client to support backwards compatibility
type LegacyClient struct {
	ZiaClient *zia.Client
	ZtwClient *ztw.Client
	ZpaClient *zpa.Client
	ZccClient *zcc.Client
	ZdxClient *zdx.Client
}

// Configuration struct holds the config for ZIA, ZPA, and common fields like HTTPClient and AuthToken.
type Configuration struct {
	Logger         logger.Logger
	HTTPClient     *http.Client
	ZPAHTTPClient  *http.Client
	ZIAHTTPClient  *http.Client
	ZTWHTTPClient  *http.Client
	ZCCHTTPClient  *http.Client
	ZDXHTTPClient  *http.Client
	DefaultHeader  map[string]string `json:"defaultHeader,omitempty"`
	UserAgent      string            `json:"userAgent,omitempty"`
	Debug          bool              `json:"debug,omitempty"`
	UserAgentExtra string
	Context        context.Context
	Zscaler        struct {
		Client struct {
			ClientID      string     `yaml:"clientId" envconfig:"ZSCALER_CLIENT_ID"`
			ClientSecret  string     `yaml:"clientSecret" envconfig:"ZSCALER_CLIENT_SECRET"`
			VanityDomain  string     `yaml:"vanityDomain" envconfig:"ZSCALER_VANITY_DOMAIN"`
			Cloud         string     `yaml:"cloud" envconfig:"ZSCALER_CLOUD"`
			CustomerID    string     `yaml:"customerId" envconfig:"ZPA_CUSTOMER_ID"`
			MicrotenantID string     `yaml:"microtenantId" envconfig:"ZPA_MICROTENANT_ID"`
			PrivateKey    []byte     `yaml:"privateKey" envconfig:"ZSCALER_PRIVATE_KEY"`
			PartnerID     string     `yaml:"partnerId" envconfig:"ZSCALER_PARTNER_ID"`
			AuthToken     *AuthToken `yaml:"authToken"`
			AccessToken   *AuthToken `yaml:"accessToken"`
			SandboxToken  string     `yaml:"sandboxToken" envconfig:"ZSCALER_SANDBOX_TOKEN"`
			SandboxCloud  string     `yaml:"sandboxCloud" envconfig:"ZSCALER_SANDBOX_CLOUD"`
			Cache         struct {
				Enabled               bool          `yaml:"enabled" envconfig:"ZSCALER_CLIENT_CACHE_ENABLED"`
				DefaultTtl            time.Duration `yaml:"defaultTtl" envconfig:"ZSCALER_CLIENT_CACHE_DEFAULT_TTL"`
				DefaultTti            time.Duration `yaml:"defaultTti" envconfig:"ZSCALER_CLIENT_CACHE_DEFAULT_TTI"`
				DefaultCacheMaxSizeMB int64         `yaml:"defaultSize" envconfig:"ZSCALER_CLIENT_CACHE_DEFAULT_SIZE"`
			} `yaml:"cache"`
			Proxy struct {
				Port     int32  `yaml:"port" envconfig:"ZSCALER_CLIENT_PROXY_PORT"`
				Host     string `yaml:"host" envconfig:"ZSCALER_CLIENT_PROXY_HOST"`
				Username string `yaml:"username" envconfig:"ZSCALER_CLIENT_PROXY_USERNAME"`
				Password string `yaml:"password" envconfig:"ZSCALER_CLIENT_PROXY_PASSWORD"`
			} `yaml:"proxy"`
			RequestTimeout time.Duration `yaml:"requestTimeout" envconfig:"ZSCALER_CLIENT_REQUEST_TIMEOUT"`
			RateLimit      struct {
				MaxRetries                int32         `yaml:"maxRetries" envconfig:"ZSCALER_CLIENT_RATE_LIMIT_MAX_RETRIES"`
				RetryWaitMin              time.Duration `yaml:"minWait" envconfig:"ZSCALER_CLIENT_RATE_LIMIT_MIN_WAIT"`
				RetryWaitMax              time.Duration `yaml:"maxWait" envconfig:"ZSCALER_CLIENT_RATE_LIMIT_MAX_WAIT"`
				RetryRemainingThreshold   int32         `yaml:"remainingThreshold" envconfig:"ZSCALER_CLIENT_REMAINING_THRESHOLD"`
				MaxSessionNotValidRetries int32         `yaml:"maxSessionNotValidRetries" envconfig:"ZSCALER_CLIENT_MAX_SESSION_NOT_VALID_RETRIES"`
			} `yaml:"rateLimit"`
		} `yaml:"client"`
		Testing struct {
			DisableHttpsCheck bool `yaml:"disableHttpsCheck" envconfig:"ZSCALER_TESTING_DISABLE_HTTPS_CHECK"`
		} `yaml:"testing"`
	} `yaml:"zscaler"`
	PrivateKeySigner jose.Signer
	CacheManager     cache.Cache
	UseLegacyClient  bool `yaml:"useLegacyClient" envconfig:"ZSCALER_USE_LEGACY_CLIENT"`
	LegacyClient     *LegacyClient
}

// NewConfiguration is the main configuration function, implementing the ConfigSetter pattern.
func NewConfiguration(conf ...ConfigSetter) (*Configuration, error) {
	logger := logger.GetDefaultLogger(loggerPrefix)
	cfg := &Configuration{
		DefaultHeader: make(map[string]string),
		Logger:        logger,
		UserAgent:     fmt.Sprintf("zscaler-sdk-go/%s golang/%s %s/%s", VERSION, runtime.Version(), runtime.GOOS, runtime.GOARCH),
		Debug:         false,
		Context:       context.Background(), // Set default context
	}

	cfg.Zscaler.Client.RateLimit.MaxRetries = MaxNumOfRetries
	cfg.Zscaler.Client.RateLimit.RetryWaitMax = time.Second * time.Duration(RetryWaitMaxSeconds)
	cfg.Zscaler.Client.RateLimit.RetryWaitMin = time.Second * time.Duration(RetryWaitMinSeconds)
	cfg.Zscaler.Client.RateLimit.MaxSessionNotValidRetries = 3 // Default to 3 consecutive SESSION_NOT_VALID retries

	cfg.Zscaler.Client.RequestTimeout = time.Duration(requestTimeout) * time.Second

	// Initialize cache
	if cfg.Zscaler.Client.Cache.DefaultTtl == 0 {
		cfg.Zscaler.Client.Cache.DefaultTtl = time.Minute * 10
	}

	if cfg.Zscaler.Client.Cache.DefaultTti == 0 {
		cfg.Zscaler.Client.Cache.DefaultTti = time.Minute * 8
	}

	cfg.CacheManager = newCache(cfg)

	cfg.Zscaler.Testing.DisableHttpsCheck = false

	cfg = readConfigFromSystem(*cfg)
	cfg = readConfigFromEnvironment(*cfg)

	setHttpClients(cfg)

	// Apply each ConfigSetter function.
	for _, confSetter := range conf {
		confSetter(cfg)
	}

	// Recheck and adjust defaults after setters are applied.
	if cfg.Zscaler.Client.RateLimit.MaxRetries == 0 {
		cfg.Zscaler.Client.RateLimit.MaxRetries = 4 // Default to 4 if user set it to zero.
	}

	if cfg.Zscaler.Client.RequestTimeout == 0 {
		cfg.Zscaler.Client.RequestTimeout = 60 * time.Second // Default to 60 seconds if user set it to zero.
	}

	// UserAgentExtra gets added if provided.
	if cfg.UserAgentExtra != "" {
		cfg.UserAgent = fmt.Sprintf("%s %s", cfg.UserAgent, cfg.UserAgentExtra)
	}

	ctx := context.WithValue(
		context.Background(),
		ContextAccessToken,
		cfg.Zscaler.Client.AuthToken.AccessToken,
	)
	cfg.Context = ctx

	return cfg, nil
}

func setHttpClients(cfg *Configuration) {
	// ZIA-specific rate limits with hourly tracking:
	// Per-second: GET: 2/sec, POST/PUT: 1/sec, DELETE: 1/sec
	// Hourly: GET: 1000/hr, POST/PUT: 1000/hr (combined), DELETE: 400/hr
	// Using conservative limits to stay well below API thresholds
	// Per-second: 20 requests per 10s (2/sec) for GET, 10 per 10s (1/sec) for POST/PUT/DELETE
	// Hourly: 950/hr for GET (buffer of 50), 950/hr for POST/PUT (buffer of 50), 380/hr for DELETE (buffer of 20)
	ziaRateLimiter := rl.NewRateLimiterWithHourly(
		20, 10, // GET: 20 per 10 seconds, POST/PUT/DELETE: 10 per 10 seconds
		10, 61, // GET frequency: 10 seconds, DELETE frequency: 61 seconds (+1 buffer)
		950, 950, 380, // Hourly: GET: 950, POST/PUT: 950, DELETE: 380 (with safety buffers)
	)

	// ZTW uses same limits as ZIA
	ztwRateLimiter := rl.NewRateLimiterWithHourly(
		20, 10,
		10, 61,
		950, 950, 380,
	)

	// ZPA-specific rate limits:
	zpaRateLimiter := rl.NewRateLimiter(20, 10, 10, 10) // GET: 20 per 10s, POST/PUT/DELETE: 10 per 10s

	// ZCC-specific rate limits:
	zccRateLimiter := rl.NewRateLimiter(100, 3, 3600, 86400) // General: 100 per hour, downloadDevices: 3 per day

	// ZDX-specific rate limits:
	zdxRateLimiter := rl.NewRateLimiter(100, 3, 3600, 86400) // General: 100 per hour, downloadDevices: 3 per day

	// Default case for unknown or unhandled services
	defaultRateLimiter := rl.NewRateLimiter(2, 1, 1, 1) // Default limits

	// Pass the config to getHTTPClient so it can access proxy settings
	cfg.HTTPClient = getHTTPClient(cfg.Logger, defaultRateLimiter, cfg)
	cfg.ZIAHTTPClient = getHTTPClient(cfg.Logger, ziaRateLimiter, cfg)
	cfg.ZTWHTTPClient = getHTTPClient(cfg.Logger, ztwRateLimiter, cfg)
	cfg.ZPAHTTPClient = getHTTPClient(cfg.Logger, zpaRateLimiter, cfg)
	cfg.ZCCHTTPClient = getHTTPClient(cfg.Logger, zccRateLimiter, cfg)
	cfg.ZDXHTTPClient = getHTTPClient(cfg.Logger, zdxRateLimiter, cfg)
}

// Authenticate performs OAuth2 authentication and retrieves an AuthToken.
func Authenticate(ctx context.Context, cfg *Configuration, l logger.Logger) (*AuthToken, error) {
	creds := cfg.Zscaler.Client

	if creds.ClientID == "" || (creds.ClientSecret == "" && len(creds.PrivateKey) == 0) {
		return nil, errors.New("no client credentials were provided")
	}

	// If private key is provided, use JWT-based authentication.
	if len(creds.PrivateKey) > 0 {
		return authenticateWithCert(cfg)
	}

	// Determine the OAuth2 provider URL based on the cloud parameter.
	var authUrl string
	if creds.Cloud == "" || strings.EqualFold(creds.Cloud, "PRODUCTION") {
		authUrl = fmt.Sprintf("https://%s.zslogin.net/oauth2/v1/token", creds.VanityDomain)
	} else {
		authUrl = fmt.Sprintf("https://%s.zslogin%s.net/oauth2/v1/token", creds.VanityDomain, strings.ToLower(creds.Cloud))
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_secret", creds.ClientSecret)
	data.Set("client_id", creds.ClientID)
	data.Set("audience", "https://api.zscaler.com")

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "POST", authUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to sign in the user %s: %v", creds.ClientID, err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", cfg.UserAgent)
	// start := time.Now()
	reqID := uuid.NewString()
	logger.LogRequest(l, req, reqID, nil, false)
	resp, err := cfg.HTTPClient.Do(req)
	// logger.LogResponse(l, resp, start, reqID)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Failed to sign in the user %s, err: %v", creds.ClientID, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Failed to sign in the user %s, err: %v", creds.ClientID, err)
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("[ERROR] Failed to sign in the user %s, got http status: %d, response body: %s", creds.ClientID, resp.StatusCode, respBody)
	}

	var token AuthToken
	if err := json.Unmarshal(respBody, &token); err != nil {
		return nil, fmt.Errorf("[ERROR] Failed to sign in: %v", err)
	}

	secondsStr := token.ExpiresIn.String()
	seconds, err := strconv.ParseInt(secondsStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] invalid expires_in value: %v", err)
	}
	token.Expiry = time.Now().Add(time.Duration(seconds) * time.Second)
	cfg.Logger.Printf("[DEBUG] parsed expires_in=%d seconds, token expiry set to: %s", seconds, token.Expiry.Format(time.RFC3339))
	return &token, nil
}

// authenticateWithCert performs JWT-based authentication using a private key.
func authenticateWithCert(cfg *Configuration) (*AuthToken, error) {
	creds := cfg.Zscaler.Client

	if creds.ClientID == "" || len(creds.PrivateKey) == 0 {
		return nil, errors.New("client ID or private key is missing")
	}

	// Create the JWT payload.
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(creds.PrivateKey))
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %v", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": creds.ClientID,
		"sub": creds.ClientID,
		"aud": "https://api.zscaler.com",
		"exp": time.Now().Add(10 * time.Minute).Unix(),
	})

	assertion, err := token.SignedString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("error signing JWT: %v", err)
	}

	formData := url.Values{
		"grant_type":            {"client_credentials"},
		"client_id":             {creds.ClientID},
		"client_assertion":      {assertion},
		"client_assertion_type": {"urn:ietf:params:oauth:client-assertion-type:jwt-bearer"},
		"audience":              {"https://api.zscaler.com"},
	}

	// Determine the OAuth2 provider URL based on the cloud parameter.
	var authUrl string
	if creds.Cloud == "" || strings.EqualFold(creds.Cloud, "PRODUCTION") {
		authUrl = fmt.Sprintf("https://%s.zslogin.net/oauth2/v1/token", creds.VanityDomain)
	} else {
		authUrl = fmt.Sprintf("https://%s.zslogin%s.net/oauth2/v1/token", creds.VanityDomain, strings.ToLower(creds.Cloud))
	}

	// Make the POST request.
	resp, err := cfg.HTTPClient.PostForm(authUrl, formData)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("auth error: %v", string(body))
	}
	// Parse the response.
	var tokenResponse AuthToken
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &tokenResponse, nil
}

// getServiceHTTPClient returns the appropriate http client for the current service
func (client *Client) getServiceHTTPClient(endpoint string) *http.Client {
	service, err := detectServiceType(endpoint)
	if err != nil {
		return client.oauth2Credentials.HTTPClient
	}
	switch service {
	case "zpa":
		return client.oauth2Credentials.ZPAHTTPClient
	case "zia":
		return client.oauth2Credentials.ZIAHTTPClient
	case "ztw":
		return client.oauth2Credentials.ZTWHTTPClient
	case "zcc":
		return client.oauth2Credentials.ZCCHTTPClient
	case "zdx":
		return client.oauth2Credentials.ZDXHTTPClient
	case "admin":
		return client.oauth2Credentials.HTTPClient // Use default client for admin endpoints
	default:
		return client.oauth2Credentials.HTTPClient
	}
}

func detectServiceType(endpoint string) (string, error) {
	path := strings.TrimPrefix(endpoint, "/")
	// Detect the service type based on the endpoint prefix
	if strings.HasPrefix(path, "zia") || strings.HasPrefix(path, "zscsb") {
		return "zia", nil
	} else if strings.HasPrefix(path, "ztw") {
		return "ztw", nil
	} else if strings.HasPrefix(path, "zpa") {
		return "zpa", nil
	} else if strings.HasPrefix(endpoint, "/zcc") {
		return "zcc", nil
	} else if strings.HasPrefix(endpoint, "/zdx") {
		return "zdx", nil
	} else if strings.HasPrefix(endpoint, "/admin") {
		return "admin", nil
	}
	return "", fmt.Errorf("unsupported service")
}

// GetAPIBaseURL gets the appropriate base url based on the cloud and sandbox mode.
func GetAPIBaseURL(cloud string) string {
	baseURL := "https://api.zsapi.net"
	if cloud != "" && !strings.EqualFold(cloud, "PRODUCTION") {
		baseURL = fmt.Sprintf("https://api.%s.zsapi.net", strings.ToLower(cloud))
	}

	return baseURL
}

func readConfigFromFile(location string, c Configuration) (*Configuration, error) {
	yamlConfig, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlConfig, &c)
	if err != nil {
		return nil, err
	}
	return &c, err
}

func readConfigFromSystem(c Configuration) *Configuration {
	currUser, err := user.Current()
	if err != nil {
		return &c
	}
	if currUser.HomeDir == "" {
		return &c
	}
	conf, err := readConfigFromFile(currUser.HomeDir+"/.zscaler/zscaler.yaml", c)
	if err != nil {
		return &c
	}
	return conf
}

func readConfigFromEnvironment(c Configuration) *Configuration {
	err := envconfig.Process("zscaler", &c)
	if err != nil {
		fmt.Println("error parsing")
		return &c
	}
	return &c
}

// AddDefaultHeader adds a new HTTP header to the default header in the request
func (c *Configuration) AddDefaultHeader(key string, value string) {
	c.DefaultHeader[key] = value
}

type ConfigSetter func(*Configuration)

// ConfigSetter type defines a function that modifies a Config struct.
// WithClientID sets the ClientID in the Config.
func WithClientID(clientID string) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Client.ClientID = clientID
	}
}

// WithClientSecret sets the ClientSecret in the Config.
func WithClientSecret(clientSecret string) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Client.ClientSecret = clientSecret
	}
}

// WithOauth2ProviderUrl sets the Oauth2ProviderUrl in the Config.
func WithVanityDomain(domain string) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Client.VanityDomain = domain
	}
}

func WithZscalerCloud(cloud string) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Client.Cloud = cloud
	}
}

// WithSandboxToken is a ConfigSetter that sets the Sandbox token for the Zscaler Client.
func WithSandboxToken(token string) ConfigSetter {
	return func(cfg *Configuration) {
		cfg.Zscaler.Client.SandboxToken = token
	}
}

func WithSandboxCloud(sandboxCloud string) ConfigSetter {
	return func(cfg *Configuration) {
		cfg.Zscaler.Client.SandboxCloud = sandboxCloud
	}
}

func WithZPACustomerID(customerID string) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Client.CustomerID = customerID
	}
}

func WithZPAMicrotenantID(microtenantID string) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Client.MicrotenantID = microtenantID
	}
}

func WithPartnerID(partnerID string) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Client.PartnerID = partnerID
	}
}

func WithCache(cache bool) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Client.Cache.Enabled = cache
	}
}

func WithCacheManager(cacheManager cache.Cache) ConfigSetter {
	return func(c *Configuration) {
		c.CacheManager = cacheManager
	}
}

func newCache(c *Configuration) cache.Cache {
	if !c.Zscaler.Client.Cache.Enabled {
		return cache.NewNopCache()
	}
	cche, err := cache.NewCache(time.Duration(c.Zscaler.Client.Cache.DefaultTtl), time.Duration(c.Zscaler.Client.Cache.DefaultTti), int(c.Zscaler.Client.Cache.DefaultCacheMaxSizeMB))
	if err != nil {
		return cache.NewNopCache()
	}
	return cche
}

func WithCacheTtl(i time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Client.Cache.DefaultTtl = i
		c.CacheManager = newCache(c)
	}
}

func WithCacheMaxSizeMB(size int64) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Client.Cache.DefaultCacheMaxSizeMB = size
		c.CacheManager = newCache(c)
	}
}

func WithCacheTti(i time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Client.Cache.DefaultTti = i
		c.CacheManager = newCache(c)
	}
}

// WithHttpClient sets the HttpClient in the Config.
func WithHttpClientPtr(httpClient *http.Client) ConfigSetter {
	return func(c *Configuration) {
		c.HTTPClient = httpClient
	}
}

func WithProxyPort(i int32) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Client.Proxy.Port = i
	}
}

func WithProxyHost(host string) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Client.Proxy.Host = host
	}
}

func WithProxyUsername(username string) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Client.Proxy.Username = username
	}
}

func WithProxyPassword(pass string) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Client.Proxy.Password = pass
	}
}

func WithTestingDisableHttpsCheck(httpsCheck bool) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Testing.DisableHttpsCheck = httpsCheck
	}
}

func WithRequestTimeout(requestTimeout time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Client.RequestTimeout = requestTimeout
		setHttpClients(c)
	}
}

func WithRateLimitMaxRetries(maxRetries int32) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Client.RateLimit.MaxRetries = maxRetries
		setHttpClients(c)
	}
}

func WithRateLimitMaxWait(maxWait time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Client.RateLimit.RetryWaitMax = maxWait
		setHttpClients(c)
	}
}

func WithRateLimitMinWait(minWait time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Client.RateLimit.RetryWaitMin = minWait
		setHttpClients(c)
	}
}

func WithRateLimitRemainingThreshold(threshold int32) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Client.RateLimit.RetryRemainingThreshold = threshold
		setHttpClients(c)
	}
}

func WithRateLimitMaxSessionNotValidRetries(maxRetries int32) ConfigSetter {
	return func(c *Configuration) {
		c.Zscaler.Client.RateLimit.MaxSessionNotValidRetries = maxRetries
	}
}

// WithUserAgent sets the UserAgent in the Config.
func WithUserAgentExtra(userAgent string) ConfigSetter {
	return func(c *Configuration) {
		c.UserAgentExtra = userAgent
	}
}

func WithDebug(debug bool) ConfigSetter {
	return func(c *Configuration) {
		c.Debug = debug
		// Automatically set the environment variables if debug is enabled
		if debug {
			_ = os.Setenv("ZSCALER_SDK_LOG", "true")
			_ = os.Setenv("ZSCALER_SDK_VERBOSE", "true")
		}
	}
}

// WithPrivateKey sets private key, privateKey can be the raw key value or a path to the pem file.
func WithPrivateKey(privateKey string) ConfigSetter {
	return func(c *Configuration) {
		if fileExists(privateKey) {
			content, err := os.ReadFile(privateKey)
			if err != nil {
				fmt.Printf("failed to read from provided private key file path: %v", err)
			}
			c.Zscaler.Client.PrivateKey = content
		} else {
			c.Zscaler.Client.PrivateKey = []byte(privateKey)
		}
	}
}

func WithPrivateKeySigner(signer jose.Signer) ConfigSetter {
	return func(c *Configuration) {
		c.PrivateKeySigner = signer
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) || errors.Is(err, syscall.ENAMETOOLONG) {
			return false
		}
		fmt.Println("can not get information about the file containing private key, using provided value as the key itself")
		return false
	}
	return !info.IsDir()
}

func WithLegacyClient(useLegacyClient bool) ConfigSetter {
	return func(c *Configuration) {
		c.UseLegacyClient = useLegacyClient
	}

}

func WithZiaLegacyClient(ziaClient *zia.Client) ConfigSetter {
	return func(c *Configuration) {
		if c.LegacyClient == nil {
			c.LegacyClient = &LegacyClient{}
		}
		c.LegacyClient.ZiaClient = ziaClient
	}
}

func WithZtwLegacyClient(ztwClient *ztw.Client) ConfigSetter {
	return func(c *Configuration) {
		if c.LegacyClient == nil {
			c.LegacyClient = &LegacyClient{}
		}
		c.LegacyClient.ZtwClient = ztwClient
	}
}

func WithZpaLegacyClient(zpaClient *zpa.Client) ConfigSetter {
	return func(c *Configuration) {
		if c.LegacyClient == nil {
			c.LegacyClient = &LegacyClient{}
		}
		c.LegacyClient.ZpaClient = zpaClient
	}
}

func WithZccLegacyClient(zccClient *zcc.Client) ConfigSetter {
	return func(c *Configuration) {
		if c.LegacyClient == nil {
			c.LegacyClient = &LegacyClient{}
		}
		c.LegacyClient.ZccClient = zccClient
	}
}

func WithZdxLegacyClient(zdxClient *zdx.Client) ConfigSetter {
	return func(c *Configuration) {
		if c.LegacyClient == nil {
			c.LegacyClient = &LegacyClient{}
		}
		c.LegacyClient.ZdxClient = zdxClient
	}
}
