package common

import (
	"api/graph/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type ValidateTodoType struct {
	Title       string
	Description string
	LabelIDs    []int
	FinishTime  string
	LabelCount  int
}

func ValidateTodo(obj ValidateTodoType) error {
	if len(obj.Title) > 50 {
		return errors.New("タイトルは50文字以下にしてください")
	}

	if len(obj.Description) > 300 {
		return errors.New("説明を300文字以下にしてください")
	}

	now := time.Now()
	finishTime, err := ChangeStringToTime(obj.FinishTime)
	if err != nil {
		return err
	}
	diff := now.Sub(finishTime)

	if diff > 0 {
		return errors.New("終了期限を現在日時以降にしてください")
	}

	if len(obj.LabelIDs) > 5 || (5-obj.LabelCount)-len(obj.LabelIDs) < 0 {
		return errors.New("labelは登録できるのは5つまでです。")
	}
	return nil
}

func ChangeStringToTime(stringFinishTime string) (time.Time, error) {
	layout := "2006-01-02 15:04"
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	finishTimeUTC, err := time.Parse(layout, stringFinishTime)
	if err != nil {
		return time.Time{}, err
	}
	finishTime := finishTimeUTC.In(jst)

	return finishTime, nil
}

func CountTodoLabel(db *gorm.DB, id int) (int64, error) {
	var labelCount int64

	db.
		Model(&[]*models.TodoLabel{}).
		Where("todo_id = ?", id).
		Count(&labelCount)
	return labelCount, nil
}
