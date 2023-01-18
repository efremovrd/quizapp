package repo_test

import (
	"context"
	"errors"
	"quiz-app/internal/auth/repo"
	"quiz-app/models"
	"quiz-app/pkg/errs"
	"quiz-app/pkg/postgres"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
)

var (
	_builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

func TestAuthRepo_Create(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewAuthRepo(&db)

	type mockBehavior func(ctx context.Context, user *models.User)

	testTable := []struct {
		nameTest     string
		ctx          context.Context
		user         models.User
		mockBehavior mockBehavior
		expectedUser models.User
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			user: models.User{
				Login:    "sdcsd",
				Password: "ecefvc",
			},
			mockBehavior: func(ctx context.Context, user *models.User) {
				pgxRows := pgxpoolmock.NewRows([]string{"id_"}).AddRow(345).ToPgxRows()
				pgxRows.Next()
				mockPool.EXPECT().QueryRow(ctx, "INSERT INTO user_ (login_, password_) VALUES ($1,$2) RETURNING \"id_\"", user.Login, user.Password).Return(pgxRows)
			},
			expectedUser: models.User{
				Id:       "345",
				Login:    "sdcsd",
				Password: "ecefvc",
			},
		},
		{
			nameTest: "invalid_inputs",
			ctx:      context.Background(),
			user: models.User{
				Id:       "5r4",
				Login:    "sdcsd",
				Password: "ecefvc",
			},
			mockBehavior: func(ctx context.Context, user *models.User) {},
		},
		{
			nameTest: "no_rows",
			ctx:      context.Background(),
			user: models.User{
				Id:       "345",
				Login:    "sdcsd",
				Password: "ecefvc",
			},
			mockBehavior: func(ctx context.Context, user *models.User) {
				pgxRows := pgxpoolmock.NewRows([]string{}).AddRow().ToPgxRows()
				mockPool.EXPECT().QueryRow(ctx, "INSERT INTO user_ (login_, password_) VALUES ($1,$2) RETURNING \"id_\"", user.Login, user.Password).Return(pgxRows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, &testCase.user)

			got, err := r.Create(testCase.ctx, &testCase.user)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedUser, *got)
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

func TestAuthRepo_GetByLogin(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewAuthRepo(&db)

	type mockBehavior func(ctx context.Context, login string)

	testTable := []struct {
		nameTest     string
		ctx          context.Context
		login        string
		mockBehavior mockBehavior
		expectedUser models.User
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			login:    "sdcsd",
			mockBehavior: func(ctx context.Context, login string) {
				pgxRows := pgxpoolmock.NewRows([]string{"id_", "password_"}).AddRow(345, "ecefvc").ToPgxRows()
				pgxRows.Next()
				mockPool.EXPECT().QueryRow(ctx, "SELECT id_, password_ FROM user_ WHERE login_ = $1", login).Return(pgxRows)
			},
			expectedUser: models.User{
				Id:       "345",
				Login:    "sdcsd",
				Password: "ecefvc",
			},
		},
		{
			nameTest: "no_rows",
			ctx:      context.Background(),
			login:    "sdcsd",
			mockBehavior: func(ctx context.Context, login string) {
				pgxRows := pgxpoolmock.NewRows([]string{}).AddRow().ToPgxRows()
				mockPool.EXPECT().QueryRow(ctx, "SELECT id_, password_ FROM user_ WHERE login_ = $1", login).Return(pgxRows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, testCase.login)

			got, err := r.GetByLogin(testCase.ctx, testCase.login)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedUser, *got)
			case "no_rows":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

func TestAuthRepo_GetById(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewAuthRepo(&db)

	type mockBehavior func(ctx context.Context, id string)

	testTable := []struct {
		nameTest     string
		ctx          context.Context
		id           string
		mockBehavior mockBehavior
		expectedUser models.User
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			id:       "345",
			mockBehavior: func(ctx context.Context, id string) {
				pgxRows := pgxpoolmock.NewRows([]string{"login_", "password_"}).AddRow("sdcsd", "ecefvc").ToPgxRows()
				pgxRows.Next()
				mockPool.EXPECT().QueryRow(ctx, "SELECT login_, password_ FROM user_ WHERE id_ = $1", gomock.Any()).Return(pgxRows)
			},
			expectedUser: models.User{
				Id:       "345",
				Login:    "sdcsd",
				Password: "ecefvc",
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
				mockPool.EXPECT().QueryRow(ctx, "SELECT login_, password_ FROM user_ WHERE id_ = $1", gomock.Any()).Return(pgxRows)
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
				assert.Equal(t, testCase.expectedUser, *got)
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
