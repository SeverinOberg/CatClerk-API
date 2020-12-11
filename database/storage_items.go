package database

import (
	"fmt"
	"time"
)

// Item structure
type Item struct {
	ID                  int       `json:"id"`
	StorageID           int       `json:"stroageID"`
	Title               string    `json:"title"`
	Image               string    `json:"image"`
	Quantity            int       `json:"quantity"`
	QuantityType        string    `json:"quantityType"`
	QuantityThreshold   int       `json:"quantityThreshold"`
	ExpirationThreshold int       `json:"expirationThreshold"`
	ExpirationDate      string    `json:"expirationDate"`
	UpdatedAt           time.Time `json:"updatedAt"`
	CreatedAt           time.Time `json:"createdAt"`
}

// CreateStorageItem creates a storage item and attaches it to an FK storageID
func (handler *Handler) CreateStorageItem(storageID int, title string, quantity int, quantityType string, quantityThreshold int, expirationThreshold int, expirationDate string) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		INSERT INTO storage_items(storage_id, title, quantity, quantity_type, quantity_threshold, expiration_threshold, expiration_date)
		VALUES(%d, "%s", %d, "%s", %d, %d, "%s");
	`, storageID, title, quantity, quantityType, quantityThreshold, expirationThreshold, expirationDate))
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return err
}

// GetStorageItems gets storage items by storageID
func (handler *Handler) GetStorageItems(storageID int) ([]Item, error) {
	items := []Item{}

	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
	SELECT *
	FROM storage_items AS si
	WHERE si.storage_id = %d
	`, storageID))
	if err != nil {
		return items, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return items, err
	}

	for rows.Next() {
		item := Item{}

		if err := rows.Scan(
			&item.ID,
			&item.StorageID,
			&item.Title,
			&item.Image,
			&item.Quantity,
			&item.QuantityType,
			&item.QuantityThreshold,
			&item.ExpirationThreshold,
			&item.ExpirationDate,
			&item.UpdatedAt,
			&item.CreatedAt,
		); err != nil {
			return items, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return items, err
}

// GetStorageItem gets a storage item by storageID and ID
func (handler *Handler) GetStorageItem(storageID, itemID int) (Item, error) {
	item := Item{}

	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
	SELECT *
	FROM storage_items
	WHERE storage_id = %d AND id = %d  
	`, storageID, itemID))
	if err != nil {
		return item, err
	}

	defer stmt.Close()

	if err := stmt.QueryRow().Scan(
		&item.ID,
		&item.StorageID,
		&item.Title,
		&item.Image,
		&item.Quantity,
		&item.QuantityType,
		&item.QuantityThreshold,
		&item.ExpirationThreshold,
		&item.ExpirationDate,
		&item.UpdatedAt,
		&item.CreatedAt,
	); err != nil {
		return item, err
	}

	return item, err
}

// GetStorageItemsCount gets the amount of storage items by username
func (handler *Handler) GetStorageItemsCount(username string) (int, error) {
	count := 0

	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		SELECT COUNT(*)
		FROM storage_items AS si
		INNER JOIN storages  AS s
		ON si.storage_id = s.id
		INNER JOIN account_storage_binder AS asb
		WHERE s.id = asb.storage_id AND asb.username="%s";
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

// UpdateStorageItem updates a storage item by ID
func (handler *Handler) UpdateStorageItem(title, image string, quantity int, quantityType string, quantityThreshold, expirationThreshold int, expirationDate string, itemID int) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		UPDATE storage_items AS si
		SET 
			si.title = "%s",
			si.image = "%s",
			si.quantity = %d,
			si.quantity_type = "%s",
			si.quantity_threshold = %d,
			si.expiration_threshold = %d,
			si.expiration_date = "%s"
		WHERE si.id = %d;
	`, title, image, quantity, quantityType, quantityThreshold, expirationThreshold, expirationDate, itemID))
	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec()
	if err != nil {
		fmt.Println(err)
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

// DecrementStorageItemQuantity decrements a storage item's quantity by ID
func (handler *Handler) DecrementStorageItemQuantity(itemID int) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		UPDATE storage_items AS si
		SET si.quantity = si.quantity - 1
		WHERE si.id = %d AND si.quantity > 0;
	`, itemID))
	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec()
	if err != nil {
		fmt.Println(err)
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

// IncrementStorageItemQuantity increments a storage item's quantity by ID
func (handler *Handler) IncrementStorageItemQuantity(itemID int) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		UPDATE storage_items AS si
		SET 
			si.quantity = si.quantity + 1
		WHERE si.id = %d;
	`, itemID))
	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec()
	if err != nil {
		fmt.Println(err)
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

// DeleteStorageItem deletes a storage item by ID
func (handler *Handler) DeleteStorageItem(itemID int) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
	DELETE FROM storage_items AS si
	WHERE si.id = %d
	`, itemID))
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
