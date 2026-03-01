package store

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/AndersBennedsgaard/msg/internal/notification"
)

type FSStore struct {
	basePath string
}

func NewFSStore(basePath string) *FSStore {
	// ensure status directories exist
	for _, status := range notification.AllStatuses {
		dirPath := filepath.Join(basePath, string(status))
		if stat, err := os.Stat(dirPath); os.IsNotExist(err) {
			fmt.Printf("Creating directory for status %s at %s\n", status, dirPath)
			inErr := os.MkdirAll(dirPath, os.ModePerm)
			if inErr != nil {
				panic(inErr)
			}
		} else if err != nil {
			panic(err)
		} else if !stat.IsDir() {
			panic(fmt.Sprintf("Path %s is not a directory", dirPath))
		}
	}

	return &FSStore{basePath: basePath}
}

func (fs *FSStore) AddMessage(msg *notification.Message) error {
	data := notification.Render(msg.Notification)
	filePath := filepath.Join(fs.basePath, string(msg.Status), string(msg.ID()))
	return os.WriteFile(filePath, []byte(data), 0644)
}

func locateFile(basePath, id string) (string, notification.MessageStatus) {
	for _, status := range notification.AllStatuses {
		filePath := filepath.Join(basePath, string(status), id)
		if _, err := os.Stat(filePath); err == nil {
			return filePath, status
		}
	}
	return "", ""
}

func (fs *FSStore) GetMessage(id string) (*notification.Message, error) {
	path, state := locateFile(fs.basePath, id)
	if path == "" {
		return nil, os.ErrNotExist
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Trim any trailing newlines
	content := strings.TrimRight(string(data), "\r\n")

	noti, err := notification.Parse(content)
	if err != nil {
		return nil, err
	}
	return &notification.Message{
		Notification: noti,
		Status:       state,
	}, nil
}

func (fs *FSStore) MoveMessage(id string, oldStatus, newStatus notification.MessageStatus) error {
	oldPath := filepath.Join(fs.basePath, string(oldStatus), id)
	newPath := filepath.Join(fs.basePath, string(newStatus), id)
	return os.Rename(oldPath, newPath)
}

type MessageFilter struct {
	Statuses   []notification.MessageStatus
	Types      []string
	Severities []notification.NotificationSeverity
}

// TODO: implement filtering logic
func (fs *FSStore) ListMessages(filter MessageFilter) ([]*notification.Message, error) {
	var allMessages []*notification.Message

	for _, status := range filter.Statuses {
		dirPath := filepath.Join(fs.basePath, string(status))
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			data, err := os.ReadFile(filepath.Join(dirPath, entry.Name()))
			if err != nil {
				return nil, err
			}
			noti, err := notification.Parse(string(data))
			if err != nil {
				return nil, err
			}
			if len(filter.Types) == 0 || !slices.Contains(filter.Types, noti.Type()) {
				continue
			}
			if len(filter.Severities) == 0 || !slices.Contains(filter.Severities, noti.Severity()) {
				continue
			}
			allMessages = append(allMessages, &notification.Message{
				Notification: noti,
				Status:       status,
			})
		}
	}
	return allMessages, nil
}
