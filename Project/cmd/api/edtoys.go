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
		Title      string       `json:"title"`
		Year       int32        `json:"year"`
		TargetAge  string       `json:"target_age"`
		Genres     []string     `json:"genres"`
		SkillFocus []string     `json:"skill_focus"`
		Runtime    data.Runtime `json:"runtime"`
	}
	// Read the JSON request body data into the input struct.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Copy the values from the request body to the appropriate fields of the movie
	// record.
	edToys.Title = input.Title
	edToys.Year = input.Year
	edToys.Runtime = input.Runtime
	edToys.Genres = input.Genres
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
		app.serverErrorResponse(w, r, err)
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
