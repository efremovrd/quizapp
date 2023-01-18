package repo_test

import (
	"context"
	"errors"
	"quiz-app/internal/answer/repo"
	"quiz-app/models"
	"quiz-app/pkg/errs"
	"quiz-app/pkg/postgres"
	"quiz-app/pkg/types"
	"strconv"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
)

var (
	_builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

func TestAnswerRepo_Create(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewAnswerRepo(&db)

	type mockBehavior func(ctx context.Context, answer *models.Answer)

	testTable := []struct {
		nameTest       string
		ctx            context.Context
		answer         models.Answer
		mockBehavior   mockBehavior
		expectedanswer models.Answer
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			answer: models.Answer{
				Question_id:    "12",
				Pool_answer_id: "14",
				Value:          "answer",
			},
			mockBehavior: func(ctx context.Context, answer *models.Answer) {
				pgxRows := pgxpoolmock.NewRows([]string{"id_"}).AddRow(345).ToPgxRows()
				pgxRows.Next()
				qintid, _ := strconv.Atoi(answer.Question_id)
				paintid, _ := strconv.Atoi(answer.Pool_answer_id)
				mockPool.EXPECT().QueryRow(ctx, "INSERT INTO answer_ (question_id_, pool_answer_id_, value_) VALUES ($1,$2,$3) RETURNING \"id_\"", qintid, paintid, answer.Value).Return(pgxRows)
			},
			expectedanswer: models.Answer{
				Id:             "345",
				Question_id:    "12",
				Pool_answer_id: "14",
				Value:          "answer",
			},
		},
		{
			nameTest: "invalid_inputs_question_id",
			ctx:      context.Background(),
			answer: models.Answer{
				Question_id:    "1r2",
				Pool_answer_id: "14",
				Value:          "answer",
			},
			mockBehavior: func(ctx context.Context, answer *models.Answer) {},
		},
		{
			nameTest: "invalid_inputs_pa_id",
			ctx:      context.Background(),
			answer: models.Answer{
				Question_id:    "12",
				Pool_answer_id: "1r4",
				Value:          "answer",
			},
			mockBehavior: func(ctx context.Context, answer *models.Answer) {},
		},
		{
			nameTest: "no_rows",
			ctx:      context.Background(),
			answer: models.Answer{
				Question_id:    "12",
				Pool_answer_id: "14",
				Value:          "answer",
			},
			mockBehavior: func(ctx context.Context, answer *models.Answer) {
				pgxRows := pgxpoolmock.NewRows([]string{}).AddRow().ToPgxRows()
				qintid, _ := strconv.Atoi(answer.Question_id)
				paintid, _ := strconv.Atoi(answer.Pool_answer_id)
				mockPool.EXPECT().QueryRow(ctx, "INSERT INTO answer_ (question_id_, pool_answer_id_, value_) VALUES ($1,$2,$3) RETURNING \"id_\"", qintid, paintid, answer.Value).Return(pgxRows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, &testCase.answer)

			got, err := r.Create(testCase.ctx, &testCase.answer)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedanswer, *got)
			case "invalid_inputs_question_id", "invalid_inputs_pa_id":
				assert.Equal(t, errs.ErrInvalidContent, err)
			case "no_rows":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestAnswerRepo_GetByPoolAnswerId(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewAnswerRepo(&db)

	type mockBehavior func(ctx context.Context, pool_answer_id string, sets types.GetSets)

	testTable := []struct {
		nameTest        string
		ctx             context.Context
		pa_id           string
		sets            types.GetSets
		mockBehavior    mockBehavior
		expectedanswers []*models.Answer
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			pa_id:    "12",
			sets:     types.GetSets{},
			mockBehavior: func(ctx context.Context, pool_answer_id string, sets types.GetSets) {
				pgxRows := pgxpoolmock.NewRows([]string{"id_", "question_id_", "value_"}).AddRow(345, 14, "ans1").AddRow(346, 15, "ans2").ToPgxRows()
				paintid, _ := strconv.Atoi(pool_answer_id)
				mockPool.EXPECT().Query(ctx, "SELECT id_, question_id_, value_ FROM answer_ WHERE pool_answer_id_ = $1 LIMIT 0 OFFSET 0", paintid).Return(pgxRows, nil)
			},
			expectedanswers: []*models.Answer{
				{
					Id:             "345",
					Pool_answer_id: "12",
					Question_id:    "14",
					Value:          "ans1",
				},
				{
					Id:             "346",
					Pool_answer_id: "12",
					Question_id:    "15",
					Value:          "ans2",
				},
			},
		},
		{
			nameTest:     "invalid_inputs",
			ctx:          context.Background(),
			pa_id:        "5r4",
			sets:         types.GetSets{},
			mockBehavior: func(ctx context.Context, pool_answer_id string, sets types.GetSets) {},
		},
		{
			nameTest: "query_error",
			ctx:      context.Background(),
			pa_id:    "12",
			sets:     types.GetSets{},
			mockBehavior: func(ctx context.Context, pool_answer_id string, sets types.GetSets) {
				paintid, _ := strconv.Atoi(pool_answer_id)
				mockPool.EXPECT().Query(ctx, "SELECT id_, question_id_, value_ FROM answer_ WHERE pool_answer_id_ = $1 LIMIT 0 OFFSET 0", paintid).Return(nil, errors.New("query_error"))
			},
		},
		{
			nameTest: "no_rows",
			ctx:      context.Background(),
			pa_id:    "12",
			sets:     types.GetSets{},
			mockBehavior: func(ctx context.Context, form_id string, sets types.GetSets) {
				pgxRows := pgxpoolmock.NewRows([]string{}).AddRow().ToPgxRows()
				pgxRows.Next()
				paintid, _ := strconv.Atoi(form_id)
				mockPool.EXPECT().Query(ctx, "SELECT id_, question_id_, value_ FROM answer_ WHERE pool_answer_id_ = $1 LIMIT 0 OFFSET 0", paintid).Return(pgxRows, nil)
			},
			expectedanswers: []*models.Answer{},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, testCase.pa_id, types.GetSets{})

			got, err := r.GetByPoolAnswerId(testCase.ctx, testCase.pa_id, types.GetSets{})

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedanswers, got)
			case "invalid_inputs":
				assert.Equal(t, errs.ErrInvalidContent, err)
			case "query_error":
				assert.NotEqual(t, nil, err)
			case "no_rows":
				assert.Equal(t, nil, err)
				assert.Equal(t, []*models.Answer{}, got)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}
