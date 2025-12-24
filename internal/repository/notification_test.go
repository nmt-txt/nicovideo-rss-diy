package repository

import (
	"errors"
	"fmt"
	"testing"
)

func TestAddNotification(t *testing.T) {

	repo := NewNotificationRepository()

	if len(repo.Notifications) != 0 {
		t.Fatalf("expected 0 notifications initially, got %d", len(repo.Notifications))
	}

	desc := errors.New("TestError")
	for i := 0; i < 4; i++ {
		dup := false
		if i >= 1 {
			dup = true
		}

		lv := NotificationInfo
		if i >= 2 {
			lv = NotificationError
		}

		repo.AddNotification(
			lv,
			fmt.Sprintf("Test%d", i),
			desc,
			dup,
		)
	}

	for i, n := range repo.Notifications {
		t.Logf("#%d [%s] %s, %s | dup:%t(%d)", i, n.Level.String(), n.Title, n.Description.Error(), n.AllowDuplication, n.DuplicateCount)
	}

	// 重複動作チェック
	if len(repo.Notifications) != 2 {
		t.Fatalf("expected 2 notifications, got %d", len(repo.Notifications))
	}

	// DuplicateCountチェック
	if repo.Notifications[1].DuplicateCount != 2 {
		t.Fatalf("expected DuplicateCount 2 for second notification, got %d", repo.Notifications[1].DuplicateCount)
	}

	// Levelチェック・上書きチェック
	if repo.Notifications[0].Level != NotificationInfo {
		t.Fatalf("expected Level NotificationInfo for first notification, got %s", repo.Notifications[0].Level.String())
	}

	if repo.Notifications[1].Level != NotificationError {
		t.Fatalf("expected Level NotificationError for second notification, got %s", repo.Notifications[1].Level.String())
	}
}
