package cmd

import (
	"database/sql"
	"fmt"

	"github.com/AndersBennedsgaard/msg/internal/logging"
	"github.com/AndersBennedsgaard/msg/internal/store"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// countCmd represents the count command
var countCmd = &cobra.Command{
	Use:   "count",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := logging.GetLogger()

		db, err := sql.Open("sqlite", path)
		if err != nil {
			return fmt.Errorf("error occurred when opening datbase: %w", err)
		}
		defer func() {
			err := db.Close()
			if err != nil {
				logger.Fatal("error occurred when closing the database: %s", zap.Error(err))
			}
		}()

		dbstore, err := store.NewSqliteStore(db)
		if err != nil {
			return fmt.Errorf("error occured when initializing the database connection: %w", err)
		}

		count, err := dbstore.CountUnreadMessages()
		if err != nil {
			return fmt.Errorf("error occured when counting unread messages: %w", err)
		}

		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Number of unread messages: %d\n", count)
		return err
	},
}

func init() {
	rootCmd.AddCommand(countCmd)
}
