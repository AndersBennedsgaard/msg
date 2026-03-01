package flags

import (
	"fmt"
	"strings"

	"github.com/AndersBennedsgaard/msg/internal/notification"
)

type SeverityListValue struct {
	values *[]notification.NotificationSeverity
}

func NewSeverityListValue(
	defaultValue []notification.NotificationSeverity,
	target *[]notification.NotificationSeverity,
) *SeverityListValue {
	*target = defaultValue
	return &SeverityListValue{values: target}
}

func (s *SeverityListValue) Set(input string) error {
	if input == "" {
		*s.values = []notification.NotificationSeverity{}
		return nil
	}

	parts := strings.Split(input, ",")
	var severities []notification.NotificationSeverity
	for _, part := range parts {
		sev := notification.NotificationSeverity(part)

		if !notification.IsValidSeverity(sev) {
			return fmt.Errorf(
				"invalid severity %q (allowed: low, medium, high, critical)",
				part,
			)
		}
		severities = append(severities, sev)
	}

	*s.values = severities
	return nil
}

func (s *SeverityListValue) String() string {
	if s.values == nil {
		return ""
	}

	out := make([]string, len(*s.values))
	for i, sev := range *s.values {
		out[i] = string(sev)
	}
	return strings.Join(out, ",")
}

func (s *SeverityListValue) Type() string {
	return "severity-list"
}
