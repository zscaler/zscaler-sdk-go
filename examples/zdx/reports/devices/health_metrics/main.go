package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
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

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Check for API key and secret in environment variables
	apiKey := os.Getenv("ZDX_API_KEY_ID")
	apiSecret := os.Getenv("ZDX_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		log.Fatalf("[ERROR] API key and secret must be set in environment variables (ZDX_API_KEY_ID, ZDX_API_SECRET)\n")
	}

	// Create configuration and client
	cfg, err := zdx.NewConfig(apiKey, apiSecret, "userAgent")
	if err != nil {
		log.Fatalf("[ERROR] creating client failed: %v\n", err)
	}
	cli := zdx.NewClient(cfg)
	service := services.New(cli)

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
		if parsedFrom > int64(int(^uint(0)>>1)) || parsedFrom < int64(-int(^uint(0)>>1)-1) {
			log.Fatalf("[ERROR] Start time is out of range for int type\n")
		}
		fromTime = parsedFrom
	}
	if toInput != "" {
		parsedTo, err := strconv.ParseInt(toInput, 10, 64)
		if err != nil {
			log.Fatalf("[ERROR] Invalid end time: %v\n", err)
		}
		if parsedTo > int64(int(^uint(0)>>1)) || parsedTo < int64(-int(^uint(0)>>1)-1) {
			log.Fatalf("[ERROR] End time is out of range for int type\n")
		}
		toTime = parsedTo
	}

	// Define filters
	filters := common.GetFromToFilters{
		From: int(fromTime),
		To:   int(toTime),
	}

	// Call GetHealthMetrics with the provided device ID and filters
	healthMetrics, resp, err := devices.GetHealthMetrics(service, deviceID, filters)
	if err != nil {
		log.Fatalf("Error getting health metrics: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(healthMetrics) == 0 {
		log.Println("No health metrics found.")
	} else {
		displayHealthMetrics(healthMetrics)
	}
}

func displayHealthMetrics(healthMetrics []devices.HealthMetrics) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Category", "Instance", "Metric", "Unit", "Value", "Timestamp"})

	for _, metric := range healthMetrics {
		for _, instance := range metric.Instances {
			for _, met := range instance.Metrics {
				for _, datapoint := range met.DataPoints {
					table.Append([]string{
						metric.Category,
						instance.Name,
						met.Metric,
						met.Unit,
						fmt.Sprintf("%f", datapoint.Value),
						time.Unix(int64(datapoint.TimeStamp), 0).Format("2006-01-02 15:04:05"),
					})
				}
			}
		}
	}
	table.Render()
}
