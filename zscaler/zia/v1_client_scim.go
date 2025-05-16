package zia

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
	rl "github.com/zscaler/zscaler-sdk-go/v3/ratelimiter"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"gopkg.in/yaml.v3"
)

const (
	ZIA_SCIM_API_TOKEN = "ZIA_SCIM_API_TOKEN"
	ZIA_SCIM_CLOUD     = "ZIA_SCIM_CLOUD"
	ZIA_SCIM_TENANT_ID = "ZIA_SCIM_TENANT_ID"
)

const loggerScimPrefix = "zia-scim-logger: "

var globalScimConfig *ScimConfiguration
var scimConfigOnce sync.Once

type ScimZiaClient struct {
	ScimConfig *ZIAScimConfig
}

type ZIAScimConfig struct {
	BaseURL     *url.URL
	HTTPClient  *http.Client
	AuthToken   string
	TenantID    string
	Logger      logger.Logger
	RateLimiter *rl.RateLimiter
	UserAgent   string
}

type ScimConfiguration struct {
	sync.Mutex
	Logger         logger.Logger
	HTTPClient     *http.Client
	BaseURL        *url.URL
	DefaultHeader  map[string]string `json:"defaultHeader,omitempty"`
	UserAgent      string            `json:"userAgent,omitempty"`
	Debug          bool              `json:"debug,omitempty"`
	UserAgentExtra string
	Context        context.Context
	ZIAScim        struct {
		Client struct {
			ZIAScimApiToken string `yaml:"zia_scim_api_token" envconfig:"ZIA_SCIM_API_TOKEN"`
			ZIAScimCloud    string `yaml:"zia_scim_cloud" envconfig:"ZIA_SCIM_CLOUD"`
			ZIAScimTenantID string `yaml:"zia_scim_tenant_id" envconfig:"ZIA_SCIM_TENANT_ID"`
		} `yaml:"client"`
	} `yaml:"ziaScim"`
}

type ScimConfigSetter func(*ScimConfiguration)

func NewScimConfig(setters ...ScimConfigSetter) (*ScimConfiguration, error) {
	scimConfigOnce.Do(func() {
		logger := logger.GetDefaultLogger(loggerScimPrefix)
		logger.Printf("[DEBUG] Initializing SCIM config")

		globalScimConfig = &ScimConfiguration{
			DefaultHeader: make(map[string]string),
			Logger:        logger,
			UserAgent:     fmt.Sprintf("zscaler-sdk-go/%s golang/%s %s/%s", VERSION, runtime.Version(), runtime.GOOS, runtime.GOARCH),
			Debug:         false,
			Context:       context.Background(),
		}

		for _, setter := range setters {
			setter(globalScimConfig)
		}

		// Read configuration from system and environment
		readScimConfigFromSystem(globalScimConfig)
		readScimConfigFromEnvironment(globalScimConfig)

		if globalScimConfig.ZIAScim.Client.ZIAScimApiToken == "" {
			globalScimConfig.ZIAScim.Client.ZIAScimApiToken = os.Getenv(ZIA_SCIM_API_TOKEN)
		}
		if globalScimConfig.ZIAScim.Client.ZIAScimTenantID == "" {
			globalScimConfig.ZIAScim.Client.ZIAScimTenantID = os.Getenv(ZIA_SCIM_TENANT_ID)
		}
		if globalScimConfig.ZIAScim.Client.ZIAScimApiToken == "" || globalScimConfig.ZIAScim.Client.ZIAScimTenantID == "" {
			logger.Printf("[ERROR] Missing SCIM API token or tenant ID")
			return
		}

		if globalScimConfig.BaseURL == nil {
			cloud := globalScimConfig.ZIAScim.Client.ZIAScimCloud
			if cloud == "" {
				cloud = os.Getenv(ZIA_SCIM_CLOUD)
			}
			if cloud == "" {
				logger.Printf("[ERROR] Missing SCIM cloud configuration")
				return
			}
			rawURL := fmt.Sprintf("https://scim.%s.net/%s/scim", cloud, globalScimConfig.ZIAScim.Client.ZIAScimTenantID)
			parsedURL, err := url.Parse(rawURL)
			if err != nil {
				logger.Printf("[ERROR] Failed to parse SCIM base URL: %v", err)
				return
			}
			globalScimConfig.BaseURL = parsedURL
			logger.Printf("[DEBUG] Constructed SCIM base URL: %s", parsedURL.String())
		}

		if globalScimConfig.HTTPClient == nil {
			globalScimConfig.HTTPClient = &http.Client{}
		}

		if globalScimConfig.Context == nil {
			globalScimConfig.Context = context.Background()
		}
	})

	if globalScimConfig == nil {
		return nil, fmt.Errorf("failed to initialize SCIM configuration")
	}
	return globalScimConfig, nil
}

// WithScimToken sets the SCIM token in the configuration.
func WithScimToken(scimToken string) ScimConfigSetter {
	return func(c *ScimConfiguration) {
		c.ZIAScim.Client.ZIAScimApiToken = scimToken
	}
}

func WithTenantID(tenantID string) ScimConfigSetter {
	return func(c *ScimConfiguration) {
		c.ZIAScim.Client.ZIAScimTenantID = tenantID
	}
}

// WithScimBaseURL sets the SCIM BaseURL in the configuration.
func WithScimCloud(input string) ScimConfigSetter {
	return func(c *ScimConfiguration) {
		// If the input does not look like a full URL, treat it as a cloud identifier
		if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
			if c.ZIAScim.Client.ZIAScimTenantID == "" {
				c.ZIAScim.Client.ZIAScimTenantID = os.Getenv(ZIA_SCIM_TENANT_ID)
			}
			if c.ZIAScim.Client.ZIAScimTenantID == "" {
				c.Logger.Printf("[ERROR] Missing tenant ID for dynamic SCIM URL construction")
				return
			}
			constructed := fmt.Sprintf("https://scim.%s.net/%s/scim", input, c.ZIAScim.Client.ZIAScimTenantID)
			parsed, err := url.Parse(constructed)
			if err != nil {
				c.Logger.Printf("[ERROR] Failed to parse constructed URL from cloud name: %v", err)
				return
			}
			c.BaseURL = parsed
			c.Logger.Printf("[DEBUG] Constructed SCIM Base URL from cloud: %s", parsed.String())
			return
		}

		// Fall back to full URL parsing
		parsed, err := url.Parse(input)
		if err == nil {
			c.BaseURL = parsed
			c.Logger.Printf("[DEBUG] Using explicitly provided full Base URL: %s", parsed.String())
		} else {
			c.Logger.Printf("[ERROR] Invalid Base URL: %v", err)
		}
	}
}

// WithScimUserAgent sets the User-Agent in the configuration.
func WithScimUserAgent(userAgent string) ScimConfigSetter {
	return func(c *ScimConfiguration) {
		c.UserAgent = userAgent
	}
}

// DoRequest performs an HTTP request specifically for SCIM endpoints with enhanced logging

func (c *ScimZiaClient) DoRequest(ctx context.Context, method, endpoint string, payload interface{}, target interface{}) (*http.Response, error) {
	// Create a copy of the base URL to avoid modifying the original
	requestURL, err := url.Parse(c.ScimConfig.BaseURL.String())
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	// Ensure the endpoint path is properly joined
	endpoint = strings.TrimPrefix(endpoint, "/")
	if endpoint != "" {
		requestURL.Path = path.Join(requestURL.Path, endpoint)
	}

	fullURL := requestURL.String()
	c.ScimConfig.Logger.Printf("[DEBUG] Making request to: %s", fullURL)

	reqID := uuid.NewString()
	start := time.Now()

	var reqBody io.Reader
	if payload != nil {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			c.ScimConfig.Logger.Printf("[ERROR] Failed to marshal payload: %v", err)
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonPayload)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		c.ScimConfig.Logger.Printf("[ERROR] Failed to create request: %v", err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/scim+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.ScimConfig.AuthToken))
	if c.ScimConfig.UserAgent != "" {
		req.Header.Add("User-Agent", c.ScimConfig.UserAgent)
	}

	logger.LogRequest(c.ScimConfig.Logger, req, reqID, nil, true)

	if c.ScimConfig.RateLimiter != nil {
		shouldWait, delay := c.ScimConfig.RateLimiter.Wait(method)
		if shouldWait {
			c.ScimConfig.Logger.Printf("[DEBUG] Rate limiter triggered. Sleeping for %v", delay)
			time.Sleep(delay)
		}
	}

	resp, err := c.ScimConfig.HTTPClient.Do(req)
	if err != nil {
		c.ScimConfig.Logger.Printf("[ERROR] Error occurred during request: %v", err)
		return nil, err
	}

	respData, err := io.ReadAll(resp.Body)
	if err == nil {
		resp.Body = io.NopCloser(bytes.NewBuffer(respData))
	}

	if err := errorx.CheckErrorInResponse(resp, err); err != nil {
		return resp, err
	}

	if target != nil {
		if err := decodeJSON(respData, target); err != nil {
			c.ScimConfig.Logger.Printf("[ERROR] Failed to decode response: %v", err)
			return resp, err
		}
	}

	logger.LogResponse(c.ScimConfig.Logger, resp, start, reqID)
	unescapeHTML(target)

	return resp, nil
}

func decodeJSON(respData []byte, v interface{}) error {
	return json.NewDecoder(bytes.NewBuffer(respData)).Decode(&v)
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
	_ = json.Unmarshal(data, entity)
}

func readScimConfigFromFile(location string, c *ScimConfiguration) (*ScimConfiguration, error) {
	yamlConfig, err := os.ReadFile(location)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlConfig, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func readScimConfigFromSystem(c *ScimConfiguration) *ScimConfiguration {
	currUser, err := user.Current()
	if err != nil {
		return c
	}
	if currUser.HomeDir == "" {
		return c
	}
	conf, err := readScimConfigFromFile(currUser.HomeDir+"/.zscaler/zscaler.yaml", c)
	if err != nil {
		return c
	}
	return conf
}

func readScimConfigFromEnvironment(c *ScimConfiguration) *ScimConfiguration {
	err := envconfig.Process("zscaler", c)
	if err != nil {
		fmt.Println("error parsing")
		return c
	}
	return c
}
