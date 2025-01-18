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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/alerts"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	apiKey := os.Getenv("ZDX_API_KEY_ID")
	apiSecret := os.Getenv("ZDX_API_SECRET")

	// Initialize ZDX configuration
	zdxCfg, err := zdx.NewConfiguration(
		zdx.WithZDXAPIKeyID(apiKey),
		zdx.WithZDXAPISecret(apiSecret),
		// Uncomment the line below if connecting to a custom ZDX cloud
		// zdx.WithZDXCloud("zdxbeta"),
		zdx.WithDebug(true),
	)
	if err != nil {
		log.Fatalf("Error creating ZDX configuration: %v", err)
	}

	// Initialize ZDX client
	zdxClient, err := zdx.NewClient(zdxCfg)
	if err != nil {
		log.Fatalf("Error creating ZDX client: %v", err)
	}

	// Wrap the ZDX client in a Service instance
	service := services.New(zdxClient)

	// Prompt the user to choose an alert type
	fmt.Println("Choose the Alert Type:")
	fmt.Println("a. Retrieve All Ongoing Alerts with Optional Filters (Defaults to the previous 2 hours)")
	fmt.Println("b. Retrieve Historical Alerts with Optional Filters (Defaults to the previous 2 hours)")
	fmt.Println("c. Retrieve Alert details including the impacted department, Zscaler locations, geolocation, and alert trigger")
	fmt.Println("d. Retrieve Alert details for affected Devices for specific AlertID")
	fmt.Print("Enter choice (a/b/c/d): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "a":
		// Prompt for optional filters
		filters := promptForFilters(reader, false)
		getOngoingAlerts(service, filters)
	case "b":
		// Prompt for optional filters
		filters := promptForFilters(reader, false)
		getHistoricalAlerts(service, filters)
	case "c":
		// Prompt for alert ID
		fmt.Print("Enter alert ID: ")
		alertID, _ := reader.ReadString('\n')
		alertID = strings.TrimSpace(alertID)
		getAlertDetails(service, alertID)
	case "d":
		// Prompt for alert ID and optional filters
		fmt.Print("Enter alert ID: ")
		alertID, _ := reader.ReadString('\n')
		alertID = strings.TrimSpace(alertID)
		filters := promptForFilters(reader, false)
		getAffectedDevices(service, alertID, filters)
	default:
		log.Fatalf("[ERROR] Invalid choice: %s\n", choice)
	}
}

func promptForFilters(reader *bufio.Reader, defaultTo14Days bool) common.GetFromToFilters {
	now := time.Now()
	var from, to int64

	if defaultTo14Days {
		from = now.Add(-14 * 24 * time.Hour).Unix() // Default to 14 days ago
	} else {
		from = now.Add(-2 * time.Hour).Unix() // Default to 2 hours ago
	}
	to = now.Unix()

	fmt.Print("Enter start time in Unix Epoch (optional: Defaults to the previous 2 hours): ")
	fromInput, _ := reader.ReadString('\n')
	fromInput = strings.TrimSpace(fromInput)
	if fromInput != "" {
		parsedFrom, err := strconv.ParseInt(fromInput, 10, 64)
		if err != nil {
			log.Fatalf("[ERROR] Invalid start time: %v\n", err)
		}
		from = parsedFrom
	}

	fmt.Print("Enter end time in Unix Epoch (optional: Defaults to the previous 2 hours): ")
	toInput, _ := reader.ReadString('\n')
	toInput = strings.TrimSpace(toInput)
	if toInput != "" {
		parsedTo, err := strconv.ParseInt(toInput, 10, 64)
		if err != nil {
			log.Fatalf("[ERROR] Invalid end time: %v\n", err)
		}
		to = parsedTo
	}

	if to-from > 14*24*60*60 {
		log.Fatalf("[ERROR] The time range cannot exceed 14 days.\n")
	}

	return common.GetFromToFilters{
		From: common.SafeCastToInt(from),
		To:   common.SafeCastToInt(to),
	}
}

func getOngoingAlerts(service *services.Service, filters common.GetFromToFilters) {
	ctx := context.Background()
	alertsResponse, _, err := alerts.GetOngoingAlerts(ctx, service, filters)
	if err != nil {
		log.Fatalf("Error getting ongoing alerts: %v", err)
	}
	displayAlerts(alertsResponse.Alerts)
}

func getHistoricalAlerts(service *services.Service, filters common.GetFromToFilters) {
	ctx := context.Background()
	alertsResponse, _, err := alerts.GetHistoricalAlerts(ctx, service, filters)
	if err != nil {
		log.Fatalf("Error getting historical alerts: %v", err)
	}
	displayAlerts(alertsResponse.Alerts)
}

func getAlertDetails(service *services.Service, alertID string) {
	ctx := context.Background()
	alertDetails, _, err := alerts.GetAlert(ctx, service, alertID)
	if err != nil {
		log.Fatalf("Error getting alert details: %v", err)
	}
	displayAlertDetails(*alertDetails) // Dereference the pointer
}

func getAffectedDevices(service *services.Service, alertID string, filters common.GetFromToFilters) {
	ctx := context.Background()
	affectedDevicesResponse, _, err := alerts.GetAffectedDevices(ctx, service, alertID, filters)
	if err != nil {
		log.Fatalf("Error getting affected devices: %v", err)
	}
	displayAffectedDevices(affectedDevicesResponse.Devices)
}

func displayAlerts(alerts []alerts.Alert) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Rule Name", "Severity", "Alert Type", "Status", "Num Geolocations", "Num Devices", "Started On", "Ended On"})

	for _, alert := range alerts {
		startedOn := time.Unix(int64(alert.StartedOn), 0).Format("2006-01-02 15:04:05")
		endedOn := ""
		if alert.EndedOn > 0 {
			endedOn = time.Unix(int64(alert.EndedOn), 0).Format("2006-01-02 15:04:05")
		} else {
			endedOn = "N/A"
		}
		table.Append([]string{
			strconv.Itoa(alert.ID),
			alert.RuleName,
			alert.Severity,
			alert.AlertType,
			alert.AlertStatus,
			strconv.Itoa(alert.NumGeolocations),
			strconv.Itoa(alert.NumDevices),
			startedOn,
			endedOn,
		})
	}
	table.Render()
}

func displayAlertDetails(alert alerts.Alert) {
	mainTable := tablewriter.NewWriter(os.Stdout)
	mainTable.SetHeader([]string{"ID", "Rule Name", "Severity", "Alert Type", "Status", "Started On", "Ended On"})

	startedOn := time.Unix(int64(alert.StartedOn), 0).Format("2006-01-02 15:04:05")
	endedOn := ""
	if alert.EndedOn > 0 {
		endedOn = time.Unix(int64(alert.EndedOn), 0).Format("2006-01-02 15:04:05")
	} else {
		endedOn = "N/A"
	}

	mainTable.Append([]string{
		strconv.Itoa(alert.ID),
		alert.RuleName,
		alert.Severity,
		alert.AlertType,
		alert.AlertStatus,
		startedOn,
		endedOn,
	})
	mainTable.Render()
}

func displayAffectedDevices(devices []alerts.Device) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Device ID", "Device Name", "User ID", "User Name", "User Email"})

	for _, device := range devices {
		table.Append([]string{
			strconv.Itoa(device.ID),
			device.Name,
			strconv.Itoa(device.UserID),
			device.UserName,
			device.UserEmail,
		})
	}
	table.Render()
}
