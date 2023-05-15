package models

import "time"

type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      int       `json:"userID"`
	StatusID    int       `json:"statusID"`
	PriorityID  int       `json:"priorityID"`
	FinishedAt  time.Time `json:"finishedAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
