package database

import (
	"fmt"
	"time"
)

// ShoppingListItem structure
type ShoppingListItem struct {
	ID             int       `json:"id"`
	ShoppingListID int       `json:"shoppingListID"`
	Title          string    `json:"title"`
	Quantity       int       `json:"quantity"`
	QuantityType   string    `json:"quantityType"`
	UpdatedAt      time.Time `json:"updatedAt"`
	CreatedAt      time.Time `json:"createdAt"`
}

// CreateShoppingListItem creates a shopping list item in the database attatched by FK to a shopping list ID
func (handler *Handler) CreateShoppingListItem(shoppingListID int, title string, quantity int, quantityType string) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		INSERT INTO shopping_list_items(shopping_list_id, title, quantity, quantity_type)
		VALUES(%d, "%s", %d, "%s");
	`, shoppingListID, title, quantity, quantityType))
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

// GetShoppingListItems gets all shopping list items by a shopping list ID
func (handler *Handler) GetShoppingListItems(shoppingListID int) ([]ShoppingListItem, error) {
	items := []ShoppingListItem{}

	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
	SELECT *
	FROM shopping_list_items AS sli
	WHERE sli.shopping_list_id = %d
	`, shoppingListID))
	if err != nil {
		return items, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return items, err
	}

	for rows.Next() {
		item := ShoppingListItem{}

		if err := rows.Scan(
			&item.ID,
			&item.ShoppingListID,
			&item.Title,
			&item.Quantity,
			&item.QuantityType,
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

// GetShoppingListItem gets a single shopping list item by ID
func (handler *Handler) GetShoppingListItem(itemID int) (ShoppingListItem, error) {
	item := ShoppingListItem{}

	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
	SELECT *
	FROM shopping_list_items AS sli
	WHERE sli.id = %d
	`, itemID))
	if err != nil {
		return item, err
	}

	defer stmt.Close()

	if err := stmt.QueryRow().Scan(
		&item.ID,
		&item.ShoppingListID,
		&item.Title,
		&item.Quantity,
		&item.QuantityType,
		&item.UpdatedAt,
		&item.CreatedAt,
	); err != nil {
		return item, err
	}

	return item, err
}

// GetShoppingListItemsCount gets the amount of shopping list items by username
func (handler *Handler) GetShoppingListItemsCount(username string) (int, error) {
	count := 0

	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		SELECT COUNT(*)
		FROM shopping_list_items AS sli
		INNER JOIN shopping_lists AS sl
		ON sli.shopping_list_id = sl.id
		INNER JOIN account_shopping_list_binder AS aslb
		WHERE sl.id = aslb.shopping_list_id AND aslb.username="%s";
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

// UpdateShoppingListItem updates a shopping list item by ID
func (handler *Handler) UpdateShoppingListItem(title string, quantity int, quantityType string, itemID int) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		UPDATE shopping_list_items AS sli
		SET 
			sli.title = "%s",
			sli.quantity = %d,
			sli.quantity_type = "%s"
		WHERE sli.id = %d;
	`, title, quantity, quantityType, itemID))
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

// UpdateShoppingListItemTitle updates a shopping list item's title by ID
func (handler *Handler) UpdateShoppingListItemTitle(title string, itemID int) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		UPDATE shopping_list_items AS sli
		SET 
			sli.title = "%s"
		WHERE sli.id = %d;
	`, title, itemID))
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

// DecrementShoppingListItemQuantity decrements a shopping list item by username
func (handler *Handler) DecrementShoppingListItemQuantity(itemID int) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		UPDATE shopping_list_items AS sli
		SET sli.quantity = sli.quantity - 1
		WHERE sli.id = %d AND sli.quantity > 0;
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

// IncrementShoppingListItemQuantity increments a shopping list item by username
func (handler *Handler) IncrementShoppingListItemQuantity(itemID int) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		UPDATE shopping_list_items AS sli
		SET 
			sli.quantity = sli.quantity + 1
		WHERE sli.id = %d;
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

// DeleteShoppingListItem deletes a shopping list item by ID
func (handler *Handler) DeleteShoppingListItem(itemID int) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
	DELETE FROM shopping_list_items AS sli
	WHERE sli.id = %d
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
