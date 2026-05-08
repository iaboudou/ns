package sqlite

import (
	"database/sql"

	"rtf/models"

	"github.com/gofrs/uuid/v5"
)

func InsertNotifInDB(db *sql.DB, notif *models.Notification, notifID uuid.UUID, now string, receivers []string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(`
		INSERT INTO notifications (id, sender_id, type, group_id, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, notifID, notif.SenderID, notif.Type, notif.GroupID, now)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO notification_users (notification_id, user_id, is_read)
		VALUES (?, ?, 0)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, uid := range receivers {
		if _, err := stmt.Exec(notifID, uid); err != nil {
			return err
		}
	}

	return tx.Commit()
}
