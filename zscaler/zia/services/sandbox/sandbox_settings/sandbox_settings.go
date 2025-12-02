package sandbox_settings

import (
	"context"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	advancedSettingsEndpoint = "/zia/api/v1/behavioralAnalysisAdvancedSettings"
	fileHashCountEndpoint    = "/fileHashCount"
)

// BaAdvancedSettings represents the legacy/deprecated structure for behavioral analysis settings.
// Kept for backward compatibility as it's unclear if the old model is fully deprecated.
type BaAdvancedSettings struct {
	FileHashesToBeBlocked []string `json:"fileHashesToBeBlocked,omitempty"`
}

// MD5HashType represents the type of MD5 hash entry
type MD5HashType string

const (
	MD5HashTypeMalware             MD5HashType = "MALWARE"
	MD5HashTypeCustomFilehashDeny  MD5HashType = "CUSTOM_FILEHASH_DENY"
	MD5HashTypeCustomFilehashAllow MD5HashType = "CUSTOM_FILEHASH_ALLOW"
)

// MD5HashValue represents an individual MD5 hash entry in the new API structure
type MD5HashValue struct {
	URL        string      `json:"url,omitempty"`
	URLComment string      `json:"urlComment,omitempty"`
	Type       MD5HashType `json:"type,omitempty"`
}

// BaAdvancedSettingsV2 represents the new structure for behavioral analysis advanced settings.
// This is the updated API model that uses md5HashValueList instead of fileHashesToBeBlocked.
// Note: MD5HashValueList does NOT use omitempty to ensure empty arrays are sent as [] (required by API to clear all hashes).
type BaAdvancedSettingsV2 struct {
	MD5HashValueList []MD5HashValue `json:"md5HashValueList"`
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

// GetV2 retrieves behavioral analysis advanced settings using the new API structure.
func GetV2(ctx context.Context, service *zscaler.Service) (*BaAdvancedSettingsV2, error) {
	var hashes BaAdvancedSettingsV2
	err := service.Client.Read(ctx, advancedSettingsEndpoint, &hashes)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning MD5 hash value list from GetV2: %v", hashes)
	return &hashes, nil
}

// UpdateV2 updates behavioral analysis advanced settings using the new API structure.
func UpdateV2(ctx context.Context, service *zscaler.Service, hashes BaAdvancedSettingsV2) (*BaAdvancedSettingsV2, error) {
	_, err := service.Client.UpdateWithPut(ctx, advancedSettingsEndpoint, hashes)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning updated MD5 hash value list from UpdateV2: %v", hashes)
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
