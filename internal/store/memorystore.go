package store

import (
	"errors"
	"os"

	"github.com/AndersBennedsgaard/msg/internal/notification"
)

type InMemoryStore struct {
	messages map[string]*notification.Message
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		messages: make(map[string]*notification.Message),
	}
}

func (store *InMemoryStore) AddMessage(msg *notification.Message) error {
	store.messages[string(msg.ID())] = msg
	return nil
}

func (store *InMemoryStore) GetMessage(id string) (*notification.Message, error) {
	msg, exists := store.messages[id]
	if !exists {
		return nil, os.ErrNotExist
	}
	return msg, nil
}

func (store *InMemoryStore) ListMessages(state notification.MessageStatus) ([]*notification.Message, error) {
	var result []*notification.Message
	for _, msg := range store.messages {
		if msg.Status == state {
			result = append(result, msg)
		}
	}
	return result, nil
}

func (store *InMemoryStore) MoveMessage(id string, oldStatus, newStatus notification.MessageStatus) error {
	msg, exists := store.messages[id]
	if !exists {
		return os.ErrNotExist
	}
	if msg.Status != oldStatus {
		return errors.New("message is not in the expected status")
	}
	msg.Status = newStatus
	return nil
}
