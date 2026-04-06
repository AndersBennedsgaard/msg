package cmd

import (
	"database/sql"
	"fmt"
	"io"
	"strings"

	"github.com/AndersBennedsgaard/msg/internal/logging"
	"github.com/AndersBennedsgaard/msg/internal/notification"
	"github.com/AndersBennedsgaard/msg/internal/store"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var limit int

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	GroupID: appGroup,
	Short:   "Print the number of unread messages",
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

		list, err := dbstore.ListUnreadMessages(limit)
		if err != nil {
			return fmt.Errorf("error occured when listing unread messages: %w", err)
		}

		_, err = io.WriteString(cmd.OutOrStdout(), formatMessages(list))
		return err
	},
}

func formatMessages(messages []notification.Notification) string {
	const timeLayout = "2006-01-02 15:04:05"

	// Headers
	hTimestamp := "Timestamp"
	hSeverity := "Severity"
	hType := "Type"

	// Initialize column widths with header lengths
	maxTimestampWidth := len(hTimestamp)
	maxSeverityWidth := len(hSeverity)
	maxTypeWidth := len(hType)

	// First pass: compute max widths
	for _, m := range messages {
		ts := m.Timestamp().Format(timeLayout)
		maxTimestampWidth = max(maxTimestampWidth, len(ts))
		maxSeverityWidth = max(maxSeverityWidth, len(m.Severity()))
		maxTypeWidth = max(maxTypeWidth, len(m.Type()))
	}

	// Helper for padding (left-align)
	pad := func(s string, width int) string {
		if len(s) >= width {
			return s
		}
		return s + strings.Repeat(" ", width-len(s))
	}

	var b strings.Builder

	// Header row
	b.WriteString(pad(hTimestamp, maxTimestampWidth))
	b.WriteString("  ")
	b.WriteString(pad(hSeverity, maxSeverityWidth))
	b.WriteString("  ")
	b.WriteString(pad(hType, maxTypeWidth))
	b.WriteString("\n")

	// Data rows
	for _, m := range messages {
		ts := m.Timestamp().Format(timeLayout)

		b.WriteString(pad(ts, maxTimestampWidth))
		b.WriteString("  ")
		b.WriteString(pad(string(m.Severity()), maxSeverityWidth))
		b.WriteString("  ")
		b.WriteString(pad(m.Type(), maxTypeWidth))
		b.WriteString("\n")
	}

	return b.String()
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().IntVarP(&limit, "limit", "l", 5, "Maximum number of messages to list")
}
