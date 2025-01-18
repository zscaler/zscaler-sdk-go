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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/applications"
)

type AppScore struct {
	Metric    string
	Unit      string
	TimeStamp int64
	Value     float64
}

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

	// Use safeCastToInt for conversion
	filters := common.GetFromToFilters{
		From: common.SafeCastToInt(fromTime),
		To:   common.SafeCastToInt(toTime),
	}

	// Get app scores
	scoresList, _, err := applications.GetAppScores(ctx, service, appID, filters)
	if err != nil {
		log.Fatalf("[ERROR] getting app scores failed: %v\n", err)
	}

	// Extract app score details and display in table format
	var scoreData []AppScore
	for _, score := range scoresList {
		for _, dp := range score.DataPoints {
			scoreData = append(scoreData, AppScore{
				Metric:    score.Metric,
				Unit:      score.Unit,
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

	for _, score := range scoreData {
		table.Append([]string{score.Metric, score.Unit, strconv.FormatInt(score.TimeStamp, 10), fmt.Sprintf("%.2f", score.Value)})
	}

	table.Render()
}
