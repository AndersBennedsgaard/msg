package cmd

import (
	"os"
	"path/filepath"

	"github.com/AndersBennedsgaard/msg/internal/logging"
	"github.com/spf13/cobra"
)

const appGroup = "app"

var (
	quiet   bool
	verbose bool
	path    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "msg",
	Short: "A messaging application CLI",
	// PersistentPreRunE is called after flags are parsed but before the
	// command's RunE function is called.
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logging.InitLogger(verbose, quiet)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Get XDG_DATA_HOME, default to ~/.local/share if not set
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome == "" {
		home, _ := os.UserHomeDir()
		xdgDataHome = filepath.Join(home, ".local", "share")
	}

	// Construct the default database path
	defaultPath := filepath.Join(xdgDataHome, "msg", "db.sql")

	// Flag configuration
	rootCmd.AddGroup(&cobra.Group{ID: appGroup, Title: "Available Commands"})

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Add additional debug logs")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Reduce the amount of logging")

	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", defaultPath, "Path to database")
}
