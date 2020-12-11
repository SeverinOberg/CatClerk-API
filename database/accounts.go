package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Account holds the account structure
type Account struct {
	ID            int       `json:"id"`
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	Salt          string    `json:"salt"`
	Email         string    `json:"email"`
	DarkTheme     bool      `json:"datkTheme"`
	Notifications bool      `json:"notifications"`
	LastLogin     time.Time `json:"lastLogin"`
	UpdatedAt     time.Time `json:"updatedAt"`
	CreatedAt     time.Time `json:"createdAt"`
}

// CreateAccount creates a new account in the database
func (handler *Handler) CreateAccount(username, email, password, salt string) (result sql.Result, err error) {
	stmt, err := handler.DB.Prepare(`INSERT INTO accounts SET username=?, email=?, password=?, salt=?`)
	if err != nil {
		return result, err
	}

	result, err = stmt.Exec(username, email, password, salt)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "for key 'username'"):
			return result, fmt.Errorf("username already taken")
		case strings.Contains(err.Error(), "for key 'email'"):
			return result, fmt.Errorf("email already taken")
		default:
			return result, err
		}
	}

	return result, nil
}

// GetAccount gets the account from the database by username
func (handler *Handler) GetAccount(username string) (Account, error) {
	acc := Account{}

	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		SELECT * FROM accounts
		WHERE username="%s"	
	`, username))
	if err != nil {
		return acc, err
	}

	defer stmt.Close()

	if err := stmt.QueryRow().Scan(
		&acc.ID,
		&acc.Username,
		&acc.Password,
		&acc.Salt,
		&acc.Email,
		&acc.DarkTheme,
		&acc.Notifications,
		&acc.LastLogin,
		&acc.UpdatedAt,
		&acc.CreatedAt,
	); err != nil {
		return acc, err
	}

	return acc, err
}

// Email struct
type Email struct {
	Email string `json:"email"`
}

// GetAccountEmail gets the account's email from the database by username
func (handler *Handler) GetAccountEmail(username string) (Email, error) {
	payload := Email{}

	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		SELECT email FROM accounts
		WHERE username="%s"	
	`, username))
	if err != nil {
		return payload, err
	}

	defer stmt.Close()

	if err := stmt.QueryRow().Scan(
		&payload.Email,
	); err != nil {
		return payload, err
	}

	return payload, err
}

// UsernameEmail -
type UsernameEmail struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// EmailExists gets the account's email from the database by username
func (handler *Handler) EmailExists(email string) (UsernameEmail, error) {
	payload := UsernameEmail{}

	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		SELECT username, email FROM accounts
		WHERE email="%s"	
	`, email))
	if err != nil {
		return payload, err
	}

	defer stmt.Close()

	if err := stmt.QueryRow().Scan(
		&payload.Username,
		&payload.Email,
	); err != nil {
		return payload, err
	}

	return payload, err
}

// GetAccounts gets all accounts from the database by username
func (handler *Handler) GetAccounts() ([]Account, error) {
	accounts := []Account{}

	stmt, err := handler.DB.Prepare(`
		SELECT * FROM accounts
	`)
	if err != nil {
		return accounts, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return accounts, err
	}

	for rows.Next() {
		acc := Account{}

		if err := rows.Scan(
			&acc.ID,
			&acc.Username,
			&acc.Password,
			&acc.Salt,
			&acc.Email,
			&acc.DarkTheme,
			&acc.Notifications,
			&acc.LastLogin,
			&acc.UpdatedAt,
			&acc.CreatedAt,
		); err != nil {
			return accounts, err
		}

		accounts = append(accounts, acc)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return accounts, err
}

// CheckAccountCredentials verifies the accounts credentials by username or email
func (handler *Handler) CheckAccountCredentials(username, email string) (Account, error) {
	login := Account{}

	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		SELECT * FROM accounts 
		WHERE (username="%s" OR email="%s")	
	`, username, email))
	if err != nil {
		return login, err
	}

	defer stmt.Close()

	if err := stmt.QueryRow().Scan(
		&login.ID,
		&login.Username,
		&login.Password,
		&login.Salt,
		&login.Email,
		&login.DarkTheme,
		&login.Notifications,
		&login.LastLogin,
		&login.UpdatedAt,
		&login.CreatedAt,
	); err != nil {
		return login, err
	}

	return login, nil
}

// UpdateAccount updates all account data fields in the database by username
func (handler *Handler) UpdateAccount(username, newUsername, password, email string, darkTheme, notifications bool) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		UPDATE accounts AS a
		SET 
			a.username = "%s",
			a.password = "%s",
			a.email = "%s",
			a.dark_theme = %t,
			a.notifications = %t
		WHERE a.username = "%s";
	`, newUsername, password, email, darkTheme, notifications, username))
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

// UpdateAccountUsername updates the accounts username by current username
func (handler *Handler) UpdateAccountUsername(username, newUsername string) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		UPDATE accounts
		SET username = "%s"
		WHERE username = "%s";
	`, newUsername, username))
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

// UpdateAccountEmail updates the accounts email by username
func (handler *Handler) UpdateAccountEmail(username, email string) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		UPDATE accounts
		SET 
			email = "%s"
		WHERE username = "%s";
	`, email, username))
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

// UpdateAccountPassword updates the accounts password by username
func (handler *Handler) UpdateAccountPassword(username, password, salt string) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		UPDATE accounts
		SET 
			password = "%s",
			salt = "%s"
		WHERE username = "%s";
	`, password, salt, username))
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

// DeleteAccount deletes the account by username
func (handler *Handler) DeleteAccount(username string) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
	DELETE FROM accounts AS a
	WHERE a.username = %s
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
