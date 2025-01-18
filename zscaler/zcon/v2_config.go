package zcon

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
	zconAPIVersion            = "api/v1"
	zconAPIAuthURL            = "/auth"
	loggerPrefix              = "zcon-logger: "
)

const (
	VERSION       = "3.0.0"
	ZCON_USERNAME = "ZCON_USERNAME"
	ZCON_PASSWORD = "ZCON_PASSWORD"
	ZCON_API_KEY  = "ZCON_API_KEY"
	ZCON_CLOUD    = "ZCON_CLOUD"
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
	stopTicker       chan bool
	refreshing       bool
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
	ZCON           struct {
		Client struct {
			ZCONUsername string `yaml:"username" envconfig:"ZCON_USERNAME"`
			ZCONPassword string `yaml:"password" envconfig:"ZCON_PASSWORD"`
			ZCONApiKey   string `yaml:"apiKey" envconfig:"ZCON_API_KEY"`
			ZCONCloud    string `yaml:"cloud" envconfig:"ZCON_CLOUD"`
			Cache        struct {
				Enabled               bool          `yaml:"enabled" envconfig:"ZCON_CLIENT_CACHE_ENABLED"`
				DefaultTtl            time.Duration `yaml:"defaultTtl" envconfig:"ZCON_CLIENT_CACHE_DEFAULT_TTL"`
				DefaultTti            time.Duration `yaml:"defaultTti" envconfig:"ZCON_CLIENT_CACHE_DEFAULT_TTI"`
				DefaultCacheMaxSizeMB int64         `yaml:"defaultTti" envconfig:"ZCON_CLIENT_CACHE_DEFAULT_SIZE"`
			} `yaml:"cache"`
			Proxy struct {
				Port     int32  `yaml:"port" envconfig:"ZCON_CLIENT_PROXY_PORT"`
				Host     string `yaml:"host" envconfig:"ZCON_CLIENT_PROXY_HOST"`
				Username string `yaml:"username" envconfig:"ZCON_CLIENT_PROXY_USERNAME"`
				Password string `yaml:"password" envconfig:"ZCON_CLIENT_PROXY_PASSWORD"`
			} `yaml:"proxy"`
			RequestTimeout time.Duration `yaml:"requestTimeout" envconfig:"ZCON_CLIENT_REQUEST_TIMEOUT"`
			RateLimit      struct {
				MaxRetries   int32         `yaml:"maxRetries" envconfig:"ZCON_CLIENT_RATE_LIMIT_MAX_RETRIES"`
				RetryWaitMin time.Duration `yaml:"minWait" envconfig:"ZCON_CLIENT_RATE_LIMIT_MIN_WAIT"`
				RetryWaitMax time.Duration `yaml:"maxWait" envconfig:"ZCON_CLIENT_RATE_LIMIT_MAX_WAIT"`
			} `yaml:"rateLimit"`
		} `yaml:"client"`
		Testing struct {
			DisableHttpsCheck bool `yaml:"disableHttpsCheck" envconfig:"ZCON_TESTING_DISABLE_HTTPS_CHECK"`
		} `yaml:"testing"`
	} `yaml:"zcon"`
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
	cfg.ZCON.Client.RateLimit.MaxRetries = MaxNumOfRetries
	cfg.ZCON.Client.RateLimit.RetryWaitMax = time.Second * time.Duration(RetryWaitMaxSeconds)
	cfg.ZCON.Client.RateLimit.RetryWaitMin = time.Second * time.Duration(RetryWaitMinSeconds)
	cfg.ZCON.Client.RequestTimeout = time.Duration(requestTimeout) * time.Second

	// Apply additional configurations from ConfigSetters
	for _, confSetter := range conf {
		confSetter(cfg)
	}

	// Read configuration from YAML and environment
	readConfigFromSystem(cfg)
	readConfigFromEnvironment(cfg)

	// Validate critical configuration fields
	if cfg.ZCON.Client.ZCONUsername == "" || cfg.ZCON.Client.ZCONPassword == "" || cfg.ZCON.Client.ZCONApiKey == "" || cfg.ZCON.Client.ZCONCloud == "" {
		logger.Printf("[ERROR] Missing client credentials (ZCON_USERNAME, ZCON_PASSWORD, ZCON_API_KEY, ZCON_CLOUD).")
		return nil, fmt.Errorf("client credentials (ZCON_USERNAME, ZCON_PASSWORD, ZCON_API_KEY, ZCON_CLOUD) are missing")
	}

	// Construct base URL with the API version
	rawURL := fmt.Sprintf("https://connector.%s.net/%s", cfg.ZCON.Client.ZCONCloud, zconAPIVersion)
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
// WithZconUsername sets the Username in the Config.
func WithZconUsername(username string) ConfigSetter {
	return func(c *Configuration) {
		c.ZCON.Client.ZCONUsername = username
	}
}

// WithZconPassword sets the Password in the Config.
func WithZconPassword(password string) ConfigSetter {
	return func(c *Configuration) {
		c.ZCON.Client.ZCONPassword = password
	}
}

// WithZconAPIKey sets the ApiKey in the Config.
func WithZconAPIKey(apiKey string) ConfigSetter {
	return func(c *Configuration) {
		c.ZCON.Client.ZCONApiKey = apiKey
	}
}

// WithZconAPIKey sets the ApiKey in the Config.
func WithZconCloud(cloud string) ConfigSetter {
	return func(c *Configuration) {
		c.ZCON.Client.ZCONCloud = cloud
	}
}

func WithCacheManager(cacheManager cache.Cache) ConfigSetter {
	return func(c *Configuration) {
		c.CacheManager = cacheManager
	}
}

func newCache(c *Configuration) cache.Cache {
	cche, err := cache.NewCache(time.Duration(c.ZCON.Client.Cache.DefaultTtl), time.Duration(c.ZCON.Client.Cache.DefaultTti), int(c.ZCON.Client.Cache.DefaultCacheMaxSizeMB))
	if err != nil {
		cche = cache.NewNopCache()
	}
	return cche
}

func WithCacheTtl(i time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZCON.Client.Cache.DefaultTtl = i
		c.CacheManager = newCache(c)
	}
}

func WithCacheMaxSizeMB(size int64) ConfigSetter {
	return func(c *Configuration) {
		c.ZCON.Client.Cache.DefaultCacheMaxSizeMB = size
		c.CacheManager = newCache(c)
	}
}

func WithCacheTti(i time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZCON.Client.Cache.DefaultTti = i
		c.CacheManager = newCache(c)
	}
}

func WithCache(cache bool) ConfigSetter {
	return func(c *Configuration) {
		c.ZCON.Client.Cache.Enabled = cache
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
		c.ZCON.Client.Proxy.Port = i
	}
}

func WithProxyHost(host string) ConfigSetter {
	return func(c *Configuration) {
		c.ZCON.Client.Proxy.Host = host
	}
}

func WithProxyUsername(username string) ConfigSetter {
	return func(c *Configuration) {
		c.ZCON.Client.Proxy.Username = username
	}
}

func WithProxyPassword(pass string) ConfigSetter {
	return func(c *Configuration) {
		c.ZCON.Client.Proxy.Password = pass
	}
}

func WithTestingDisableHttpsCheck(httpsCheck bool) ConfigSetter {
	return func(c *Configuration) {
		c.ZCON.Testing.DisableHttpsCheck = httpsCheck
	}
}

func WithRequestTimeout(requestTimeout time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZCON.Client.RequestTimeout = requestTimeout
		setHttpClients(c)
	}
}

func WithRateLimitMaxRetries(maxRetries int32) ConfigSetter {
	return func(c *Configuration) {
		c.ZCON.Client.RateLimit.MaxRetries = maxRetries
		setHttpClients(c)
	}
}

func WithRateLimitMaxWait(maxWait time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZCON.Client.RateLimit.RetryWaitMax = maxWait
		setHttpClients(c)
	}
}

func WithRateLimitMinWait(minWait time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZCON.Client.RateLimit.RetryWaitMin = minWait
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

	// ZCON-specific rate limits:
	// GET: 20 requests per 10s (2/sec), POST/PUT: 10 requests per 10s (1/sec), DELETE: 1 request per 61s
	zconRateLimiter := rl.NewRateLimiter(20, 10, 10, 61) // Adjusted for ZCON based on official limits and +1 sec buffer

	// Configure the ZCON HTTP client
	cfg.HTTPClient = getHTTPClient(log, zconRateLimiter, cfg)
	if cfg.HTTPClient == nil {
		log.Printf("[ERROR] Failed to initialize ZCON HTTP client.")
	} else {
		// log.Printf("[DEBUG] ZCON HTTP client initialized successfully.")
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
