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

	// Use SafeCastToInt for conversion with error handling
	fromInt, err := common.SafeCastToInt(fromTime)
	if err != nil {
		log.Fatalf("[ERROR] %v\n", err)
	}
	toInt, err := common.SafeCastToInt(toTime)
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
	geoLocations, resp, err := devices.GetGeoLocations(ctx, service, filters)
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
