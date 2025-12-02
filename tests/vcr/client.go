package vcr

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

// bodyBufferingTransport wraps an http.RoundTripper and ensures request bodies are buffered
// before being passed to the wrapped transport. This is necessary because go-vcr reads
// the body during recording, which would leave an empty body for the actual HTTP request.
// By buffering here, we ensure the body can be read multiple times.
type bodyBufferingTransport struct {
	wrapped http.RoundTripper
}

func (t *bodyBufferingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Buffer the request body if present
	if req.Body != nil && req.Body != http.NoBody {
		bodyBytes, err := io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to buffer request body: %w", err)
		}
		// Replace with a new reader
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		req.ContentLength = int64(len(bodyBytes))
		// Also set GetBody so the body can be re-read on redirects or retries
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(bodyBytes)), nil
		}
	}
	return t.wrapped.RoundTrip(req)
}

// VCRClient wraps a Zscaler service client with VCR recorder
type VCRClient struct {
	Service  *zscaler.Service
	Recorder *recorder.Recorder
}

// Stop stops the VCR recorder and saves the cassette
func (c *VCRClient) Stop() error {
	if c.Recorder != nil {
		return c.Recorder.Stop()
	}
	return nil
}

// NewVCRClient creates a new Zscaler client with VCR recording/playback
func NewVCRClient(t *testing.T, cassetteName string, service string) (*VCRClient, error) {
	// Create VCR recorder
	rec, err := NewVCRRecorder(t, VCRConfig{
		CassetteName: cassetteName,
		Service:      service,
		Mode:         GetRecordMode(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create VCR recorder: %w", err)
	}

	// Create HTTP client with body-buffering transport wrapping the VCR recorder
	// The body buffering transport ensures request bodies are buffered BEFORE
	// go-vcr processes them, so the body can be read by both VCR and the real HTTP request
	httpClient := &http.Client{
		Transport: &bodyBufferingTransport{wrapped: rec},
	}

	// Get credentials - in playback mode, we can use dummy credentials
	clientID := os.Getenv("ZSCALER_CLIENT_ID")
	clientSecret := os.Getenv("ZSCALER_CLIENT_SECRET")
	vanityDomain := os.Getenv("ZSCALER_VANITY_DOMAIN")
	customerID := os.Getenv("ZPA_CUSTOMER_ID")
	zscalerCloud := os.Getenv("ZSCALER_CLOUD")

	// In playback mode, provide dummy credentials if not set
	if IsPlaybackMode() {
		if clientID == "" {
			clientID = "dummy_client_id"
		}
		if clientSecret == "" {
			clientSecret = "dummy_client_secret"
		}
		if vanityDomain == "" {
			vanityDomain = "dummy_vanity_domain"
		}
		if customerID == "" {
			// Use the real customer ID from recordings for ZPA
			customerID = "216196257331281920"
		}
	} else {
		// In recording mode, require real credentials
		var missing []string
		if clientID == "" {
			missing = append(missing, "ZSCALER_CLIENT_ID")
		}
		if clientSecret == "" {
			missing = append(missing, "ZSCALER_CLIENT_SECRET")
		}
		if vanityDomain == "" {
			missing = append(missing, "ZSCALER_VANITY_DOMAIN")
		}
		if len(missing) > 0 {
			rec.Stop()
			return nil, fmt.Errorf("recording mode requires environment variables: %s", strings.Join(missing, ", "))
		}
	}

	// Build configuration with VCR HTTP client
	config, err := zscaler.NewConfiguration(
		zscaler.WithClientID(clientID),
		zscaler.WithClientSecret(clientSecret),
		zscaler.WithVanityDomain(vanityDomain),
		zscaler.WithZPACustomerID(customerID),
		zscaler.WithZscalerCloud(zscalerCloud),
		zscaler.WithTestingDisableHttpsCheck(true),
	)
	if err != nil {
		rec.Stop()
		return nil, fmt.Errorf("failed to create configuration: %w", err)
	}

	// Override ALL HTTP clients with VCR transport
	// This is necessary because the SDK uses separate HTTP clients per service
	config.HTTPClient = httpClient
	config.ZPAHTTPClient = httpClient
	config.ZIAHTTPClient = httpClient
	config.ZTWHTTPClient = httpClient
	config.ZCCHTTPClient = httpClient
	config.ZDXHTTPClient = httpClient

	// Create the client
	client, err := zscaler.NewOneAPIClient(config)
	if err != nil {
		rec.Stop()
		return nil, fmt.Errorf("failed to create OneAPI client: %w", err)
	}

	return &VCRClient{
		Service:  client,
		Recorder: rec,
	}, nil
}

// NewZPAVCRClient is a convenience function for ZPA tests
func NewZPAVCRClient(t *testing.T, cassetteName string) (*VCRClient, error) {
	return NewVCRClient(t, cassetteName, "zpa")
}

// NewZIAVCRClient is a convenience function for ZIA tests
func NewZIAVCRClient(t *testing.T, cassetteName string) (*VCRClient, error) {
	return NewVCRClient(t, cassetteName, "zia")
}

// NewZCCVCRClient is a convenience function for ZCC tests
func NewZCCVCRClient(t *testing.T, cassetteName string) (*VCRClient, error) {
	return NewVCRClient(t, cassetteName, "zcc")
}

// NewZDXVCRClient is a convenience function for ZDX tests
func NewZDXVCRClient(t *testing.T, cassetteName string) (*VCRClient, error) {
	return NewVCRClient(t, cassetteName, "zdx")
}

// NewZTWVCRClient is a convenience function for ZTW tests
func NewZTWVCRClient(t *testing.T, cassetteName string) (*VCRClient, error) {
	return NewVCRClient(t, cassetteName, "ztw")
}
