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
	service := services.New(cli)

	// Define filters
	fromInt, err := safeIntConversion(fromTime)
	if err != nil {
		log.Fatalf("[ERROR] %v\n", err)
	}

	toInt, err := safeIntConversion(toTime)
	if err != nil {
		log.Fatalf("[ERROR] %v\n", err)
	}

	filters := devices.GeoLocationFilter{
		GetFromToFilters: common.GetFromToFilters{
			From: fromInt,
			To:   toInt,
		},
	}

	// Get geolocations
	geoLocations, resp, err := devices.GetGeoLocations(service, filters)
	if err != nil {
		log.Fatalf("Error getting geo locations: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(geoLocations) == 0 {
		log.Println("No geolocations found.")
	} else {
		displayGeoLocations(geoLocations)
	}
}

func safeIntConversion(value int64) (int, error) {
	if value > int64(int(^uint(0)>>1)) || value < int64(-int(^uint(0)>>1)-1) {
		return 0, fmt.Errorf("value %d is out of range for int type", value)
	}
	return int(value), nil
}

func displayGeoLocations(geoLocations []devices.GeoLocation) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "GeoType", "Description", "Child ID", "Child GeoType", "Child Description"})

	for _, geoLocation := range geoLocations {
		if len(geoLocation.Children) == 0 {
			table.Append([]string{
				geoLocation.ID,
				geoLocation.Name,
				geoLocation.GeoType,
				geoLocation.Description,
				"",
				"",
				"",
			})
		} else {
			for _, child := range geoLocation.Children {
				table.Append([]string{
					geoLocation.ID,
					geoLocation.Name,
					geoLocation.GeoType,
					geoLocation.Description,
					child.ID,
					child.GeoType,
					child.Description,
				})
			}
		}
	}
	table.Render()
}
