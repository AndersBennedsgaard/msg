package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/AndersBennedsgaard/msg/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfg     config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "msg",
	Short: "A messaging application CLI",
	// PersistentPreRunE is called after flags are parsed but before the
	// command's RunE function is called.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initializeConfig(cmd)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/msg/config.yaml)")

	rootCmd.PersistentFlags().StringP("path", "p", "./data", "Base path for storing messages")
	cobra.CheckErr(viper.BindPFlag("basePath", rootCmd.PersistentFlags().Lookup("path")))
}

func initializeConfig(cmd *cobra.Command) error {
	viper.SetEnvPrefix("MSG")
	// Allow for nested keys in environment variables (e.g. `MSG_DATA_PATH`).
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "*", "-", "*"))
	viper.AutomaticEnv()

	if cfgFile == "" {
		// Find the user's config directory.
		configdir, err := os.UserConfigDir()
		if err != nil {
			return err
		}

		// Search for a config file with the name "config"
		cfgFile = path.Join(configdir, "msg", "config.yaml")
	}

	viper.SetConfigFile(cfgFile)

	// Read the configuration file.
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if the config file doesn't exist.
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	// Bind Cobra flags to Viper.
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}

	// This is an optional but useful step to debug your config.
	fmt.Println("Configuration initialized. Using config file:", viper.ConfigFileUsed())
	return nil
}
