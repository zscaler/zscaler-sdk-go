package lssconfigcontroller

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
)

func TestLSSConfigController(t *testing.T) {
	policyType := "SIEM_POLICY"
	ipAddress, _ := acctest.RandIpAddress("192.168.0.0/24")
	rPort := strconv.Itoa(acctest.RandIntRange(1000, 9999))
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	accessPolicySet, _, err := policysetcontroller.GetByPolicyType(context.Background(), service, policyType)
	if err != nil {
		t.Errorf("Error getting access inspection policy set: %v", err)
		return
	}
	// create app connector group for testing
	appConnGroup, _, err := appconnectorgroup.Create(context.Background(), service, appconnectorgroup.AppConnectorGroup{
		Name:                     name,
		Description:              name,
		Enabled:                  true,
		CityCountry:              "San Jose, US",
		Latitude:                 "37.33874",
		Longitude:                "-121.8852525",
		Location:                 "San Jose, CA, USA",
		UpgradeDay:               "SUNDAY",
		UpgradeTimeInSecs:        "66600",
		OverrideVersionProfile:   true,
		VersionProfileName:       "Default",
		VersionProfileID:         "0",
		DNSQueryType:             "IPV4_IPV6",
		PRAEnabled:               false,
		WAFDisabled:              true,
		TCPQuickAckApp:           true,
		TCPQuickAckAssistant:     true,
		TCPQuickAckReadAssistant: true,
	})
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error creating app connector group for testing server group: %v", err)
	}
	defer func() {
		time.Sleep(time.Second * 2) // Sleep for 2 seconds before deletion
		_, _, getErr := appconnectorgroup.Get(context.Background(), service, appConnGroup.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := appconnectorgroup.Delete(context.Background(), service, appConnGroup.ID)
			if err != nil {
				t.Errorf("Error deleting app connector group: %v", err)
			}
		}
	}()

	lssConfig := &LSSResource{
		LSSConfig: &LSSConfig{
			Name:          name,
			Description:   name,
			Enabled:       true,
			LSSHost:       ipAddress,
			LSSPort:       rPort,
			Format:        "json",
			SourceLogType: "zpn_trans_log",
			Filter: []string{
				"BRK_MT_SETUP_FAIL_BIND_TO_AST_LOCAL_OWNER",
				"CLT_INVALID_DOMAIN",
				"AST_MT_SETUP_ERR_HASH_TBL_FULL",
				"AST_MT_SETUP_ERR_CONN_PEER",
				"BRK_MT_SETUP_FAIL_REJECTED_BY_POLICY_APPROVAL",
				"BRK_MT_SETUP_FAIL_ICMP_RATE_LIMIT_NUM_APP_EXCEEDED",
				"EXPTR_MT_TLS_SETUP_FAIL_VERSION_MISMATCH",
				"BRK_MT_SETUP_FAIL_RATE_LIMIT_LOOP_DETECTED",
				"CLT_INVALID_TAG",
				"AST_MT_SETUP_ERR_NO_SYSTEM_FD",
				"AST_MT_SETUP_ERR_NO_PROCESS_FD",
				"BROKER_NOT_ENABLED",
				"AST_MT_SETUP_ERR_AST_CFG_DISABLED",
				"BRK_MT_SETUP_FAIL_TOO_MANY_FAILED_ATTEMPTS",
				"BRK_MT_AUTH_NO_SAML_ASSERTION_IN_MSG",
				"BRK_MT_SETUP_FAIL_CTRL_BRK_CANNOT_FIND_CONNECTOR",
				"INVALID_DOMAIN",
				"BRK_MT_TERMINATED_BRK_SWITCHED",
				"AST_MT_SETUP_ERR_OPEN_SERVER_CLOSE",
				"AST_MT_SETUP_ERR_BIND_TO_AST_LOCAL_OWNER",
				"NO_CONNECTOR_AVAILABLE",
				"BRK_MT_AUTH_SAML_CANNOT_ADD_ATTR_TO_HEAP",
				"EXPTR_MT_TLS_SETUP_FAIL_NOT_TRUSTED_CA",
				"AST_MT_SETUP_TIMEOUT_NO_ACK_TO_BIND",
				"CLT_PORT_UNREACHABLE",
				"C2C_CLIENT_CONN_EXPIRED",
				"BRK_MT_SETUP_FAIL_BIND_TO_CLIENT_LOCAL_OWNER",
				"BRK_MT_AUTH_SAML_CANNOT_ADD_ATTR_TO_HASH",
				"BRK_MT_SETUP_FAIL_REPEATED_DISPATCH",
				"AST_MT_SETUP_ERR_OPEN_SERVER_ERROR",
				"DSP_MT_SETUP_FAIL_DISCOVERY_TIMEOUT",
				"CUSTOMER_NOT_ENABLED",
				"BRK_CONN_UPGRADE_REQUEST_FAILED",
				"C2C_MTUNNEL_FAILED_FORWARD",
				"EXPTR_MT_TLS_SETUP_FAIL_CERT_CHAIN_ISSUE",
				"AST_MT_SETUP_ERR_RATE_LIMIT_REACHED",
				"BRK_MT_SETUP_FAIL_RATE_LIMIT_NUM_APP_EXCEEDED",
				"CLT_WRONG_PORT",
				"AST_MT_SETUP_TIMEOUT_CANNOT_CONN_TO_SERVER",
				"BRK_MT_AUTH_SAML_FINGER_PRINT_FAIL",
				"AST_MT_SETUP_ERR_NO_EPHEMERAL_PORT",
				"BRK_CONN_UPGRADE_REQUEST_FORBIDDEN",
				"AST_MT_SETUP_ERR_OPEN_SERVER_CONN",
				"CLT_PROBE_FAILED",
				"AST_MT_SETUP_ERR_APP_NOT_FOUND",
				"AST_MT_SETUP_ERR_OPEN_BROKER_CONN",
				"BRK_MT_SETUP_FAIL_ICMP_RATE_LIMIT_EXCEEDED",
				"AST_MT_SETUP_ERR_OPEN_SERVER_TIMEOUT",
				"C2C_MTUNNEL_BAD_STATE",
				"CLT_DUPLICATE_TAG",
				"AST_MT_SETUP_TIMEOUT",
				"CLT_DOUBLEENCRYPT_NOT_SUPPORTED",
				"BRK_MT_SETUP_FAIL_CANNOT_SEND_MT_COMPLETE",
				"BRK_MT_SETUP_FAIL_BIND_RECV_IN_BAD_STATE",
				"APP_NOT_AVAILABLE",
				"BRK_MT_AUTH_SAML_NO_USER_ID",
				"AST_MT_SETUP_TIMEOUT_CANNOT_CONN_TO_BROKER",
				"DSP_MT_SETUP_FAIL_MISSING_HEALTH",
				"AST_MT_SETUP_ERR_DUP_MT_ID",
				"AST_MT_SETUP_ERR_BIND_GLOBAL_OWNER",
				"BRK_MT_TERMINATED_APPROVAL_TIMEOUT",
				"AST_MT_SETUP_ERR_BIND_ACK",
				"CLT_CONN_FAILED",
				"BRK_MT_SETUP_FAIL_ACCESS_DENIED",
				"AST_MT_SETUP_ERR_INIT_FOHH_MCONN",
				"AST_MT_SETUP_ERR_MEM_LIMIT_REACHED",
				"BRK_MT_SETUP_FAIL_DUPLICATE_TAG_ID",
				"BRK_MT_AUTH_SAML_FAILURE",
				"AST_MT_SETUP_ERR_PRA_UNAVAILABLE",
				"C2C_MTUNNEL_NOT_FOUND",
				"MT_CLOSED_INTERNAL_ERROR",
				"DSP_MT_SETUP_FAIL_CANNOT_SEND_TO_BROKER",
				"CLT_READ_FAILED",
				"BRK_MT_SETUP_FAIL_CANNOT_SEND_TO_DISPATCHER",
				"AST_MT_SETUP_ERR_BROKER_BIND_FAIL",
				"BRK_MT_SETUP_FAIL_RATE_LIMIT_EXCEEDED",
				"CLT_INVALID_CLIENT",
				"BRK_MT_SETUP_FAIL_APP_NOT_FOUND",
				"C2C_NOT_AVAILABLE",
				"AST_MT_SETUP_ERR_MAX_SESSIONS_REACHED",
				"BRK_MT_AUTH_TWO_SAML_ASSERTION_IN_MSG",
				"AST_MT_SETUP_ERR_CPU_LIMIT_REACHED",
				"AST_MT_SETUP_ERR_NO_DNS_TO_SERVER",
				"CLT_PROTOCOL_NOT_SUPPORTED",
				"BRK_MT_AUTH_ALREADY_FAILED",
				"BRK_MT_SETUP_FAIL_CONNECTOR_GROUPS_MISSING",
				"BRK_MT_SETUP_FAIL_SCIM_INACTIVE",
				"EXPTR_MT_TLS_SETUP_FAIL_PEER",
				"BRK_MT_AUTH_SAML_DECODE_FAIL",
				"AST_MT_SETUP_ERR_BRK_HASH_TBL_FULL",
				"APP_NOT_REACHABLE",
				"BRK_MT_SETUP_TIMEOUT",
				"BRK_MT_TERMINATED_IDLE_TIMEOUT",
				"MT_CLOSED_DTLS_CONN_GONE_CLIENT_CLOSED",
				"MT_CLOSED_DTLS_CONN_GONE",
				"MT_CLOSED_DTLS_CONN_GONE_AST_CLOSED",
				"MT_CLOSED_TLS_CONN_GONE_SCIM_USER_DISABLE",
				"MT_CLOSED_TLS_CONN_GONE_CLIENT_CLOSED",
				"MT_CLOSED_TLS_CONN_GONE",
				"OPEN_OR_ACTIVE_CONNECTION",
				"MT_CLOSED_TLS_CONN_GONE_AST_CLOSED",
				"ZPN_ERR_SCIM_INACTIVE",
				"BRK_MT_CLOSED_FROM_ASSISTANT",
				"MT_CLOSED_TERMINATED",
				"AST_MT_TERMINATED",
				"BRK_MT_CLOSED_FROM_CLIENT",
				"BRK_MT_TERMINATED",
				"BRK_MT_SETUP_FAIL_NO_POLICY_FOUND",
				"BRK_MT_SETUP_FAIL_REJECTED_BY_POLICY",
				"BRK_MT_SETUP_FAIL_SAML_EXPIRED",
			},
			UseTLS: true,
		},
		ConnectorGroups: []ConnectorGroups{
			{
				ID: appConnGroup.ID,
			},
		},
		PolicyRuleResource: &PolicyRuleResource{
			Name:        name,
			Description: name,
			Action:      "LOG",
			PolicySetID: accessPolicySet.ID,
			Conditions: []PolicyRuleResourceConditions{
				{
					Negated:  false,
					Operator: "OR",
					Operands: &[]PolicyRuleResourceOperands{
						{
							ObjectType: "CLIENT_TYPE",
							Values:     []string{"zpn_client_type_exporter"},
						},
						{
							ObjectType: "CLIENT_TYPE",
							Values:     []string{"zpn_client_type_machine_tunnel"},
						},
						{
							ObjectType: "CLIENT_TYPE",
							Values:     []string{"zpn_client_type_ip_anchoring"},
						},
						{
							ObjectType: "CLIENT_TYPE",
							Values:     []string{"zpn_client_type_edge_connector"},
						},
						{
							ObjectType: "CLIENT_TYPE",
							Values:     []string{"zpn_client_type_zapp"},
						},
						{
							ObjectType: "CLIENT_TYPE",
							Values:     []string{"zpn_client_type_slogger"},
						},
						{
							ObjectType: "CLIENT_TYPE",
							Values:     []string{"zpn_client_type_branch_connector"},
						},
						// {
						// 	ObjectType: "CLIENT_TYPE",
						// 	Values:     []string{"zpn_client_type_browser_isolation"},
						// },
					},
				},
			},
		},
	}

	// Test resource creation
	createdResource, _, err := Create(context.Background(), service, lssConfig)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making POST request: %v", err)
	}

	if createdResource.ID == "" {
		t.Error("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.LSSConfig.Name != name {
		t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.LSSConfig.Name)
	}
	// Test resource retrieval
	retrievedResource, _, err := Get(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.LSSConfig.Name != name {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", name, createdResource.LSSConfig.Name)
	}
	// Test resource update
	retrievedResource.LSSConfig.Name = updateName
	_, err = Update(context.Background(), service, createdResource.ID, retrievedResource)
	if err != nil {
		t.Errorf("Error updating resource: %v", err)
	}
	updatedResource, _, err := Get(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%s', but got '%s'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.LSSConfig.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.LSSConfig.Name)
	}

	// Test resource retrieval by name
	retrievedResource, _, err = GetByName(context.Background(), service, updateName)
	if err != nil {
		t.Errorf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.LSSConfig.Name != updateName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, createdResource.LSSConfig.Name)
	}
	// Test resources retrieval
	// resources, _, err := GetAll(context.Background(), service)
	// if err != nil {
	// 	t.Errorf("Error retrieving resources: %v", err)
	// }
	// if len(resources) == 0 {
	// 	t.Error("Expected retrieved resources to be non-empty, but got empty slice")
	// }
	// // check if the created resource is in the list
	// found := false
	// for _, resource := range resources {
	// 	if resource.ID == createdResource.ID {
	// 		found = true
	// 		break
	// 	}
	// }
	// if !found {
	// 	t.Errorf("Expected retrieved resources to contain created resource '%s', but it didn't", createdResource.ID)
	// }
	// Test resource removal
	_, err = Delete(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Errorf("Error deleting resource: %v", err)
		return
	}

	// Test resource retrieval after deletion
	_, resp, err := Get(context.Background(), service, createdResource.ID)
	if err == nil || (resp != nil && resp.StatusCode == http.StatusOK) {
		t.Errorf("Expected deletion to remove resource ID '%s', but it still exists", createdResource.ID)
	} else {
		t.Logf("Confirmed deletion of resource ID '%s'", createdResource.ID)
	}
}

func TestRetrieveNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	_, _, err = Get(context.Background(), service, "non_existent_id")
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	_, err = Delete(context.Background(), service, "non_existent_id")
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }
	_, err = Update(context.Background(), service, "non_existent_id", &LSSResource{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	_, _, err = GetByName(context.Background(), service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
