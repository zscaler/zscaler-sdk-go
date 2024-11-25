package sandbox_settings

import (
	"context"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
)

const (
	advancedSettingsEndpoint = "/zia/api/v1/behavioralAnalysisAdvancedSettings"
	fileHashCountEndpoint    = "/fileHashCount"
)

type BaAdvancedSettings struct {
	FileHashesToBeBlocked []string `json:"fileHashesToBeBlocked,omitempty"`
}

type FileHashCount struct {
	BlockedFileHashesCount int `json:"blockedFileHashesCount,omitempty"`
	RemainingFileHashes    int `json:"remainingFileHashes,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service) (*BaAdvancedSettings, error) {
	var hashes BaAdvancedSettings
	err := service.Client.Read(ctx, advancedSettingsEndpoint, &hashes)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning custom list of MD5 hashes from Get: %v", hashes)
	return &hashes, nil
}

func Update(ctx context.Context, service *zscaler.Service, hashes BaAdvancedSettings) (*BaAdvancedSettings, error) {
	_, err := service.Client.UpdateWithPut(ctx, advancedSettingsEndpoint, hashes)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning updated custom list of MD5 hashes from Get: %v", hashes)
	return &hashes, nil
}

func GetFileHashCount(ctx context.Context, service *zscaler.Service) (*FileHashCount, error) {
	var hashes FileHashCount
	err := service.Client.Read(ctx, advancedSettingsEndpoint+fileHashCountEndpoint, &hashes)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning used andd unused quota for blocking MD5 file hashes from Get: %v", hashes)
	return &hashes, nil
}
