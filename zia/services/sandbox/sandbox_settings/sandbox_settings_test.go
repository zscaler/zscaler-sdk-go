package sandbox_settings

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestUpdateBaAdvancedSettings(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := &Service{Client: client}

	// Define the desired settings for the update
	desiredSettings := BaAdvancedSettings{
		FileHashesToBeBlocked: []string{
			"42914d6d213a20a2684064be5c80ffa9",
			"c0202cf6aeab8437c638533d14563d35",
			"1ca31319721740ecb79f4b9ee74cd9b0",
			"2c373a7e86d0f3469849971e053bcc82",
			"40858748e03a544f6b562a687777397a",
			"465e89654a72256e7d1fb066388cc2a3",
			"47e7b297f020d53f7de7dc0f450e262d",
			"53d9af8829a9c7f6f177178885901c01",
			"9578c2be6437dcc8517e78a5de1fa975",
			"dfb689196faa945217a8929131f1d670",
			"8f9b7c1c2b84b8c71318b6776d31c9af",
			"a24bb61df75034769ffdda61c7a25926",
			"e5aea3b998644e394f506ac1f0f2f107",
			"1727de1b3d5636f1817d68ba0208fb50",
			"383498f810f0a992b964c19fc21ca398",
			"64990a45cf6b1b900c6b284bb54a1402",
			"97835760aa696d8ab7acbb5a78a5b013",
			"a8ab5aca96d260e649026e7fc05837bf",
			"c63a7c559870873133a84f0eb6ca54cd",
			"cc89100f20002801fa401b77dab0c512",
			"f8c110929606dca4c08ecaa9f9baf140",
			"f3dcf80b6251cfba1cd754006f693a73",
			"2c50efc0fef1601ce1b96b1b7cf991fb",
		},
	}

	updatedSettings, err := service.Update(desiredSettings)
	if err != nil {
		t.Errorf("Error updating BA Advanced Settings: %v", err)
	}
	if updatedSettings == nil {
		t.Error("Expected updated BA Advanced Settings, got nil")
	}
}

/*
	func TestUpdateBaAdvancedSettings(t *testing.T) {
		client, err := tests.NewZiaClient()
		if err != nil {
			t.Fatalf("Error creating client: %v", err)
		}

		service := &Service{Client: client}

		// Read hashes from a file
		file, err := os.Open("hashes.txt")
		if err != nil {
			t.Fatalf("Error opening hashes file: %v", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var hashes []string
		for scanner.Scan() {
			hashes = append(hashes, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			t.Fatalf("Error reading hashes from file: %v", err)
		}

		// Define the desired settings for the update
		desiredSettings := BaAdvancedSettings{
			FileHashesToBeBlocked: hashes,
		}

		updatedSettings, err := service.Update(desiredSettings)
		if err != nil {
			t.Errorf("Error updating BA Advanced Settings: %v", err)
		}
		if updatedSettings == nil {
			t.Error("Expected updated BA Advanced Settings, got nil")
		}
	}
*/
func TestGetBaAdvancedSettings(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := &Service{Client: client}

	settings, err := service.Get()
	if err != nil {
		t.Errorf("Error getting BA Advanced Settings: %v", err)
	}
	if settings == nil {
		t.Error("Expected BA Advanced Settings, got nil")
	}
}

func TestGetFileHashCount(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := &Service{Client: client}

	hashCount, err := service.GetFileHashCount()
	if err != nil {
		t.Errorf("Error getting file hash count: %v", err)
	}
	if hashCount == nil {
		t.Error("Expected file hash count, got nil")
	}
}

func TestEmptyHashList(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := &Service{Client: client}

	// Define the desired settings for the update
	desiredSettings := BaAdvancedSettings{
		FileHashesToBeBlocked: []string{},
	}

	updatedSettings, err := service.Update(desiredSettings)
	if err != nil {
		t.Errorf("Error updating BA Advanced Settings: %v", err)
	}
	if updatedSettings == nil {
		t.Error("Expected updated BA Advanced Settings, got nil")
	}
}
