package store

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/AndersBennedsgaard/msg/internal/notification"

	_ "modernc.org/sqlite"
)

type SqliteStore struct {
	db *sql.DB
}

func NewSqliteStore(db *sql.DB) (*SqliteStore, error) {
	query := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		message TEXT,
		severity TEXT,
		type TEXT,
		created_at DATETIME,
		read INTEGER DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_unread ON events(read, created_at);
	`
	_, err := db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("error initializing database: %w", err)
	}

	return &SqliteStore{db: db}, nil
}

func (store *SqliteStore) AddMessage(msg *notification.Message) (int64, error) {
	res, err := store.db.Exec(
		`INSERT INTO events (message, severity, type, created_at, read)
		 VALUES (?, ?, ?, ?, 0)`,
		msg.Message(), msg.Severity(), msg.Type(), time.Now(),
	)
	if err != nil {
		return 0, fmt.Errorf("error occurred when storing a message: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error occurred when retrieving the id of the last message: %w", err)
	}

	return id, err
}

func readNextMessage(tx *sql.Tx) (*notification.Notification, error) {
	row := tx.QueryRow(`
		SELECT id, message, severity, type, created_at
		FROM events
		WHERE read = 0
		ORDER BY created_at
		LIMIT 1
	`)

	var id int64
	var msg, severity, eventType string
	var created time.Time

	err := row.Scan(&id, &msg, &severity, &eventType, &created)
	if err == sql.ErrNoRows {
		return nil, errors.New("no unread messages")
	}
	if err != nil {
		return nil, fmt.Errorf("unknown error occurred during read from database: %w", err)
	}

	_, err = tx.Exec(`UPDATE events SET read = 1 WHERE id = ?`, id)
	if err != nil {
		return nil, fmt.Errorf("error updating notification status: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing the database transaction: %w", err)
	}

	notif, err := notification.NewNotification(eventType, created, notification.NotificationSeverity(severity), msg)
	return &notif, err
}

func (store *SqliteStore) ReadNextMessage() (*notification.Message, error) {
	tx, err := store.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}

	notif, err := readNextMessage(tx)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, fmt.Errorf("error occurred during rollback: %s", err)
		}
	}
	if notif == nil {
		// if no message is found, return no message
		return nil, nil
	}

	message := notification.Message{
		Notification: *notif,
		Status:       notification.StatusRead,
	}
	return &message, nil
}

func (store *SqliteStore) ListUnreadMessages(limit int) ([]notification.Notification, error) {
	rows, err := store.db.Query(`
		SELECT id, message, severity, type, created_at
		FROM events
		WHERE read = 0
		ORDER BY created_at
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Fatalf("error closing database rows: %s", err)
		}
	}()

	var messages []notification.Notification

	for rows.Next() {
		var id int64
		var msg, severity, eventType string
		var created time.Time

		err := rows.Scan(&id, &msg, &severity, &eventType, &created)
		if err != nil {
			return nil, err
		}

		notif, err := notification.NewNotification(eventType, created, notification.NotificationSeverity(severity), msg)
		if err != nil {
			return nil, fmt.Errorf("error when creating new notification: %w", err)
		}
		messages = append(messages, notif)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (store *SqliteStore) CountUnreadMessages() (int, error) {
	var total int
	err := store.db.QueryRow(
		`SELECT COUNT(*) FROM events WHERE read = 0`,
	).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("error occurred when counting unread messages: %w", err)
	}

	return total, nil
}
