package zpa

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/zscaler/zscaler-sdk-go/v3/cache"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
	rl "github.com/zscaler/zscaler-sdk-go/v3/ratelimiter"
	"gopkg.in/yaml.v3"
)

const (
	maxIdleConnections  int = 40
	requestTimeout      int = 60
	contentTypeJSON         = "application/json"
	MaxNumOfRetries         = 100
	RetryWaitMaxSeconds     = 20
	RetryWaitMinSeconds     = 5
	loggerPrefix            = "zpa-logger: "
)

const (
	VERSION           = "3.0.0"
	ZPA_CLIENT_ID     = "ZPA_CLIENT_ID"
	ZPA_CLIENT_SECRET = "ZPA_CLIENT_SECRET"
	ZPA_CUSTOMER_ID   = "ZPA_CUSTOMER_ID"
	ZPA_CLOUD         = "ZPA_CLOUD"
	defaultBaseURL    = "https://config.private.zscaler.com"
	betaBaseURL       = "https://config.zpabeta.net"
	zpaTwoBaseUrl     = "https://config.zpatwo.net"
	govBaseURL        = "https://config.zpagov.net"
	govUsBaseURL      = "https://config.zpagov.us"
	previewBaseUrl    = "https://config.zpapreview.net"
	devBaseUrl        = "https://public-api.dev.zpath.net"
	devAuthUrl        = "https://authn1.dev.zpath.net/authn/v1/oauth/token?grant_type=CLIENT_CREDENTIALS"
	qaBaseUrl         = "https://config.qa.zpath.net"
	qa2BaseUrl        = "https://pdx2-zpa-config.qa2.zpath.net"
)

type contextKey string

func (c contextKey) String() string {
	return "zscaler " + string(c)
}

var (
	// ContextAccessToken takes a string OAuth2 access token as authentication for the request.
	ContextAccessToken = contextKey("access_token")
)

type AuthToken struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"-"` // Skip direct JSON unmarshaling
	Expiry       time.Time
	RawExpiresIn string `json:"expires_in"`
}

type Configuration struct {
	sync.Mutex
	Logger         logger.Logger
	HTTPClient     *http.Client
	BaseURL        *url.URL
	DefaultHeader  map[string]string `json:"defaultHeader,omitempty"`
	UserAgent      string            `json:"userAgent,omitempty"`
	Debug          bool              `json:"debug,omitempty"`
	UserAgentExtra string
	Context        context.Context
	ZPA            struct {
		Client struct {
			ZPAClientID      string     `yaml:"clientId" envconfig:"ZPA_CLIENT_ID"`
			ZPAClientSecret  string     `yaml:"clientSecret" envconfig:"ZPA_CLIENT_SECRET"`
			ZPACustomerID    string     `yaml:"customerId" envconfig:"ZPA_CUSTOMER_ID"`
			ZPACloud         string     `yaml:"cloud" envconfig:"ZPA_CLOUD"`
			ZPAMicrotenantID string     `yaml:"microtenantId" envconfig:"ZPA_MICROTENANT_ID"`
			AuthToken        *AuthToken `yaml:"authToken"`
			AccessToken      *AuthToken `yaml:"accessToken"`
			Cache            struct {
				Enabled               bool          `yaml:"enabled" envconfig:"ZPA_CLIENT_CACHE_ENABLED"`
				DefaultTtl            time.Duration `yaml:"defaultTtl" envconfig:"ZPA_CLIENT_CACHE_DEFAULT_TTL"`
				DefaultTti            time.Duration `yaml:"defaultTti" envconfig:"ZPA_CLIENT_CACHE_DEFAULT_TTI"`
				DefaultCacheMaxSizeMB int64         `yaml:"defaultTti" envconfig:"ZPA_CLIENT_CACHE_DEFAULT_SIZE"`
			} `yaml:"cache"`
			Proxy struct {
				Port     int32  `yaml:"port" envconfig:"ZPA_CLIENT_PROXY_PORT"`
				Host     string `yaml:"host" envconfig:"ZPA_CLIENT_PROXY_HOST"`
				Username string `yaml:"username" envconfig:"ZPA_CLIENT_PROXY_USERNAME"`
				Password string `yaml:"password" envconfig:"ZPA_CLIENT_PROXY_PASSWORD"`
			} `yaml:"proxy"`
			RequestTimeout time.Duration `yaml:"requestTimeout" envconfig:"ZPA_CLIENT_REQUEST_TIMEOUT"`
			RateLimit      struct {
				MaxRetries   int32         `yaml:"maxRetries" envconfig:"ZPA_CLIENT_RATE_LIMIT_MAX_RETRIES"`
				RetryWaitMin time.Duration `yaml:"minWait" envconfig:"ZPA_CLIENT_RATE_LIMIT_MIN_WAIT"`
				RetryWaitMax time.Duration `yaml:"maxWait" envconfig:"ZPA_CLIENT_RATE_LIMIT_MAX_WAIT"`
			} `yaml:"rateLimit"`
		} `yaml:"client"`
		Testing struct {
			DisableHttpsCheck bool `yaml:"disableHttpsCheck" envconfig:"ZPA_TESTING_DISABLE_HTTPS_CHECK"`
		} `yaml:"testing"`
	} `yaml:"zpa"`
	CacheManager cache.Cache
}

func NewConfiguration(conf ...ConfigSetter) (*Configuration, error) {
	logger := logger.GetDefaultLogger(loggerPrefix)
	cfg := &Configuration{
		DefaultHeader: make(map[string]string),
		Logger:        logger,
		UserAgent:     fmt.Sprintf("zscaler-sdk-go/%s golang/%s %s/%s", VERSION, runtime.Version(), runtime.GOOS, runtime.GOARCH),
		Debug:         false,
		Context:       context.Background(),
	}

	// logger.Printf("[DEBUG] Initializing configuration with default values.")

	// Set default rate limit and request timeout values
	cfg.ZPA.Client.RateLimit.MaxRetries = MaxNumOfRetries
	cfg.ZPA.Client.RateLimit.RetryWaitMax = time.Second * time.Duration(RetryWaitMaxSeconds)
	cfg.ZPA.Client.RateLimit.RetryWaitMin = time.Second * time.Duration(RetryWaitMinSeconds)
	cfg.ZPA.Client.RequestTimeout = time.Duration(requestTimeout) * time.Second

	// Initialize cache with defaults
	if cfg.ZPA.Client.Cache.DefaultTtl == 0 {
		cfg.ZPA.Client.Cache.DefaultTtl = time.Minute * 10
	}
	if cfg.ZPA.Client.Cache.DefaultTti == 0 {
		cfg.ZPA.Client.Cache.DefaultTti = time.Minute * 8
	}
	cfg.CacheManager = newCache(cfg)
	// logger.Printf("[DEBUG] Cache initialized with TTL: %s, TTI: %s", cfg.ZPA.Client.Cache.DefaultTtl, cfg.ZPA.Client.Cache.DefaultTti)

	// Read configuration from YAML (lowest precedence)
	readConfigFromSystem(cfg)

	// Read environment variables (medium precedence)
	readConfigFromEnvironment(cfg)

	// Apply ConfigSetter functions
	for _, confSetter := range conf {
		confSetter(cfg)
	}
	// logger.Printf("[DEBUG] Configuration setters applied.")

	// Validate credentials after both setters and environment variables
	if cfg.ZPA.Client.ZPAClientID == "" || cfg.ZPA.Client.ZPAClientSecret == "" || cfg.ZPA.Client.ZPACustomerID == "" {
		logger.Printf("[ERROR] Missing ZPA credentials. Ensure they are provided via setters or environment variables.")
		return nil, errors.New("missing required ZPA credentials")
	}

	// Determine the base URL based on the ZPACloud value
	rawURL := defaultBaseURL
	switch strings.ToUpper(cfg.ZPA.Client.ZPACloud) {
	case "PRODUCTION", "":
		rawURL = defaultBaseURL
	case "ZPATWO":
		rawURL = zpaTwoBaseUrl
	case "BETA":
		rawURL = betaBaseURL
	case "GOV":
		rawURL = govBaseURL
	case "GOVUS":
		rawURL = govUsBaseURL
	case "PREVIEW":
		rawURL = previewBaseUrl
	case "DEV":
		rawURL = devBaseUrl
	case "QA":
		rawURL = qaBaseUrl
	case "QA2":
		rawURL = qa2BaseUrl
	}
	// logger.Printf("[DEBUG] Selected base URL: %s", rawURL)

	// Parse and validate the base URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		logger.Printf("[ERROR] Error occurred while parsing the base URL: %v", err)
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	cfg.BaseURL = parsedURL
	// logger.Printf("[DEBUG] Base URL parsed successfully: %s", parsedURL)

	// Set up HTTP clients
	setHttpClients(cfg)
	if cfg.HTTPClient == nil {
		logger.Printf("[ERROR] HTTP clients not initialized")
		return nil, errors.New("HTTP clients not initialized")
	}
	// logger.Printf("[DEBUG] HTTP clients configured.")

	// Authenticate the client and populate the AuthToken
	authToken, err := Authenticate(cfg.Context, cfg, logger)
	if err != nil {
		logger.Printf("[ERROR] Authentication failed: %v", err)
		return nil, fmt.Errorf("authentication failed: %w", err)
	}
	cfg.ZPA.Client.AuthToken = authToken

	// Add the AuthToken to the context
	if cfg.ZPA.Client.AuthToken != nil && cfg.ZPA.Client.AuthToken.AccessToken != "" {
		cfg.Context = context.WithValue(context.Background(), ContextAccessToken, cfg.ZPA.Client.AuthToken.AccessToken)
		// logger.Printf("[DEBUG] AuthToken added to context.")
	} else {
		logger.Printf("[ERROR] Failed to set AuthToken in context.")
		return nil, errors.New("AuthToken is missing or invalid after authentication")
	}

	// logger.Printf("[DEBUG] Configuration successfully initialized.")
	return cfg, nil
}

type ConfigSetter func(*Configuration)

// ConfigSetter type defines a function that modifies a Config struct.
// WithClientID sets the ClientID in the Config.
func WithZPAClientID(clientID string) ConfigSetter {
	return func(c *Configuration) {
		c.ZPA.Client.ZPAClientID = clientID
	}
}

// WithClientSecret sets the ClientSecret in the Config.
func WithZPAClientSecret(clientSecret string) ConfigSetter {
	return func(c *Configuration) {
		c.ZPA.Client.ZPAClientSecret = clientSecret
	}
}

func WithZPACustomerID(customerID string) ConfigSetter {
	return func(c *Configuration) {
		c.ZPA.Client.ZPACustomerID = customerID
	}
}

func WithZPAMicrotenantID(microtenantID string) ConfigSetter {
	return func(c *Configuration) {
		c.ZPA.Client.ZPAMicrotenantID = microtenantID
	}
}

func WithZPACloud(cloud string) ConfigSetter {
	return func(c *Configuration) {
		c.ZPA.Client.ZPACloud = cloud
	}
}

func WithCache(cache bool) ConfigSetter {
	return func(c *Configuration) {
		c.ZPA.Client.Cache.Enabled = cache
	}
}

func WithCacheManager(cacheManager cache.Cache) ConfigSetter {
	return func(c *Configuration) {
		c.CacheManager = cacheManager
	}
}

func newCache(c *Configuration) cache.Cache {
	cche, err := cache.NewCache(time.Duration(c.ZPA.Client.Cache.DefaultTtl), time.Duration(c.ZPA.Client.Cache.DefaultTti), int(c.ZPA.Client.Cache.DefaultCacheMaxSizeMB))
	if err != nil {
		cche = cache.NewNopCache()
	}
	return cche
}

func WithCacheTtl(i time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZPA.Client.Cache.DefaultTtl = i
		c.CacheManager = newCache(c)
	}
}

func WithCacheMaxSizeMB(size int64) ConfigSetter {
	return func(c *Configuration) {
		c.ZPA.Client.Cache.DefaultCacheMaxSizeMB = size
		c.CacheManager = newCache(c)
	}
}

func WithCacheTti(i time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZPA.Client.Cache.DefaultTti = i
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
		c.ZPA.Client.Proxy.Port = i
	}
}

func WithProxyHost(host string) ConfigSetter {
	return func(c *Configuration) {
		c.ZPA.Client.Proxy.Host = host
	}
}

func WithProxyUsername(username string) ConfigSetter {
	return func(c *Configuration) {
		c.ZPA.Client.Proxy.Username = username
	}
}

func WithProxyPassword(pass string) ConfigSetter {
	return func(c *Configuration) {
		c.ZPA.Client.Proxy.Password = pass
	}
}

func WithTestingDisableHttpsCheck(httpsCheck bool) ConfigSetter {
	return func(c *Configuration) {
		c.ZPA.Testing.DisableHttpsCheck = httpsCheck
	}
}

func WithRequestTimeout(requestTimeout time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZPA.Client.RequestTimeout = requestTimeout
		setHttpClients(c)
	}
}

func WithRateLimitMaxRetries(maxRetries int32) ConfigSetter {
	return func(c *Configuration) {
		c.ZPA.Client.RateLimit.MaxRetries = maxRetries
		setHttpClients(c)
	}
}

func WithRateLimitMaxWait(maxWait time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZPA.Client.RateLimit.RetryWaitMax = maxWait
		setHttpClients(c)
	}
}

func WithRateLimitMinWait(minWait time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZPA.Client.RateLimit.RetryWaitMin = minWait
		setHttpClients(c)
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

func setHttpClients(cfg *Configuration) {
	// Use a temporary logger if cfg or cfg.Logger is nil
	var log logger.Logger
	if cfg == nil || cfg.Logger == nil {
		log = logger.GetDefaultLogger(loggerPrefix) // Use default logger
		log.Printf("[ERROR] Configuration is nil. Cannot initialize HTTP clients.")
		return
	} else {
		log = cfg.Logger
	}

	// ZPA-specific rate limits:
	zpaRateLimiter := rl.NewRateLimiter(20, 10, 10, 10) // GET: 20 per 10s, POST/PUT/DELETE: 10 per 10s

	// Configure the ZPA HTTP client
	cfg.HTTPClient = getHTTPClient(log, zpaRateLimiter, cfg)
	if cfg.HTTPClient == nil {
		log.Printf("[ERROR] Failed to initialize ZPA HTTP client.")
	} else {
		// log.Printf("[DEBUG] ZPA HTTP client initialized successfully.")
	}

	// Configure the generic HTTP client
	cfg.HTTPClient = getHTTPClient(log, nil, cfg)
	if cfg.HTTPClient == nil {
		log.Printf("[ERROR] Failed to initialize generic HTTP client.")
	} else {
		// log.Printf("[DEBUG] Generic HTTP client initialized successfully.")
	}
}

func readConfigFromFile(location string, c *Configuration) (*Configuration, error) {
	yamlConfig, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlConfig, &c)
	if err != nil {
		return nil, err
	}
	return c, err
}

func readConfigFromSystem(c *Configuration) *Configuration {
	currUser, err := user.Current()
	if err != nil {
		return c
	}
	if currUser.HomeDir == "" {
		return c
	}
	conf, err := readConfigFromFile(currUser.HomeDir+"/.zscaler/zscaler.yaml", c)
	if err != nil {
		return c
	}
	return conf
}

func readConfigFromEnvironment(c *Configuration) *Configuration {
	err := envconfig.Process("zscaler", c)
	if err != nil {
		c.Logger.Printf("[ERROR] Error parsing environment variables: %v", err)
		return c
	}
	// c.Logger.Printf("[DEBUG] Successfully parsed environment variables.")
	return c
}

// AddDefaultHeader adds a new HTTP header to the default header in the request
func (c *Configuration) AddDefaultHeader(key string, value string) {
	c.DefaultHeader[key] = value
}
