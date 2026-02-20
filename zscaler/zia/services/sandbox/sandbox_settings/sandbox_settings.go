package sandbox_settings

import (
	"context"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
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

// Md5HashValue represents a single MD5 hash entry in the sandbox settings
type Md5HashValue struct {
	URL        string `json:"url,omitempty"`
	URLComment string `json:"urlComment,omitempty"`
	Type       string `json:"type,omitempty"` // e.g. "MALWARE"
}

// Md5HashValueListPayload is the request/response payload for MD5 hash value list operations
type Md5HashValueListPayload struct {
	// Do not use omitempty - API requires md5HashValueList to be present (even as []) when clearing the list
	Md5HashValueList []Md5HashValue `json:"md5HashValueList"`
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

func Getv2(ctx context.Context, service *zscaler.Service) (*Md5HashValueListPayload, error) {
	var payload Md5HashValueListPayload
	err := service.Client.Read(ctx, advancedSettingsEndpoint, &payload)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning MD5 hash value list from Getv2: %v", payload)
	return &payload, nil
}

func Updatev2(ctx context.Context, service *zscaler.Service, payload Md5HashValueListPayload) (*Md5HashValueListPayload, error) {
	_, err := service.Client.UpdateWithPut(ctx, advancedSettingsEndpoint, payload)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning updated MD5 hash value list from Updatev2: %v", payload)
	return &payload, nil
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
