package database

import (
	"fmt"
	"time"
)

// Folder structure
type Folder struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
	Count     int       `json:"count"`
}

// CreateStorage creates a storage and attaches the account to it by username
func (handler *Handler) CreateStorage(username, title string, owner bool) (int64, error) {
	lastInsertID := int64(0)
	stmtS, err := handler.DB.Prepare(fmt.Sprintf(`
		INSERT INTO storages(title)
		VALUES("%s");
	`, title))
	if err != nil {
		return lastInsertID, err
	}

	defer stmtS.Close()

	result, err := stmtS.Exec()
	if err != nil {
		return lastInsertID, err
	}

	lastInsertID, err = result.LastInsertId()
	if err != nil {
		return lastInsertID, err
	}

	stmtASB, err := handler.DB.Prepare(fmt.Sprintf(`
	INSERT INTO account_storage_binder(username, storage_id, owner)
	VALUES("%s", %d, %t);
	`, username, lastInsertID, owner))
	if err != nil {
		return lastInsertID, err
	}

	defer stmtASB.Close()

	_, err = stmtASB.Exec()
	if err != nil {
		return lastInsertID, err
	}

	return lastInsertID, err
}

// GetStorages gets all storages by username
func (handler *Handler) GetStorages(username string) ([]Folder, error) {
	folders := []Folder{}

	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		SELECT s.*, COUNT(si.id) AS "count"
		FROM storages AS s
		LEFT JOIN storage_items AS si
		ON s.id = si.storage_id
		INNER JOIN account_storage_binder AS asb
		ON asb.storage_id = s.id AND asb.username = "%s"
		GROUP BY s.id
	`, username))
	if err != nil {
		return folders, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return folders, err
	}

	for rows.Next() {
		folder := Folder{}

		if err := rows.Scan(
			&folder.ID,
			&folder.Title,
			&folder.UpdatedAt,
			&folder.CreatedAt,
			&folder.Count,
		); err != nil {
			return folders, err
		}

		folders = append(folders, folder)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return folders, err
}

// GetStorageFoldersCount gets the amount of storages by username
func (handler *Handler) GetStoragesCount(username string) (int, error) {
	count := 0

	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		SELECT COUNT(*)
		FROM account_storage_binder
		WHERE username="%s"; 
	`, username))
	if err != nil {
		return count, err
	}

	defer stmt.Close()

	if err := stmt.QueryRow().Scan(
		&count,
	); err != nil {
		return count, err
	}

	return count, err
}

// UpdateStorage updates a storage by ID
func (handler *Handler) UpdateStorage(title string, storageID int) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		UPDATE storages AS s
		SET s.title = "%s"
		WHERE s.id = %d
	`, title, storageID))
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

// DeleteStorage deletes a storage by ID
func (handler *Handler) DeleteStorage(storageID int) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
	DELETE FROM storages AS s
	WHERE s.id = %d
	`, storageID))
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

// Share structure
type Share struct {
	ID        int `json:"id"`
	Username  int `json:"username"`
	StorageID int `json:"storageID"`
	Owner     int `json:"owner"`
}

// ShareStorage attaches a storage to an account by username and ID
func (handler *Handler) ShareStorage(username string, storageID int) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		INSERT INTO account_storage_binder(username, storage_id)
		VALUES("%s", %d)
	`, username, storageID))
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

// RemoveShareStorage removes an accounts attachment to a storage by username and ID
func (handler *Handler) RemoveShareStorage(usernameRequest string, storageID int) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		DELETE FROM account_storage_binder
		WHERE username="%s" AND storage_id=%d
	`, usernameRequest, storageID))
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

// GetStorageOwner returns true or false whether it is the storage owner by ID
func (handler *Handler) GetStorageOwner(owner string, storageID int) (Share, error) {
	payload := Share{}

	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
	SELECT owner
	FROM account_storage_binder
	WHERE username="%s" AND storage_id=%d
	`, owner, storageID))
	if err != nil {
		return payload, err
	}

	defer stmt.Close()

	if err := stmt.QueryRow().Scan(
		&payload.Owner,
	); err != nil {
		return payload, err
	}

	return payload, err
}
