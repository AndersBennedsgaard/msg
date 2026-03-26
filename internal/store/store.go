package store

import (
	"github.com/AndersBennedsgaard/msg/internal/notification"
)

type Store interface {
	AddMessage(msg *notification.Message) (int64, error)
	ReadNextMessage() (*notification.Message, error)
	CountUnreadMessages() (int, error)
}
