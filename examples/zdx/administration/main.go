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
	"github.com/zscaler/zscaler-sdk-go/v2/zdx"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services/administration"
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

	// Prompt the user to choose a resource type
	fmt.Println("Choose the Resource Type:")
	fmt.Println("a. Retrieve Departments with Optional Filters")
	fmt.Println("b. Retrieve Locations with Optional Filters")
	fmt.Print("Enter choice (a/b): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "a":
		// Prompt for optional filters
		filters := promptForFilters(reader)
		getDepartments(service, filters)
	case "b":
		// Prompt for optional filters
		filters := promptForFilters(reader)
		getLocations(service, filters)
	default:
		log.Fatalf("[ERROR] Invalid choice: %s\n", choice)
	}
}

func promptForFilters(reader *bufio.Reader) administration.GetDepartmentsFilters {
	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix() // Default to 2 hours ago
	to := now.Unix()

	fmt.Print("Enter start time in Unix Epoch (optional: Defaults to the previous 2 hours): ")
	fromInput, _ := reader.ReadString('\n')
	fromInput = strings.TrimSpace(fromInput)
	if fromInput != "" {
		parsedFrom, err := strconv.ParseInt(fromInput, 10, 64)
		if err == nil {
			if parsedFrom > int64(int(^uint(0)>>1)) || parsedFrom < int64(-int(^uint(0)>>1)-1) {
				log.Fatalf("[ERROR] Start time is out of range for int type\n")
			}
			from = parsedFrom
		} else {
			log.Fatalf("[ERROR] Invalid start time: %v\n", err)
		}
	}

	fmt.Print("Enter end time in Unix Epoch (optional: Defaults to the previous 2 hours): ")
	toInput, _ := reader.ReadString('\n')
	toInput = strings.TrimSpace(toInput)
	if toInput != "" {
		parsedTo, err := strconv.ParseInt(toInput, 10, 64)
		if err == nil {
			if parsedTo > int64(int(^uint(0)>>1)) || parsedTo < int64(-int(^uint(0)>>1)-1) {
				log.Fatalf("[ERROR] End time is out of range for int type\n")
			}
			to = parsedTo
		} else {
			log.Fatalf("[ERROR] Invalid end time: %v\n", err)
		}
	}

	fmt.Print("Enter search query (optional): ")
	search, _ := reader.ReadString('\n')
	search = strings.TrimSpace(search)

	return administration.GetDepartmentsFilters{
		From:   int(from),
		To:     int(to),
		Search: search,
	}
}

func getDepartments(service *services.Service, filters administration.GetDepartmentsFilters) {
	departments, _, err := administration.GetDepartments(service, filters)
	if err != nil {
		log.Fatalf("Error getting departments: %v", err)
	}
	displayDepartments(departments)
}

func getLocations(service *services.Service, filters administration.GetDepartmentsFilters) {
	locations, _, err := administration.GetLocations(service, administration.GetLocationsFilters(filters))
	if err != nil {
		log.Fatalf("Error getting locations: %v", err)
	}
	displayLocations(locations)
}

func displayDepartments(departments []administration.Department) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Department ID", "Name"})

	for _, department := range departments {
		table.Append([]string{
			strconv.Itoa(department.ID),
			department.Name,
		})
	}
	table.Render()
}

func displayLocations(locations []administration.Location) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Location ID", "Name"})

	for _, location := range locations {
		table.Append([]string{
			strconv.Itoa(location.ID),
			location.Name,
		})
	}
	table.Render()
}
