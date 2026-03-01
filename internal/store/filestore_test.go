package store_test

import (
	"os"
	"testing"
	"time"

	"github.com/AndersBennedsgaard/msg/internal/notification"
	"github.com/AndersBennedsgaard/msg/internal/store"
)

func TestFSStore_SaveAndGet(t *testing.T) {
	basePath := "./testdata"
	defer func() {
		if t.Failed() {
			t.Logf("Test failed, preserving test data at %s for inspection", basePath)
			return
		}

		err := os.RemoveAll(basePath)
		if err != nil {
			t.Fatalf("Failed to clean up test data: %v", err)
		}
	}()

	fsStore := store.NewFSStore(basePath)

	now := time.Now()
	noti, err := notification.NewNotification(
		"test-id",
		"type1",
		now,
		"high",
		"notification message content",
	)
	if err != nil {
		t.Fatalf("Failed to create notification: '%v'", err)
	}

	msg := &notification.Message{
		Notification: noti,
		Status:       notification.StatusUnread,
	}

	err = fsStore.AddMessage(msg)
	if err != nil {
		t.Fatalf("Failed to save message: '%v'", err)
	}

	retrievedMsg, err := fsStore.GetMessage(string(msg.ID()))
	if err != nil {
		t.Fatalf("Failed to get message: '%v'", err)
	}

	// Only checking ID and Status.
	// The render and parse process are tested in the notification package tests.
	if retrievedMsg.ID() != msg.ID() {
		t.Errorf("Expected ID '%s', got '%s'", msg.ID(), retrievedMsg.ID())
	}
	if retrievedMsg.Status != msg.Status {
		t.Errorf("Expected Status '%s', got '%s'", msg.Status, retrievedMsg.Status)
	}
}
