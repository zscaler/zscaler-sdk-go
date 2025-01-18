package main

import (
	"bufio"
	"context"
	"encoding/csv"
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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/troubleshooting/deeptrace"
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
	// Prompt for device ID
	fmt.Print("Enter device ID: ")
	deviceIDInput, _ := reader.ReadString('\n')
	deviceIDInput = strings.TrimSpace(deviceIDInput)
	deviceID, err := strconv.Atoi(deviceIDInput)
	if err != nil {
		log.Fatalf("[ERROR] Invalid device ID: %v\n", err)
	}

	// Prompt for trace ID (optional)
	fmt.Print("Enter trace ID (optional): ")
	traceID, _ := reader.ReadString('\n')
	traceID = strings.TrimSpace(traceID)

	if traceID == "" {
		// Get all deep trace sessions for the device
		deepTraces, _, err := deeptrace.GetDeepTraces(ctx, service, deviceID)
		if err != nil {
			log.Fatalf("[ERROR] getting deep traces failed: %v\n", err)
		}

		// Display the deep trace sessions in a table
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Trace ID", "Session Name", "User ID", "Username", "App ID", "App Name", "Device ID", "Device Name", "Web Probe ID", "Web Probe Name", "Session Length", "Probe Device", "Status", "Expected Time", "Created At", "Started At", "Ended At"})

		for _, trace := range deepTraces {
			table.Append([]string{
				trace.TraceID,
				trace.TraceDetails.SessionName,
				trace.TraceDetails.UserID,
				trace.TraceDetails.Username,
				trace.TraceDetails.AppID,
				trace.TraceDetails.AppName,
				trace.TraceDetails.DeviceID,
				trace.TraceDetails.DeviceName,
				trace.TraceDetails.WebProbeID,
				trace.TraceDetails.WebProbeName,
				strconv.Itoa(trace.TraceDetails.SessionLength),
				strconv.FormatBool(trace.TraceDetails.ProbeDevice),
				trace.Status,
				strconv.Itoa(trace.ExpectedTimeMinutes),
				time.Unix(int64(trace.CreatedAt), 0).Format("2006-01-02"),
				time.Unix(int64(trace.StartedAt), 0).Format("2006-01-02"),
				time.Unix(int64(trace.EndedAt), 0).Format("2006-01-02"),
			})
		}

		table.Render()

		// Ask the user if they want to save the table to a CSV file
		fmt.Print("Do you want to save the table to a CSV file? (yes/no): ")
		saveToCSV, _ := reader.ReadString('\n')
		saveToCSV = strings.TrimSpace(saveToCSV)

		if strings.ToLower(saveToCSV) == "yes" {
			// Prompt for the CSV file name
			fmt.Print("Enter the CSV file name: ")
			csvFileName, _ := reader.ReadString('\n')
			csvFileName = strings.TrimSpace(csvFileName)

			file, err := os.Create(csvFileName)
			if err != nil {
				log.Fatalf("[ERROR] creating CSV file failed: %v\n", err)
			}
			defer file.Close()

			writer := csv.NewWriter(file)
			defer writer.Flush()

			// Write the table header
			writer.Write([]string{"Trace ID", "Session Name", "User ID", "Username", "App ID", "App Name", "Device ID", "Device Name", "Web Probe ID", "Web Probe Name", "Session Length", "Probe Device", "Status", "Expected Time", "Created At", "Started At", "Ended At"})

			// Write the table rows
			for _, trace := range deepTraces {
				writer.Write([]string{
					trace.TraceID,
					trace.TraceDetails.SessionName,
					trace.TraceDetails.UserID,
					trace.TraceDetails.Username,
					trace.TraceDetails.AppID,
					trace.TraceDetails.AppName,
					trace.TraceDetails.DeviceID,
					trace.TraceDetails.DeviceName,
					trace.TraceDetails.WebProbeID,
					trace.TraceDetails.WebProbeName,
					strconv.Itoa(trace.TraceDetails.SessionLength),
					strconv.FormatBool(trace.TraceDetails.ProbeDevice),
					trace.Status,
					strconv.Itoa(trace.ExpectedTimeMinutes),
					time.Unix(int64(trace.CreatedAt), 0).Format("2006-01-02"),
					time.Unix(int64(trace.StartedAt), 0).Format("2006-01-02"),
					time.Unix(int64(trace.EndedAt), 0).Format("2006-01-02"),
				})
			}

			log.Printf("Table saved to %s\n", csvFileName)
		}
	} else {
		// Get specific deep trace session details
		resp, err := deeptrace.GetDeepTraceSession(ctx, service, deviceID, traceID)
		if err != nil {
			log.Fatalf("[ERROR] getting deep trace session failed: %v\n", err)
		}

		if resp.StatusCode == http.StatusOK {
			log.Printf("Retrieved details for trace ID: %s\n", traceID)
		} else {
			log.Printf("Failed to retrieve details for trace ID: %s, Status Code: %d\n", traceID, resp.StatusCode)
		}
	}
}
