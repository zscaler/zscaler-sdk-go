package manage_pass

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	managePassEndpoint = "/zcc/papi/public/v1/managePass"
)

type ManagePass struct {
	CompanyID      int    `json:"companyId"`
	DeviceType     int    `json:"deviceType"`
	ExitPass       string `json:"exitPass"`
	LogoutPass     string `json:"logoutPass"`
	PolicyName     string `json:"policyName"`
	UninstallPass  string `json:"uninstallPass"`
	ZadDisablePass string `json:"zadDisablePass"`
	ZdpDisablePass string `json:"zdpDisablePass"`
	ZdxDisablePass string `json:"zdxDisablePass"`
	ZiaDisablePass string `json:"ziaDisablePass"`
	ZpaDisablePass string `json:"zpaDisablePass"`
}

type ManagePassResponseContract struct {
	ErrorMessage string `json:"errorMessage"`
}

func UpdateManagePass(ctx context.Context, service *zscaler.Service, managePass *ManagePass) (*ManagePassResponseContract, error) {
	body, err := json.Marshal(managePass)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal manage pass request: %w", err)
	}

	resp, err := service.Client.NewZccRequestDo(ctx, "POST", managePassEndpoint, nil, bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update manage pass: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update manage pass: received status code %d", resp.StatusCode)
	}

	var response ManagePassResponseContract
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
