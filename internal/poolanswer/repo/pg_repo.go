package repo

import (
	"context"
	"errors"
	"quizapp/internal/poolanswer"
	"quizapp/models"
	"quizapp/pkg/errs"
	"quizapp/pkg/postgres"
	"quizapp/pkg/types"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type PoolAnswerDB struct {
	ID     int
	UserID int
	FormID int
}

type poolAnswerRepo struct {
	*postgres.Postgres
}

func NewPoolAnswerRepo(db *postgres.Postgres) poolanswer.Repo {
	return &poolAnswerRepo{db}
}

func (p *poolAnswerRepo) Create(ctx context.Context, poolanswer *models.PoolAnswer) (*models.PoolAnswer, error) {
	poolanswerDB, err := paBLToDB(poolanswer)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	sql, args, err := p.Builder.
		Insert("pool_answer_").
		Columns("user_id_, form_id_").
		Values(poolanswerDB.UserID, poolanswerDB.FormID).
		Suffix("RETURNING \"id_\"").
		ToSql()
	if err != nil {
		return nil, err
	}

	err = p.Pool.QueryRow(ctx, sql, args...).Scan(&poolanswerDB.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == postgres.PermDenied {
			return nil, errs.ErrForbidden
		}

		return nil, err
	}

	return paDBToBL(poolanswerDB)
}

func (p *poolAnswerRepo) GetByFormId(ctx context.Context, form_id string, sets types.GetSets) ([]*models.PoolAnswer, error) {
	intid, err := strconv.Atoi(form_id)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	sql, args, err := p.Builder.
		Select("id_, user_id_").
		From("pool_answer_").
		Where(squirrel.Eq{"form_id_": intid}).
		Limit(sets.Limit).
		Offset(sets.Offset).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*models.PoolAnswer, 0)

	for rows.Next() {
		paDB := PoolAnswerDB{FormID: intid}

		err = rows.Scan(&paDB.ID, &paDB.UserID)
		if err != nil {
			return nil, err
		}

		formBL, err := paDBToBL(&paDB)
		if err != nil {
			return nil, err
		}

		res = append(res, formBL)
	}

	return res, nil
}

func (p *poolAnswerRepo) GetById(ctx context.Context, id string) (*models.PoolAnswer, error) {
	intid, err := strconv.Atoi(id)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	sql, args, err := p.Builder.
		Select("user_id_, form_id_").
		From("pool_answer_").
		Where(squirrel.Eq{"id_": intid}).
		ToSql()
	if err != nil {
		return nil, err
	}

	modelDB := PoolAnswerDB{ID: intid}
	err = p.Pool.QueryRow(ctx, sql, args...).Scan(&modelDB.UserID, &modelDB.FormID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.ErrContentNotFound
		}

		return nil, err
	}

	return paDBToBL(&modelDB)
}

func (p *poolAnswerRepo) Delete(ctx context.Context, id string) error {
	intid, err := strconv.Atoi(id)
	if err != nil {
		return errs.ErrInvalidContent
	}

	sql, args, err := p.Builder.
		Delete("pool_answer_").
		Where(squirrel.Eq{"id_": intid}).
		ToSql()
	if err != nil {
		return err
	}

	res, err := p.Pool.Exec(ctx, sql, args...)
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

func paDBToBL(paDB *PoolAnswerDB) (*models.PoolAnswer, error) {
	return &models.PoolAnswer{
		Id:      strconv.Itoa(paDB.ID),
		User_id: strconv.Itoa(paDB.UserID),
		Form_id: strconv.Itoa(paDB.FormID),
	}, nil
}

func paBLToDB(paBL *models.PoolAnswer) (*PoolAnswerDB, error) {
	var (
		err error
		id  int
	)

	if paBL.Id != "" {
		id, err = strconv.Atoi(paBL.Id)
		if err != nil {
			return nil, err
		}
	}

	var uid int
	if paBL.User_id != "" {
		uid, err = strconv.Atoi(paBL.User_id)
		if err != nil {
			return nil, err
		}
	}

	var fid int
	if paBL.Form_id != "" {
		fid, err = strconv.Atoi(paBL.Form_id)
		if err != nil {
			return nil, err
		}
	}

	return &PoolAnswerDB{
		ID:     id,
		UserID: uid,
		FormID: fid,
	}, nil
}
