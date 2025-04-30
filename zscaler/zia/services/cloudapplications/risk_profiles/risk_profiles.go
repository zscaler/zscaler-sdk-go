package risk_profiles

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	riskProfilesEndpoint = "/zia/api/v1/riskProfiles"
)

type RiskProfiles struct {
	ID                        int                       `json:"id,omitempty"`
	ProfileName               string                    `json:"profileName,omitempty"`
	ProfileType               string                    `json:"profileType,omitempty"`
	Status                    string                    `json:"status,omitempty"`
	ExcludeCertificates       int                       `json:"excludeCertificates,omitempty"`
	PoorItemsOfService        string                    `json:"poorItemsOfService,omitempty"`
	AdminAuditLogs            string                    `json:"adminAuditLogs,omitempty"`
	DataBreach                string                    `json:"dataBreach,omitempty"`
	SourceIpRestrictions      string                    `json:"sourceIpRestrictions,omitempty"`
	MfaSupport                string                    `json:"mfaSupport,omitempty"`
	SslPinned                 string                    `json:"sslPinned,omitempty"`
	HttpSecurityHeaders       string                    `json:"httpSecurityHeaders,omitempty"`
	Evasive                   string                    `json:"evasive,omitempty"`
	DnsCaaPolicy              string                    `json:"dnsCaaPolicy,omitempty"`
	WeakCipherSupport         string                    `json:"weakCipherSupport,omitempty"`
	PasswordStrength          string                    `json:"passwordStrength,omitempty"`
	SslCertValidity           string                    `json:"sslCertValidity,omitempty"`
	Vulnerability             string                    `json:"vulnerability,omitempty"`
	MalwareScanningForContent string                    `json:"malwareScanningForContent,omitempty"`
	FileSharing               string                    `json:"fileSharing,omitempty"`
	SslCertKeySize            string                    `json:"sslCertKeySize,omitempty"`
	VulnerableToHeartBleed    string                    `json:"vulnerableToHeartBleed,omitempty"`
	VulnerableToLogJam        string                    `json:"vulnerableToLogJam,omitempty"`
	VulnerableToPoodle        string                    `json:"vulnerableToPoodle,omitempty"`
	VulnerabilityDisclosure   string                    `json:"vulnerabilityDisclosure,omitempty"`
	SupportForWaf             string                    `json:"supportForWaf,omitempty"`
	RemoteScreenSharing       string                    `json:"remoteScreenSharing,omitempty"`
	SenderPolicyFramework     string                    `json:"senderPolicyFramework,omitempty"`
	DomainKeysIdentifiedMail  string                    `json:"domainKeysIdentifiedMail,omitempty"`
	DomainBasedMessageAuth    string                    `json:"domainBasedMessageAuth,omitempty"`
	LastModTime               int                       `json:"lastModTime,omitempty"`
	CreateTime                int                       `json:"createTime,omitempty"`
	Certifications            []string                  `json:"certifications,omitempty"`
	DataEncryptionInTransit   []string                  `json:"dataEncryptionInTransit,omitempty"`
	RiskIndex                 []int                     `json:"riskIndex,omitempty"`
	ModifiedBy                *common.IDNameExtensions  `json:"modifiedBy,omitempty"`
	CustomTags                []common.IDNameExtensions `json:"customTags,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, profileID int) (*RiskProfiles, error) {
	var riskProfiles RiskProfiles
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", riskProfilesEndpoint, profileID), &riskProfiles)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning risk profile from Get: %d", riskProfiles.ID)
	return &riskProfiles, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, profileName string) (*RiskProfiles, error) {
	var riskProfiles []RiskProfiles
	err := common.ReadAllPages(ctx, service.Client, riskProfilesEndpoint, &riskProfiles)
	if err != nil {
		return nil, err
	}
	for _, riskProfile := range riskProfiles {
		if strings.EqualFold(riskProfile.ProfileName, profileName) {
			return &riskProfile, nil
		}
	}
	return nil, fmt.Errorf("no risk profiles found with name: %s", profileName)
}

func Create(ctx context.Context, service *zscaler.Service, profileID *RiskProfiles) (*RiskProfiles, *http.Response, error) {
	resp, err := service.Client.Create(ctx, riskProfilesEndpoint, *profileID)
	if err != nil {
		return nil, nil, err
	}

	createdRiskProfile, ok := resp.(*RiskProfiles)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a risk profile pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new risk profile from create: %d", createdRiskProfile.ID)
	return createdRiskProfile, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, profileID int, profiles *RiskProfiles) (*RiskProfiles, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", riskProfilesEndpoint, profileID), *profiles)
	if err != nil {
		return nil, nil, err
	}
	updatedRiskProfile, _ := resp.(*RiskProfiles)

	service.Client.GetLogger().Printf("[DEBUG]returning updates risk profile from update: %d", updatedRiskProfile.ID)
	return updatedRiskProfile, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, profileID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", riskProfilesEndpoint, profileID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAllLite(ctx context.Context, service *zscaler.Service) ([]RiskProfiles, error) {
	var profiles []RiskProfiles
	err := common.ReadAllPages(ctx, service.Client, riskProfilesEndpoint+"/lite", &profiles)
	return profiles, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]RiskProfiles, error) {
	var profiles []RiskProfiles
	err := common.ReadAllPages(ctx, service.Client, riskProfilesEndpoint, &profiles)
	return profiles, err
}
