package sandbox_settings

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

func (service *Service) Get() (*BaAdvancedSettings, error) {
	var hashes BaAdvancedSettings
	err := service.Client.Read(advancedSettingsEndpoint, &hashes)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG] Returning custom list of MD5 hashes from Get: %v", hashes)
	return &hashes, nil
}

func (service *Service) Update() (*BaAdvancedSettings, error) {
	var hashes BaAdvancedSettings
	err := service.Client.Read(advancedSettingsEndpoint, &hashes)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG] Returning updated custom list of MD5 hashes from Get: %v", hashes)
	return &hashes, nil
}

func (service *Service) GetFileHashCount() (*FileHashCount, error) {
	var hashes FileHashCount
	err := service.Client.Read(advancedSettingsEndpoint+fileHashCountEndpoint, &hashes)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG] Returning used andd unusedd quota for blocking MD5 file hashes from Get: %v", hashes)
	return &hashes, nil
}
