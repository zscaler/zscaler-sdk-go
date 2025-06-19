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
	"time"

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

	// Prompt the user for optional from and to timestamps
	fmt.Print("Enter start time in Unix Epoch (optional, defaults to 2 hours ago): ")
	fromInput, _ := reader.ReadString('\n')
	fromInput = strings.TrimSpace(fromInput)

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

	// Convert times to safe int values
	fromInt, err := safeIntConversion(fromTime)
	if err != nil {
		log.Fatalf("[ERROR] %v\n", err)
	}

	toInt, err := safeIntConversion(toTime)
	if err != nil {
		log.Fatalf("[ERROR] %v\n", err)
	}

	// Define filters
	filters := common.GetFromToFilters{
		From: fromInt,
		To:   toInt,
	}

	// Call GetEvents with the provided device ID and filters
	deviceEvents, resp, err := devices.GetEvents(ctx, service, deviceID, filters)
	if err != nil {
		log.Fatalf("Error getting events: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(deviceEvents) == 0 {
		log.Println("No events found.")
	} else {
		displayDeviceEvents(deviceEvents)
	}
}

func safeIntConversion(value int64) (int, error) {
	if value > int64(int(^uint(0)>>1)) || value < int64(-int(^uint(0)>>1)-1) {
		return 0, fmt.Errorf("value %d is out of range for int type", value)
	}
	return int(value), nil
}

func displayDeviceEvents(deviceEvents []devices.DeviceEvents) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Event Category", "Name", "DisplayName", "Prev", "Curr"})

	for _, deviceEvent := range deviceEvents {
		for _, event := range deviceEvent.Events {
			table.Append([]string{
				event.Category,
				event.Name,
				event.DisplayName,
				event.Prev,
				event.Curr,
			})
		}
	}
	table.Render()
}
