package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/AndersBennedsgaard/msg/cmd/flags"
	"github.com/AndersBennedsgaard/msg/internal/notification"
	"github.com/AndersBennedsgaard/msg/internal/store"
	"github.com/spf13/cobra"
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
	DisableFlagsInUseLine: true,
	Example: ` # Add a new notification message:
 msg add -t alert -s high -m "Disk space low"

 # Add a new notification message by reading content from stdin:
 echo "CPU usage high" | msg add -t alert -s critical`,
	Short: "Add a new notification message",
	Long:  `Add a new notification message to the store with specified type, severity, and content.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputReader := cmd.InOrStdin()

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

		fsstore := store.NewFSStore(cfg.BasePath)

		now := time.Now()
		id := fmt.Sprintf("%d", now.UnixNano())

		noti, err := notification.NewNotification(
			notification.NotificationId(id),
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
		err = fsstore.AddMessage(msg)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Added message with ID: %s\n", msg.ID())
		return err
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&msgType, "type", "t", "info", "Type of the notification")
	addCmd.Flags().StringVarP(&msgContent, "message", "m", "", "Content of the notification")
	addCmd.Flags().VarP(flags.NewSeverityValue(notification.NotificationSeverityLow, &msgSeverity), "severity", "s", "Severity of the notification. Possible values: low, medium, high, critical")
}
