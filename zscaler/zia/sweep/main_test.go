package sweep

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/adminuserrolemgmt/admins"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/filteringrules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/ipdestinationgroups"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/ipsourcegroups"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/networkapplicationgroups"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/networkservicegroups"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/networkservices"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/location/locationmanagement"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/rule_labels"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sandbox/sandbox_settings"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/security_policy_settings"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/gretunnels"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/staticips"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/vpncredentials"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlcategories"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlfilteringpolicies"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/user_authentication_settings"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

var sweepFlag = flag.Bool("sweep", false, "Perform resource sweep")

func TestMain(m *testing.M) {
	flag.Parse() // Parse any flags that are defined.

	// Check if the sweep flag is set and the environment variable is true.
	if *sweepFlag && os.Getenv("ZIA_SDK_TEST_SWEEP") == "true" {
		log.Println("Sweep flag is set and ZIA_SDK_TEST_SWEEP is true. Starting sweep.")
		err := sweep()
		if err != nil {
			log.Printf("Failed to clean up resources: %v", err)
			os.Exit(1)
		}
	} else if *sweepFlag {
		log.Println("Sweep flag is set but ZIA_SDK_TEST_SWEEP environment variable is not set to true. Skipping sweep.")
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
	log.Println("[INFO] Sweeping ZIA test resources")
	service, err := tests.NewOneAPIClient() // This returns a *zscaler.Service
	if err != nil {
		log.Printf("[ERROR] Failed to instantiate OneAPI client: %v", err)
		return err
	}

	client := service.Client // Extract the *zscaler.Client from the Service

	// List of all sweep functions to execute
	sweepFunctions := []func(*zscaler.Client) error{
		sweepFirewallFilteringRules,
		sweepURLFilteringPolicies,
		sweepLocationManagement,
		sweepAdminUsers,
		// sweepDLPEngines,
		// sweepDLPNotificationTemplates,
		// sweepADLPWebRules,
		// sweepDLPDictionaries,
		sweepIPDestinationGroup,
		sweepIPSourceGroup,
		sweepNetworkAplicationGroups,
		sweepNetworkServiceGroups,
		sweepNetworkServices,
		// sweepForwardingControlRules,
		// sweepZPAGateways,
		sweepRuleLabels,
		sweepGRETunnels,
		sweepStaticIP,
		sweepVPNCredentials,
		sweepURLCategories,
		// sweepUserManagement,
		sweepSandboxSettings,
		sweepSecurityPolicySettings,
		sweepUserAuthenticationSettings,
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

func sweepAdminUsers(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, err := admins.GetAllAdminUsers(service)
	if err != nil {
		log.Printf("[ERROR] Failed to get admin users: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.UserName, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %d, Name: %s", r.ID, r.UserName)
		_, err := admins.DeleteAdminUser(service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete app connector group with ID: %d, Name: %s: %v", r.ID, r.UserName, err)
		}
	}
	return nil
}

/*
	func sweepDLPEngines(client *zscaler.Client) error {
		service := zscaler.NewService(client)
		resources, err := dlp_engines.GetAll(service)
		if err != nil {
			log.Printf("[ERROR] Failed to get dlp engines: %v", err)
			return err
		}

		for _, r := range resources {
			if !strings.HasPrefix(r.Name, "tests-") {
				continue
			}
			log.Printf("Deleting resource with ID: %d, Name: %s", r.ID, r.Name)
			_, err := dlp_engines.Delete(service, r.ID)
			if err != nil {
				log.Printf("[ERROR] Failed to delete dlp engines with ID: %d, Name: %s: %v", r.ID, r.Name, err)
			}
		}
		return nil
	}

	func sweepDLPNotificationTemplates(client *zscaler.Client) error {
		service := zscaler.NewService(client)
		resources, err := dlp_notification_templates.GetAll(service)
		if err != nil {
			log.Printf("[ERROR] Failed to get dlp notification templates: %v", err)
			return err
		}

		for _, r := range resources {
			if !strings.HasPrefix(r.Name, "tests-") {
				continue
			}
			log.Printf("Deleting resource with ID: %d, Name: %s", r.ID, r.Name)
			_, err := dlp_notification_templates.Delete(service, r.ID)
			if err != nil {
				log.Printf("[ERROR] Failed to delete application segment with ID: %d, Name: %s: %v", r.ID, r.Name, err)
			}
		}
		return nil
	}

	func sweepADLPWebRules(client *zscaler.Client) error {
		service := zscaler.NewService(client)
		resources, err := dlp_web_rules.GetAll(service)
		if err != nil {
			log.Printf("[ERROR] Failed to get dlp web rules: %v", err)
			return err
		}

		for _, r := range resources {
			if !strings.HasPrefix(r.Name, "tests-") {
				continue
			}
			log.Printf("Deleting resource with ID: %d, Name: %s", r.ID, r.Name)
			_, err := dlp_web_rules.Delete(service, r.ID)
			if err != nil {
				log.Printf("[ERROR] Failed to delete dlp web rules with ID: %d, Name: %s: %v", r.ID, r.Name, err)
			}
		}
		return nil
	}

	func sweepDLPDictionaries(client *zscaler.Client) error {
		service := zscaler.NewService(client)
		resources, err := dlpdictionaries.GetAll(service)
		if err != nil {
			log.Printf("[ERROR] Failed to get dlp dictionaries: %v", err)
			return err
		}

		for _, r := range resources {
			if !strings.HasPrefix(r.Name, "tests-") {
				continue
			}
			log.Printf("Deleting resource with ID: %d, Name: %s", r.ID, r.Name)
			_, err := dlpdictionaries.DeleteDlpDictionary(service, r.ID)
			if err != nil {
				log.Printf("[ERROR] Failed to delete dlp dictionaries with ID: %d, Name: %s: %v", r.ID, r.Name, err)
			}
		}
		return nil
	}
*/
func sweepFirewallFilteringRules(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, err := filteringrules.GetAll(service)
	if err != nil {
		log.Printf("[ERROR] Failed to get Firewall filtering rule: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %d, Name: %s", r.ID, r.Name)
		_, err := filteringrules.Delete(service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete Firewall filtering rule with ID: %d, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func sweepIPDestinationGroup(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, err := ipdestinationgroups.GetAll(service)
	if err != nil {
		log.Printf("[ERROR] Failed to get ip destination group: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %d, Name: %s", r.ID, r.Name)
		_, err := ipdestinationgroups.Delete(service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete ip destination group with ID: %d, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func sweepIPSourceGroup(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, err := ipsourcegroups.GetAll(service)
	if err != nil {
		log.Printf("[ERROR] Failed to get ip source group: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %d, Name: %s", r.ID, r.Name)
		_, err := ipsourcegroups.Delete(service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete ip source group with ID: %d, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func sweepNetworkAplicationGroups(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, err := networkapplicationgroups.GetAllNetworkApplicationGroups(service)
	if err != nil {
		log.Printf("[ERROR] Failed to get network application groups: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %d, Name: %s", r.ID, r.Name)
		_, err := networkapplicationgroups.Delete(service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete network application groups with ID: %d, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func sweepNetworkServiceGroups(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, err := networkservicegroups.GetAllNetworkServiceGroups(service)
	if err != nil {
		log.Printf("[ERROR] Failed to get network service groups: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %d, Name: %s", r.ID, r.Name)
		_, err := networkservicegroups.DeleteNetworkServiceGroups(service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete network service groupsnetwork service groups with ID: %d, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func sweepNetworkServices(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, err := networkservices.GetAllNetworkServices(service)
	if err != nil {
		log.Printf("[ERROR] Failed to get network services: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %d, Name: %s", r.ID, r.Name)
		_, err := networkservices.Delete(service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete network services with ID: %d, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

/*
func sweepForwardingControlRules(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, err := forwarding_rules.GetAll(service)
	if err != nil {
		log.Printf("[ERROR] Failed to get forwarding control rules: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %d, Name: %s", r.ID, r.Name)
		_, err := forwarding_rules.Delete(service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete forwarding control rules with ID: %d, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}
*/
/*
func sweepZPAGateways(client *zscaler.Client) error {
	service := zpa_gateways.New(client)
	resources, err := service.GetAll()
	if err != nil {
		log.Printf("[ERROR] Failed to get zpa gateways: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %d, Name: %s", r.ID, r.Name)
		_, err := service.Delete(r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete zpa gateways with ID: %d, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}
*/

func sweepLocationManagement(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, err := locationmanagement.GetAll(service)
	if err != nil {
		log.Printf("[ERROR] Failed to get location management: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %d, Name: %s", r.ID, r.Name)
		_, err := locationmanagement.Delete(service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete location management with ID: %d, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func sweepRuleLabels(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, err := rule_labels.GetAll(service)
	if err != nil {
		log.Printf("[ERROR] Failed to rule labels: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %d, Name: %s", r.ID, r.Name)
		_, err := rule_labels.Delete(service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete rule labels with ID: %d, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

// TODO: Need to review method calls.
func sweepSandboxSettings(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	// First, fetch the current list of MD5 hashes
	currentSettings, err := sandbox_settings.Get(service)
	if err != nil {
		log.Printf("[ERROR] Failed to get current sandbox settings: %v", err)
		return err
	}

	// Check if the list of FileHashesToBeBlocked contains any items
	if len(currentSettings.FileHashesToBeBlocked) > 0 {
		// Since the list is not empty, proceed with clearing the MD5 hashes
		emptyHashes := sandbox_settings.BaAdvancedSettings{
			FileHashesToBeBlocked: []string{}, // Explicitly setting an empty slice
		}

		// Use the Update function with the emptyHashes object to clear the MD5 hashes
		_, err := sandbox_settings.Update(service, emptyHashes)
		if err != nil {
			log.Printf("[ERROR] Failed to clear MD5 hashes in sandbox settings: %v", err)
			return err
		}

		log.Println("[INFO] Successfully cleared MD5 hashes in sandbox settings")
	} else {
		// The list is already empty, so no need to send an update request
		log.Println("[INFO] No MD5 hashes to clear in sandbox settings")
	}

	return nil
}

// TODO: Need to review method calls.
func sweepSecurityPolicySettings(client *zscaler.Client) error {
	service := zscaler.NewService(client)

	// First, fetch the current lists of whitelist and blacklist URLs
	currentSettings, err := security_policy_settings.GetListUrls(service)
	if err != nil {
		log.Printf("[ERROR] Failed to get current security policy settings: %v", err)
		return err
	}

	// Check if either the whitelist or blacklist URLs contain any items
	if len(currentSettings.White) > 0 || len(currentSettings.Black) > 0 {
		// Since at least one list is not empty, proceed with clearing the URLs
		emptyUrls := security_policy_settings.ListUrls{
			White: []string{}, // Explicitly setting an empty slice for whitelist
			Black: []string{}, // Explicitly setting an empty slice for blacklist
		}

		// Use the UpdateListUrls function with the emptyUrls object to clear the URLs
		_, err := security_policy_settings.UpdateListUrls(service, emptyUrls)
		if err != nil {
			log.Printf("[ERROR] Failed to clear URLs in security policy settings: %v", err)
			return err
		}

		log.Println("[INFO] Successfully cleared URLs in security policy settings")
	} else {
		// Both lists are already empty, so no need to send an update request
		log.Println("[INFO] No URLs to clear in security policy settings")
	}

	return nil
}

func sweepUserAuthenticationSettings(client *zscaler.Client) error {
	service := zscaler.NewService(client)

	currentSettings, err := user_authentication_settings.Get(service)
	if err != nil {
		log.Printf("[ERROR] Failed to get current sandbox settings: %v", err)
		return err
	}

	if len(currentSettings.URLs) > 0 {
		emptyURLs := user_authentication_settings.ExemptedUrls{
			URLs: []string{},
		}
		_, err := user_authentication_settings.Update(service, emptyURLs)
		if err != nil {
			log.Printf("[ERROR] Failed to clear URLs from user authentication settings: %v", err)
			return err
		}

		log.Println("[INFO] Successfully cleared URLs from user authentication settings")
	} else {
		log.Println("[INFO] No URLs to clear in user authentication settings")
	}

	return nil
}

func sweepGRETunnels(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, err := gretunnels.GetAll(service)
	if err != nil {
		log.Printf("[ERROR] Failed to get gre tunnels: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.SourceIP, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %d, Name: %s", r.ID, r.SourceIP)
		_, err := gretunnels.DeleteGreTunnels(service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete gre tunnels with ID: %d, Name: %s: %v", r.ID, r.SourceIP, err)
		}
	}
	return nil
}

func sweepStaticIP(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, err := staticips.GetAll(service)
	if err != nil {
		log.Printf("[ERROR] Failed to get static ip: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Comment, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %d, Name: %s", r.ID, r.Comment)
		_, err := staticips.Delete(service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete static ip with ID: %d, Name: %s: %v", r.ID, r.Comment, err)
		}
	}
	return nil
}

func sweepVPNCredentials(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, err := vpncredentials.GetAll(service)
	if err != nil {
		log.Printf("[ERROR] Failed to get vpn credentials: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.FQDN, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %d, FQDN: %s", r.ID, r.FQDN)
		err := vpncredentials.Delete(service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete vpn credentials with ID: %d, Name: %s: %v", r.ID, r.Comments, err)
		}
	}
	return nil
}

func sweepURLCategories(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, err := urlcategories.GetAll(service)
	if err != nil {
		log.Printf("[ERROR] Failed to get url categories: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.ConfiguredName, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.ConfiguredName)
		err, _ := urlcategories.DeleteURLCategories(service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete url categories with ID: %s, Name: %s: %v", r.ID, r.ConfiguredName, err)
		}
	}
	return nil
}

func sweepURLFilteringPolicies(client *zscaler.Client) error {
	service := zscaler.NewService(client)
	resources, err := urlfilteringpolicies.GetAll(service)
	if err != nil {
		log.Printf("[ERROR] Failed to get url filtering policies: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %d, Name: %s", r.ID, r.Name)
		err, _ := urlfilteringpolicies.Delete(service, r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete url filtering policies with ID: %d, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

/*
func sweepUserManagement(client *zscaler.Client) error {
	service := users.New(client)
	resources, err := service.GetAllUsers()
	if err != nil {
		log.Printf("[ERROR] Failed to get users: %v", err)
		return err
	}

	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %d, Name: %s", r.ID, r.Name)
		err, _ := service.Delete(r.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete users with ID: %d, Name: %s: %v", r.ID, r.Name, err)
		}
	}
	return nil
}
*/
