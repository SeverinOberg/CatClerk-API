package database

import (
	"fmt"
)

// Settings structure
type Settings struct {
	Notifications bool `json:"notifications"`
}

// GetNotificationSetting gets the account's notification setting preference from the database by username
func (handler *Handler) GetNotificationSetting(username string) (Settings, error) {
	response := Settings{}

	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		SELECT notifications FROM accounts
		WHERE username="%s"	
	`, username))
	if err != nil {
		return response, err
	}

	defer stmt.Close()

	if err := stmt.QueryRow().Scan(
		&response.Notifications,
	); err != nil {
		return response, err
	}

	return response, err
}

// ToggleNotificationSetting toggles the notification setting by username
func (handler *Handler) ToggleNotificationSetting(username string) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		UPDATE accounts
		SET 
			notifications = !notifications
		WHERE username = "%s";
	`, username))
	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec()
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected < 1 {
		return fmt.Errorf("no rows affected")
	}

	return err
}
