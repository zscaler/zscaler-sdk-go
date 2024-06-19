```go
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services/reports/devices"
)

type Device struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Platform string `json:"platform"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Check for API key and secret in environment variables
	apiKey := os.Getenv("ZDX_API_KEY_ID")
	apiSecret := os.Getenv("ZDX_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		log.Fatalf("[ERROR] API key and secret must be set in environment variables (ZDX_API_KEY_ID, ZDX_API_SECRET)\n")
	}

	// Prompt for from time
	fmt.Print("Enter start time in Unix Epoch (optional, defaults to 2 hours ago): ")
	fromInput, _ := reader.ReadString('\n')
	fromInput = strings.TrimSpace(fromInput)

	// Prompt for to time
	fmt.Print("Enter end time in Unix Epoch (optional, defaults to now): ")
	toInput, _ := reader.ReadString('\n')
	toInput = strings.TrimSpace(toInput)

	// Set default time range to last 2 hours if not provided
	now := time.Now()
	fromTime := now.Add(-2 * time.Hour).Unix()
	toTime := now.Unix()

	if fromInput != "" {
		parsedFrom, err := strconv.ParseInt(fromInput, 10, 64)
		if err != nil {
			log.Fatalf("[ERROR] Invalid start time: %v\n", err)
		}
		fromTime = parsedFrom
	}
	if toInput != "" {
		parsedTo, err := strconv.ParseInt(toInput, 10, 64)
		if err != nil {
			log.Fatalf("[ERROR] Invalid end time: %v\n", err)
		}
		toTime = parsedTo
	}

	// Create configuration and client
	cfg, err := zdx.NewConfig(apiKey, apiSecret, "userAgent")
	if err != nil {
		log.Fatalf("[ERROR] creating client failed: %v\n", err)
	}
	cli := zdx.NewClient(cfg)
	deviceService := services.New(cli)

	// Define filters
	filters := devices.GetDevicesFilters{
		GetFromToFilters: common.GetFromToFilters{
			From: int(fromTime),
			To:   int(toTime),
		},
	}

	// Get all devices
	devicesList, _, err := devices.GetAllDevices(deviceService, filters)
	if err != nil {
		log.Fatalf("[ERROR] getting all devices failed: %v\n", err)
	}

	// Extract device details and display in table format
	var deviceData []Device
	for _, device := range devicesList {
		// Extract platform information from device name
		parts := strings.Split(device.Name, "(")
		platform := ""
		if len(parts) > 1 {
			platform = strings.TrimSuffix(parts[1], ")")
		}
		deviceData = append(deviceData, Device{
			ID:       device.ID,
			Name:     parts[0],
			Platform: platform,
		})
	}

	// Display the data in a formatted table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"device_id", "device_name", "platform"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	for _, device := range deviceData {
		table.Append([]string{strconv.Itoa(device.ID), device.Name, device.Platform})
	}

	table.Render()
}
```