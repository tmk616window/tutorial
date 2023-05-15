package update

import (
	"api/graph/model"
	"api/graph/models"
	"api/graph/services/common"

	"gorm.io/gorm"
)

func UpdateTodo(db *gorm.DB, input model.UpdateTodo) (*models.Todo, error) {

	// フォーマット　"2022/6/28 13:00"
	finishTime, err := common.ChangeStringToTime(input.FinishedAt)
	if err != nil {
		return nil, err
	}

	var todo models.Todo
	err = db.Model(&todo).First(&todo, input.ID).Updates(models.Todo{
		Title:       input.Title,
		Description: input.Description,
		FinishedAt:  finishTime,
	}).Error
	if err != nil {
		return nil, err
	}

	return &todo, nil
}
