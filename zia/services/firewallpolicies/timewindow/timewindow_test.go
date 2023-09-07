package timewindow

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestTimeWindow_data(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	tWindows, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting time windows: %v", err)
		return
	}
	if len(tWindows) == 0 {
		t.Errorf("No time windows found")
		return
	}
	tWindowName := tWindows[0].Name
	t.Log("Getting time window by name:" + tWindowName)
	tWindow, err := service.GetTimeWindowByName(tWindowName)
	if err != nil {
		t.Errorf("Error getting time windows by name: %v", err)
		return
	}
	if tWindow.Name != tWindowName {
		t.Errorf("time window name does not match: expected %s, got %s", tWindowName, tWindow.Name)
		return
	}
}
