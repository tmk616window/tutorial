package factory

import (
	"api/graph/model"
	"api/graph/models"
	createTodoService "api/graph/services/todo/create"
	createTodoLabelService "api/graph/services/todoLabel/create"

	"gorm.io/gorm"
)

func NewTodo(db *gorm.DB) *models.Todo {
	title := "testTitle"
	description := "testDescription"
	userID := 1
	priorityID := 1
	finishedTimeString := "2024-01-02 15:04"
	labelIDs := []int{1, 2, 3, 4, 5, 6}

	object := model.NewTodo{
		Title:       title,
		Description: description,
		UserID:      userID,
		PriorityID:  priorityID,
		FinishedAt:  finishedTimeString,
	}

	todo, _ := createTodoService.CreateTodo(db, object)
	_ = createTodoLabelService.CreateTodoLabel(db, labelIDs, todo.ID)

	return todo
}
