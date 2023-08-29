package policysetcontroller

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/tests"
)

var supportedPolicyTypes = []string{
	"ACCESS_POLICY",
	"TIMEOUT_POLICY",
	"CLIENT_FORWARDING_POLICY",
	"ISOLATION_POLICY",
}

// clean all resources
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	cleanResources() // clean up at the beginning
}

func teardown() {
	cleanResources() // clean up at the end
}

func shouldClean() bool {
	val, present := os.LookupEnv("ZSCALER_SDK_TEST_SWEEP")
	return !present || (present && (val == "" || val == "true")) // simplified for clarity
}

func cleanResources() {
	if !shouldClean() {
		return
	}

	client, err := tests.NewZpaClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	for _, policyType := range supportedPolicyTypes {
		resources, _, err := service.GetAllByType(policyType)
		if err != nil {
			log.Printf("Error fetching resources of type %s: %v", policyType, err)
			continue
		}

		for _, r := range resources {
			if !strings.HasPrefix(r.Name, "tests-") {
				continue
			}
			log.Printf("Deleting resource with ID: %s, Name: %s, Type: %s", r.ID, r.Name, policyType)

			// Assuming that the Delete function needs both policySetID and policyRuleID
			_, _ = service.Delete(r.PolicySetID, r.ID)
		}
	}
}
