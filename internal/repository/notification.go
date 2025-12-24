package repository

import (
	"errors"
	"time"
)

type NotificationLevel int

const (
	NotificationInfo NotificationLevel = iota
	NotificationError
)

func (n NotificationLevel) String() string {
	switch n {
	case NotificationInfo:
		return "INFO"
	case NotificationError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type Notification struct {
	Level            NotificationLevel
	Title            string
	Description      error
	Date             time.Time
	AllowDuplication bool
	DuplicateCount   int
}

// NotificationRepository RSSフィード上で通知したい項目を保持する
type NotificationRepository struct {
	Notifications []Notification
}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{
		Notifications: make([]Notification, 0, 10),
	}
}

// AddNotification 通知項目を追加する。
// allowDuplicationがtrueの場合、同じdescriptionの通知が存在する場合かつ既存の通知もDuplicationを許可している場合は引数の方のもので上書きする。
func (r *NotificationRepository) AddNotification(
	level NotificationLevel,
	title string,
	description error,
	allowDuplication bool,
) {

	new := Notification{
		Level:            level,
		Title:            title,
		Description:      description,
		AllowDuplication: allowDuplication,
		Date:             time.Now(),
		DuplicateCount:   0,
	}

	if allowDuplication {
		for i, n := range r.Notifications {
			if errors.Is(n.Description, description) && n.AllowDuplication {
				new.DuplicateCount = n.DuplicateCount + 1
				r.Notifications[i] = new
				return
			}
		}
	}

	r.Notifications = append(r.Notifications, new)
}

func (r *NotificationRepository) ClearNotifications() {
	r.Notifications = make([]Notification, 0, 10)
}
