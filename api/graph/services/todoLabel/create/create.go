package create

import (
	"api/graph/models"

	"gorm.io/gorm"
)

func CreateTodoLabel(db *gorm.DB, LabelIDs []int, todoID int) error {
	var todoLabels []*models.TodoLabel

	for _, labelID := range LabelIDs {
		todoLabel := &models.TodoLabel{
			TodoID:  todoID,
			LabelID: labelID,
		}
		todoLabels = append(todoLabels, todoLabel)
	}

	err := db.Create(&todoLabels).Error
	if err != nil {
		return err
	}
	return nil
}
