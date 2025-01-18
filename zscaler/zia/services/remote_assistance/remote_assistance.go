package remote_assistance

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	remoteAssistanceEndpoint = "/zia/api/v1/remoteAssistance"
)

type RemoteAssistance struct {
	// The time until when view-only access is granted. Unix time is used.
	ViewOnlyUntil int64 `json:"viewOnlyUntil,omitempty"`

	// A Boolean value that indicates whether the user names for single sign-on users should be obfuscated or visible
	// NOTE: This is incorrect: The API swagger description indicates Boolean but the value is int.
	FullAccessUntil int64 `json:"fullAccessUntil,omitempty"`

	// A Boolean value that indicates whether the user names for single sign-on users should be obfuscated or visible
	UsernameObfuscated bool `json:"usernameObfuscated"`

	// A Boolean value that indicates whether the device information (Device Hostname, Device Name, and Device Owner) should be obfuscated or visible on the Dashboard and Analytics pages
	DeviceInfoObfuscate bool `json:"deviceInfoObfuscate"`
}

// GetAll returns the all rules.
func GetRemoteAssistance(ctx context.Context, service *zscaler.Service) (*RemoteAssistance, error) {
	var remoteAssistance RemoteAssistance
	err := service.Client.Read(ctx, remoteAssistanceEndpoint, &remoteAssistance)
	if err != nil {
		return nil, err
	}
	return &remoteAssistance, nil
}

func UpdateRemoteAssistance(ctx context.Context, service *zscaler.Service, remoteAssistance RemoteAssistance) (*RemoteAssistance, *http.Response, error) {
	updated, err := service.Client.UpdateWithPut(ctx, remoteAssistanceEndpoint, remoteAssistance)
	if err != nil {
		return nil, nil, err
	}

	if updated == nil {
		service.Client.GetLogger().Printf("[DEBUG] Remote Assistance Settings updated successfully with no content")
		return nil, nil, nil
	}

	updatedRemoteAssistance, ok := updated.(*RemoteAssistance)
	if !ok {
		return nil, nil, fmt.Errorf("unexpected response type")
	}

	service.Client.GetLogger().Printf("[DEBUG] Updated Remote Assistance Settings: %+v", updatedRemoteAssistance)
	return updatedRemoteAssistance, nil, nil
}
