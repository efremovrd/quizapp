package repo

import (
	"context"
	"errors"
	"quiz-app/internal/auth"
	"quiz-app/models"
	"quiz-app/pkg/errs"
	"quiz-app/pkg/postgres"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type UserDB struct {
	Id              int
	Login, Password string
}

type authRepo struct {
	*postgres.Postgres
}

func NewAuthRepo(db *postgres.Postgres) auth.Repo {
	return &authRepo{db}
}

func (a *authRepo) Create(ctx context.Context, user *models.User) (*models.User, error) {
	userDB, err := userBLToDB(user)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	sql, args, err := a.Builder.
		Insert("user_").
		Columns("login_, password_").
		Values(userDB.Login, userDB.Password).
		Suffix("RETURNING \"id_\"").
		ToSql()
	if err != nil {
		return nil, err
	}

	err = a.Pool.QueryRow(ctx, sql, args...).Scan(&userDB.Id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == postgres.PermDenied {
			return nil, errs.ErrForbidden
		}

		return nil, err
	}

	return userDBToBL(userDB)
}

func (a *authRepo) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	sql, args, err := a.Builder.
		Select("id_, password_").
		From("user_").
		Where(squirrel.Eq{"login_": login}).
		ToSql()
	if err != nil {
		return nil, err
	}

	userDB := UserDB{Login: login}
	err = a.Pool.QueryRow(ctx, sql, args...).Scan(&userDB.Id, &userDB.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.ErrContentNotFound
		}

		return nil, err
	}

	return userDBToBL(&userDB)
}

func (a *authRepo) GetById(ctx context.Context, id string) (*models.User, error) {
	intid, err := strconv.Atoi(id)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	sql, args, err := a.Builder.
		Select("login_, password_").
		From("user_").
		Where(squirrel.Eq{"id_": intid}).
		ToSql()
	if err != nil {
		return nil, err
	}

	userDB := UserDB{Id: intid}
	err = a.Pool.QueryRow(ctx, sql, args...).Scan(&userDB.Login, &userDB.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.ErrContentNotFound
		}

		return nil, err
	}

	return userDBToBL(&userDB)
}

func userDBToBL(userDB *UserDB) (*models.User, error) {
	return &models.User{
		Id:       strconv.Itoa(userDB.Id),
		Login:    userDB.Login,
		Password: userDB.Password,
	}, nil
}

func userBLToDB(userBL *models.User) (*UserDB, error) {
	var (
		err error
		id  int
	)

	if userBL.Id != "" {
		id, err = strconv.Atoi(userBL.Id)
		if err != nil {
			return nil, err
		}
	}

	return &UserDB{
		Id:       id,
		Login:    userBL.Login,
		Password: userBL.Password,
	}, nil
}
