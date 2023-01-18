package repo_test

import (
	"context"
	"errors"
	"quiz-app/internal/form/repo"
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

	r := repo.NewFormRepo("any", &db)

	type mockBehavior func(ctx context.Context, form *models.Form)

	testTable := []struct {
		nameTest     string
		ctx          context.Context
		form         models.Form
		mockBehavior mockBehavior
		expectedForm models.Form
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			form: models.Form{
				User_id:     "12",
				Title:       "sdcsd",
				Description: "ecefvc",
			},
			mockBehavior: func(ctx context.Context, form *models.Form) {
				pgxRows := pgxpoolmock.NewRows([]string{"id_"}).AddRow(345).ToPgxRows()
				pgxRows.Next()
				useridint, _ := strconv.Atoi(form.User_id)
				mockPool.EXPECT().QueryRow(ctx, "INSERT INTO form_ (user_id_, title_, description_) VALUES ($1,$2,$3) RETURNING \"id_\"", useridint, form.Title, form.Description).Return(pgxRows)
			},
			expectedForm: models.Form{
				Id:          "345",
				User_id:     "12",
				Title:       "sdcsd",
				Description: "ecefvc",
			},
		},
		{
			nameTest: "invalid_inputs",
			ctx:      context.Background(),
			form: models.Form{
				User_id:     "5r4",
				Title:       "sdcsd",
				Description: "ecefvc",
			},
			mockBehavior: func(ctx context.Context, form *models.Form) {},
		},
		{
			nameTest: "no_rows",
			ctx:      context.Background(),
			form: models.Form{
				User_id:     "12",
				Title:       "sdcsd",
				Description: "ecefvc",
			},
			mockBehavior: func(ctx context.Context, form *models.Form) {
				pgxRows := pgxpoolmock.NewRows([]string{}).AddRow().ToPgxRows()
				useridint, _ := strconv.Atoi(form.User_id)
				mockPool.EXPECT().QueryRow(ctx, "INSERT INTO form_ (user_id_, title_, description_) VALUES ($1,$2,$3) RETURNING \"id_\"", useridint, form.Title, form.Description).Return(pgxRows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, &testCase.form)

			got, err := r.Create(testCase.ctx, &testCase.form)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedForm, *got)
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

func TestFormRepo_GetById(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewFormRepo("any", &db)

	type mockBehavior func(ctx context.Context, id string)

	testTable := []struct {
		nameTest     string
		ctx          context.Context
		id           string
		mockBehavior mockBehavior
		expectedForm models.Form
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			id:       "345",
			mockBehavior: func(ctx context.Context, id string) {
				pgxRows := pgxpoolmock.NewRows([]string{"user_id_", "title_", "description_"}).AddRow(12, "sdcsd", "ecefvc").ToPgxRows()
				pgxRows.Next()
				idint, _ := strconv.Atoi(id)
				mockPool.EXPECT().QueryRow(ctx, "SELECT user_id_, title_, description_ FROM form_ WHERE id_ = $1", idint).Return(pgxRows)
			},
			expectedForm: models.Form{
				Id:          "345",
				User_id:     "12",
				Title:       "sdcsd",
				Description: "ecefvc",
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
				mockPool.EXPECT().QueryRow(ctx, "SELECT user_id_, title_, description_ FROM form_ WHERE id_ = $1", idint).Return(pgxRows)
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
				assert.Equal(t, testCase.expectedForm, *got)
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

func TestFormRepo_GetByUserId(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewFormRepo("any", &db)

	type mockBehavior func(ctx context.Context, user_id string, sets types.GetSets)

	testTable := []struct {
		nameTest      string
		ctx           context.Context
		user_id       string
		sets          types.GetSets
		mockBehavior  mockBehavior
		expectedForms []*models.Form
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			user_id:  "12",
			sets:     types.GetSets{},
			mockBehavior: func(ctx context.Context, user_id string, sets types.GetSets) {
				pgxRows := pgxpoolmock.NewRows([]string{"id_", "title_", "description_"}).AddRow(345, "sdcsd", "ecefvc").AddRow(346, "qwer", "ty").ToPgxRows()
				useridint, _ := strconv.Atoi(user_id)
				mockPool.EXPECT().Query(ctx, "SELECT id_, title_, description_ FROM form_ WHERE user_id_ = $1 LIMIT 0 OFFSET 0", useridint).Return(pgxRows, nil)
			},
			expectedForms: []*models.Form{
				{
					Id:          "345",
					User_id:     "12",
					Title:       "sdcsd",
					Description: "ecefvc",
				},
				{
					Id:          "346",
					User_id:     "12",
					Title:       "qwer",
					Description: "ty",
				},
			},
		},
		{
			nameTest:     "invalid_inputs",
			ctx:          context.Background(),
			user_id:      "5r4",
			sets:         types.GetSets{},
			mockBehavior: func(ctx context.Context, user_id string, sets types.GetSets) {},
		},
		{
			nameTest: "query_error",
			ctx:      context.Background(),
			user_id:  "12",
			sets:     types.GetSets{},
			mockBehavior: func(ctx context.Context, user_id string, sets types.GetSets) {
				useridint, _ := strconv.Atoi(user_id)
				mockPool.EXPECT().Query(ctx, "SELECT id_, title_, description_ FROM form_ WHERE user_id_ = $1 LIMIT 0 OFFSET 0", useridint).Return(nil, errors.New("query_error"))
			},
		},
		{
			nameTest: "no_rows",
			ctx:      context.Background(),
			user_id:  "12",
			sets:     types.GetSets{},
			mockBehavior: func(ctx context.Context, user_id string, sets types.GetSets) {
				pgxRows := pgxpoolmock.NewRows([]string{}).AddRow().ToPgxRows()
				pgxRows.Next()
				useridint, _ := strconv.Atoi(user_id)
				mockPool.EXPECT().Query(ctx, "SELECT id_, title_, description_ FROM form_ WHERE user_id_ = $1 LIMIT 0 OFFSET 0", useridint).Return(pgxRows, nil)
			},
			expectedForms: []*models.Form{},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, testCase.user_id, types.GetSets{})

			got, err := r.GetByUserId(testCase.ctx, testCase.user_id, types.GetSets{})

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedForms, got)
			case "invalid_inputs":
				assert.Equal(t, errs.ErrInvalidContent, err)
			case "query_error":
				assert.NotEqual(t, nil, err)
			case "no_rows":
				assert.Equal(t, nil, err)
				assert.Equal(t, []*models.Form{}, got)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestFormRepo_Update(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewFormRepo("any", &db)

	type mockBehavior func(ctx context.Context, form *models.Form)

	testTable := []struct {
		nameTest     string
		ctx          context.Context
		form         models.Form
		mockBehavior mockBehavior
		expectedForm models.Form
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			form: models.Form{
				Id:          "345",
				User_id:     "12",
				Title:       "sdcsd",
				Description: "ecefvc",
			},
			mockBehavior: func(ctx context.Context, form *models.Form) {
				idint, _ := strconv.Atoi(form.Id)
				mockPool.EXPECT().Exec(ctx, "UPDATE form_ SET title_ = $1, description_ = $2 WHERE id_ = $3", form.Title, form.Description, idint).Return(pgxmock.NewResult("UPDATE", 1), nil)
			},
			expectedForm: models.Form{
				Id:          "345",
				User_id:     "12",
				Title:       "sdcsd",
				Description: "ecefvc",
			},
		},
		{
			nameTest: "invalid_inputs_user_id",
			ctx:      context.Background(),
			form: models.Form{
				Id:          "345",
				User_id:     "5r4",
				Title:       "sdcsd",
				Description: "ecefvc",
			},
			mockBehavior: func(ctx context.Context, form *models.Form) {},
		},
		{
			nameTest: "invalid_inputs_id",
			ctx:      context.Background(),
			form: models.Form{
				Id:          "5r4",
				User_id:     "45",
				Title:       "sdcsd",
				Description: "ecefvc",
			},
			mockBehavior: func(ctx context.Context, form *models.Form) {},
		},
		{
			nameTest: "no_form_to_update",
			ctx:      context.Background(),
			form: models.Form{
				Id:          "345",
				User_id:     "12",
				Title:       "sdcsd",
				Description: "ecefvc",
			},
			mockBehavior: func(ctx context.Context, form *models.Form) {
				idint, _ := strconv.Atoi(form.Id)
				mockPool.EXPECT().Exec(ctx, "UPDATE form_ SET title_ = $1, description_ = $2 WHERE id_ = $3", form.Title, form.Description, idint).Return(pgxmock.NewResult("UPDATE", 0), nil)
			},
		},
		{
			nameTest: "exec_error",
			ctx:      context.Background(),
			form: models.Form{
				Id:          "345",
				User_id:     "12",
				Title:       "sdcsd",
				Description: "ecefvc",
			},
			mockBehavior: func(ctx context.Context, form *models.Form) {
				idint, _ := strconv.Atoi(form.Id)
				mockPool.EXPECT().Exec(ctx, "UPDATE form_ SET title_ = $1, description_ = $2 WHERE id_ = $3", form.Title, form.Description, idint).Return(nil, errors.New("exec_error"))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, &testCase.form)

			got, err := r.Update(testCase.ctx, &testCase.form)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedForm, *got)
			case "invalid_inputs_id", "invalid_inputs_user_id":
				assert.Equal(t, errs.ErrInvalidContent, err)
			case "no_form_to_update":
				assert.Equal(t, errs.ErrContentNotFound, err)
			case "exec_error":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestFormRepo_Delete(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewFormRepo("any", &db)

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
				mockPool.EXPECT().Exec(ctx, "DELETE FROM form_ WHERE id_ = $1", idint).Return(pgxmock.NewResult("DELETE", 1), nil)
			},
		},
		{
			nameTest:     "invalid_inputs",
			ctx:          context.Background(),
			id:           "5r4",
			mockBehavior: func(ctx context.Context, id string) {},
		},
		{
			nameTest: "no_form_to_delete",
			ctx:      context.Background(),
			id:       "345",
			mockBehavior: func(ctx context.Context, id string) {
				idint, _ := strconv.Atoi(id)
				mockPool.EXPECT().Exec(ctx, "DELETE FROM form_ WHERE id_ = $1", idint).Return(pgxmock.NewResult("DELETE", 0), nil)
			},
		},
		{
			nameTest: "exec_error",
			ctx:      context.Background(),
			id:       "345",
			mockBehavior: func(ctx context.Context, id string) {
				idint, _ := strconv.Atoi(id)
				mockPool.EXPECT().Exec(ctx, "DELETE FROM form_ WHERE id_ = $1", idint).Return(nil, errors.New("exec_error"))
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
			case "no_form_to_delete":
				assert.Equal(t, errs.ErrContentNotFound, err)
			case "exec_error":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestFormRepo_ValidateIsOwner(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	ctxuserkey := "ctxuserkey"

	r := repo.NewFormRepo(ctxuserkey, &db)

	type mockBehavior func(ctx context.Context, form_id string)

	user := models.User{
		Id:       "23",
		Login:    "qwe",
		Password: "rty",
	}

	testTable := []struct {
		nameTest     string
		ctx          context.Context
		form_id      string
		mockBehavior mockBehavior
	}{
		{
			nameTest: "ok",
			ctx:      context.WithValue(context.Background(), ctxuserkey, &user),
			form_id:  "345",
			mockBehavior: func(ctx context.Context, id string) {
				pgxrows := pgxpoolmock.NewRows([]string{"user_id_"}).AddRow(23).ToPgxRows()
				pgxrows.Next()
				idint, _ := strconv.Atoi(id)
				mockPool.EXPECT().QueryRow(ctx, "SELECT user_id_ FROM form_ WHERE id_ = $1", idint).Return(pgxrows)
			},
		},
		{
			nameTest:     "invalid_inputs",
			ctx:          context.Background(),
			form_id:      "5r4",
			mockBehavior: func(ctx context.Context, id string) {},
		},
		{
			nameTest: "no_form",
			ctx:      context.Background(),
			form_id:  "345",
			mockBehavior: func(ctx context.Context, id string) {
				pgxRows := pgxpoolmock.NewRows([]string{}).AddRow().ToPgxRows()
				idint, _ := strconv.Atoi(id)
				mockPool.EXPECT().QueryRow(ctx, "SELECT user_id_ FROM form_ WHERE id_ = $1", idint).Return(pgxRows)
			},
		},
		{
			nameTest: "unauthorized",
			ctx:      context.Background(),
			form_id:  "345",
			mockBehavior: func(ctx context.Context, id string) {
				pgxrows := pgxpoolmock.NewRows([]string{"user_id_"}).AddRow(23).ToPgxRows()
				pgxrows.Next()
				idint, _ := strconv.Atoi(id)
				mockPool.EXPECT().QueryRow(ctx, "SELECT user_id_ FROM form_ WHERE id_ = $1", idint).Return(pgxrows)
			},
		},
		{
			nameTest: "not_an_owner",
			ctx:      context.WithValue(context.Background(), ctxuserkey, &user),
			form_id:  "345",
			mockBehavior: func(ctx context.Context, id string) {
				pgxrows := pgxpoolmock.NewRows([]string{"user_id_"}).AddRow(50).ToPgxRows()
				pgxrows.Next()
				idint, _ := strconv.Atoi(id)
				mockPool.EXPECT().QueryRow(ctx, "SELECT user_id_ FROM form_ WHERE id_ = $1", idint).Return(pgxrows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, testCase.form_id)

			err := r.ValidateIsOwner(testCase.ctx, testCase.form_id)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
			case "invalid_inputs":
				assert.Equal(t, errs.ErrInvalidContent, err)
			case "unauthorized":
				assert.Equal(t, errs.ErrUnauthorized, err)
			case "not_an_owner":
				assert.Equal(t, errs.ErrForbidden, err)
			case "no_form":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}
