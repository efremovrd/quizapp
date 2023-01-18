package repo_test

import (
	"context"
	"errors"
	"quiz-app/internal/question/repo"
	"quiz-app/models"
	"quiz-app/pkg/errs"
	"quiz-app/pkg/postgres"
	"quiz-app/pkg/types"
	"strconv"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/pashagolub/pgxmock"

	"github.com/stretchr/testify/assert"
)

var (
	_builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

func TestFormRepo_Create(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewQuestionRepo(&db)

	type mockBehavior func(ctx context.Context, question *models.Question)

	testTable := []struct {
		nameTest         string
		ctx              context.Context
		question         models.Question
		mockBehavior     mockBehavior
		expectedQuestion models.Question
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			question: models.Question{
				Form_id: "12",
				Header:  "sdcsd",
			},
			mockBehavior: func(ctx context.Context, question *models.Question) {
				pgxRows := pgxpoolmock.NewRows([]string{"id_"}).AddRow(345).ToPgxRows()
				pgxRows.Next()
				formidint, _ := strconv.Atoi(question.Form_id)
				mockPool.EXPECT().QueryRow(ctx, "INSERT INTO question_ (form_id_, header_) VALUES ($1,$2) RETURNING \"id_\"", formidint, question.Header).Return(pgxRows)
			},
			expectedQuestion: models.Question{
				Id:      "345",
				Form_id: "12",
				Header:  "sdcsd",
			},
		},
		{
			nameTest: "invalid_inputs",
			ctx:      context.Background(),
			question: models.Question{
				Form_id: "5r4",
				Header:  "sdcsd",
			},
			mockBehavior: func(ctx context.Context, question *models.Question) {},
		},
		{
			nameTest: "no_rows",
			ctx:      context.Background(),
			question: models.Question{
				Form_id: "12",
				Header:  "sdcsd",
			},
			mockBehavior: func(ctx context.Context, question *models.Question) {
				pgxRows := pgxpoolmock.NewRows([]string{}).AddRow().ToPgxRows()
				formidint, _ := strconv.Atoi(question.Form_id)
				mockPool.EXPECT().QueryRow(ctx, "INSERT INTO question_ (form_id_, header_) VALUES ($1,$2) RETURNING \"id_\"", formidint, question.Header).Return(pgxRows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, &testCase.question)

			got, err := r.Create(testCase.ctx, &testCase.question)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedQuestion, *got)
			case "invalid_inputs":
				assert.Equal(t, errs.ErrInvalidContent, err)
			case "no_rows":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestQuestionRepo_GetById(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewQuestionRepo(&db)

	type mockBehavior func(ctx context.Context, id string)

	testTable := []struct {
		nameTest         string
		ctx              context.Context
		id               string
		mockBehavior     mockBehavior
		expectedQuestion models.Question
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			id:       "345",
			mockBehavior: func(ctx context.Context, id string) {
				pgxRows := pgxpoolmock.NewRows([]string{"form_id_", "header_"}).AddRow(12, "sdcsd").ToPgxRows()
				pgxRows.Next()
				idint, _ := strconv.Atoi(id)
				mockPool.EXPECT().QueryRow(ctx, "SELECT form_id_, header_ FROM question_ WHERE id_ = $1", idint).Return(pgxRows)
			},
			expectedQuestion: models.Question{
				Id:      "345",
				Form_id: "12",
				Header:  "sdcsd",
			},
		},
		{
			nameTest:     "invalid_inputs",
			ctx:          context.Background(),
			id:           "5r4",
			mockBehavior: func(ctx context.Context, id string) {},
		},
		{
			nameTest: "no_rows",
			ctx:      context.Background(),
			id:       "345",
			mockBehavior: func(ctx context.Context, id string) {
				pgxRows := pgxpoolmock.NewRows([]string{}).AddRow().ToPgxRows()
				idint, _ := strconv.Atoi(id)
				mockPool.EXPECT().QueryRow(ctx, "SELECT form_id_, header_ FROM question_ WHERE id_ = $1", idint).Return(pgxRows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, testCase.id)

			got, err := r.GetById(testCase.ctx, testCase.id)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedQuestion, *got)
			case "invalid_inputs":
				assert.Equal(t, errs.ErrInvalidContent, err)
			case "no_rows":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestQuestionRepo_GetByFormId(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewQuestionRepo(&db)

	type mockBehavior func(ctx context.Context, form_id string, sets types.GetSets)

	testTable := []struct {
		nameTest          string
		ctx               context.Context
		form_id           string
		sets              types.GetSets
		mockBehavior      mockBehavior
		expectedQuestions []*models.Question
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			form_id:  "12",
			sets:     types.GetSets{},
			mockBehavior: func(ctx context.Context, form_id string, sets types.GetSets) {
				pgxRows := pgxpoolmock.NewRows([]string{"id_", "header_"}).AddRow(345, "ecefvc").AddRow(346, "ty").ToPgxRows()
				formidint, _ := strconv.Atoi(form_id)
				mockPool.EXPECT().Query(ctx, "SELECT id_, header_ FROM question_ WHERE form_id_ = $1 LIMIT 0 OFFSET 0", formidint).Return(pgxRows, nil)
			},
			expectedQuestions: []*models.Question{
				{
					Id:      "345",
					Form_id: "12",
					Header:  "ecefvc",
				},
				{
					Id:      "346",
					Form_id: "12",
					Header:  "ty",
				},
			},
		},
		{
			nameTest:     "invalid_inputs",
			ctx:          context.Background(),
			form_id:      "5r4",
			sets:         types.GetSets{},
			mockBehavior: func(ctx context.Context, form_id string, sets types.GetSets) {},
		},
		{
			nameTest: "query_error",
			ctx:      context.Background(),
			form_id:  "12",
			sets:     types.GetSets{},
			mockBehavior: func(ctx context.Context, form_id string, sets types.GetSets) {
				formidint, _ := strconv.Atoi(form_id)
				mockPool.EXPECT().Query(ctx, "SELECT id_, header_ FROM question_ WHERE form_id_ = $1 LIMIT 0 OFFSET 0", formidint).Return(nil, errors.New("query_error"))
			},
		},
		{
			nameTest: "no_rows",
			ctx:      context.Background(),
			form_id:  "12",
			sets:     types.GetSets{},
			mockBehavior: func(ctx context.Context, form_id string, sets types.GetSets) {
				pgxRows := pgxpoolmock.NewRows([]string{}).AddRow().ToPgxRows()
				pgxRows.Next()
				formidint, _ := strconv.Atoi(form_id)
				mockPool.EXPECT().Query(ctx, "SELECT id_, header_ FROM question_ WHERE form_id_ = $1 LIMIT 0 OFFSET 0", formidint).Return(pgxRows, nil)
			},
			expectedQuestions: []*models.Question{},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, testCase.form_id, types.GetSets{})

			got, err := r.GetByFormId(testCase.ctx, testCase.form_id, types.GetSets{})

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedQuestions, got)
			case "invalid_inputs":
				assert.Equal(t, errs.ErrInvalidContent, err)
			case "query_error":
				assert.NotEqual(t, nil, err)
			case "no_rows":
				assert.Equal(t, nil, err)
				assert.Equal(t, []*models.Question{}, got)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestQuestionRepo_Update(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewQuestionRepo(&db)

	type mockBehavior func(ctx context.Context, question *models.Question)

	testTable := []struct {
		nameTest         string
		ctx              context.Context
		question         models.Question
		mockBehavior     mockBehavior
		expectedQuestion models.Question
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			question: models.Question{
				Id:      "345",
				Form_id: "12",
				Header:  "sdcsd",
			},
			mockBehavior: func(ctx context.Context, question *models.Question) {
				idint, _ := strconv.Atoi(question.Id)
				mockPool.EXPECT().Exec(ctx, "UPDATE question_ SET header_ = $1 WHERE id_ = $2", question.Header, idint).Return(pgxmock.NewResult("UPDATE", 1), nil)
			},
			expectedQuestion: models.Question{
				Id:      "345",
				Form_id: "12",
				Header:  "sdcsd",
			},
		},
		{
			nameTest: "invalid_inputs_form_id",
			ctx:      context.Background(),
			question: models.Question{
				Id:      "345",
				Form_id: "5r4",
				Header:  "sdcsd",
			},
			mockBehavior: func(ctx context.Context, question *models.Question) {},
		},
		{
			nameTest: "invalid_inputs_id",
			ctx:      context.Background(),
			question: models.Question{
				Id:      "5r4",
				Form_id: "45",
				Header:  "sdcsd",
			},
			mockBehavior: func(ctx context.Context, question *models.Question) {},
		},
		{
			nameTest: "no_question_to_update",
			ctx:      context.Background(),
			question: models.Question{
				Id:      "345",
				Form_id: "12",
				Header:  "sdcsd",
			},
			mockBehavior: func(ctx context.Context, question *models.Question) {
				idint, _ := strconv.Atoi(question.Id)
				mockPool.EXPECT().Exec(ctx, "UPDATE question_ SET header_ = $1 WHERE id_ = $2", question.Header, idint).Return(pgxmock.NewResult("UPDATE", 0), nil)
			},
		},
		{
			nameTest: "exec_error",
			ctx:      context.Background(),
			question: models.Question{
				Id:      "345",
				Form_id: "12",
				Header:  "sdcsd",
			},
			mockBehavior: func(ctx context.Context, question *models.Question) {
				idint, _ := strconv.Atoi(question.Id)
				mockPool.EXPECT().Exec(ctx, "UPDATE question_ SET header_ = $1 WHERE id_ = $2", question.Header, idint).Return(nil, errors.New("exec_error"))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, &testCase.question)

			got, err := r.Update(testCase.ctx, &testCase.question)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedQuestion, *got)
			case "invalid_inputs_id", "invalid_inputs_user_id":
				assert.Equal(t, errs.ErrInvalidContent, err)
			case "no_question_to_update":
				assert.Equal(t, errs.ErrContentNotFound, err)
			case "exec_error":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestQuestionRepo_Delete(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewQuestionRepo(&db)

	type mockBehavior func(ctx context.Context, id string)

	testTable := []struct {
		nameTest     string
		ctx          context.Context
		id           string
		mockBehavior mockBehavior
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			id:       "345",
			mockBehavior: func(ctx context.Context, id string) {
				idint, _ := strconv.Atoi(id)
				mockPool.EXPECT().Exec(ctx, "DELETE FROM question_ WHERE id_ = $1", idint).Return(pgxmock.NewResult("DELETE", 1), nil)
			},
		},
		{
			nameTest:     "invalid_inputs",
			ctx:          context.Background(),
			id:           "5r4",
			mockBehavior: func(ctx context.Context, id string) {},
		},
		{
			nameTest: "no_question_to_delete",
			ctx:      context.Background(),
			id:       "345",
			mockBehavior: func(ctx context.Context, id string) {
				idint, _ := strconv.Atoi(id)
				mockPool.EXPECT().Exec(ctx, "DELETE FROM question_ WHERE id_ = $1", idint).Return(pgxmock.NewResult("DELETE", 0), nil)
			},
		},
		{
			nameTest: "exec_error",
			ctx:      context.Background(),
			id:       "345",
			mockBehavior: func(ctx context.Context, id string) {
				idint, _ := strconv.Atoi(id)
				mockPool.EXPECT().Exec(ctx, "DELETE FROM question_ WHERE id_ = $1", idint).Return(nil, errors.New("exec_error"))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, testCase.id)

			err := r.Delete(testCase.ctx, testCase.id)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
			case "invalid_inputs":
				assert.Equal(t, errs.ErrInvalidContent, err)
			case "no_question_to_delete":
				assert.Equal(t, errs.ErrContentNotFound, err)
			case "exec_error":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}
