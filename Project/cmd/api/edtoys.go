package main

import (
	"Project/internal/data"
	"fmt"
	"net/http"
)

func (app *application) createEdtoysHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new edtoys")
}

func (app *application) showEdtoysHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	edtoys := data.Edtoys{
		ID:         id,
		Title:      "Educational Game 1",
		Year:       2023,
		TargetAge:  "5-10",
		Genres:     []string{"educational", "puzzle"},
		SkillFocus: []string{"mathematics", "logic"},
		Runtime:    30,
	}
	// Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, envelope{"edtoys": edtoys}, nil)
	if err != nil {
		// Use the new serverErrorResponse() helper.
		app.serverErrorResponse(w, r, err)
	}
}
