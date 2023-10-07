package data

import (
	"Project/internal/validator" // New import
	"time"
)

type Edtoys struct {
	ID         int64    `json:"id"`
	Title      string   `json:"title"`
	Year       int32    `json:"year,omitempty"`
	TargetAge  string   `json:"target_age"`
	Genres     []string `json:"genres,omitempty"`
	SkillFocus []string `json:"skill_focus"`
	Runtime    Runtime  `json:"runtime,omitempty"`
}

func ValidateMovie(v *validator.Validator, edtoys *Edtoys) {
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
