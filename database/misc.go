package database

// Foods structure
type Foods struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetFoods gets all food varieties from the foods database
func (handler *Handler) GetFoods() ([]Foods, error) {
	foods := []Foods{}

	stmt, err := handler.DB.Prepare("SELECT * FROM foods")
	if err != nil {
		return foods, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return foods, err
	}

	for rows.Next() {
		food := Foods{}

		if err := rows.Scan(
			&food.ID,
			&food.Name,
		); err != nil {
			return foods, err
		}

		foods = append(foods, food)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return foods, err
}
