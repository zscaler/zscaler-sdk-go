package adminroles

/*
import (
	"context"
	"log"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestAdminRole(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "adminroles", "ztw")
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	roles, err := GetAllAdminRoles(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting admin roles: %v", err)
		return
	}
	defer client.Stop()
	service := client.Service
	if len(roles) == 0 {
		t.Errorf("No admin roles found")
		return
	}
	name := roles[0].Name
	t.Log("Getting admin roles by name:" + name)
	ecgroup, err := GetByName(context.Background(), service, name)
	if err != nil {
		t.Errorf("Error getting admin roles by name: %v", err)
		return
	}
	defer client.Stop()
	service := client.Service
	if ecgroup.Name != name {
		t.Errorf("admin roles name does not match: expected %s, got %s", name, ecgroup.Name)
		return
	}

	adminRole, err := GetByName(context.Background(), service, name)
	if err != nil {
		t.Errorf("Error getting admin roles by name: %v", err)
		return
	}
	defer client.Stop()
	service := client.Service
	if adminRole.Name != name {
		t.Errorf("admin roles name does not match: expected %s, got %s", name, adminRole.Name)
		return
	}
	// Negative Test: Try to retrieve a admin role with a non-existent name
	nonExistentName := "ThisAdminRoleNotExist"
	_, err = GetByName(context.Background(), service, nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}
}
*/
