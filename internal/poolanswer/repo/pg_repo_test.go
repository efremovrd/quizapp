package repo_test

import (
	"context"
	"errors"
	"quizapp/internal/poolanswer/repo"
	"quizapp/models"
	"quizapp/pkg/errs"
	"quizapp/pkg/postgres"
	"quizapp/pkg/types"
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

func TestPoolAnswerRepo_Create(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewPoolAnswerRepo(&db)

	type mockBehavior func(ctx context.Context, poolanswer *models.PoolAnswer)

	testTable := []struct {
		nameTest            string
		ctx                 context.Context
		pool_answer         models.PoolAnswer
		mockBehavior        mockBehavior
		expectedpoolsanswer models.PoolAnswer
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			pool_answer: models.PoolAnswer{
				Form_id: "12",
				User_id: "14",
			},
			mockBehavior: func(ctx context.Context, pool_answer *models.PoolAnswer) {
				pgxRows := pgxpoolmock.NewRows([]string{"id_"}).AddRow(345).ToPgxRows()
				pgxRows.Next()
				formidint, _ := strconv.Atoi(pool_answer.Form_id)
				useridint, _ := strconv.Atoi(pool_answer.User_id)
				mockPool.EXPECT().QueryRow(ctx, "INSERT INTO pool_answer_ (user_id_, form_id_) VALUES ($1,$2) RETURNING \"id_\"", useridint, formidint).Return(pgxRows)
			},
			expectedpoolsanswer: models.PoolAnswer{
				Id:      "345",
				Form_id: "12",
				User_id: "14",
			},
		},
		{
			nameTest: "invalid_inputs_form_id",
			ctx:      context.Background(),
			pool_answer: models.PoolAnswer{
				Form_id: "5r4",
				User_id: "14",
			},
			mockBehavior: func(ctx context.Context, pool_answer *models.PoolAnswer) {},
		},
		{
			nameTest: "invalid_inputs_user_id",
			ctx:      context.Background(),
			pool_answer: models.PoolAnswer{
				Form_id: "54",
				User_id: "1r4",
			},
			mockBehavior: func(ctx context.Context, pool_answer *models.PoolAnswer) {},
		},
		{
			nameTest: "no_rows",
			ctx:      context.Background(),
			pool_answer: models.PoolAnswer{
				Form_id: "12",
				User_id: "14",
			},
			mockBehavior: func(ctx context.Context, pool_answer *models.PoolAnswer) {
				pgxRows := pgxpoolmock.NewRows([]string{}).AddRow().ToPgxRows()
				formidint, _ := strconv.Atoi(pool_answer.Form_id)
				useridint, _ := strconv.Atoi(pool_answer.User_id)
				mockPool.EXPECT().QueryRow(ctx, "INSERT INTO pool_answer_ (user_id_, form_id_) VALUES ($1,$2) RETURNING \"id_\"", useridint, formidint).Return(pgxRows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, &testCase.pool_answer)

			got, err := r.Create(testCase.ctx, &testCase.pool_answer)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedpoolsanswer, *got)
			case "invalid_inputs_form_id", "invalid_inputs_user_id":
				assert.Equal(t, errs.ErrInvalidContent, err)
			case "no_rows":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestPoolAnswerRepo_GetById(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewPoolAnswerRepo(&db)

	type mockBehavior func(ctx context.Context, id string)

	testTable := []struct {
		nameTest            string
		ctx                 context.Context
		id                  string
		mockBehavior        mockBehavior
		expectedpoolsanswer models.PoolAnswer
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			id:       "345",
			mockBehavior: func(ctx context.Context, id string) {
				pgxRows := pgxpoolmock.NewRows([]string{"user_id_", "form_id_"}).AddRow(12, 14).ToPgxRows()
				pgxRows.Next()
				idint, _ := strconv.Atoi(id)
				mockPool.EXPECT().QueryRow(ctx, "SELECT user_id_, form_id_ FROM pool_answer_ WHERE id_ = $1", idint).Return(pgxRows)
			},
			expectedpoolsanswer: models.PoolAnswer{
				Id:      "345",
				Form_id: "14",
				User_id: "12",
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
				mockPool.EXPECT().QueryRow(ctx, "SELECT user_id_, form_id_ FROM pool_answer_ WHERE id_ = $1", idint).Return(pgxRows)
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
				assert.Equal(t, testCase.expectedpoolsanswer, *got)
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

func TestPoolAnswerRepo_GetByFormId(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewPoolAnswerRepo(&db)

	type mockBehavior func(ctx context.Context, form_id string, sets types.GetSets)

	testTable := []struct {
		nameTest             string
		ctx                  context.Context
		form_id              string
		sets                 types.GetSets
		mockBehavior         mockBehavior
		expectedpoolsanswers []*models.PoolAnswer
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			form_id:  "12",
			sets:     types.GetSets{},
			mockBehavior: func(ctx context.Context, form_id string, sets types.GetSets) {
				pgxRows := pgxpoolmock.NewRows([]string{"id_", "user_id_"}).AddRow(345, 14).AddRow(346, 15).ToPgxRows()
				formidint, _ := strconv.Atoi(form_id)
				mockPool.EXPECT().Query(ctx, "SELECT id_, user_id_ FROM pool_answer_ WHERE form_id_ = $1 LIMIT 0 OFFSET 0", formidint).Return(pgxRows, nil)
			},
			expectedpoolsanswers: []*models.PoolAnswer{
				{
					Id:      "345",
					Form_id: "12",
					User_id: "14",
				},
				{
					Id:      "346",
					Form_id: "12",
					User_id: "15",
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
				mockPool.EXPECT().Query(ctx, "SELECT id_, user_id_ FROM pool_answer_ WHERE form_id_ = $1 LIMIT 0 OFFSET 0", formidint).Return(nil, errors.New("query_error"))
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
				mockPool.EXPECT().Query(ctx, "SELECT id_, user_id_ FROM pool_answer_ WHERE form_id_ = $1 LIMIT 0 OFFSET 0", formidint).Return(pgxRows, nil)
			},
			expectedpoolsanswers: []*models.PoolAnswer{},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, testCase.form_id, types.GetSets{})

			got, err := r.GetByFormId(testCase.ctx, testCase.form_id, types.GetSets{})

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedpoolsanswers, got)
			case "invalid_inputs":
				assert.Equal(t, errs.ErrInvalidContent, err)
			case "query_error":
				assert.NotEqual(t, nil, err)
			case "no_rows":
				assert.Equal(t, nil, err)
				assert.Equal(t, []*models.PoolAnswer{}, got)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestPoolAnswerRepo_Delete(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewPoolAnswerRepo(&db)

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
				mockPool.EXPECT().Exec(ctx, "DELETE FROM pool_answer_ WHERE id_ = $1", idint).Return(pgxmock.NewResult("DELETE", 1), nil)
			},
		},
		{
			nameTest:     "invalid_inputs",
			ctx:          context.Background(),
			id:           "5r4",
			mockBehavior: func(ctx context.Context, id string) {},
		},
		{
			nameTest: "no_pool_answer_to_delete",
			ctx:      context.Background(),
			id:       "345",
			mockBehavior: func(ctx context.Context, id string) {
				idint, _ := strconv.Atoi(id)
				mockPool.EXPECT().Exec(ctx, "DELETE FROM pool_answer_ WHERE id_ = $1", idint).Return(pgxmock.NewResult("DELETE", 0), nil)
			},
		},
		{
			nameTest: "exec_error",
			ctx:      context.Background(),
			id:       "345",
			mockBehavior: func(ctx context.Context, id string) {
				idint, _ := strconv.Atoi(id)
				mockPool.EXPECT().Exec(ctx, "DELETE FROM pool_answer_ WHERE id_ = $1", idint).Return(nil, errors.New("exec_error"))
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
			case "no_pool_answer_to_delete":
				assert.Equal(t, errs.ErrContentNotFound, err)
			case "exec_error":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}
