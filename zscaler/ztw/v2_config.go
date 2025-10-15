package ztw

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"runtime"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/zscaler/zscaler-sdk-go/v3/cache"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
	rl "github.com/zscaler/zscaler-sdk-go/v3/ratelimiter"
	"gopkg.in/yaml.v3"
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
	ztwAPIVersion             = "api/v1"
	ztwAPIAuthURL             = "/auth"
	loggerPrefix              = "ztw-logger: "
)

const (
	VERSION      = "3.7.5"
	ZTW_USERNAME = "ZTW_USERNAME"
	ZTW_PASSWORD = "ZTW_PASSWORD"
	ZTW_API_KEY  = "ZTW_API_KEY"
	ZTW_CLOUD    = "ZTW_CLOUD"
)

type contextKey string

func (c contextKey) String() string {
	return "zscaler " + string(c)
}

var (
	// ContextAccessToken takes a string OAuth2 access token as authentication for the request.
	ContextAccessToken = contextKey("access_token")
)

type Client struct {
	sync.Mutex
	userName         string
	password         string
	cloud            string
	apiKey           string
	session          *Session
	sessionRefreshed time.Time     // Also indicates last usage
	sessionTimeout   time.Duration // in minutes
	URL              string
	HTTPClient       *http.Client
	Logger           logger.Logger
	UserAgent        string
	freshCache       bool
	cacheEnabled     bool
	cache            cache.Cache
	cacheTtl         time.Duration
	cacheCleanwindow time.Duration
	cacheMaxSizeMB   int
	rateLimiter      *rl.RateLimiter
	sessionTicker    *time.Ticker
	// stopTicker       chan bool
	ctx        context.Context
	cancelFunc context.CancelFunc
	refreshing bool
}

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
	ZTW            struct {
		Client struct {
			ZTWUsername string `yaml:"username" envconfig:"ZTW_USERNAME"`
			ZTWPassword string `yaml:"password" envconfig:"ZTW_PASSWORD"`
			ZTWApiKey   string `yaml:"apiKey" envconfig:"ZTW_API_KEY"`
			ZTWCloud    string `yaml:"cloud" envconfig:"ZTW_CLOUD"`
			Cache       struct {
				Enabled               bool          `yaml:"enabled" envconfig:"ZTW_CLIENT_CACHE_ENABLED"`
				DefaultTtl            time.Duration `yaml:"defaultTtl" envconfig:"ZTW_CLIENT_CACHE_DEFAULT_TTL"`
				DefaultTti            time.Duration `yaml:"defaultTti" envconfig:"ZTW_CLIENT_CACHE_DEFAULT_TTI"`
				DefaultCacheMaxSizeMB int64         `yaml:"defaultSize" envconfig:"ZTW_CLIENT_CACHE_DEFAULT_SIZE"`
			} `yaml:"cache"`
			Proxy struct {
				Port     int32  `yaml:"port" envconfig:"ZTW_CLIENT_PROXY_PORT"`
				Host     string `yaml:"host" envconfig:"ZTW_CLIENT_PROXY_HOST"`
				Username string `yaml:"username" envconfig:"ZTW_CLIENT_PROXY_USERNAME"`
				Password string `yaml:"password" envconfig:"ZTW_CLIENT_PROXY_PASSWORD"`
			} `yaml:"proxy"`
			RequestTimeout time.Duration `yaml:"requestTimeout" envconfig:"ZTW_CLIENT_REQUEST_TIMEOUT"`
			RateLimit      struct {
				MaxRetries   int32         `yaml:"maxRetries" envconfig:"ZTW_CLIENT_RATE_LIMIT_MAX_RETRIES"`
				RetryWaitMin time.Duration `yaml:"minWait" envconfig:"ZTW_CLIENT_RATE_LIMIT_MIN_WAIT"`
				RetryWaitMax time.Duration `yaml:"maxWait" envconfig:"ZTW_CLIENT_RATE_LIMIT_MAX_WAIT"`
			} `yaml:"rateLimit"`
		} `yaml:"client"`
		Testing struct {
			DisableHttpsCheck bool `yaml:"disableHttpsCheck" envconfig:"ZTW_TESTING_DISABLE_HTTPS_CHECK"`
		} `yaml:"testing"`
	} `yaml:"ztw"`
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

	logger.Printf("[DEBUG] Initializing configuration with default values.")

	// Set default rate limit and request timeout values
	cfg.ZTW.Client.RateLimit.MaxRetries = MaxNumOfRetries
	cfg.ZTW.Client.RateLimit.RetryWaitMax = time.Second * time.Duration(RetryWaitMaxSeconds)
	cfg.ZTW.Client.RateLimit.RetryWaitMin = time.Second * time.Duration(RetryWaitMinSeconds)
	cfg.ZTW.Client.RequestTimeout = time.Duration(requestTimeout) * time.Second

	// Apply additional configurations from ConfigSetters
	for _, confSetter := range conf {
		confSetter(cfg)
	}

	// Read configuration from YAML and environment
	readConfigFromSystem(cfg)
	readConfigFromEnvironment(cfg)

	// Validate critical configuration fields
	if cfg.ZTW.Client.ZTWUsername == "" || cfg.ZTW.Client.ZTWPassword == "" || cfg.ZTW.Client.ZTWApiKey == "" || cfg.ZTW.Client.ZTWCloud == "" {
		logger.Printf("[ERROR] Missing client credentials (ZTW_USERNAME, ZTW_PASSWORD, ZTW_API_KEY, ZTW_CLOUD).")
		return nil, fmt.Errorf("client credentials (ZTW_USERNAME, ZTW_PASSWORD, ZTW_API_KEY, ZTW_CLOUD) are missing")
	}

	// Construct base URL with the API version
	rawURL := fmt.Sprintf("https://connector.%s.net/%s", cfg.ZTW.Client.ZTWCloud, ztwAPIVersion)
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		logger.Printf("[ERROR] Error occurred while parsing the base URL: %v", err)
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	cfg.BaseURL = parsedURL

	// Set up HTTP clients
	setHttpClients(cfg)
	if cfg.HTTPClient == nil {
		logger.Printf("[ERROR] HTTP clients not initialized")
		return nil, errors.New("HTTP clients not initialized")
	}

	logger.Printf("[DEBUG] Configuration successfully initialized.")
	return cfg, nil
}

type ConfigSetter func(*Configuration)

// ConfigSetter type defines a function that modifies a Config struct.
// WithZtwUsername sets the Username in the Config.
func WithZtwUsername(username string) ConfigSetter {
	return func(c *Configuration) {
		c.ZTW.Client.ZTWUsername = username
	}
}

// WithZtwPassword sets the Password in the Config.
func WithZtwPassword(password string) ConfigSetter {
	return func(c *Configuration) {
		c.ZTW.Client.ZTWPassword = password
	}
}

// func WithZtwAPIKey(apiKey string) ConfigSetter {
func WithZtwAPIKey(apiKey string) ConfigSetter {
	return func(c *Configuration) {
		c.ZTW.Client.ZTWApiKey = apiKey
	}
}

// WithZtwCloud sets the ApiKey in the Config.
func WithZtwCloud(cloud string) ConfigSetter {
	return func(c *Configuration) {
		c.ZTW.Client.ZTWCloud = cloud
	}
}

func WithCacheManager(cacheManager cache.Cache) ConfigSetter {
	return func(c *Configuration) {
		c.CacheManager = cacheManager
	}
}

func newCache(c *Configuration) cache.Cache {
	cche, err := cache.NewCache(time.Duration(c.ZTW.Client.Cache.DefaultTtl), time.Duration(c.ZTW.Client.Cache.DefaultTti), int(c.ZTW.Client.Cache.DefaultCacheMaxSizeMB))
	if err != nil {
		cche = cache.NewNopCache()
	}
	return cche
}

func WithCacheTtl(i time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZTW.Client.Cache.DefaultTtl = i
		c.CacheManager = newCache(c)
	}
}

func WithCacheMaxSizeMB(size int64) ConfigSetter {
	return func(c *Configuration) {
		c.ZTW.Client.Cache.DefaultCacheMaxSizeMB = size
		c.CacheManager = newCache(c)
	}
}

func WithCacheTti(i time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZTW.Client.Cache.DefaultTti = i
		c.CacheManager = newCache(c)
	}
}

func WithCache(cache bool) ConfigSetter {
	return func(c *Configuration) {
		c.ZTW.Client.Cache.Enabled = cache
	}
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

// WithHttpClient sets the HttpClient in the Config.
func WithHttpClientPtr(httpClient *http.Client) ConfigSetter {
	return func(c *Configuration) {
		c.HTTPClient = httpClient
	}
}

func WithProxyPort(i int32) ConfigSetter {
	return func(c *Configuration) {
		c.ZTW.Client.Proxy.Port = i
	}
}

func WithProxyHost(host string) ConfigSetter {
	return func(c *Configuration) {
		c.ZTW.Client.Proxy.Host = host
	}
}

func WithProxyUsername(username string) ConfigSetter {
	return func(c *Configuration) {
		c.ZTW.Client.Proxy.Username = username
	}
}

func WithProxyPassword(pass string) ConfigSetter {
	return func(c *Configuration) {
		c.ZTW.Client.Proxy.Password = pass
	}
}

func WithTestingDisableHttpsCheck(httpsCheck bool) ConfigSetter {
	return func(c *Configuration) {
		c.ZTW.Testing.DisableHttpsCheck = httpsCheck
	}
}

func WithRequestTimeout(requestTimeout time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZTW.Client.RequestTimeout = requestTimeout
		setHttpClients(c)
	}
}

func WithRateLimitMaxRetries(maxRetries int32) ConfigSetter {
	return func(c *Configuration) {
		c.ZTW.Client.RateLimit.MaxRetries = maxRetries
		setHttpClients(c)
	}
}

func WithRateLimitMaxWait(maxWait time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZTW.Client.RateLimit.RetryWaitMax = maxWait
		setHttpClients(c)
	}
}

func WithRateLimitMinWait(minWait time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZTW.Client.RateLimit.RetryWaitMin = minWait
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

	// ZTW-specific rate limits:
	// GET: 20 requests per 10s (2/sec), POST/PUT: 10 requests per 10s (1/sec), DELETE: 1 request per 61s
	ztwRateLimiter := rl.NewRateLimiter(20, 10, 10, 61) // Adjusted for ZTW based on official limits and +1 sec buffer

	// Configure the ZTW HTTP client
	cfg.HTTPClient = getHTTPClient(log, ztwRateLimiter, cfg)
	if cfg.HTTPClient == nil {
		log.Printf("[ERROR] Failed to initialize ZTW HTTP client.")
	} else {
		// log.Printf("[DEBUG] ZTW HTTP client initialized successfully.")
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
		fmt.Println("error parsing")
		return c
	}
	return c
}

// AddDefaultHeader adds a new HTTP header to the default header in the request
func (c *Configuration) AddDefaultHeader(key string, value string) {
	c.DefaultHeader[key] = value
}
