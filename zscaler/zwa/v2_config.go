package zwa

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
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
	rl "github.com/zscaler/zscaler-sdk-go/v3/ratelimiter"
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
	maxIdleConnections    int = 40
	requestTimeout        int = 60
	JSessionIDTimeout         = 30 // minutes.
	jSessionTimeoutOffset     = 5 * time.Minute
	contentTypeJSON           = "application/json"
	cookieName                = "JSESSIONID"
	MaxNumOfRetries           = 100
	RetryWaitMaxSeconds       = 20
	RetryWaitMinSeconds       = 5
	loggerPrefix              = "zwa-logger: "
)

const (
	VERSION        = "3.7.5"
	ZWA_API_KEY_ID = "ZWA_API_KEY_ID"
	ZWA_API_SECRET = "ZWA_API_SECRET"
)

type BackoffConfig struct {
	Enabled             bool // Set to true to enable backoff and retry mechanism
	RetryWaitMinSeconds int  // Minimum time to wait
	RetryWaitMaxSeconds int  // Maximum time to wait
	MaxNumOfRetries     int  // Maximum number of retries
}

type AuthRequest struct {
	APIKeyID     string `json:"key_id"`
	APIKeySecret string `json:"key_secret"`
	Timestamp    int64  `json:"timestamp"`
}

type AuthToken struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"token"`
	ExpiresIn   int    `json:"expires_in"`
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
	ZWA            struct {
		Client struct {
			ZWAAPIKeyID  string     `yaml:"key_id" envconfig:"ZWA_API_KEY_ID"`
			ZWAAPISecret string     `yaml:"key_secret" envconfig:"ZWA_API_SECRET"`
			ZWACloud     string     `yaml:"cloud" envconfig:"ZWA_CLOUD"`
			PartnerID    string     `yaml:"partnerId" envconfig:"ZSCALER_PARTNER_ID"`
			AuthToken    *AuthToken `yaml:"authToken"`
			AccessToken  *AuthToken `yaml:"accessToken"`
			Proxy        struct {
				Port     int32  `yaml:"port" envconfig:"ZWA_CLIENT_PROXY_PORT"`
				Host     string `yaml:"host" envconfig:"ZWA_CLIENT_PROXY_HOST"`
				Username string `yaml:"username" envconfig:"ZWA_CLIENT_PROXY_USERNAME"`
				Password string `yaml:"password" envconfig:"ZWA_CLIENT_PROXY_PASSWORD"`
			} `yaml:"proxy"`
			RequestTimeout time.Duration `yaml:"requestTimeout" envconfig:"ZWA_CLIENT_REQUEST_TIMEOUT"`
			RateLimit      struct {
				MaxRetries   int32         `yaml:"maxRetries" envconfig:"ZWA_CLIENT_RATE_LIMIT_MAX_RETRIES"`
				RetryWaitMin time.Duration `yaml:"minWait" envconfig:"ZWA_CLIENT_RATE_LIMIT_MIN_WAIT"`
				RetryWaitMax time.Duration `yaml:"maxWait" envconfig:"ZZWA_CLIENT_RATE_LIMIT_MAX_WAIT"`
				BackoffConf  *BackoffConfig
			} `yaml:"rateLimit"`
		} `yaml:"client"`
		Testing struct {
			DisableHttpsCheck bool `yaml:"disableHttpsCheck" envconfig:"ZWA_TESTING_DISABLE_HTTPS_CHECK"`
		} `yaml:"testing"`
	} `yaml:"zwa"`
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
	cfg.ZWA.Client.RateLimit.MaxRetries = MaxNumOfRetries
	cfg.ZWA.Client.RateLimit.RetryWaitMax = time.Second * time.Duration(RetryWaitMaxSeconds)
	cfg.ZWA.Client.RateLimit.RetryWaitMin = time.Second * time.Duration(RetryWaitMinSeconds)
	cfg.ZWA.Client.RequestTimeout = time.Duration(requestTimeout) * time.Second

	// Read configuration from YAML (lowest precedence)
	readConfigFromSystem(cfg)
	// logger.Printf("[DEBUG] Configuration loaded from system configuration.")

	// Read environment variables (medium precedence)
	readConfigFromEnvironment(cfg)
	// logger.Printf("[DEBUG] Configuration loaded from environment variables.")

	// Apply ConfigSetter functions (highest precedence)
	for _, confSetter := range conf {
		confSetter(cfg)
	}
	// logger.Printf("[DEBUG] Configuration setters applied.")

	// Validate credentials after all sources (YAML, env, setters)
	if cfg.ZWA.Client.ZWAAPIKeyID == "" || cfg.ZWA.Client.ZWAAPISecret == "" {
		logger.Printf("[ERROR] Missing ZWA credentials. Ensure they are provided via setters, environment variables, or YAML configuration.")
		return nil, errors.New("missing required ZWA credentials")
	}

	// Construct and validate the base URL
	var rawBaseURL string
	if cfg.ZWA.Client.ZWACloud != "" {
		rawBaseURL = fmt.Sprintf("https://api.%s.zsworkflow.net", cfg.ZWA.Client.ZWACloud)
	} else {
		rawBaseURL = "https://api.us1.zsworkflow.net" // Default to "us1" if no cloud is specified
	}

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		logger.Printf("[ERROR] Error occurred while configuring the base URL: %v", err)
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	cfg.BaseURL = baseURL
	// logger.Printf("[DEBUG] Base URL configured successfully: %s", cfg.BaseURL.String())

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
	cfg.ZWA.Client.AuthToken = authToken

	// Add the AuthToken to the context
	if cfg.ZWA.Client.AuthToken != nil && cfg.ZWA.Client.AuthToken.AccessToken != "" {
		cfg.Context = context.WithValue(context.Background(), ContextAccessToken, cfg.ZWA.Client.AuthToken.AccessToken)
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
func WithZWAAPIKeyID(keyID string) ConfigSetter {
	return func(c *Configuration) {
		c.ZWA.Client.ZWAAPIKeyID = keyID
	}
}

// WithClientSecret sets the ClientSecret in the Config.
func WithZWAAPISecret(apiSecret string) ConfigSetter {
	return func(c *Configuration) {
		c.ZWA.Client.ZWAAPISecret = apiSecret
	}
}

func WithZWACloud(cloud string) ConfigSetter {
	return func(c *Configuration) {
		c.ZWA.Client.ZWACloud = cloud
	}
}

func WithPartnerID(partnerID string) ConfigSetter {
	return func(c *Configuration) {
		c.ZWA.Client.PartnerID = partnerID
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
		c.ZWA.Client.Proxy.Port = i
	}
}

func WithProxyHost(host string) ConfigSetter {
	return func(c *Configuration) {
		c.ZWA.Client.Proxy.Host = host
	}
}

func WithProxyUsername(username string) ConfigSetter {
	return func(c *Configuration) {
		c.ZWA.Client.Proxy.Username = username
	}
}

func WithProxyPassword(pass string) ConfigSetter {
	return func(c *Configuration) {
		c.ZWA.Client.Proxy.Password = pass
	}
}

func WithTestingDisableHttpsCheck(httpsCheck bool) ConfigSetter {
	return func(c *Configuration) {
		c.ZWA.Testing.DisableHttpsCheck = httpsCheck
	}
}

func WithRequestTimeout(requestTimeout time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZWA.Client.RequestTimeout = requestTimeout
		setHttpClients(c)
	}
}

func WithRateLimitMaxRetries(maxRetries int32) ConfigSetter {
	return func(c *Configuration) {
		c.ZWA.Client.RateLimit.MaxRetries = maxRetries
		setHttpClients(c)
	}
}

func WithRateLimitMaxWait(maxWait time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZWA.Client.RateLimit.RetryWaitMax = maxWait
		setHttpClients(c)
	}
}

func WithRateLimitMinWait(minWait time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZWA.Client.RateLimit.RetryWaitMin = minWait
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
	var log logger.Logger
	if cfg == nil || cfg.Logger == nil {
		log = logger.GetDefaultLogger(loggerPrefix)
		log.Printf("[ERROR] Configuration is nil. Cannot initialize HTTP clients.")
		return
	} else {
		log = cfg.Logger
	}

	// Initialize the global rate limiter (example: 100 requests/min)
	globalLimiter := rl.NewGlobalRateLimiter(100, 60)

	// Configure the HTTP client with rate limiting
	httpClient := &http.Client{
		Transport: &rl.RateLimitTransport{
			GlobalLimiter:   globalLimiter,
			WaitFunc:        globalLimiter.Wait, // Pass the method reference of the limiter
			Logger:          log,
			AdditionalDelay: 5 * time.Second,
		},
	}

	// Assign the rate-limited HTTP client to the configuration
	cfg.HTTPClient = httpClient

	// log.Printf("[DEBUG] HTTP client initialized with global rate limiting.")
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
