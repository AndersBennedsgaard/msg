package cmd

import (
	"fmt"

	"github.com/AndersBennedsgaard/msg/cmd/flags"
	"github.com/AndersBennedsgaard/msg/internal/notification"
	"github.com/AndersBennedsgaard/msg/internal/store"
	"github.com/spf13/cobra"
)

var (
	listRead   bool
	severities []notification.NotificationSeverity
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List notification messages",
	Long:  `Prints a list of notification messages stored in the system, optionally filtered by type or severity.`,
	Run: func(cmd *cobra.Command, args []string) {
		fsstore := store.NewFSStore(cfg.BasePath)

		states := []notification.MessageStatus{notification.StatusUnread}
		if listRead {
			states = append(states, notification.StatusRead)
		}

		typeFilter := ""

		messages, err := fsstore.ListMessages(store.MessageFilter{
			Statuses:   states,
			Types:      []string{typeFilter},
			Severities: severities,
		})
		cobra.CheckErr(err)

		if len(messages) == 0 {
			fmt.Println("No messages found.")
			return
		} else {
			fmt.Printf("Found %d message(s):\n", len(messages))
		}

		for _, msg := range messages {
			fmt.Printf("- [%s] (%s) %s: %s\n", msg.ID(), msg.Severity(), msg.Type(), msg.Message())
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&listRead, "read", "r", false, "List read messages")
	listCmd.Flags().VarP(flags.NewSeverityListValue([]notification.NotificationSeverity{}, &severities), "severity", "s", "Filter messages by severity (comma-separated list of low, medium, high, critical)")
}
