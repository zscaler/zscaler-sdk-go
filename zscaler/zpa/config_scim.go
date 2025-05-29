package zpa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/user"
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
	ZPA_SCIM_TOKEN     = "ZPA_SCIM_TOKEN"
	ZPA_IDP_ID         = "ZPA_IDP_ID"
	ZPA_SCIM_CLOUD     = "ZPA_SCIM_CLOUD"
	scimdefaultBaseURL = "https://scim1.private.zscaler.com/scim/1/"
	scimbetaBaseURL    = "https://scim1.zpabeta.ne/scim/1/"
	scimzpaTwoBaseUrl  = "https://scim1.zpatwo.net/scim/1/"
	scimgovBaseURL     = "https://scim1.zpagov.net/scim/1/"
	scimgovUsBaseURL   = "https://scim1.zpagov.us/scim/1/"
	scimpreviewBaseUrl = "https://scim1.zpapreview.net/scim/1/"
	defaultTimeout     = 240 * time.Second
)

const loggerScimPrefix = "zpa-scim-logger: "

var globalScimConfig *ScimConfiguration
var scimConfigOnce sync.Once

type ScimZpaClient struct {
	ScimConfig *ZPAScimConfig
}

type ZPAScimConfig struct {
	BaseURL     *url.URL
	HTTPClient  *http.Client
	AuthToken   string
	IDPId       string
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
	ZPAScim        struct {
		Client struct {
			ZPAScimToken string `yaml:"zpa_scim_token" envconfig:"ZPA_SCIM_TOKEN"`
			ZPAIdPID     string `yaml:"zpa_idp_id" envconfig:"ZPA_IDP_ID"`
			ZPAScimCloud string `yaml:"zpa_scim_cloud" envconfig:"ZPA_SCIM_CLOUD"`
		} `yaml:"client"`
	} `yaml:"zpaScim"`
}

type ScimConfigSetter func(*ScimConfiguration)

func NewScimConfig(setters ...ScimConfigSetter) (*ScimConfiguration, error) {
	var initErr error

	scimConfigOnce.Do(func() {
		logger := logger.GetDefaultLogger(loggerScimPrefix)

		globalScimConfig = &ScimConfiguration{
			DefaultHeader: make(map[string]string),
			Logger:        logger,
			HTTPClient:    &http.Client{Timeout: defaultTimeout},
			UserAgent: fmt.Sprintf("zscaler-sdk-go/%s golang/%s %s/%s",
				VERSION, runtime.Version(), runtime.GOOS, runtime.GOARCH),
			Context: context.Background(),
		}

		for _, setter := range setters {
			setter(globalScimConfig)
		}

		readScimConfigFromSystem(globalScimConfig)
		readScimConfigFromEnvironment(globalScimConfig)

		// Fallback to environment variables
		if globalScimConfig.ZPAScim.Client.ZPAScimToken == "" {
			globalScimConfig.ZPAScim.Client.ZPAScimToken = os.Getenv(ZPA_SCIM_TOKEN)
		}
		if globalScimConfig.ZPAScim.Client.ZPAIdPID == "" {
			globalScimConfig.ZPAScim.Client.ZPAIdPID = os.Getenv(ZPA_IDP_ID)
		}

		// Validate required fields
		if globalScimConfig.ZPAScim.Client.ZPAScimToken == "" || globalScimConfig.ZPAScim.Client.ZPAIdPID == "" {
			initErr = fmt.Errorf("scim token and idp id are required")
			globalScimConfig = nil
			return
		}

		// Only set BaseURL if not already set by WithScimCloud
		if globalScimConfig.BaseURL == nil {
			cloud := globalScimConfig.ZPAScim.Client.ZPAScimCloud
			if cloud == "" {
				cloud = os.Getenv(ZPA_SCIM_CLOUD)
				if cloud == "" {
					cloud = "PRODUCTION"
				}
			}

			var baseURL string
			switch strings.ToUpper(cloud) {
			case "BETA":
				baseURL = scimbetaBaseURL
			case "ZPATWO":
				baseURL = scimzpaTwoBaseUrl
			case "GOV":
				baseURL = scimgovBaseURL
			case "GOVUS":
				baseURL = scimgovUsBaseURL
			case "PREVIEW":
				baseURL = scimpreviewBaseUrl
			default:
				baseURL = scimdefaultBaseURL
			}

			// Incorporate the IDP ID into the URL path
			fullURL := fmt.Sprintf("%s%s/", baseURL, globalScimConfig.ZPAScim.Client.ZPAIdPID)
			parsedURL, err := url.Parse(fullURL)
			if err != nil {
				initErr = fmt.Errorf("failed to parse SCIM base URL: %w", err)
				globalScimConfig = nil
				return
			}
			globalScimConfig.BaseURL = parsedURL
			logger.Printf("[DEBUG] Constructed SCIM base URL: %s", parsedURL.String())
		}
	})

	if initErr != nil {
		return nil, initErr
	}
	if globalScimConfig == nil {
		return nil, fmt.Errorf("failed to initialize SCIM configuration")
	}
	return globalScimConfig, nil
}

// WithScimToken sets the SCIM token in the configuration.
func WithScimToken(scimToken string) ScimConfigSetter {
	return func(c *ScimConfiguration) {
		c.ZPAScim.Client.ZPAScimToken = scimToken
	}
}

// WithIDPId sets the IDP ID in the configuration.
func WithIDPId(idpId string) ScimConfigSetter {
	return func(c *ScimConfiguration) {
		c.ZPAScim.Client.ZPAIdPID = idpId
	}
}

// WithScimBaseURL sets the SCIM BaseURL in the configuration.
func WithScimCloud(env string) ScimConfigSetter {
	return func(c *ScimConfiguration) {
		var raw string
		// full URL?
		if strings.HasPrefix(env, "http://") || strings.HasPrefix(env, "https://") {
			raw = env
		} else {
			// cloud name → constant
			switch strings.ToUpper(env) {
			case "BETA":
				raw = scimbetaBaseURL
			case "ZPATWO":
				raw = scimzpaTwoBaseUrl
			case "GOV":
				raw = scimgovBaseURL
			case "GOVUS":
				raw = scimgovUsBaseURL
			case "PREVIEW":
				raw = scimpreviewBaseUrl
			default:
				raw = scimdefaultBaseURL
			}
		}
		u, err := url.Parse(raw)
		if err != nil {
			c.Logger.Printf("[ERROR] invalid SCIM cloud/url %q: %v", env, err)
			return
		}
		c.BaseURL = u
		c.Logger.Printf("[DEBUG] SCIM BaseURL set to %s", u)
	}
}

// WithScimUserAgent sets the User-Agent in the configuration.
func WithScimUserAgent(userAgent string) ScimConfigSetter {
	return func(c *ScimConfiguration) {
		c.UserAgent = userAgent
	}
}

// DoRequest performs an HTTP request specifically for SCIM endpoints with enhanced logging
func (c *ScimZpaClient) DoRequest(ctx context.Context, method, endpoint string, payload interface{}, target interface{}) (*http.Response, error) {
	fullURL := fmt.Sprintf("%s%s", c.ScimConfig.BaseURL.String(), endpoint)
	reqID := uuid.NewString() // Generate a unique request ID
	start := time.Now()

	// Marshal payload if provided
	var reqBody io.Reader
	if payload != nil {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			c.ScimConfig.Logger.Printf("[ERROR] Failed to marshal payload: %v", err)
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonPayload)
	}

	// Create the HTTP request with context
	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		c.ScimConfig.Logger.Printf("[ERROR] Failed to create request: %v", err)
		return nil, err
	}

	// Add headers, including the Authorization token
	req.Header.Add("Content-Type", "application/scim+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.ScimConfig.AuthToken))
	if c.ScimConfig.UserAgent != "" {
		req.Header.Add("User-Agent", c.ScimConfig.UserAgent)
	}

	// Log the request
	logger.LogRequest(c.ScimConfig.Logger, req, reqID, nil, true)

	// Send the HTTP request
	resp, err := c.ScimConfig.HTTPClient.Do(req)
	if err != nil {
		c.ScimConfig.Logger.Printf("[ERROR] Error occurred during request: %v", err)
		return nil, err
	}

	// Read and log the response data
	respData, err := io.ReadAll(resp.Body)
	if err == nil {
		resp.Body = io.NopCloser(bytes.NewBuffer(respData))
	}

	// Check for errors in the response status or body
	if err := errorx.CheckErrorInResponse(resp, err); err != nil {
		return resp, err
	}

	// Decode JSON into target if provided and no error
	if target != nil {
		if err := decodeJSON(respData, target); err != nil {
			c.ScimConfig.Logger.Printf("[ERROR] Failed to decode response: %v", err)
			return resp, err
		}
	}

	// Log the response details with the same reqID
	logger.LogResponse(c.ScimConfig.Logger, resp, start, reqID)

	// Optional: Unescape HTML in target, if necessary
	unescapeHTML(target)

	c.ScimConfig.Logger.Printf("[DEBUG] Successfully completed request to %s", fullURL)
	return resp, nil
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
