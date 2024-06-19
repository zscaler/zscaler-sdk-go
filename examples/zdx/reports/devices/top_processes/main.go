package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services/reports/devices"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Check for API key and secret in environment variables
	apiKey := os.Getenv("ZDX_API_KEY_ID")
	apiSecret := os.Getenv("ZDX_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		log.Fatalf("[ERROR] API key and secret must be set in environment variables (ZDX_API_KEY_ID, ZDX_API_SECRET)\n")
	}

	// Prompt the user for device ID
	fmt.Print("Enter device ID: ")
	deviceIDInput, _ := reader.ReadString('\n')
	deviceIDInput = strings.TrimSpace(deviceIDInput)
	deviceID, err := strconv.Atoi(deviceIDInput)
	if err != nil {
		log.Fatalf("[ERROR] Invalid device ID: %v\n", err)
	}

	// Prompt the user for trace ID
	fmt.Print("Enter trace ID: ")
	traceIDInput, _ := reader.ReadString('\n')
	traceID := strings.TrimSpace(traceIDInput)

	// Create configuration and client
	cfg, err := zdx.NewConfig(apiKey, apiSecret, "userAgent")
	if err != nil {
		log.Fatalf("[ERROR] creating client failed: %v\n", err)
	}
	cli := zdx.NewClient(cfg)
	service := services.New(cli)

	// Call GetDeviceTopProcesses with the provided device ID and trace ID
	topProcesses, resp, err := devices.GetDeviceTopProcesses(service, deviceID, traceID, common.GetFromToFilters{})
	if err != nil {
		log.Fatalf("Error getting device top processes: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Expected status code %d, got %d", resp.StatusCode)
	}

	if len(topProcesses) == 0 {
		log.Println("No top processes found.")
	} else {
		displayTopProcesses(topProcesses)
	}
}

func displayTopProcesses(topProcesses []devices.DeviceTopProcesses) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Category", "Process ID", "Process Name"})

	for _, process := range topProcesses {
		for _, topProcess := range process.TopProcesses {
			for _, proc := range topProcess.Processes {
				table.Append([]string{
					topProcess.Category,
					strconv.Itoa(proc.ID),
					proc.Name,
				})
			}
		}
	}
	table.Render()
}
