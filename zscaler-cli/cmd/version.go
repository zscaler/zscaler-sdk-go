package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Injected at build time
	Version   = "unknown"
	GitSHA    = "unknown"
	BuildTime = "unknown"

	// versionCmd represents the version command
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the Zscaler CLI version",
		Long: `Prints out the installed version of the Zscaler CLI and checks for newer
versions available for update.`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("zscaler-cli v%s (sha:%s) (built:%s)\n", Version, GitSHA, BuildTime)
			checkForUpdates()
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

func checkForUpdates() {
	// Placeholder function to implement version checking logic
	fmt.Println("Checking for updates...")
	// Insert logic to check for the latest version and notify the user
}
