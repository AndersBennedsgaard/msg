package notification

type MessageStatus string

const (
	StatusUnread MessageStatus = "unread"
	StatusRead   MessageStatus = "read"
)

var AllStatuses = []MessageStatus{
	StatusUnread,
	StatusRead,
}

type Message struct {
	Notification

	Status MessageStatus
}
