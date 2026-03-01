package flags

import (
	"fmt"

	"github.com/AndersBennedsgaard/msg/internal/notification"
)

type SeverityValue struct {
	value *notification.NotificationSeverity
}

func NewSeverityValue(
	defaultValue notification.NotificationSeverity,
	target *notification.NotificationSeverity,
) *SeverityValue {
	*target = defaultValue
	return &SeverityValue{value: target}
}

func (s *SeverityValue) Set(input string) error {
	sev := notification.NotificationSeverity(input)

	if !notification.IsValidSeverity(sev) {
		return fmt.Errorf(
			"invalid severity %q (allowed: low, medium, high, critical)",
			input,
		)
	}

	*s.value = sev
	return nil
}

func (s *SeverityValue) String() string {
	if s.value == nil {
		return ""
	}
	return string(*s.value)
}

func (s *SeverityValue) Type() string {
	return "severity"
}
