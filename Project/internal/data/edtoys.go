package data

import (
	"Project/internal/validator" // New import
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"time"
)

type Edtoys struct {
	ID         int64     `json:"id"`
	CreatedAt  time.Time `json:"-"`
	Title      string    `json:"title"`
	Year       int32     `json:"year,omitempty"`
	TargetAge  string    `json:"target_age"`
	Genres     []string  `json:"genres,omitempty"`
	SkillFocus []string  `json:"skill_focus"`
	Runtime    Runtime   `json:"runtime,omitempty"`
	Version    int32     `json:"version"`
}

func ValidateEdtoys(v *validator.Validator, edtoys *Edtoys) {
	v.Check(edtoys.Title != "", "title", "must be provided")
	v.Check(len(edtoys.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(edtoys.Year != 0, "year", "must be provided")
	v.Check(edtoys.Year >= 1888, "year", "must be greater than 1888")
	v.Check(edtoys.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(edtoys.Runtime != 0, "runtime", "must be provided")
	v.Check(edtoys.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(edtoys.Genres != nil, "genres", "must be provided")
	v.Check(len(edtoys.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(edtoys.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(edtoys.Genres), "genres", "must not contain duplicate values")
}

type EdtoysModel struct {
	DB *sql.DB
}

func (m EdtoysModel) Insert(edtoys *Edtoys) error {

	query := `
		INSERT INTO edToys (title, year, target_age, genres, skill_focus, runtime)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, version
		`

	args := []interface{}{edtoys.Title, edtoys.Year, edtoys.TargetAge, pq.Array(edtoys.Genres), pq.Array(edtoys.SkillFocus), edtoys.Runtime}

	return m.DB.QueryRow(query, args...).Scan(&edtoys.ID, &edtoys.CreatedAt, &edtoys.Version)

}

// Add a placeholder method for fetching a specific record from the Edtoyss table.
func (m EdtoysModel) Get(id int64) (*Edtoys, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
SELECT id, created_at, title, year, target_age,genres, skill_focus,runtime, version
FROM edtoys
WHERE id = $1`
	// Declare a Movie struct to hold the data returned by the query.
	var edToy Edtoys
	// Execute the query using the QueryRow() method, passing in the provided id value
	// as a placeholder parameter, and scan the response data into the fields of the
	// Movie struct. Importantly, notice that we need to convert the scan target for the
	// genres column using the pq.Array() adapter function again.
	err := m.DB.QueryRow(query, id).Scan(
		&edToy.ID,
		&edToy.CreatedAt,
		&edToy.Title,
		&edToy.Year,
		&edToy.TargetAge,
		pq.Array(&edToy.Genres),
		pq.Array(&edToy.SkillFocus),
		&edToy.Runtime,
		&edToy.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &edToy, nil
}

// Add a placeholder method for updating a specific record in the Edtoyss table.
func (m EdtoysModel) Update(edtoys *Edtoys) error {
	query := `
UPDATE edtoys
SET title = $1, year = $2, target_age = $3, genres = $4, skill_focus = $5, runtime = $6, version = version + 1
WHERE id = $7
RETURNING version`
	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{
		edtoys.Title,
		edtoys.Year,
		edtoys.TargetAge,
		pq.Array(edtoys.Genres),
		pq.Array(edtoys.SkillFocus),
		edtoys.Runtime,
		edtoys.ID,
	}
	// Use the QueryRow() method to execute the query, passing in the args slice as a
	// variadic parameter and scanning the new version value into the movie struct.
	return m.DB.QueryRow(query, args...).Scan(&edtoys.Version)
}

// Add a placeholder method for deleting a specific record from the Edtoyss table.
func (m EdtoysModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
DELETE FROM edtoys
WHERE id = $1`
	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil

}
