package graph

import (
	"api/graph/generated"
	"api/graph/model"
	"api/graph/models"
	createTodoService "api/graph/services/todo/create"
	"api/graph/test/util"
	"api/graph/test/util/factory"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GraphTestSuite struct {
	util.Suite
	mutationResolver generated.MutationResolver
	queryResolver    generated.QueryResolver
	todoResolver     generated.TodoResolver
	todoLabel        generated.TodoLabelResolver
	resolver         Resolver
}

func (s *GraphTestSuite) SetupTest() {
	tx := s.DB().Begin()

	resolver := NewResolver(
		tx,
	)
	s.resolver = *resolver
	s.mutationResolver = s.resolver.Mutation()
	s.queryResolver = s.resolver.Query()
	s.todoResolver = s.resolver.Todo()
	s.todoLabel = s.resolver.TodoLabel()
}

func TestMain(t *testing.T) {
	suite.Run(t, new(GraphTestSuite))
}

// 各テストメソッドの実行後
func (s *GraphTestSuite) TearDownTest() {
	s.resolver.DB.Rollback()
	s.CloseDB()
}

// スイートの実行前
func (s *GraphTestSuite) SetupSuite() {
	s.SetupDB()
}

func (s *GraphTestSuite) TestGetTodo() {
	db := s.resolver.DB
	for i := 0; i < 5; i++ {
		_ = factory.NewTodo(db)
	}
	s.Run("正常系", func() {
		var sortInput *model.SortTodo
		var searchInput *model.SearchTodo

		result, _ := s.queryResolver.GqlgenTodos(context.Background(), sortInput, searchInput)

		var todos []*models.Todo
		db.Find(&todos)
		assert.Equal(s.T(), result, todos)
	})

	s.Run("並び替え", func() {
		sortInput := &model.SortTodo{
			Column: "id",
			Sort:   model.SortDesc,
		}
		var searchInput *model.SearchTodo
		results, _ := s.queryResolver.GqlgenTodos(context.Background(), sortInput, searchInput)

		var todos []*models.Todo
		db.Order(sortInput.Column + " " + string(sortInput.Sort)).Find(&todos)

		assert.Equal(s.T(), results, todos)

	})

	s.Run("検索", func() {
		var sortInput *model.SortTodo
		searchInput := &model.SearchTodo{
			Column: "title",
			Value:  "time",
		}
		results, _ := s.queryResolver.GqlgenTodos(context.Background(), sortInput, searchInput)

		var todos []*models.Todo
		db.Where(searchInput.Column+" = ?", searchInput.Value).Find(&todos)
		for _, result := range results {
			assert.Equal(s.T(), result.Title, "time")
		}
		assert.Equal(s.T(), results, todos)
	})
}

func (s *GraphTestSuite) TestCreateTodo() {
	db := s.resolver.DB
	title := "testTitle"
	description := "testDescription"
	userID := 1
	priorityID := 1
	finishedTimeString := "2024-01-02 15:04"
	s.Run("正常系", func() {
		object := model.NewTodo{
			Title:       title,
			Description: description,
			UserID:      userID,
			PriorityID:  priorityID,
			FinishedAt:  finishedTimeString,
		}
		result, _ := createTodoService.CreateTodo(db, object)

		var todo models.Todo
		db.Last(&todo)

		assert.Equal(s.T(), result.Title, title)
		assert.Equal(s.T(), result.Description, description)
		assert.Equal(s.T(), result.UserID, userID)
		assert.Equal(s.T(), result.PriorityID, priorityID)
		assert.Equal(s.T(), result.Title, todo.Title)
		assert.Equal(s.T(), result.Description, todo.Description)
		assert.Equal(s.T(), result.UserID, todo.UserID)
		assert.Equal(s.T(), result.PriorityID, todo.PriorityID)
		assert.Equal(s.T(), result.StatusID, todo.StatusID)
	})

	s.Run("異常系", func() {
		s.Run("タイトルが50字以上であればエラーを返す", func() {
			title := "testTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitle"

			object := model.NewTodo{
				Title:       title,
				Description: description,
				UserID:      userID,
				PriorityID:  priorityID,
				FinishedAt:  finishedTimeString,
			}
			_, err := s.mutationResolver.CreateTodo(context.Background(), object)
			assert.Equal(s.T(), err.Error(), "タイトルは50文字以下にしてください")

		})

		s.Run("説明文が300字以上であればエラーを返す", func() {
			description := "testDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescription"

			object := model.NewTodo{
				Title:       title,
				Description: description,
				UserID:      userID,
				PriorityID:  priorityID,
				FinishedAt:  finishedTimeString,
			}
			_, err := s.mutationResolver.CreateTodo(context.Background(), object)
			assert.Equal(s.T(), err.Error(), "説明を300文字以下にしてください")
		})

		s.Run("ラベルの登録が5つ以上の時にエラーが発生する", func() {
			labelIDs := []int{1, 2, 3, 4, 5, 6}

			object := model.NewTodo{
				Title:       title,
				Description: description,
				UserID:      userID,
				LabelIDs:    labelIDs,
				PriorityID:  priorityID,
				FinishedAt:  finishedTimeString,
			}

			_, err := s.mutationResolver.CreateTodo(context.Background(), object)
			assert.Equal(s.T(), err.Error(), "labelは登録できるのは5つまでです。")
		})

		s.Run("終了時間は現在時刻以前にするとエラーが発生する", func() {
			finishedTimeString := "2020-01-02 15:04"

			object := model.NewTodo{
				Title:       title,
				Description: description,
				UserID:      userID,
				PriorityID:  priorityID,
				FinishedAt:  finishedTimeString,
			}
			_, err := s.mutationResolver.CreateTodo(context.Background(), object)
			assert.Equal(s.T(), err.Error(), "終了期限を現在日時以降にしてください")
		})
	})
}

func (s *GraphTestSuite) TestUpdateTodo() {
	db := s.resolver.DB
	todo, _ := createTodoService.CreateTodo(db, model.NewTodo{
		Title:       "testTitle",
		Description: "testDescription",
		UserID:      1,
		PriorityID:  1,
		FinishedAt:  "2024-01-02 15:04",
	})
	s.Run("正常系", func() {
		id := todo.ID
		title := "updateTitle"
		description := "updateDescription"
		finishTimeString := "2025-01-02 15:04"

		object := model.UpdateTodo{
			ID:          id,
			Title:       title,
			Description: description,
			FinishedAt:  finishTimeString,
		}
		result, _ := s.mutationResolver.UpdateTodo(context.Background(), object)

		var todo models.Todo
		db.Find(&todo, id)

		assert.Equal(s.T(), result.ID, id)
		assert.Equal(s.T(), result.Title, title)
		assert.Equal(s.T(), result.Description, description)
		assert.Equal(s.T(), result.ID, todo.ID)
		assert.Equal(s.T(), result.Title, todo.Title)
		assert.Equal(s.T(), result.Description, todo.Description)
	})

	s.Run("異常系", func() {
		s.Run("タイトルが50字以上であればエラーを返す", func() {
			id := todo.ID
			title := "testTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitletestTitle"

			object := model.UpdateTodo{
				ID:    id,
				Title: title,
			}
			_, err := s.mutationResolver.UpdateTodo(context.Background(), object)

			assert.Equal(s.T(), err.Error(), "タイトルは50文字以下にしてください")
		})

		s.Run("説明文が300字以上であればエラーを返す", func() {
			id := todo.ID
			description := "testDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescriptiontestDescription"

			object := model.UpdateTodo{
				ID:          id,
				Description: description,
			}
			_, err := s.mutationResolver.UpdateTodo(context.Background(), object)

			assert.Equal(s.T(), err.Error(), "説明を300文字以下にしてください")
		})

		s.Run("ラベルの登録が5つ以上の時にエラーが発生する", func() {
			todo := factory.NewTodo(db)

			id := todo.ID
			addLabelIDs := []int{4, 5, 6}

			object := model.UpdateTodo{
				ID:          id,
				Title:       "updateTitle",
				Description: "updateDescription",
				AddLabelIDs: addLabelIDs,
				FinishedAt:  "2025-01-02 15:04",
			}
			_, err := s.mutationResolver.UpdateTodo(context.Background(), object)

			assert.Equal(s.T(), err.Error(), "labelは登録できるのは5つまでです。")
		})

		s.Run("終了時間は現在時刻以前にするとエラーが発生する", func() {
			id := todo.ID
			finishTimeString := "2020-01-02 15:04"

			object := model.UpdateTodo{
				ID:         id,
				FinishedAt: finishTimeString,
			}
			_, err := s.mutationResolver.UpdateTodo(context.Background(), object)

			assert.Equal(s.T(), err.Error(), "終了期限を現在日時以降にしてください")
		})
	})
}

func (s *GraphTestSuite) TestDeleteTodo() {
	db := s.resolver.DB
	s.Run("正常系", func() {
		todo := factory.NewTodo(db)
		id := todo.ID

		result, _ := s.mutationResolver.DeleteTodo(context.Background(), id)

		var deleteTodo models.Todo
		err := db.First(&deleteTodo, id).Error

		assert.Equal(s.T(), err.Error(), "record not found")
		assert.Equal(s.T(), result, "削除が完了しました")
	})
}

func (s *GraphTestSuite) TestGetStatus() {
	db := s.resolver.DB
	s.Run("正常系", func() {
		todo := factory.NewTodo(db)

		result, _ := s.todoResolver.Status(context.Background(), todo)

		var status model.Status
		db.First(&status, todo.StatusID)
		assert.Equal(s.T(), result.ID, status.ID)
		assert.Equal(s.T(), result.Name, status.Name)
	})
}

func (s *GraphTestSuite) TestGetPriority() {
	db := s.resolver.DB
	s.Run("正常系", func() {
		todo := factory.NewTodo(db)

		result, _ := s.todoResolver.Priority(context.Background(), todo)

		var priority model.Priority
		db.First(&priority, todo.PriorityID)
		assert.Equal(s.T(), result.ID, priority.ID)
		assert.Equal(s.T(), result.Name, priority.Name)
	})
}

func (s *GraphTestSuite) TestGetTodoLabels() {
	db := s.resolver.DB
	s.Run("正常系", func() {
		todo := factory.NewTodo(db)

		result, _ := s.todoResolver.TodoLabels(context.Background(), todo)

		var todoLabel []*models.TodoLabel
		db.Where("todo_id = ?", todo.ID).Find(&todoLabel)
		assert.Equal(s.T(), result, todoLabel)
	})
}

func (s *GraphTestSuite) TestGetLabel() {
	db := s.resolver.DB
	s.Run("正常系", func() {
		todo := factory.NewTodo(db)

		results, _ := s.todoResolver.TodoLabels(context.Background(), todo)

		for _, result := range results {
			var label model.Label
			db.First(&label, result.LabelID)
			assert.Equal(s.T(), result.LabelID, label.ID)
		}
	})
}
