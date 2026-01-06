/*
Copyright © 2025 CypherGoat <contact@cyphergoat.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cyphergoat",
	Short: "CypherGoat CLI - Cryptocurrency swap tool",
	Long: `CypherGoat CLI is a tool that helps you perform cryptocurrency swaps from the command line. 

CypherGoat is an instant swap exchange aggregator.`,
}

var verbose bool

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose debug output")
	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
