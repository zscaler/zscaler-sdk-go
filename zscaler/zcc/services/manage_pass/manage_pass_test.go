package manage_pass

/*
import (
	"testing"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zcc/services"
)

func TestUpdateManagePass(t *testing.T) {
	client, err := tests.NewZccClient()
	if err != nil {
		t.Fatalf("Failed to create ZCC client: %v", err)
	}
	service := services.New(client)

	managePass := &ManagePass{
		CompanyID:      12345,
		DeviceType:     2,
		ExitPass:       "exitPassword",
		LogoutPass:     "logoutPassword",
		PolicyName:     "DefaultPolicy",
		UninstallPass:  "uninstallPassword",
		ZadDisablePass: "zadDisablePassword",
		ZdpDisablePass: "zdpDisablePassword",
		ZdxDisablePass: "zdxDisablePassword",
		ZiaDisablePass: "ziaDisablePassword",
		ZpaDisablePass: "zpaDisablePassword",
	}

	response, err := UpdateManagePass(service, managePass)
	if err != nil {
		t.Fatalf("Error updating manage pass: %v", err)
	}

	if response.ErrorMessage != "" {
		t.Errorf("Expected no error message, but got: %s", response.ErrorMessage)
	}
}
*/
