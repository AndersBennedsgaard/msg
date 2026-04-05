package cmd

import (
	"database/sql"
	"fmt"

	"github.com/AndersBennedsgaard/msg/internal/logging"
	"github.com/AndersBennedsgaard/msg/internal/store"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:     "read",
	GroupID: appGroup,
	Short:   "Print the next unread message and mark it as read",
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

		message, err := dbstore.ReadNextMessage()
		if err != nil {
			return fmt.Errorf("error occured when reading unread message: %w", err)
		}

		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Time: %s\nSeverity: %s\nType: %s\n\n%s", message.Timestamp(), message.Severity(), message.Type(), message.Message())
		return err
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}
