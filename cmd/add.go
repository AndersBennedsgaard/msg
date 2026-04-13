package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/AndersBennedsgaard/msg/cmd/flags"
	"github.com/AndersBennedsgaard/msg/internal/logging"
	"github.com/AndersBennedsgaard/msg/internal/notification"
	"github.com/AndersBennedsgaard/msg/internal/store"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	msgType     string
	msgContent  string
	msgSeverity notification.NotificationSeverity
)

func hasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:                   "add [flags] [message]",
	GroupID:               appGroup,
	DisableFlagsInUseLine: true,
	Short:                 "Add a new notification message",
	Long:                  `Add a new notification message to the store with specified type, severity, and content.`,
	Example: ` # Add a new notification message:
 msg add -t alert -s high -m "Disk space low"

 # Add a new notification message by reading content from stdin:
 echo "CPU usage high" | msg add -t alert -s critical`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputReader := cmd.InOrStdin()
		logger := logging.GetLogger()

		if msgContent == "" {
			// read from inputReader
			if hasStdin() {
				data, err := io.ReadAll(inputReader)
				if err != nil {
					return err
				}
				msgContent = string(data)
			} else if len(args) == 1 {
				msgContent = args[0]
			} else {
				return errors.New("message content must be provided via --message flag, as an argument, or through stdin")
			}
		}

		db, err := sql.Open("sqlite", path)
		if err != nil {
			return fmt.Errorf("an error occurred when opening database: %w", err)
		}
		defer func() {
			err := db.Close()
			if err != nil {
				logger.Fatal("an error occurred when closing the database: %s", zap.Error(err))
			}
		}()

		dbstore, err := store.NewSqliteStore(db)
		if err != nil {
			return fmt.Errorf("an error occured when initializing the database connection: %w", err)
		}

		now := time.Now()

		noti, err := notification.NewNotification(
			msgType,
			now,
			msgSeverity,
			msgContent,
		)
		if err != nil {
			return err
		}

		msg := &notification.Message{
			Notification: noti,
			Status:       notification.StatusUnread,
		}
		id, err := dbstore.AddMessage(msg)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Added message with ID: %d\n", id)
		return err
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&msgType, "type", "t", "info", "Type of the notification")
	addCmd.Flags().StringVarP(&msgContent, "message", "m", "", "Content of the notification")
	addCmd.Flags().VarP(flags.NewSeverityValue(notification.NotificationSeverityLow, &msgSeverity), "severity", "s", "Severity of the notification. Possible values: low, medium, high, critical")
}
