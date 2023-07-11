package dlp_engines

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/tests"
)

func TestDLPEngines_data(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	engines, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting dlp engines: %v", err)
		return
	}
	if len(engines) == 0 {
		t.Errorf("No dlp engines found")
		return
	}
	name := engines[0].Name
	t.Log("Getting dlp engines by name:" + name)
	engine, err := service.GetByName(name)
	if err != nil {
		t.Errorf("Error getting dlp engines by name: %v", err)
		return
	}
	if engine.Name != name {
		t.Errorf("dlp engine name does not match: expected %s, got %s", name, engine.Name)
		return
	}
}
