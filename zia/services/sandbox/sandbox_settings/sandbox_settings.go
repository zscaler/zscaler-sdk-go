package sandbox_settings

import "github.com/zscaler/zscaler-sdk-go/v2/zia/services"

const (
	advancedSettingsEndpoint = "/behavioralAnalysisAdvancedSettings"
	fileHashCountEndpoint    = "/fileHashCount"
)

type BaAdvancedSettings struct {
	FileHashesToBeBlocked []string `json:"fileHashesToBeBlocked,omitempty"`
}

type FileHashCount struct {
	BlockedFileHashesCount int `json:"blockedFileHashesCount,omitempty"`
	RemainingFileHashes    int `json:"remainingFileHashes,omitempty"`
}

func Get(service *services.Service) (*BaAdvancedSettings, error) {
	var hashes BaAdvancedSettings
	err := service.Client.Read(advancedSettingsEndpoint, &hashes)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG] Returning custom list of MD5 hashes from Get: %v", hashes)
	return &hashes, nil
}

func Update(service *services.Service, hashes BaAdvancedSettings) (*BaAdvancedSettings, error) {
	_, err := service.Client.UpdateWithPut(advancedSettingsEndpoint, hashes)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG] Returning updated custom list of MD5 hashes from Get: %v", hashes)
	return &hashes, nil
}

func GetFileHashCount(service *services.Service) (*FileHashCount, error) {
	var hashes FileHashCount
	err := service.Client.Read(advancedSettingsEndpoint+fileHashCountEndpoint, &hashes)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG] Returning used andd unused quota for blocking MD5 file hashes from Get: %v", hashes)
	return &hashes, nil
}
