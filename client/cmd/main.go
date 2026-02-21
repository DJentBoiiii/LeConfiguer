package main

import (
	"fmt"
	"os"

	"client/internal/commands"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "le",
		Short: "LeConfiguer - Configuration management CLI",
		Long:  `A CLI tool for managing configurations through the API Gateway`,
	}

	// Get the API Gateway URL from environment or use default
	apiURL := os.Getenv("API_GATEWAY_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	// Add all commands
	rootCmd.AddCommand(commands.NewAddCommand(apiURL))
	rootCmd.AddCommand(commands.NewRemoveCommand(apiURL))
	rootCmd.AddCommand(commands.NewUpdateCommand(apiURL))
	rootCmd.AddCommand(commands.NewDiffCommand(apiURL))
	rootCmd.AddCommand(commands.NewUseCommand(apiURL))
	rootCmd.AddCommand(commands.NewRollbackCommand(apiURL))
	rootCmd.AddCommand(commands.NewVersionsCommand(apiURL))

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
