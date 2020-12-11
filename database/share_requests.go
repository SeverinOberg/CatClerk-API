package database

import (
	"fmt"
	"time"
)

// ShareRequest structure
type ShareRequest struct {
	ID           int       `json:"id"`
	FromUsername string    `json:"fromUsername"`
	ToUsername   string    `json:"toUsername"`
	ShareType    string    `json:"shareType"`
	Title        string    `json:"title"`
	IDRequest    int       `json:"idRequest"`
	CreatedAt    time.Time `json:"createdAt"`
}

// CreateShareRequest creates a share request in the database
func (handler *Handler) CreateShareRequest(fromUsername, toUsername, shareType, title string, idRequest int) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		INSERT INTO share_requests(from_username, to_username, share_type, title, id_request)
		VALUES("%s", "%s", "%s", "%s", %d)
	`, fromUsername, toUsername, shareType, title, idRequest))
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

// GetShareRequests gets all share requests from the database by username
func (handler *Handler) GetShareRequests(username string) ([]ShareRequest, error) {
	shareRequests := []ShareRequest{}

	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
		SELECT *
		FROM share_requests AS sr
		WHERE sr.to_username = "%s"	
	`, username))
	if err != nil {
		return shareRequests, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return shareRequests, err
	}

	for rows.Next() {
		shareRequest := ShareRequest{}

		if err := rows.Scan(
			&shareRequest.ID,
			&shareRequest.FromUsername,
			&shareRequest.ToUsername,
			&shareRequest.ShareType,
			&shareRequest.Title,
			&shareRequest.IDRequest,
			&shareRequest.CreatedAt,
		); err != nil {
			return shareRequests, err
		}

		shareRequests = append(shareRequests, shareRequest)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return shareRequests, err
}

// DeleteShareRequest deletes a share request by ID
func (handler *Handler) DeleteShareRequest(shareID int) error {
	stmt, err := handler.DB.Prepare(fmt.Sprintf(`
	DELETE FROM share_requests AS sr
	WHERE sr.id = %d
	`, shareID))
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
