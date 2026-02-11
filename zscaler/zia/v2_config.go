package zia

import (
	"context"
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
	JSessionIDTimeout         = 5               // minutes - updated per API provider requirements
	jSessionTimeoutOffset     = 1 * time.Minute // reduced to 1 minute for 5-minute timeout
	contentTypeJSON           = "application/json"
	cookieName                = "JSESSIONID"
	MaxNumOfRetries           = 100
	RetryWaitMaxSeconds       = 20
	RetryWaitMinSeconds       = 5
	ziaAPIVersion             = "api/v1"
	ziaAPIAuthURL             = "/authenticatedSession"
	loggerPrefix              = "zia-logger: "
)

const (
	VERSION      = "3.7.5"
	ZIA_USERNAME = "ZIA_USERNAME"
	ZIA_PASSWORD = "ZIA_PASSWORD"
	ZIA_API_KEY  = "ZIA_API_KEY"
	ZIA_CLOUD    = "ZIA_CLOUD"
)

var globalConfig *Configuration
var configOnce sync.Once

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
	partnerID        string
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
	ZIA            struct {
		Client struct {
			ZIAUsername string `yaml:"username" envconfig:"ZIA_USERNAME"`
			ZIAPassword string `yaml:"password" envconfig:"ZIA_PASSWORD"`
			ZIAApiKey   string `yaml:"apiKey" envconfig:"ZIA_API_KEY"`
			ZIACloud    string `yaml:"cloud" envconfig:"ZIA_CLOUD"`
			PartnerID   string `yaml:"partnerId" envconfig:"ZSCALER_PARTNER_ID"`
			Cache       struct {
				Enabled               bool          `yaml:"enabled" envconfig:"ZIA_CLIENT_CACHE_ENABLED"`
				DefaultTtl            time.Duration `yaml:"defaultTtl" envconfig:"ZIA_CLIENT_CACHE_DEFAULT_TTL"`
				DefaultTti            time.Duration `yaml:"defaultTti" envconfig:"ZIA_CLIENT_CACHE_DEFAULT_TTI"`
				DefaultCacheMaxSizeMB int64         `yaml:"defaultSize" envconfig:"ZIA_CLIENT_CACHE_DEFAULT_SIZE"`
			} `yaml:"cache"`
			Proxy struct {
				Port     int32  `yaml:"port" envconfig:"ZIA_CLIENT_PROXY_PORT"`
				Host     string `yaml:"host" envconfig:"ZIA_CLIENT_PROXY_HOST"`
				Username string `yaml:"username" envconfig:"ZIA_CLIENT_PROXY_USERNAME"`
				Password string `yaml:"password" envconfig:"ZIA_CLIENT_PROXY_PASSWORD"`
			} `yaml:"proxy"`
			RequestTimeout time.Duration `yaml:"requestTimeout" envconfig:"ZIA_CLIENT_REQUEST_TIMEOUT"`
			RateLimit      struct {
				MaxRetries   int32         `yaml:"maxRetries" envconfig:"ZIA_CLIENT_RATE_LIMIT_MAX_RETRIES"`
				RetryWaitMin time.Duration `yaml:"minWait" envconfig:"ZIA_CLIENT_RATE_LIMIT_MIN_WAIT"`
				RetryWaitMax time.Duration `yaml:"maxWait" envconfig:"ZIA_CLIENT_RATE_LIMIT_MAX_WAIT"`
			} `yaml:"rateLimit"`
		} `yaml:"client"`
		Testing struct {
			DisableHttpsCheck bool `yaml:"disableHttpsCheck" envconfig:"ZIA_TESTING_DISABLE_HTTPS_CHECK"`
		} `yaml:"testing"`
	} `yaml:"zia"`
	CacheManager cache.Cache
}

func NewConfiguration(conf ...ConfigSetter) (*Configuration, error) {
	configOnce.Do(func() {
		logger := logger.GetDefaultLogger(loggerPrefix)
		globalConfig = &Configuration{
			DefaultHeader: make(map[string]string),
			Logger:        logger,
			UserAgent:     fmt.Sprintf("zscaler-sdk-go/%s golang/%s %s/%s", VERSION, runtime.Version(), runtime.GOOS, runtime.GOARCH),
			Debug:         false,
			Context:       context.Background(),
		}

		// Apply ConfigSetters first to ensure all values are set before use
		for _, confSetter := range conf {
			confSetter(globalConfig)
		}

		// Parse and validate the base URL
		if globalConfig.ZIA.Client.ZIACloud == "" {
			logger.Printf("[ERROR] Missing ZIA cloud configuration.")
			return
		}

		// Debug log to ensure ZIACloud is set correctly
		logger.Printf("[DEBUG] ZIACloud: %s", globalConfig.ZIA.Client.ZIACloud)

		rawURL := fmt.Sprintf("https://zsapi.%s.net", globalConfig.ZIA.Client.ZIACloud)
		if globalConfig.ZIA.Client.ZIACloud == "zspreview" {
			rawURL = fmt.Sprintf("https://admin.%s.net", globalConfig.ZIA.Client.ZIACloud)
		}

		parsedURL, err := url.Parse(rawURL)
		if err != nil {
			logger.Printf("[ERROR] Error occurred while parsing the base URL: %v", err)
			return
		}
		globalConfig.BaseURL = parsedURL

		// Initialize cache with defaults
		if globalConfig.ZIA.Client.Cache.DefaultTtl == 0 {
			globalConfig.ZIA.Client.Cache.DefaultTtl = time.Minute * 10
		}
		if globalConfig.ZIA.Client.Cache.DefaultTti == 0 {
			globalConfig.ZIA.Client.Cache.DefaultTti = time.Minute * 8
		}
		globalConfig.CacheManager = newCache(globalConfig)

		// Set default rate limit and request timeout values
		globalConfig.ZIA.Client.RateLimit.MaxRetries = MaxNumOfRetries
		globalConfig.ZIA.Client.RateLimit.RetryWaitMax = time.Second * time.Duration(RetryWaitMaxSeconds)
		globalConfig.ZIA.Client.RateLimit.RetryWaitMin = time.Second * time.Duration(RetryWaitMinSeconds)
		globalConfig.ZIA.Client.RequestTimeout = time.Duration(requestTimeout) * time.Second

		// Initialize testing configurations
		globalConfig.ZIA.Testing.DisableHttpsCheck = false

		// Read configuration from system and environment
		readConfigFromSystem(globalConfig)
		readConfigFromEnvironment(globalConfig)

		// Set up HTTP clients
		setHttpClients(globalConfig)
		if globalConfig.HTTPClient == nil {
			logger.Printf("[ERROR] HTTP clients not initialized")
			return
		}

		// Validate critical configuration fields
		if globalConfig.ZIA.Client.ZIAUsername == "" ||
			globalConfig.ZIA.Client.ZIAPassword == "" ||
			globalConfig.ZIA.Client.ZIAApiKey == "" {
			logger.Printf("[ERROR] Missing client credentials (ZIA_USERNAME, ZIA_PASSWORD, ZIA_API_KEY).")
			return
		}

		// Log success initialization
		logger.Printf("[INFO] Configuration successfully initialized.")
	})

	if globalConfig.UserAgentExtra != "" {
		globalConfig.UserAgent = fmt.Sprintf("%s %s", globalConfig.UserAgent, globalConfig.UserAgentExtra)
	}

	// Return the global configuration
	if globalConfig == nil {
		return nil, fmt.Errorf("failed to initialize configuration")
	}
	return globalConfig, nil
}

type ConfigSetter func(*Configuration)

// ConfigSetter type defines a function that modifies a Config struct.
// WithZiaUsername sets the Username in the Config.
func WithZiaUsername(username string) ConfigSetter {
	return func(c *Configuration) {
		c.ZIA.Client.ZIAUsername = username
	}
}

// WithZiaPassword sets the Password in the Config.
func WithZiaPassword(password string) ConfigSetter {
	return func(c *Configuration) {
		c.ZIA.Client.ZIAPassword = password
	}
}

// WithZiaAPIKey sets the ApiKey in the Config.
func WithZiaAPIKey(apiKey string) ConfigSetter {
	return func(c *Configuration) {
		c.ZIA.Client.ZIAApiKey = apiKey
	}
}

// WithZiaAPIKey sets the ApiKey in the Config.
func WithZiaCloud(cloud string) ConfigSetter {
	return func(c *Configuration) {
		c.ZIA.Client.ZIACloud = cloud
	}
}

// WithPartnerID sets the PartnerID in the Config.
func WithPartnerID(partnerID string) ConfigSetter {
	return func(c *Configuration) {
		c.ZIA.Client.PartnerID = partnerID
	}
}

func WithCache(cache bool) ConfigSetter {
	return func(c *Configuration) {
		c.ZIA.Client.Cache.Enabled = cache
	}
}

func WithCacheManager(cacheManager cache.Cache) ConfigSetter {
	return func(c *Configuration) {
		c.CacheManager = cacheManager
	}
}

func newCache(c *Configuration) cache.Cache {
	if !c.ZIA.Client.Cache.Enabled {
		return cache.NewNopCache()
	}
	cche, err := cache.NewCache(
		time.Duration(c.ZIA.Client.Cache.DefaultTtl),
		time.Duration(c.ZIA.Client.Cache.DefaultTti),
		int(c.ZIA.Client.Cache.DefaultCacheMaxSizeMB),
	)
	if err != nil {
		return cache.NewNopCache()
	}
	return cche
}

func WithCacheTtl(i time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZIA.Client.Cache.DefaultTtl = i
		c.CacheManager = newCache(c)
	}
}

func WithCacheMaxSizeMB(size int64) ConfigSetter {
	return func(c *Configuration) {
		c.ZIA.Client.Cache.DefaultCacheMaxSizeMB = size
		c.CacheManager = newCache(c)
	}
}

func WithCacheTti(i time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZIA.Client.Cache.DefaultTti = i
		c.CacheManager = newCache(c)
	}
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

// WithHttpClient sets the HttpClient in the Config.
func WithHttpClientPtr(httpClient *http.Client) ConfigSetter {
	return func(c *Configuration) {
		c.HTTPClient = httpClient
	}
}

func WithProxyPort(i int32) ConfigSetter {
	return func(c *Configuration) {
		c.ZIA.Client.Proxy.Port = i
	}
}

func WithProxyHost(host string) ConfigSetter {
	return func(c *Configuration) {
		c.ZIA.Client.Proxy.Host = host
	}
}

func WithProxyUsername(username string) ConfigSetter {
	return func(c *Configuration) {
		c.ZIA.Client.Proxy.Username = username
	}
}

func WithProxyPassword(pass string) ConfigSetter {
	return func(c *Configuration) {
		c.ZIA.Client.Proxy.Password = pass
	}
}

func WithTestingDisableHttpsCheck(httpsCheck bool) ConfigSetter {
	return func(c *Configuration) {
		c.ZIA.Testing.DisableHttpsCheck = httpsCheck
	}
}

func WithRequestTimeout(requestTimeout time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZIA.Client.RequestTimeout = requestTimeout
		setHttpClients(c)
	}
}

func WithRateLimitMaxRetries(maxRetries int32) ConfigSetter {
	return func(c *Configuration) {
		c.ZIA.Client.RateLimit.MaxRetries = maxRetries
		setHttpClients(c)
	}
}

func WithRateLimitMaxWait(maxWait time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZIA.Client.RateLimit.RetryWaitMax = maxWait
		setHttpClients(c)
	}
}

func WithRateLimitMinWait(minWait time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZIA.Client.RateLimit.RetryWaitMin = minWait
		setHttpClients(c)
	}
}

// WithUserAgent sets the UserAgent in the Config.
func WithUserAgentExtra(userAgent string) ConfigSetter {
	return func(c *Configuration) {
		c.UserAgentExtra = userAgent
	}
}

func WithUserAgent(userAgent string) ConfigSetter {
	return func(cfg *Configuration) {
		cfg.UserAgent = userAgent
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

	// ZIA-specific rate limits:
	// GET: 20 requests per 10s (2/sec), POST/PUT: 10 requests per 10s (1/sec), DELETE: 1 request per 61s
	ziaRateLimiter := rl.NewRateLimiter(20, 10, 10, 61) // Adjusted for ZIA based on official limits and +1 sec buffer

	// Configure the ZIA HTTP client
	cfg.HTTPClient = getHTTPClient(log, ziaRateLimiter, cfg)
	if cfg.HTTPClient == nil {
		log.Printf("[ERROR] Failed to initialize ZIA HTTP client.")
	} else {
		// log.Printf("[DEBUG] ZIA HTTP client initialized successfully.")
	}

	// Remove the call that overwrote the ZIA client:
	// cfg.HTTPClient = getHTTPClient(log, nil, cfg)
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
