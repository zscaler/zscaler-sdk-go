package security_policy_settings

const (
	securityEndpoint         = "/security"
	securityAdvancedEndpoint = "/security/advanced"
)

// TODO: because there isn't an endpoint to get all Urls, we need to have all action types here.
var AddRemoveURLFromList []string = []string{
	"ADD_TO_LIST",
	"REMOVE_FROM_LIST",
}

type ListUrls struct {
	// Allowlist URLs whose contents will not be scanned. Allows up to 255 URLs. There may be trusted websites the content of which might be blocked due to anti-virus, anti-spyware, or anti-malware policies. Enter the URLs of sites you do not want scanned. The service allows users to download content from these URLs without inspecting the traffic. The allowlist applies to the Malware Protection, Advanced Threats Protection, and Sandbox policies.
	White []string `json:"whitelistUrls,omitempty"`

	// URLs on the denylist for your organization. Allow up to 25000 URLs.
	Black []string `json:"blacklistUrls,omitempty"`
}

func (service *Service) GetListUrls() (*ListUrls, error) {
	whitelist, err := service.GetWhiteListUrls()
	if err != nil {
		return nil, err
	}
	blacklist, err := service.GetBlackListUrls()
	if err != nil {
		return nil, err
	}
	return &ListUrls{
		White: whitelist.White,
		Black: blacklist.Black,
	}, nil
}

func (service *Service) UpdateListUrls(listUrls ListUrls) (*ListUrls, error) {
	whitelist, err := service.UpdateWhiteListUrls(ListUrls{White: listUrls.White})
	if err != nil {
		return nil, err
	}
	blacklist, err := service.UpdateBlackListUrls(ListUrls{Black: listUrls.Black})
	if err != nil {
		return nil, err
	}
	return &ListUrls{
		White: whitelist.White,
		Black: blacklist.Black,
	}, nil
}

func (service *Service) UpdateWhiteListUrls(list ListUrls) (*ListUrls, error) {
	_, err := service.Client.UpdateWithPut(securityEndpoint, list)
	if err != nil {
		return nil, err
	}

	return &list, nil
}

func (service *Service) UpdateBlackListUrls(list ListUrls) (*ListUrls, error) {
	_, err := service.Client.UpdateWithPut(securityAdvancedEndpoint, list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (service *Service) GetWhiteListUrls() (*ListUrls, error) {
	var whitelist ListUrls
	err := service.Client.Read(securityEndpoint, &whitelist)
	if err != nil {
		return nil, err
	}
	return &whitelist, nil
}

func (service *Service) GetBlackListUrls() (*ListUrls, error) {
	var blacklist ListUrls
	err := service.Client.Read(securityAdvancedEndpoint, &blacklist)
	if err != nil {
		return nil, err
	}
	return &blacklist, nil
}
