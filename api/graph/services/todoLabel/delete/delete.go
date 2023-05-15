package delete

import (
	"api/graph/models"

	"gorm.io/gorm"
)

func DeleteTodoLabel(db *gorm.DB, LabelIDs []int, todoID int) error {
	var todoLabels []*models.TodoLabel

	err := db.
		Model(&[]*models.TodoLabel{}).
		Where("todo_id", todoID).
		Where("label_id", LabelIDs).
		Delete(&todoLabels).Error
	if err != nil {
		return err
	}
	return nil
}
