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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/applications"
)

type AppMetric struct {
	Metric    string
	Unit      string
	TimeStamp int64
	Value     float64
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Check for API key and secret in environment variables
	apiKey := os.Getenv("ZDX_API_KEY_ID")
	apiSecret := os.Getenv("ZDX_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		log.Fatalf("[ERROR] API key and secret must be set in environment variables (ZDX_API_KEY_ID, ZDX_API_SECRET)\n")
	}

	// Prompt for application ID
	fmt.Print("Enter application ID: ")
	appIDInput, _ := reader.ReadString('\n')
	appIDInput = strings.TrimSpace(appIDInput)
	appID, err := strconv.Atoi(appIDInput)
	if err != nil {
		log.Fatalf("[ERROR] Invalid application ID: %v\n", err)
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

	// Create configuration and client
	cfg, err := zdx.NewConfig(apiKey, apiSecret, "userAgent")
	if err != nil {
		log.Fatalf("[ERROR] creating client failed: %v\n", err)
	}
	cli := zdx.NewClient(cfg)
	appService := services.New(cli)

	// Define filters
	filters := common.GetFromToFilters{
		From: int(fromTime),
		To:   int(toTime),
	}

	// Get app metrics
	metricsList, _, err := applications.GetAppMetrics(appService, appID, filters)
	if err != nil {
		log.Fatalf("[ERROR] getting app metrics failed: %v\n", err)
	}

	// Extract app metric details and display in table format
	var metricData []AppMetric
	for _, metric := range metricsList {
		for _, dp := range metric.DataPoints {
			metricData = append(metricData, AppMetric{
				Metric:    metric.Metric,
				Unit:      metric.Unit,
				TimeStamp: int64(dp.TimeStamp),
				Value:     dp.Value,
			})
		}
	}

	// Display the data in a formatted table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Metric", "Unit", "Timestamp", "Value"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	for _, metric := range metricData {
		table.Append([]string{metric.Metric, metric.Unit, strconv.FormatInt(metric.TimeStamp, 10), fmt.Sprintf("%.2f", metric.Value)})
	}

	table.Render()
}
