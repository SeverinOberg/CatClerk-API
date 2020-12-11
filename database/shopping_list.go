package database

import (
	"fmt"
	"time"
)

// ShoppingList structure
type ShoppingList struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
	Count     int       `json:"count"`
}

// CreateShoppingList creates a shopping list and attaches the account to it by username
func (handler *Handler) CreateShoppingList(username, title string, owner bool) error {
	stmtSL, err := handler.DB.Prepare(fmt.Sprintf(`
		INSERT INTO shopping_lists(title)
		VALUES("%s");
	`, title))
	if err != nil {
		return err
	}

	defer stmtSL.Close()

	result, err := stmtSL.Exec()
	if err != nil {
		return err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	stmtASLB, err := handler.DB.Prepare(fmt.Sprintf(`
	INSERT INTO account_shopping_list_binder(username, shopping_list_id, owner)
	VALUES("%s", %d, %t);
	`, username, lastInsertID, owner))
	if err != nil {
		return err
	}

	defer stmtASLB.Close()

	_, err = stmtASLB.Exec()
	if err != nil {
		return err
	}

	return err
}

// GetShoppingLists gets shopping lists by username
func (handler *Handler) GetShoppingLists(username string) ([]ShoppingList, error) {
	shoppingLists := []ShoppingList{}

	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		SELECT sl.*, COUNT(sli.id) AS "count"
		FROM shopping_lists AS sl
		LEFT JOIN shopping_list_items AS sli
		ON sl.id = sli.shopping_list_id
		INNER JOIN account_shopping_list_binder AS aslb
		ON aslb.shopping_list_id = sl.id AND aslb.username = "%s"
		GROUP BY sl.id
	`, username))
	if err != nil {
		return shoppingLists, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return shoppingLists, err
	}

	for rows.Next() {
		sl := ShoppingList{}

		if err := rows.Scan(
			&sl.ID,
			&sl.Title,
			&sl.UpdatedAt,
			&sl.CreatedAt,
			&sl.Count,
		); err != nil {
			return shoppingLists, err
		}

		shoppingLists = append(shoppingLists, sl)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return shoppingLists, err
}

// GetShoppingListsCount gets the amount of shopping lists by username
func (handler *Handler) GetShoppingListsCount(username string) (int, error) {
	count := 0

	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		SELECT COUNT(*)
		FROM account_shopping_list_binder
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

// UpdateShoppingListTitle updates a shopping list's title by ID
func (handler *Handler) UpdateShoppingListTitle(title string, shoppingListID int) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		UPDATE shopping_lists AS sl
		SET sl.title = "%s"
		WHERE sl.id = %d
	`, title, shoppingListID))
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

// DeleteShoppingList deletes a shopping list by ID
func (handler *Handler) DeleteShoppingList(shoppingListID int) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
	DELETE FROM shopping_lists AS sl
	WHERE sl.id = %d
	`, shoppingListID))
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

// ShareShoppingList attaches a shopping list to an account by username and ID
func (handler *Handler) ShareShoppingList(username string, shoppingListID int) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		INSERT INTO account_shopping_list_binder(username, shopping_list_id)
		VALUES("%s", %d)
	`, username, shoppingListID))
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

// RemoveShareShoppingList removes an accounts attachment to a shopping list by username and ID
func (handler *Handler) RemoveShareShoppingList(username string, shoppingListID int) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		DELETE FROM account_shopping_list_binder
		WHERE username="%s" AND shopping_list_id=%d
	`, username, shoppingListID))
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

// GetShoppingListOwner returns true or false whether it is the shopping list owner by ID
func (handler *Handler) GetShoppingListOwner(owner string, shoppingListID int) (bool, error) {
	payload := false

	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
	SELECT owner
	FROM account_shopping_list_binder
	WHERE username="%s" AND shopping_list_id=%d 
	`, owner, shoppingListID))
	if err != nil {
		return payload, err
	}

	defer stmt.Close()

	if err := stmt.QueryRow().Scan(
		&payload,
	); err != nil {
		return payload, err
	}

	return payload, err
}
