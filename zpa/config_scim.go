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
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zscaler/zscaler-sdk-go/v2/logger"
	rl "github.com/zscaler/zscaler-sdk-go/v2/ratelimiter"
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
)

type ScimClient struct {
	ScimConfig *ScimConfig
}

type ScimConfig struct {
	BaseURL     *url.URL
	httpClient  *http.Client
	AuthToken   string
	IDPId       string
	Logger      logger.Logger
	rateLimiter *rl.RateLimiter
	BackoffConf *BackoffConfig
	UserAgent   string
}

// NewScimConfig initializes a configuration specifically for SCIM-based API endpoints
func NewScimConfig(scimToken, idpId, scimCloud, userAgent string) (*ScimClient, error) {
	var logger logger.Logger = logger.GetDefaultLogger(loggerPrefix)

	// Load from environment variables if not provided
	if scimToken == "" {
		scimToken = os.Getenv(ZPA_SCIM_TOKEN)
	}
	if idpId == "" {
		idpId = os.Getenv(ZPA_IDP_ID)
	}
	if scimCloud == "" {
		scimCloud = os.Getenv(ZPA_SCIM_CLOUD)
	}

	// Ensure that both scimToken and idpId are provided
	if scimToken == "" || idpId == "" {
		return nil, fmt.Errorf("scim token and idp id are required for SCIM-based configuration")
	}

	// Select the base URL based on the provided SCIM cloud environment
	rawUrl := scimdefaultBaseURL
	switch strings.ToUpper(scimCloud) {
	case "BETA":
		rawUrl = scimbetaBaseURL
	case "ZPATWO":
		rawUrl = scimzpaTwoBaseUrl
	case "GOV":
		rawUrl = scimgovBaseURL
	case "GOVUS":
		rawUrl = scimgovUsBaseURL
	case "PRODUCTION":
		rawUrl = scimdefaultBaseURL
	case "PREVIEW":
		rawUrl = scimpreviewBaseUrl
	}

	// Parse the SCIM base URL
	baseURL, err := url.Parse(rawUrl)
	if err != nil {
		logger.Printf("[ERROR] error occurred while configuring the SCIM client: %v", err)
		return nil, err
	}

	// Create the ScimConfig
	scimConfig := &ScimConfig{
		BaseURL:     baseURL,
		AuthToken:   scimToken,
		IDPId:       idpId,
		Logger:      logger,
		httpClient:  &http.Client{Timeout: defaultTimeout},
		BackoffConf: defaultBackoffConf,
		UserAgent:   userAgent,
		rateLimiter: rl.NewRateLimiter(20, 10, 10, 10),
	}

	// Wrap ScimConfig inside ScimClient and return it
	return &ScimClient{ScimConfig: scimConfig}, nil
}

// DoRequest performs an HTTP request specifically for SCIM endpoints with enhanced logging
func (c *ScimClient) DoRequest(ctx context.Context, method, endpoint string, payload interface{}, target interface{}) (*http.Response, error) {
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
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.ScimConfig.AuthToken))
	if c.ScimConfig.UserAgent != "" {
		req.Header.Add("User-Agent", c.ScimConfig.UserAgent)
	}

	// Log the request
	logger.LogRequest(c.ScimConfig.Logger, req, reqID, nil, true)

	// Send the HTTP request
	resp, err := c.ScimConfig.httpClient.Do(req)
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
	if err := checkErrorInResponse(resp, respData); err != nil {
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
