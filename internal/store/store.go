package store

import "github.com/AndersBennedsgaard/msg/internal/notification"

type Store interface {
	AddMessage(notification *notification.Message) error
	GetMessage(id string) (*notification.Message, error)
	ListMessages(state notification.MessageStatus) ([]*notification.Message, error)
	MoveMessage(id string, oldStatus, newStatus notification.MessageStatus) error
}
