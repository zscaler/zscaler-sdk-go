package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/devices"
)

type Device struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Platform string `json:"platform"`
}

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

	var fromInt, toInt int

	if fromInput != "" {
		parsedFrom, err := strconv.ParseInt(fromInput, 10, 64)
		if err != nil {
			log.Fatalf("[ERROR] Invalid start time: %v\n", err)
		}
		fromInt, err = safeIntConversion(parsedFrom)
		if err != nil {
			log.Fatalf("[ERROR] %v\n", err)
		}
	} else {
		fromInt, err = safeIntConversion(fromTime)
		if err != nil {
			log.Fatalf("[ERROR] %v\n", err)
		}
	}

	if toInput != "" {
		parsedTo, err := strconv.ParseInt(toInput, 10, 64)
		if err != nil {
			log.Fatalf("[ERROR] Invalid end time: %v\n", err)
		}
		toInt, err = safeIntConversion(parsedTo)
		if err != nil {
			log.Fatalf("[ERROR] %v\n", err)
		}
	} else {
		toInt, err = safeIntConversion(toTime)
		if err != nil {
			log.Fatalf("[ERROR] %v\n", err)
		}
	}

	// Define filters
	filters := devices.GetDevicesFilters{
		GetFromToFilters: common.GetFromToFilters{
			From: fromInt,
			To:   toInt,
		},
	}

	// Get all devices
	devicesList, _, err := devices.GetAllDevices(ctx, service, filters)
	if err != nil {
		log.Fatalf("[ERROR] getting all devices failed: %v\n", err)
	}

	// Extract device details and display in table format
	var deviceData []Device
	for _, device := range devicesList {
		// Extract platform information from device name
		parts := strings.Split(device.Name, "(")
		name := parts[0]
		platform := ""
		if len(parts) > 1 {
			platform = strings.TrimSuffix(parts[1], ")")
		}
		deviceData = append(deviceData, Device{
			ID:       device.ID,
			Name:     name,
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

func safeIntConversion(value int64) (int, error) {
	if value > int64(int(^uint(0)>>1)) || value < int64(-int(^uint(0)>>1)-1) {
		return 0, fmt.Errorf("value %d is out of range for int type", value)
	}
	return int(value), nil
}
