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

`msg` uses SQLite for storage.

By default it uses `$XDG_DATA_HOME/msg/db.sql` (usually `~/.local/share/msg/db.sql`) for the path of the database.

---

## Core Actions

The system is controlled via the `msg` command:

* **`add`**: Create a new notification.
  ```bash
  msg add --type ci_fail --message "Pipeline #402 failed"
  ```

* **`count`**: Return the number of unread messages (ideal for tmux/polybar).
  ```bash
  msg count  # Returns "5"
  ```

* **`read`**: Read the next unread message.
  ```bash
  msg read
  ```

* **`list`**: Print a list of unread messages for a quick overview.
  ```bash
  msg list
  ```
