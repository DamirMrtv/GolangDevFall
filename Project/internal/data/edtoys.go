package data

type Edtoys struct {
	ID         int64    `json:"id"`
	Title      string   `json:"title"`
	Year       int32    `json:"year,omitempty"`
	TargetAge  string   `json:"target_age"`
	Genres     []string `json:"genres,omitempty"`
	SkillFocus []string `json:"skill_focus"`
	Runtime    Runtime  `json:"runtime,omitempty"`
}
