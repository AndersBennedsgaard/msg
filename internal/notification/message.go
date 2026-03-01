package notification

type MessageStatus string

const (
	StatusUnread   MessageStatus = "unread"
	StatusRead     MessageStatus = "read"
	StatusArchived MessageStatus = "archived"
)

var AllStatuses = []MessageStatus{
	StatusUnread,
	StatusRead,
	StatusArchived,
}

type Message struct {
	Notification

	Status MessageStatus
}
