# Local Notification Inbox (`msg`)

A persistent, filesystem-backed notification manager for Unix-like environments.
It provides a "stateful inbox" for command-line events, bridging the gap between ephemeral desktop notifications and infinite logs.

---

## Why?

Most system events are either **too loud** (desktop popups that vanish), **too heavy** (email), or **too static** (logs).
`msg` provides a middle ground:

* **Stateful:** Tracks unread, read, and dismissed states.
* **Persistent:** Notifications survive reboots and shell sessions.
* **Decoupled:** Any script can "fire and forget" a message without a running daemon.
* **Human-Centric:** Designed for users who live in the terminal and want an "Inbox Zero" workflow for their local machine.

---

## Data Storage & Performance

### RFC-822 Format

To maintain Unix-friendliness and human readability, notifications are stored as **RFC-822–style headers** (similar to email or Debian packages) rather than JSON.
This allows for easy inspection using `grep`, `sed`, or `cat`.

```rfc822
ID: 1694524800-7f2a
Timestamp: 2023-09-12T14:40:00Z
Type: process_crash
Source: backup-daemon
Severity: error

The backup script failed with exit code 1. Check /var/log/backup.log.
```

### Performance & Concurrency

* **One File Per Message:** Counting unread messages is a simple `ls | wc -l` operation.
  No need to parse a massive database file to get a status update.
* **Atomic Operations:** Notifications are written to a `tmp` directory and then moved (`mv`) to `unread`.
  On Unix, `rename()` is atomic, ensuring no partial reads or race conditions between concurrent writers and readers.

---

## Storage Structure

The system uses a simple directory tree under `$XDG_DATA_HOME/msg/` (usually `~/.local/share/msg/`):

```bash
.local/share/msg/
├── unread/      # New notifications
├── read/        # Viewed notifications
└── dismissed/   # Acknowledged/Archived notifications

```

---

## Core Actions

The system is controlled via the `msg` command:

* **`add`**: Create a new notification.
  ```bash
  msg add --type ci_fail "Pipeline #402 failed"
  ```

* **`count`**: Return the number of unread messages (ideal for tmux/polybar).
  ```bash
  msg count  # Returns "5"
  ```

* **`list`**: Show a summary of notifications.
  ```bash
  msg list [--unread | --read | --all]
  ```

* **`show <id>`**: Display the full message and move it from `unread/` to `read/`.
* **`unread <id>`**: Move a message from `read/` back to `unread/`.
* **`dismiss <id>`**: Move a message to `dismissed/` without reading it.
* **`clear`**: Permanently delete all dismissed messages.
