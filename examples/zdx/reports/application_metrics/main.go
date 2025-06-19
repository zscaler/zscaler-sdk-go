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
		fromTime = parsedFrom
	}
	if toInput != "" {
		parsedTo, err := strconv.ParseInt(toInput, 10, 64)
		if err != nil {
			log.Fatalf("[ERROR] Invalid end time: %v\n", err)
		}
		toTime = parsedTo
	}

	// Convert int64 values safely to int
	fromInt, err := common.SafeCastToInt(fromTime)
	if err != nil {
		log.Fatalf("[ERROR] %v\n", err)
	}
	toInt, err := common.SafeCastToInt(toTime)
	if err != nil {
		log.Fatalf("[ERROR] %v\n", err)
	}

	filters := common.GetFromToFilters{
		From: fromInt,
		To:   toInt,
	}

	// Get app metrics
	metricsList, _, err := applications.GetAppMetrics(ctx, service, appID, filters)
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
