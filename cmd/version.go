package cmd

import (
	"fmt"

	"github.com/AndersBennedsgaard/msg/internal/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	GroupID: appGroup,
	Short:   "Version prints the build information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:    %s\n", version.Version)
		fmt.Printf("Commit:     %s\n", version.Commit)
		fmt.Printf("Build Date: %s\n", version.Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
