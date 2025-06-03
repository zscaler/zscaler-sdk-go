package zscaler

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserAgent(t *testing.T) {
	configuration, err := NewConfiguration()
	require.NoError(t, err, "Creating a new config should not error")
	userAgent := "zscaler-sdk-go/" + VERSION + " golang/" + runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH
	require.Equal(t, userAgent, configuration.UserAgent)
}

func TestUserAgentWithExtra(t *testing.T) {
	configuration, err := NewConfiguration(
		WithUserAgentExtra("extra/info"),
	)
	require.NoError(t, err, "Creating a new config should not error")
	userAgent := "zscaler-sdk-go/" + VERSION + " golang/" + runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH + " extra/info"
	require.Equal(t, userAgent, configuration.UserAgent)
}

func TestDetectServiceTypeUnknown(t *testing.T) {
	_, err := detectServiceType("/foo")
	require.Error(t, err)
}

func TestGetServiceHTTPClientUnknown(t *testing.T) {
	cfg, err := NewConfiguration()
	require.NoError(t, err)

	svc, err := NewOneAPIClient(cfg)
	require.NoError(t, err)

	generic := cfg.HTTPClient

	httpClient := svc.Client.getServiceHTTPClient("/foo")
	require.Equal(t, generic, httpClient)
}
