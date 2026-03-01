package notification_test

import (
	"errors"
	"testing"
	"time"

	"github.com/AndersBennedsgaard/msg/internal/notification"
)

func TestRender(t *testing.T) {
	utcLoc, err := time.LoadLocation("")
	if err != nil {
		t.Fatal("Error loading location")
	}

	timestamp := time.Date(2025, time.November, 10, 8, 15, 0, 0, utcLoc)
	n, err := notification.NewNotification("someId", "type", timestamp, "high", "some message")
	if err != nil {
		t.Fatalf("Error creating notification: %v", err)
	}

	expected := `ID: someId
Type: type
Timestamp: 10 Nov 25 08:15 UTC
Severity: high

some message
`

	actual := notification.Render(n)
	if actual != expected {
		t.Fatalf("Rendered notification is not equal to the expected. Expected: %s. Actual: %s", expected, actual)
	}
}

func TestParse(t *testing.T) {
	utcLoc, err := time.LoadLocation("")
	if err != nil {
		t.Fatal("Error loading location")
	}

	timestamp := time.Date(2025, time.November, 10, 8, 15, 0, 0, utcLoc)
	validNotification, err := notification.NewNotification("someId", "type", timestamp, "high", "some message")
	if err != nil {
		t.Fatalf("Error creating notification: %v", err)
	}

	testCases := []struct {
		name        string
		input       string
		expected    *notification.Notification
		expectedErr error
	}{
		{
			name: "valid input",
			input: `ID: someId
Type: type
Timestamp: 10 Nov 25 08:15 UTC
Severity: high

some message`,
			expected:    &validNotification,
			expectedErr: nil,
		},
		{
			name: "invalid timestamp input",
			input: `ID: someId
Type: type
Timestamp: invalid-timestamp
Severity: high

some message
`,
			expected:    nil,
			expectedErr: notification.ErrInvalidTimestampFormat,
		},
		{
			name:        "invalid input",
			input:       "wrong input",
			expected:    nil,
			expectedErr: notification.ErrInvalidFormat,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual, err := notification.Parse(tc.input)

			if tc.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error '%v', got nil", tc.expectedErr)
				}
				if !errors.Is(err, tc.expectedErr) {
					t.Fatalf("expected error '%v', got '%v'", tc.expectedErr, err)
				}
				return // no value assertions when an error is expected
			}

			if err != nil {
				t.Fatalf("unexpected error: '%v'", err)
			}

			if actual.ID() != tc.expected.ID() {
				t.Fatalf("expected ID '%s', got '%s'", tc.expected.ID(), actual.ID())
			}
			if actual.Type() != tc.expected.Type() {
				t.Fatalf("expected Type '%s', got '%s'", tc.expected.Type(), actual.Type())
			}
			if !actual.Timestamp().Equal(tc.expected.Timestamp()) {
				t.Fatalf("expected Timestamp '%v', got '%v'", tc.expected.Timestamp(), actual.Timestamp())
			}
			if actual.Severity() != tc.expected.Severity() {
				t.Fatalf("expected Severity '%s', got '%s'", tc.expected.Severity(), actual.Severity())
			}
			if actual.Message() != tc.expected.Message() {
				t.Fatalf("expected Message '%s', got '%s'", tc.expected.Message(), actual.Message())
			}

			if !actual.Equals(*tc.expected) {
				t.Fatalf("expected notifications to be equal. Actual: %+v, Expected: %+v", actual, tc.expected)
			}
		})
	}
}
