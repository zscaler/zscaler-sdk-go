package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/devices"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	clientID := os.Getenv("ZSCALER_CLIENT_ID")
	clientSecret := os.Getenv("ZSCALER_CLIENT_SECRET")
	vanityDomain := os.Getenv("ZSCALER_VANITY_DOMAIN")

	// Initialize ZDX configuration
	zdxCfg, err := zscaler.NewConfiguration(
		zscaler.WithClientID(clientID),
		zscaler.WithClientSecret(clientSecret),
		zscaler.WithVanityDomain(vanityDomain),
		zscaler.WithDebug(true),
	)
	if err != nil {
		log.Fatalf("Error creating ZDX configuration: %v", err)
	}

	service, err := zscaler.NewOneAPIClient(zdxCfg)
	if err != nil {
		log.Fatalf("Error creating OneAPI client: %v", err)
	}

	ctx := context.Background()

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

	// Call GetDeviceTopProcesses with the provided device ID and trace ID
	topProcesses, resp, err := devices.GetDeviceTopProcesses(ctx, service, deviceID, traceID, common.GetFromToFilters{})
	if err != nil {
		log.Fatalf("Error getting device top processes: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
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
