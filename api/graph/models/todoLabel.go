package models

type TodoLabel struct {
	ID      int `json:"id"`
	TodoID  int `json:"todoID"`
	LabelID int `json:"labelID"`
}
