package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	version = "v1.0.0"
	commit  = "unknown"
	date    = "unknown"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  `Print the version number, build information, and Go runtime version.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("CypherGoat CLI %s\n", version)
			fmt.Printf("Commit: %s\n", commit)
			fmt.Printf("Date: %s\n", date)
			fmt.Printf("Go: %s\n", runtime.Version())
		},
	}
}

func GetVersion() string {
	return version
}
