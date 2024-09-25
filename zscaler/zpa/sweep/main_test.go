package sweep

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appservercontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/bacertificate"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbibannercontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbicertificatecontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbiprofilecontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_custom_controls"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_profile"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/lssconfigcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/microtenants"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praapproval"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praconsole"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/pracredential"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praportal"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/provisioningkey"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/servergroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgegroup"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

var sweepFlag = flag.Bool("sweep", false, "Perform resource sweep")

func TestMain(m *testing.M) {
	flag.Parse() // Parse any flags that are defined.

	// Check if the sweep flag is set and the environment variable is true.
	if *sweepFlag && os.Getenv("ZPA_SDK_TEST_SWEEP") == "true" {
		log.Println("Sweep flag is set and ZPA_SDK_TEST_SWEEP is true. Starting sweep.")
		err := sweep()
		if err != nil {
			log.Printf("Failed to clean up resources: %v", err)
			os.Exit(1)
		}
	} else if *sweepFlag {
		log.Println("Sweep flag is set but ZPA_SDK_TEST_SWEEP environment variable is not set to true. Skipping sweep.")
		// Optionally, exit if you require the sweep to run.
		// os.Exit(1)
	} else {
		log.Println("Sweep flag not set. Proceeding with tests.")
	}

	// Proceed with normal testing.
	exitVal := m.Run()
	os.Exit(exitVal)
}

// sweep the resources before running integration tests
func sweep() error {
	log.Println("[INFO] Sweeping ZPA test resources")
	service, err := tests.NewOneAPIClient() // This returns a *zscaler.Service
	if err != nil {
		log.Printf("[ERROR] Failed to instantiate OneAPI client: %v", err)
		return err
	}

	client := service.Client // Extract the *zscaler.Client from the Service

	// List of all sweep functions to execute
	sweepFunctions := []func(*zscaler.Client) error{
		sweepPrivilegedApproval,
		sweepApplicationSegment,
		sweepSegmentGroup,
		sweepServerGroup,
		sweepAppConnectorGroups,
		sweepApplicationServers,
		sweepBaCertificateController,
		sweepCBIBannerController,
		sweepCBICertificateController,
		sweepCBIProfileController,
		sweepInspectionCustomControl,
		sweepInspectionProfile,
		sweepLSSController,
		sweepMicrotenants,
		sweepServiceEdgeGroup,
		sweepProvisioningKey,
		sweepPolicySetController,
		sweeppracredential,
		sweepPRAConsole,
		sweepPRAPortal,
	}

	// Execute each sweep function
	for _, fn := range sweepFunctions {
		// Get the function name using reflection
		fnName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		// Extracting the short function name from the full package path
		shortFnName := fnName[strings.LastIndex(fnName, ".")+1:]
		log.Printf("[INFO] Starting sweep: %s", shortFnName)

		if err := fn(client); err != nil {
			log.Printf("[ERROR] %s function error: %v", shortFnName, err)
			return err
		}
	}

	log.Println("[INFO] Sweep concluded successfully")
	return nil
}

func sweepAppConnectorGroups(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, _, err := appconnectorgroup.GetAll(context.Background(), service)
	if err != nil {
		log.Printf("[ERROR] Failed to get app connector groups: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, err := appconnectorgroup.Delete(context.Background(), service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete app connector group with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func sweepApplicationServers(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, _, err := appservercontroller.GetAll(context.Background(), service)
	if err != nil {
		log.Printf("[ERROR] Failed to get application server: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, err := appservercontroller.Delete(context.Background(), service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete application server with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func sweepApplicationSegment(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, _, err := applicationsegment.GetAll(context.Background(), service)
	if err != nil {
		log.Printf("[ERROR] Failed to get application segment: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, err := applicationsegment.Delete(context.Background(), service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete application segment with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func sweepBaCertificateController(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, _, err := bacertificate.GetAll(context.Background(), service)
	if err != nil {
		log.Printf("[ERROR] Failed to get browser access certificate: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, err := bacertificate.Delete(context.Background(), service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete browser access certificate with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func sweepCBIBannerController(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, _, err := cbibannercontroller.GetAll(context.Background(), service)
	if err != nil {
		log.Printf("[ERROR] Failed to get cbi banner controller: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") || strings.HasPrefix(r.Name, "updated-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, err := cbibannercontroller.Delete(context.Background(), service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete cbi banner controller with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func sweepCBICertificateController(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, _, err := cbicertificatecontroller.GetAll(context.Background(), service)
	if err != nil {
		log.Printf("[ERROR] Failed to get cbi certificate controller: %v", err)
		return err
	}

	for _, r := range resources {
		// Check if the resource's name starts with "tests-" or "updated-"
		if strings.HasPrefix(r.Name, "tests-") || strings.HasPrefix(r.Name, "updated-") {
			log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
			_, err := cbicertificatecontroller.Delete(context.Background(), service, r.ID)
			if err != nil {
				log.Printf("[ERROR] Failed to delete cbi certificate controller with ID: %s, Name: %s: %v", r.ID, r.Name, err)
			}
		}
	}
	return nil
}

func sweepCBIProfileController(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, _, err := cbiprofilecontroller.GetAll(context.Background(), service)
	if err != nil {
		log.Printf("[ERROR] Failed to get cbi profile controller: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, err := cbiprofilecontroller.Delete(context.Background(), service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete cbi profile controller with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func sweepInspectionCustomControl(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, _, err := inspection_custom_controls.GetAll(context.Background(), service)
	if err != nil {
		log.Printf("[ERROR] Failed to get inspection custom control: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, err := inspection_custom_controls.Delete(context.Background(), service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete inspection custom control with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func sweepInspectionProfile(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, _, err := inspection_profile.GetAll(context.Background(), service)
	if err != nil {
		log.Printf("[ERROR] Failed to get inspection profile: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, err := inspection_profile.Delete(context.Background(), service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete inspection profile with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func sweepLSSController(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, _, err := lssconfigcontroller.GetAll(context.Background(), service)
	if err != nil {
		log.Printf("[ERROR] Failed to get lss config controller: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.LSSConfig.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.LSSConfig.Name)
		_, err := lssconfigcontroller.Delete(context.Background(), service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete lss config controller with ID: %s, Name: %s: %v", r.ID, r.LSSConfig.Name, err)
		}
	}
	return nil
}

func sweepMicrotenants(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, _, err := microtenants.GetAll(context.Background(), service)
	if err != nil {
		log.Printf("[ERROR] Failed to get microtenants: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, err := microtenants.Delete(context.Background(), service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete microtenants with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func sweepSegmentGroup(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, _, err := segmentgroup.GetAll(context.Background(), service)
	if err != nil {
		log.Printf("[ERROR] Failed to get segment group: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, err := segmentgroup.Delete(context.Background(), service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete segment group with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func sweepServerGroup(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, _, err := servergroup.GetAll(context.Background(), service)
	if err != nil {
		log.Printf("[ERROR] Failed to get server group: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, err := servergroup.Delete(context.Background(), service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete server group with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func sweepServiceEdgeGroup(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, _, err := serviceedgegroup.GetAll(context.Background(), service)
	if err != nil {
		log.Printf("[ERROR] Failed to get service edge group: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, err := serviceedgegroup.Delete(context.Background(), service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete service edge group with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func sweepProvisioningKey(client *zscaler.Client) error {
	service := zscaler.NewService(client)

	// Define the association types to iterate over
	associationTypes := []string{"CONNECTOR_GRP", "SERVICE_EDGE_GRP"}

	for _, associationType := range associationTypes {
		resources, err := provisioningkey.GetAllByAssociationType(context.Background(), service, associationType)
		if err != nil {
			log.Printf("[ERROR] Failed to get provisioning keys for association type %s: %v", associationType, err)
			return err
		}

		for _, r := range resources {
			if !strings.HasPrefix(r.Name, "tests-") {
				continue
			}
			log.Printf("Deleting provisioning key with ID: %s, Name: %s, AssociationType: %s", r.ID, r.Name, associationType)
			_, err := provisioningkey.Delete(context.Background(), service, associationType, r.ID) // Assuming Delete method requires ID and associationType
			if err != nil {
				log.Printf("[ERROR] Failed to delete provisioning key with ID: %s, Name: %s, AssociationType: %s: %v", r.ID, r.Name, associationType, err)
			}
		}
	}
	return nil
}

func sweepPolicySetController(client *zscaler.Client) error {
	service := zscaler.NewService(client)

	policyTypes := []string{"ACCESS_POLICY", "TIMEOUT_POLICY", "CLIENT_FORWARDING_POLICY", "ISOLATION_POLICY", "INSPECTION_POLICY", "CREDENTIAL_POLICY", "CAPABILITIES_POLICY", "CLIENTLESS_SESSION_PROTECTION_POLICY", "REDIRECTION_POLICY"}

	for _, policyType := range policyTypes {
		// Fetch the global policy set ID for the current policy type
		globalPolicySet, _, err := policysetcontroller.GetByPolicyType(context.Background(), service, policyType)
		if err != nil {
			log.Printf("[ERROR] Failed to get global policy set for policy type %s: %v", policyType, err)
			return err
		}

		// Fetch all rules for the current policy type
		resources, _, err := policysetcontroller.GetAllByType(context.Background(), service, policyType)
		if err != nil {
			log.Printf("[ERROR] Failed to get access rules for policy type %s: %v", policyType, err)
			return err
		}

		for _, r := range resources {
			if !strings.HasPrefix(r.Name, "tests-") {
				continue
			}
			log.Printf("Deleting access rule with ID: %s, Name: %s, PolicyType: %s, PolicySetID: %s", r.ID, r.Name, policyType, globalPolicySet.ID)
			_, err := policysetcontroller.Delete(context.Background(), service, globalPolicySet.ID, r.ID) // Use the fetched policySetID for deletion
			if err != nil {
				log.Printf("[ERROR] Failed to delete access rule with ID: %s, Name: %s, PolicyType: %s, PolicySetID: %s: %v", r.ID, r.Name, policyType, globalPolicySet.ID, err)
			}
		}
	}
	return nil
}

func sweeppracredential(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, _, err := pracredential.GetAll(context.Background(), service)
	if err != nil {
		log.Printf("[ERROR] Failed to get credential controller: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, err := pracredential.Delete(context.Background(), service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete credential controller with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func sweepPRAConsole(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, _, err := praconsole.GetAll(context.Background(), service)
	if err != nil {
		log.Printf("[ERROR] Failed to get pra console: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, err := praconsole.Delete(context.Background(), service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete pra console with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func sweepPRAPortal(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, _, err := praportal.GetAll(context.Background(), service)
	if err != nil {
		log.Printf("[ERROR] Failed to get pra portal: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, err := praportal.Delete(context.Background(), service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete pra portal with ID: %s, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func sweepPrivilegedApproval(client *zscaler.Client) error {
	service := zscaler.NewService(client)

	// Retrieve all privileged approvals
	approvals, _, err := praapproval.GetAll(context.Background(), service)
	if err != nil {
		log.Printf("[ERROR] Failed to get all privileged approvals: %v", err)
		return err
	}

	// Delete each privileged approval by ID
	for _, approval := range approvals {
		log.Printf("Deleting privileged approval with ID: %s", approval.ID)
		resp, err := praapproval.Delete(context.Background(), service, approval.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete privileged approval with ID: %s: %v", approval.ID, err)
			return err
		} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			log.Printf("[ERROR] Unexpected status code when deleting privileged approval with ID: %s: %d", approval.ID, resp.StatusCode)
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
	}

	log.Printf("[INFO] Successfully deleted all privileged approvals")
	return nil
}
