package main

import (
	"Project/internal/data"
	"Project/internal/validator" // New import
	"errors"
	"fmt"
	"net/http"
)

func (app *application) createEdtoysHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title      string       `json:"title"`
		Year       int32        `json:"year"`
		TargetAge  string       `json:"target_age"`
		Genres     []string     `json:"genres"`
		SkillFocus []string     `json:"skill_focus"`
		Runtime    data.Runtime `json:"runtime"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	edtoys := &data.Edtoys{
		Title:      input.Title,
		Year:       input.Year,
		TargetAge:  input.TargetAge,
		Genres:     input.Genres,
		SkillFocus: input.SkillFocus,
		Runtime:    input.Runtime,
	}
	// Initialize a new Validator.
	v := validator.New()
	// Call the ValidateMovie() function and return a response containing the errors if
	// any of the checks fail.
	if data.ValidateEdtoys(v, edtoys); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.EdToys.Insert(edtoys)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/edToys/%d", edtoys.ID))
	err = app.writeJSON(w, http.StatusCreated, envelope{"educational_toys": edtoys}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showEdtoysHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	edToy, err := app.models.EdToys.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"educational_toys": edToy}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) updateEdToysHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Fetch the existing movie record from the database, sending a 404 Not Found
	// response to the client if we couldn't find a matching record.
	edToys, err := app.models.EdToys.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		Title      *string       `json:"title"`
		Year       *int32        `json:"year"`
		TargetAge  *string       `json:"target_age"`
		Genres     []string      `json:"genres"`
		SkillFocus []string      `json:"skill_focus"`
		Runtime    *data.Runtime `json:"runtime"`
	}
	// Read the JSON request body data into the input struct.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Copy the values from the request body to the appropriate fields of the movie
	// record.
	if input.Title != nil {
		edToys.Title = *input.Title
	}
	if input.Year != nil {
		edToys.Year = *input.Year
	}
	if input.TargetAge != nil {
		edToys.TargetAge = *input.TargetAge
	}
	if input.Genres != nil {
		edToys.Genres = input.Genres // Note that we don't need to dereference a slice.
	}
	if input.SkillFocus != nil {
		edToys.SkillFocus = input.SkillFocus
	}
	if input.Runtime != nil {
		edToys.Runtime = *input.Runtime
	}
	// Validate the updated movie record, sending the client a 422 Unprocessable Entity
	// response if any checks fail.
	v := validator.New()
	if data.ValidateEdtoys(v, edToys); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Pass the updated movie record to our new Update() method.
	err = app.models.EdToys.Update(edToys)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Write the updated movie record in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"educational_toys": edToys}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteEdToysHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the movie ID from the URL.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Delete the movie from the database, sending a 404 Not Found response to the
	// client if there isn't a matching record.
	err = app.models.EdToys.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return a 200 OK status code along with a success message.
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "educational toy successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listEdToysHandler(w http.ResponseWriter, r *http.Request) {
	// To keep things consistent with our other handlers, we'll define an input struct
	// to hold the expected values from the request query string.
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}
	// Initialize a new Validator instance.
	v := validator.New()
	// Call r.URL.Query() to get the url.Values map containing the query string data.
	qs := r.URL.Query()
	// Use our helpers to extract the title and genres query string values, falling back
	// to defaults of an empty string and an empty slice respectively if they are not
	// provided by the client.
	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})
	// Get the page and page_size query string values as integers. Notice that we set
	// the default page value to 1 and default page_size to 20, and that we pass the
	// validator instance as the final argument here.
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	// Extract the sort query string value, falling back to "id" if it is not provided
	// by the client (which will imply a ascending sort on movie ID).
	input.Filters.Sort = app.readString(qs, "sort", "id")
	// Check the Validator instance for any errors and use the failedValidationResponse()
	// helper to send the client a response if necessary.
	input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	edToys, metadata, err := app.models.EdToys.GetAll(input.Title, input.Genres, input.Filters)
	// Dump the contents of the input struct in a HTTP response.
	err = app.writeJSON(w, http.StatusOK, envelope{"educational_toys": edToys, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
