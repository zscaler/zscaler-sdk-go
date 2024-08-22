package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zscaler/zscaler-sdk-go/v2/zscaler-cli/cmd/zdx"
)

var (
	// Global CLI state
	cli = NewDefaultState()

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:               "zscaler",
		Short:             "A tool to manage the Zscaler platform.",
		DisableAutoGenTag: true,
		SilenceErrors:     true,
		Long: `The Zscaler Command Line Interface is a tool that helps you manage the
Zscaler platform. Use it to interact with different services like ZDX, ZPA, ZIA, and others.

Start by configuring the Zscaler CLI with the command:

    zscaler configure

This will prompt you for your Zscaler account and a set of API access keys.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return cliPersistentPreRun(cmd, args)
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)

func cliPersistentPreRun(cmd *cobra.Command, args []string) error {
	cli.Log.Debugw("starting command", "command", cmd.CommandPath(), "args", args)

	switch cmd.Use {
	case "help [command]", "configure", "version":
		return nil
	default:
		// Initialize the client for each command
		if err := cli.NewClient(); err != nil {
			return err
		}
	}

	cli.Log.Debugw("finished pre-run setup", "command", cmd.CommandPath())
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().Bool("debug", false, "turn on debug logging")
	rootCmd.PersistentFlags().Bool("nocolor", false, "turn off colors")
	rootCmd.PersistentFlags().Bool("json", false, "output in JSON format")

	// Bind global flags to Viper
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("nocolor", rootCmd.PersistentFlags().Lookup("nocolor"))
	viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))

	// Initialize ZDX command and add it to the root command
	rootCmd.AddCommand(zdx.NewZDXCommand(cli))
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigName(".zscaler")

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
