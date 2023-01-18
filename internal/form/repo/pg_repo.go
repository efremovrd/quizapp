package repo

import (
	"context"
	"errors"
	"quiz-app/internal/form"
	"quiz-app/models"
	"quiz-app/pkg/errs"
	"quiz-app/pkg/postgres"
	"quiz-app/pkg/types"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type formDB struct {
	Id, UserId         int
	Title, Description string
}

type formRepo struct {
	ctxUserKey string
	*postgres.Postgres
}

func NewFormRepo(ctx_user_key string, db *postgres.Postgres) form.Repo {
	return &formRepo{ctx_user_key, db}
}

func (f *formRepo) Create(ctx context.Context, modelBL *models.Form) (*models.Form, error) {
	modelDB, err := formBLToDB(modelBL)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	sql, args, err := f.Builder.
		Insert("form_").
		Columns("user_id_, title_, description_").
		Values(modelDB.UserId, modelDB.Title, modelDB.Description).
		Suffix("RETURNING \"id_\"").
		ToSql()
	if err != nil {
		return nil, err
	}

	err = f.Pool.QueryRow(ctx, sql, args...).Scan(&modelDB.Id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == postgres.PermDenied {
			return nil, errs.ErrForbidden
		}

		return nil, err
	}

	return formDBToBL(modelDB)
}

func (f *formRepo) GetById(ctx context.Context, id string) (*models.Form, error) {
	intid, err := strconv.Atoi(id)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	sql, args, err := f.Builder.
		Select("user_id_, title_, description_").
		From("form_").
		Where(squirrel.Eq{"id_": intid}).
		ToSql()
	if err != nil {
		return nil, err
	}

	modelDB := formDB{Id: intid}
	err = f.Pool.QueryRow(ctx, sql, args...).Scan(&modelDB.UserId, &modelDB.Title, &modelDB.Description)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.ErrContentNotFound
		}

		return nil, err
	}

	return formDBToBL(&modelDB)
}

func (f *formRepo) GetByUserId(ctx context.Context, user_id string, sets types.GetSets) ([]*models.Form, error) {
	intuserid, err := strconv.Atoi(user_id)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	sql, args, err := f.Builder.
		Select("id_, title_, description_").
		From("form_").
		Where(squirrel.Eq{"user_id_": intuserid}).
		Limit(sets.Limit).
		Offset(sets.Offset).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := f.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*models.Form, 0)

	for rows.Next() {
		modelDB := formDB{UserId: intuserid}

		err = rows.Scan(&modelDB.Id, &modelDB.Title, &modelDB.Description)
		if err != nil {
			return nil, err
		}

		formBL, err := formDBToBL(&modelDB)
		if err != nil {
			return nil, err
		}

		res = append(res, formBL)
	}

	return res, nil
}

func (f *formRepo) Update(ctx context.Context, modelBL *models.Form) (*models.Form, error) {
	modelDB, err := formBLToDB(modelBL)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	builder := f.Builder.
		Update("form_")

	if modelDB.Title != "" {
		builder = builder.
			Set("title_", modelDB.Title)
	}

	if modelDB.Description != "" {
		builder = builder.
			Set("description_", modelDB.Description)
	}

	sql, args, err := builder.
		Where(squirrel.Eq{"id_": modelDB.Id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	res, err := f.Pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == postgres.PermDenied {
			return nil, errs.ErrForbidden
		}

		return nil, err
	}

	if res.RowsAffected() == 0 {
		return nil, errs.ErrContentNotFound
	}

	return modelBL, nil
}

func (f *formRepo) Delete(ctx context.Context, id string) error {
	intid, err := strconv.Atoi(id)
	if err != nil {
		return errs.ErrInvalidContent
	}

	sql, args, err := f.Builder.
		Delete("form_").
		Where(squirrel.Eq{"id_": intid}).
		ToSql()
	if err != nil {
		return err
	}

	res, err := f.Pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == postgres.PermDenied {
			return errs.ErrForbidden
		}

		return err
	}

	if res.RowsAffected() == 0 {
		return errs.ErrContentNotFound
	}

	return nil
}

func (f *formRepo) ValidateIsOwner(ctx context.Context, form_id string) error {
	intid, err := strconv.Atoi(form_id)
	if err != nil {
		return errs.ErrInvalidContent
	}

	sql, args, err := f.Builder.
		Select("user_id_").
		From("form_").
		Where(squirrel.Eq{"id_": intid}).
		ToSql()
	if err != nil {
		return err
	}

	var owner_id int
	err = f.Pool.QueryRow(ctx, sql, args...).Scan(&owner_id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errs.ErrContentNotFound
		}

		return err
	}

	currentuser, ok := ctx.Value(f.ctxUserKey).(*models.User)
	if !ok {
		return errs.ErrUnauthorized
	}

	if currentuser.Id != strconv.Itoa(owner_id) {
		return errs.ErrForbidden
	}

	return nil
}

func formDBToBL(modelDB *formDB) (*models.Form, error) {
	return &models.Form{
		Id:          strconv.Itoa(modelDB.Id),
		User_id:     strconv.Itoa(modelDB.UserId),
		Title:       modelDB.Title,
		Description: modelDB.Description,
	}, nil
}

func formBLToDB(modelBL *models.Form) (*formDB, error) {
	var (
		err error
		id  int
	)

	if modelBL.Id != "" {
		id, err = strconv.Atoi(modelBL.Id)
		if err != nil {
			return nil, err
		}
	}

	var uid int
	if modelBL.User_id != "" {
		uid, err = strconv.Atoi(modelBL.User_id)
		if err != nil {
			return nil, err
		}
	}

	return &formDB{
		Id:          id,
		UserId:      uid,
		Title:       modelBL.Title,
		Description: modelBL.Description,
	}, nil
}
