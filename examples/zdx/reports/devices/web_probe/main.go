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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/devices"
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

	ctx := context.Background()

	// Prompt the user to choose a report type
	fmt.Println("Choose Device Web Probe Report:")
	fmt.Println("a. List all Web Probe metrics on a device for an application")
	fmt.Println("b. List all active web probes on a device")
	fmt.Print("Enter choice (a/b): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	// Define filters
	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix()
	to := now.Unix()
	filters := common.GetFromToFilters{
		From: int(from),
		To:   int(to),
	}

	switch choice {
	case "a":
		// Prompt for device ID
		fmt.Print("Enter device ID: ")
		deviceIDInput, _ := reader.ReadString('\n')
		deviceIDInput = strings.TrimSpace(deviceIDInput)
		deviceID, err := strconv.Atoi(deviceIDInput)
		if err != nil {
			log.Fatalf("[ERROR] Invalid device ID: %v\n", err)
		}

		// Prompt for app ID
		fmt.Print("Enter app ID: ")
		appIDInput, _ := reader.ReadString('\n')
		appIDInput = strings.TrimSpace(appIDInput)
		appID, err := strconv.Atoi(appIDInput)
		if err != nil {
			log.Fatalf("[ERROR] Invalid app ID: %v\n", err)
		}

		// Prompt for probe ID
		fmt.Print("Enter probe ID: ")
		probeIDInput, _ := reader.ReadString('\n')
		probeIDInput = strings.TrimSpace(probeIDInput)
		probeID, err := strconv.Atoi(probeIDInput)
		if err != nil {
			log.Fatalf("[ERROR] Invalid probe ID: %v\n", err)
		}

		// Get Web Probes metrics
		webProbeMetrics, resp, err := devices.GetWebProbes(ctx, service, deviceID, appID, probeID, filters)
		if err != nil {
			log.Fatalf("Error getting web probe metrics: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		if len(webProbeMetrics) == 0 {
			log.Println("No web probe metrics found.")
		} else {
			displayWebProbeMetrics(webProbeMetrics)
		}
	case "b":
		// Prompt for device ID
		fmt.Print("Enter device ID: ")
		deviceIDInput, _ := reader.ReadString('\n')
		deviceIDInput = strings.TrimSpace(deviceIDInput)
		deviceID, err := strconv.Atoi(deviceIDInput)
		if err != nil {
			log.Fatalf("[ERROR] Invalid device ID: %v\n", err)
		}

		// Prompt for app ID
		fmt.Print("Enter app ID: ")
		appIDInput, _ := reader.ReadString('\n')
		appIDInput = strings.TrimSpace(appIDInput)
		appID, err := strconv.Atoi(appIDInput)
		if err != nil {
			log.Fatalf("[ERROR] Invalid app ID: %v\n", err)
		}

		// Get all active Web Probes
		webProbes, resp, err := devices.GetAllWebProbes(ctx, service, deviceID, appID, filters)
		if err != nil {
			log.Fatalf("Error getting web probes: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		if len(webProbes) == 0 {
			log.Println("No web probes found.")
		} else {
			displayAllWebProbes(webProbes)
		}
	default:
		log.Fatalf("Invalid choice. Please enter 'a' or 'b'.")
	}
}

func displayWebProbeMetrics(webProbeMetrics []common.Metric) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Metric", "Unit", "Timestamp", "Value"})

	for _, metric := range webProbeMetrics {
		for _, dp := range metric.DataPoints {
			table.Append([]string{
				metric.Metric,
				metric.Unit,
				time.Unix(int64(dp.TimeStamp), 0).Format("2006-01-02 15:04:05"),
				fmt.Sprintf("%.2f", dp.Value),
			})
		}
	}
	table.Render()
}

func displayAllWebProbes(webProbes []devices.DeviceWebProbe) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "NumProbes", "AvgScore", "AvgPFT"})

	for _, probe := range webProbes {
		table.Append([]string{
			strconv.Itoa(probe.ID),
			probe.Name,
			strconv.Itoa(probe.NumProbes),
			fmt.Sprintf("%.2f", probe.AvgScore),
			fmt.Sprintf("%.2f", probe.AvgPFT),
		})
	}
	table.Render()
}
