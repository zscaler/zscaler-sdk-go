package zcc

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
	maxIdleConnections  int = 40
	requestTimeout      int = 60
	contentTypeJSON         = "application/json"
	MaxNumOfRetries         = 50
	RetryWaitMaxSeconds     = 20
	RetryWaitMinSeconds     = 5
	loggerPrefix            = "zcc-logger: "
)

const (
	VERSION           = "3.5.0"
	ZCC_CLIENT_ID     = "ZCC_CLIENT_ID"
	ZCC_CLIENT_SECRET = "ZCC_CLIENT_SECRET"
	ZCC_CLOUD         = "ZCC_CLOUD"
)

type AuthRequest struct {
	APIKey    string `json:"apiKey"`
	SecretKey string `json:"secretKey"`
}

type BackoffConfig struct {
	Enabled             bool // Set to true to enable backoff and retry mechanism
	RetryWaitMinSeconds int  // Minimum time to wait
	RetryWaitMaxSeconds int  // Maximum time to wait
	MaxNumOfRetries     int  // Maximum number of retries
}

type AuthToken struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"jwtToken"`
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
	ZCC            struct {
		Client struct {
			ZCCClientID     string     `yaml:"apiKey" envconfig:"ZCC_CLIENT_ID"`
			ZCCClientSecret string     `yaml:"secretKey" envconfig:"ZCC_CLIENT_SECRET"`
			ZCCCloud        string     `yaml:"cloud" envconfig:"ZCC_CLOUD"`
			AuthToken       *AuthToken `yaml:"authToken"`
			AccessToken     *AuthToken `yaml:"accessToken"`
			Proxy           struct {
				Port     int32  `yaml:"port" envconfig:"ZCC_CLIENT_PROXY_PORT"`
				Host     string `yaml:"host" envconfig:"ZCC_CLIENT_PROXY_HOST"`
				Username string `yaml:"username" envconfig:"ZCC_CLIENT_PROXY_USERNAME"`
				Password string `yaml:"password" envconfig:"ZCC_CLIENT_PROXY_PASSWORD"`
			} `yaml:"proxy"`
			RequestTimeout time.Duration `yaml:"requestTimeout" envconfig:"ZCC_CLIENT_REQUEST_TIMEOUT"`
			RateLimit      struct {
				MaxRetries   int32         `yaml:"maxRetries" envconfig:"ZCC_CLIENT_RATE_LIMIT_MAX_RETRIES"`
				RetryWaitMin time.Duration `yaml:"minWait" envconfig:"ZCC_CLIENT_RATE_LIMIT_MIN_WAIT"`
				RetryWaitMax time.Duration `yaml:"maxWait" envconfig:"ZCC_CLIENT_RATE_LIMIT_MAX_WAIT"`
				BackoffConf  *BackoffConfig
			} `yaml:"rateLimit"`
		} `yaml:"client"`
		Testing struct {
			DisableHttpsCheck bool `yaml:"disableHttpsCheck" envconfig:"ZCC_TESTING_DISABLE_HTTPS_CHECK"`
		} `yaml:"testing"`
	} `yaml:"zcc"`
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
	cfg.ZCC.Client.RateLimit.MaxRetries = MaxNumOfRetries
	cfg.ZCC.Client.RateLimit.RetryWaitMax = time.Second * time.Duration(RetryWaitMaxSeconds)
	cfg.ZCC.Client.RateLimit.RetryWaitMin = time.Second * time.Duration(RetryWaitMinSeconds)
	cfg.ZCC.Client.RequestTimeout = time.Duration(requestTimeout) * time.Second

	// Read configuration from YAML (lowest precedence)
	readConfigFromSystem(cfg)

	// Read environment variables (medium precedence)
	readConfigFromEnvironment(cfg)

	// Apply ConfigSetter functions (highest precedence)
	for _, confSetter := range conf {
		confSetter(cfg)
	}
	logger.Printf("[DEBUG] Configuration setters applied.")

	// Validate credentials after all sources (YAML, env, setters)
	if cfg.ZCC.Client.ZCCClientID == "" || cfg.ZCC.Client.ZCCClientSecret == "" || cfg.ZCC.Client.ZCCCloud == "" {
		logger.Printf("[ERROR] Missing ZCC credentials. Ensure they are provided via setters, environment variables, or YAML configuration.")
		return nil, errors.New("missing required ZCC credentials")
	}

	// Construct and validate the base URL
	rawBaseURL := fmt.Sprintf("https://api-mobile.%s.net/papi", cfg.ZCC.Client.ZCCCloud)
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		logger.Printf("[ERROR] Error occurred while configuring the base URL: %v", err)
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	cfg.BaseURL = baseURL
	logger.Printf("[DEBUG] Base URL configured successfully: %s", cfg.BaseURL.String())

	// Set up HTTP clients
	setHttpClients(cfg)
	if cfg.HTTPClient == nil {
		logger.Printf("[ERROR] HTTP clients not initialized")
		return nil, errors.New("HTTP clients not initialized")
	}
	logger.Printf("[DEBUG] HTTP clients configured.")

	// Authenticate the client and populate the AuthToken
	authToken, err := Authenticate(cfg.Context, cfg, logger)
	if err != nil {
		logger.Printf("[ERROR] Authentication failed: %v", err)
		return nil, fmt.Errorf("authentication failed: %w", err)
	}
	cfg.ZCC.Client.AuthToken = authToken

	// Add the AuthToken to the context
	if cfg.ZCC.Client.AuthToken != nil && cfg.ZCC.Client.AuthToken.AccessToken != "" {
		cfg.Context = context.WithValue(context.Background(), ContextAccessToken, cfg.ZCC.Client.AuthToken.AccessToken)
		logger.Printf("[DEBUG] AuthToken added to context.")
	} else {
		logger.Printf("[ERROR] Failed to set AuthToken in context.")
		return nil, errors.New("AuthToken is missing or invalid after authentication")
	}

	logger.Printf("[DEBUG] Configuration successfully initialized.")
	return cfg, nil
}

type ConfigSetter func(*Configuration)

// ConfigSetter type defines a function that modifies a Config struct.
// WithClientID sets the ClientID in the Config.
func WithZCCClientID(clientID string) ConfigSetter {
	return func(c *Configuration) {
		c.ZCC.Client.ZCCClientID = clientID
	}
}

// WithClientSecret sets the ClientSecret in the Config.
func WithZCCClientSecret(clientSecret string) ConfigSetter {
	return func(c *Configuration) {
		c.ZCC.Client.ZCCClientSecret = clientSecret
	}
}

func WithZCCCloud(cloud string) ConfigSetter {
	return func(c *Configuration) {
		c.ZCC.Client.ZCCCloud = cloud
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
		c.ZCC.Client.Proxy.Port = i
	}
}

func WithProxyHost(host string) ConfigSetter {
	return func(c *Configuration) {
		c.ZCC.Client.Proxy.Host = host
	}
}

func WithProxyUsername(username string) ConfigSetter {
	return func(c *Configuration) {
		c.ZCC.Client.Proxy.Username = username
	}
}

func WithProxyPassword(pass string) ConfigSetter {
	return func(c *Configuration) {
		c.ZCC.Client.Proxy.Password = pass
	}
}

func WithTestingDisableHttpsCheck(httpsCheck bool) ConfigSetter {
	return func(c *Configuration) {
		c.ZCC.Testing.DisableHttpsCheck = httpsCheck
	}
}

func WithRequestTimeout(requestTimeout time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZCC.Client.RequestTimeout = requestTimeout
		setHttpClients(c)
	}
}

func WithRateLimitMaxRetries(maxRetries int32) ConfigSetter {
	return func(c *Configuration) {
		c.ZCC.Client.RateLimit.MaxRetries = maxRetries
		setHttpClients(c)
	}
}

func WithRateLimitMaxWait(maxWait time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZCC.Client.RateLimit.RetryWaitMax = maxWait
		setHttpClients(c)
	}
}

func WithRateLimitMinWait(minWait time.Duration) ConfigSetter {
	return func(c *Configuration) {
		c.ZCC.Client.RateLimit.RetryWaitMin = minWait
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

	// ZCC-specific rate limits:
	zccRateLimiter := rl.NewRateLimiter(100, 3, 3600, 86400) // General: 100 per hour, downloadDevices: 3 per day

	// Configure the ZCC HTTP client
	cfg.HTTPClient = getHTTPClient(log, zccRateLimiter, cfg)
	if cfg.HTTPClient == nil {
		log.Printf("[ERROR] Failed to initialize ZCC HTTP client.")
	} else {
		log.Printf("[DEBUG] ZCC HTTP client initialized successfully.")
	}

	// Remove the call that overwrote the ZCC client:
	// cfg.HTTPClient = getHTTPClient(log, nil, cfg)
	if cfg.HTTPClient == nil {
		log.Printf("[ERROR] Failed to initialize generic HTTP client.")
	} else {
		log.Printf("[DEBUG] Generic HTTP client initialized successfully.")
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
	c.Logger.Printf("[DEBUG] Successfully parsed environment variables.")
	return c
}

// AddDefaultHeader adds a new HTTP header to the default header in the request
func (c *Configuration) AddDefaultHeader(key string, value string) {
	c.DefaultHeader[key] = value
}
