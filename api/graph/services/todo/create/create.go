package create

import (
	"api/graph/model"
	"api/graph/models"
	"api/graph/services/common"
	"time"

	"gorm.io/gorm"
)

func CreateTodo(db *gorm.DB, input model.NewTodo) (*models.Todo, error) {
	const defaultStatus = 1

	// フォーマット　"2022-6-28 13:00"
	finishTime, err := common.ChangeStringToTime(input.FinishedAt)
	if err != nil {
		return nil, err
	}

	todo := models.Todo{
		Title:       input.Title,
		Description: input.Description,
		UserID:      input.UserID,
		StatusID:    defaultStatus,
		PriorityID:  input.PriorityID,
		FinishedAt:  finishTime,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = db.Create(&todo).Error
	if err != nil {
		return nil, err
	}

	return &todo, nil
}
