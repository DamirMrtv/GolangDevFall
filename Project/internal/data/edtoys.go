package data

import (
	"Project/internal/validator" // New import
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&edtoys.ID, &edtoys.CreatedAt, &edtoys.Version)

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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Execute the query using the QueryRow() method, passing in the provided id value
	// as a placeholder parameter, and scan the response data into the fields of the
	// Movie struct. Importantly, notice that we need to convert the scan target for the
	// genres column using the pq.Array() adapter function again.
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
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
WHERE id = $7 AND version = $8
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
		edtoys.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use the QueryRow() method to execute the query, passing in the args slice as a
	// variadic parameter and scanning the new version value into the movie struct.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&edtoys.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.
	result, err := m.DB.ExecContext(ctx, query, id)
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

func (m EdtoysModel) GetAll(title string, genres []string, filters Filters) ([]*Edtoys, Metadata, error) {
	// Construct the SQL query to retrieve all movie records.
	query := fmt.Sprintf(`
		SELECT  count(*) OVER(), id, created_at, title, year, target_age, genres, skill_focus, runtime, version
		FROM edtoys
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (genres @> $2 or $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.
	args := []interface{}{title, pq.Array(genres), filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
	// before GetAll() returns.
	defer rows.Close()
	// Initialize an empty slice to hold the movie data.
	totalRecords := 0
	edToys := []*Edtoys{}
	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var edtoys Edtoys
		// Scan the values from the row into the Movie struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&totalRecords,
			&edtoys.ID,
			&edtoys.CreatedAt,
			&edtoys.Title,
			&edtoys.Year,
			&edtoys.TargetAge,
			pq.Array(&edtoys.Genres),
			pq.Array(&edtoys.SkillFocus),
			&edtoys.Runtime,
			&edtoys.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		// Add the Movie struct to the slice.
		edToys = append(edToys, &edtoys)
	}
	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// If everything went OK, then return the slice of movies.
	return edToys, metadata, nil
}
