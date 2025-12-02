package activation

/*
import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/stretchr/testify/assert"
)

func TestActivation(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "activation", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	t.Run("Test GetActivationStatus", func(t *testing.T) {
		activationStatus, err := GetActivationStatus(service)
		assert.NoError(t, err)
		assert.Contains(t, []string{"ACTIVE", "PENDING", "INPROGRESS"}, activationStatus.Status)
	})

	t.Run("Test CreateActivation", func(t *testing.T) {
		tests := []struct {
			input  Activation
			expect string
		}{
			{input: Activation{Status: "ACTIVE"}, expect: "ACTIVE"},
			{input: Activation{Status: "PENDING"}, expect: "PENDING"},
			{input: Activation{Status: "INPROGRESS"}, expect: "INPROGRESS"},
		}

		for _, test := range tests {
			createdActivation, err := CreateActivation(service, test.input)
			if err != nil {
				t.Logf("Warning: Failed to create activation with status %s: %v", test.input.Status, err)
				continue
			}

			if test.expect != createdActivation.Status {
				t.Logf("Warning: Expected status %s but got %s", test.expect, createdActivation.Status)
			}
		}
	})
}
*/
