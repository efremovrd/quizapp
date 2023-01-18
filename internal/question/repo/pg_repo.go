package repo

import (
	"context"
	"errors"
	"quiz-app/internal/question"
	"quiz-app/models"
	"quiz-app/pkg/errs"
	"quiz-app/pkg/postgres"
	"quiz-app/pkg/types"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type QuestionDB struct {
	Id, FormId int
	Header     string
}

type questionRepo struct {
	*postgres.Postgres
}

func NewQuestionRepo(db *postgres.Postgres) question.Repo {
	return &questionRepo{db}
}

func (qr *questionRepo) Create(ctx context.Context, modelBL *models.Question) (*models.Question, error) {
	questionDB, err := questionBLToDB(modelBL)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	sql, args, err := qr.Builder.
		Insert("question_").
		Columns("form_id_, header_").
		Values(questionDB.FormId, questionDB.Header).
		Suffix("RETURNING \"id_\"").
		ToSql()
	if err != nil {
		return nil, err
	}

	err = qr.Pool.QueryRow(ctx, sql, args...).Scan(&questionDB.Id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == postgres.PermDenied {
			return nil, errs.ErrForbidden
		}

		return nil, err
	}

	return questionDBToBL(questionDB)
}

func (q *questionRepo) GetByFormId(ctx context.Context, form_id string, sets types.GetSets) ([]*models.Question, error) {
	intformid, err := strconv.Atoi(form_id)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	sql, args, err := q.Builder.
		Select("id_, header_").
		From("question_").
		Where(squirrel.Eq{"form_id_": intformid}).
		Limit(sets.Limit).
		Offset(sets.Offset).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := q.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*models.Question, 0)

	for rows.Next() {
		modelDB := QuestionDB{FormId: intformid}

		err = rows.Scan(&modelDB.Id, &modelDB.Header)
		if err != nil {
			return nil, err
		}

		questionBL, err := questionDBToBL(&modelDB)
		if err != nil {
			return nil, err
		}

		res = append(res, questionBL)
	}

	return res, nil
}

func (q *questionRepo) Update(ctx context.Context, modelBL *models.Question) (*models.Question, error) {
	modelDB, err := questionBLToDB(modelBL)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	sql, args, err := q.Builder.
		Update("question_").
		Set("header_", modelDB.Header).
		Where(squirrel.Eq{"id_": modelDB.Id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	res, err := q.Pool.Exec(ctx, sql, args...)
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

func (q *questionRepo) Delete(ctx context.Context, id string) error {
	intid, err := strconv.Atoi(id)
	if err != nil {
		return errs.ErrInvalidContent
	}

	sql, args, err := q.Builder.
		Delete("question_").
		Where(squirrel.Eq{"id_": intid}).
		ToSql()
	if err != nil {
		return err
	}

	res, err := q.Pool.Exec(ctx, sql, args...)
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

func (q *questionRepo) GetById(ctx context.Context, id string) (*models.Question, error) {
	intid, err := strconv.Atoi(id)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	sql, args, err := q.Builder.
		Select("form_id_, header_").
		From("question_").
		Where(squirrel.Eq{"id_": intid}).
		ToSql()
	if err != nil {
		return nil, err
	}

	modelDB := QuestionDB{Id: intid}
	err = q.Pool.QueryRow(ctx, sql, args...).Scan(&modelDB.FormId, &modelDB.Header)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.ErrContentNotFound
		}

		return nil, err
	}

	return questionDBToBL(&modelDB)
}

func questionDBToBL(questionDB *QuestionDB) (*models.Question, error) {
	return &models.Question{
		Id:      strconv.Itoa(questionDB.Id),
		Form_id: strconv.Itoa(questionDB.FormId),
		Header:  questionDB.Header,
	}, nil
}

func questionBLToDB(questionBL *models.Question) (*QuestionDB, error) {
	var (
		err error
		id  int
	)

	if questionBL.Id != "" {
		id, err = strconv.Atoi(questionBL.Id)
		if err != nil {
			return nil, err
		}
	}

	var fid int
	if questionBL.Form_id != "" {
		fid, err = strconv.Atoi(questionBL.Form_id)
		if err != nil {
			return nil, err
		}
	}

	return &QuestionDB{
		Id:     id,
		FormId: fid,
		Header: questionBL.Header,
	}, nil
}
