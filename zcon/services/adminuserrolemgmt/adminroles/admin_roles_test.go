package adminroles

/*
import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

// const (
// 	maxRetries    = 3
// 	retryInterval = 2 * time.Second
// )

// // Constants for conflict retries
// const (
// 	maxConflictRetries    = 5
// 	conflictRetryInterval = 1 * time.Second
// )

// func retryOnConflict(operation func() error) error {
// 	var lastErr error
// 	for i := 0; i < maxConflictRetries; i++ {
// 		lastErr = operation()
// 		if lastErr == nil {
// 			return nil
// 		}

// 		if strings.Contains(lastErr.Error(), `"code":"EDIT_LOCK_NOT_AVAILABLE"`) {
// 			log.Printf("Conflict error detected, retrying in %v... (Attempt %d/%d)", conflictRetryInterval, i+1, maxConflictRetries)
// 			time.Sleep(conflictRetryInterval)
// 			continue
// 		}

// 		return lastErr
// 	}
// 	return lastErr
// }

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	cleanResources()
}

func teardown() {
	cleanResources()
}

func shouldClean() bool {
	val, present := os.LookupEnv("ZSCALER_SDK_TEST_SWEEP")
	return !present || (present && (val == "" || val == "true")) // simplified for clarity
}

func cleanResources() {
	if !shouldClean() {
		return
	}

	client, err := tests.NewZConClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	service := New(client)
	resources, err := service.GetAllAdminRoles()
	if err != nil {
		log.Printf("Error retrieving resources during cleanup: %v", err)
		return
	}

	for _, r := range resources {
		if strings.HasPrefix(r.Name, "tests-") {
			_, err := service.Delete(r.ID)
			if err != nil {
				log.Printf("Error deleting resource %d: %v", r.ID, err)
			}
		}
	}
}

func TestZCONAdminRoles(t *testing.T) {
	client, err := tests.NewZConClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	service := New(client)

	testRoles := []AdminRoles{
		// {
		// 	Name:             "test-role-any",
		// 	PolicyAccess:     "READ_WRITE",
		// 	AlertingAccess:   "READ_WRITE",
		// 	AnalysisAccess:   "READ_WRITE",
		// 	DashboardAccess:  "READ_WRITE",
		// 	ReportAccess:     "READ_WRITE",
		// 	UsernameAccess:   "READ_WRITE",
		// 	DeviceInfoAccess: "READ_WRITE",
		// 	AdminAcctAccess:  "READ_WRITE",
		// 	LogsLimit:        "UNRESTRICTED",
		// 	RoleType:         "ANY",
		// },
		// {
		// 	Name:             "test-role-public-api",
		// 	PolicyAccess:     "READ_ONLY",
		// 	Rank:             7,
		// 	AlertingAccess:   "NONE",
		// 	AnalysisAccess:   "NONE",
		// 	DashboardAccess:  "NONE",
		// 	ReportAccess:     "NONE",
		// 	UsernameAccess:   "NONE",
		// 	DeviceInfoAccess: "NONE",
		// 	AdminAcctAccess:  "READ_ONLY",
		// 	RoleType:         "PUBLIC_API",
		// 	LogsLimit:        "UNRESTRICTED",
		// 	FeaturePermissions: map[string]interface{}{
		// 		"EDGE_CONNECTOR_LOCATION_MANAGEMENT": "READ_WRITE",
		// 		"REMOTE_ASSISTANCE_MANAGEMENT":       "READ_WRITE",
		// 		"EDGE_CONNECTOR_CLOUD_PROVISIONING":  "READ_WRITE",
		// 		"APIKEY_MANAGEMENT":                  "READ_WRITE",
		// 		"EDGE_CONNECTOR_NSS_CONFIGURATION":   "READ_WRITE",
		// 		"EDGE_CONNECTOR_TEMPLATE":            "READ_WRITE",
		// 		"EDGE_CONNECTOR_ADMIN_MANAGEMENT":    "READ_WRITE",
		// 		"EDGE_CONNECTOR_DASHBOARD":           "READ_ONLY",
		// 		"EDGE_CONNECTOR_FORWARDING":          "READ_WRITE",
		// 	},
		// },
		// {
		// 	Name:             "test-role-sdwan",
		// 	PolicyAccess:     "READ_WRITE",
		// 	AlertingAccess:   "READ_WRITE",
		// 	AnalysisAccess:   "READ_WRITE",
		// 	DashboardAccess:  "READ_WRITE",
		// 	ReportAccess:     "READ_WRITE",
		// 	UsernameAccess:   "READ_WRITE",
		// 	DeviceInfoAccess: "READ_WRITE",
		// 	AdminAcctAccess:  "READ_WRITE",
		// 	RoleType:         "SDWAN",
		// },
		{
			Name:             "test-role-org-admin",
			PolicyAccess:     "READ_WRITE",
			AlertingAccess:   "READ_WRITE",
			AnalysisAccess:   "READ_WRITE",
			DashboardAccess:  "READ_WRITE",
			ReportAccess:     "READ_WRITE",
			UsernameAccess:   "READ_WRITE",
			DeviceInfoAccess: "READ_WRITE",
			AdminAcctAccess:  "READ_WRITE",
			RoleType:         "ORG_ADMIN",
		},
	}

	// Create test roles and verify creation
	for _, role := range testRoles {
		if err := createAndVerifyRole(service, &role, t); err != nil {
			t.Fatalf("Error creating test role '%s': %v", role.Name, err)
		}

		// Test search mechanisms
		if role.RoleType == "PUBLIC_API" {
			_, err := service.GetAPIRole(role.Name)
			if err != nil {
				t.Errorf("Error retrieving API role '%s': %v", role.Name, err)
			}
		} else if role.RoleType == "ORG_ADMIN" {
			_, err := service.GetPartnerRole(role.Name)
			if err != nil {
				t.Errorf("Error retrieving partner role '%s': %v", role.Name, err)
			}
		}
	}

	// Clean up created roles
	for _, role := range testRoles {
		if err := deleteRole(service, role.Name, t); err != nil {
			t.Errorf("Error cleaning up role '%s': %v", role.Name, err)
		}
	}
}

func createAndVerifyRole(service *Service, role *AdminRoles, t *testing.T) error {
	createdRole, err := service.Create(role)
	if err != nil {
		return err
	}
	if createdRole.ID == 0 {
		t.Errorf("Expected non-zero ID for created role '%s'", role.Name)
	}
	return nil
}

func deleteRole(service *Service, roleName string, t *testing.T) error {
	role, err := service.GetByName(roleName)
	if err != nil {
		return err
	}
	_, err = service.Delete(role.ID)
	return err
}
*/
