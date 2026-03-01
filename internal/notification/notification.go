package notification

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type NotificationId string
type NotificationSeverity string

const (
	NotificationSeverityLow      NotificationSeverity = "low"
	NotificationSeverityMedium   NotificationSeverity = "medium"
	NotificationSeverityHigh     NotificationSeverity = "high"
	NotificationSeverityCritical NotificationSeverity = "critical"
)

type Notification struct {
	id        NotificationId
	_type     string
	timestamp time.Time
	severity  NotificationSeverity
	message   string
}

var ErrInvalidFormat = errors.New("invalid notification format")
var ErrInvalidTimestampFormat = errors.New("invalid timestamp format")

type ErrInvalidSeverity struct {
	Value string
}

func (e *ErrInvalidSeverity) Error() string {
	return fmt.Sprintf("invalid notification severity: %q", e.Value)
}

func IsValidSeverity(sev NotificationSeverity) bool {
	return sev == NotificationSeverityLow ||
		sev == NotificationSeverityMedium ||
		sev == NotificationSeverityHigh ||
		sev == NotificationSeverityCritical
}

func NewNotification(id NotificationId, _type string, timestamp time.Time, severity NotificationSeverity, message string) (Notification, error) {
	if !IsValidSeverity(severity) {
		return Notification{}, &ErrInvalidSeverity{Value: string(severity)}
	}

	return Notification{
		id:        id,
		_type:     _type,
		timestamp: timestamp,
		severity:  severity,
		message:   message,
	}, nil
}

func (n Notification) ID() NotificationId {
	return n.id
}

func (n Notification) Type() string {
	return n._type
}

func (n Notification) Timestamp() time.Time {
	return n.timestamp
}

func (n Notification) Severity() NotificationSeverity {
	return n.severity
}

func (n Notification) Message() string {
	return n.message
}

func (n Notification) Equals(other Notification) bool {
	if n.ID() != other.ID() {
		return false
	}
	if n.Type() != other.Type() {
		return false
	}
	if !n.Timestamp().Equal(other.Timestamp()) {
		return false
	}
	if n.Severity() != other.Severity() {
		return false
	}
	if n.Message() != other.Message() {
		return false
	}

	return true
}

func Render(n Notification) string {
	timestampStr := n.timestamp.Format(time.RFC822)

	return "ID: " + string(n.id) + "\n" +
		"Type: " + string(n._type) + "\n" +
		"Timestamp: " + timestampStr + "\n" +
		"Severity: " + string(n.severity) + "\n\n" +
		n.message + "\n"
}

func Parse(input string) (Notification, error) {
	lines := strings.Split(input, "\n")

	if len(lines) < 5 {
		return Notification{}, ErrInvalidFormat
	}

	id := NotificationId(strings.TrimPrefix(lines[0], "ID: "))
	_type := string(strings.TrimPrefix(lines[1], "Type: "))
	timestampStr := strings.TrimPrefix(lines[2], "Timestamp: ")
	severity := NotificationSeverity(strings.TrimPrefix(lines[3], "Severity: "))
	message := strings.Join(lines[5:], "\n")

	timestamp, err := time.Parse(time.RFC822, timestampStr)
	if err != nil {
		return Notification{}, ErrInvalidTimestampFormat
	}

	return NewNotification(id, _type, timestamp, severity, message)
}
