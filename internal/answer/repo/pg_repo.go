package repo

import (
	"context"
	"errors"
	"quiz-app/internal/answer"
	"quiz-app/models"
	"quiz-app/pkg/errs"
	"quiz-app/pkg/postgres"
	"quiz-app/pkg/types"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
)

type AnswerDB struct {
	Id, QuestionId, PoolAnswerId int
	Value                        string
}

type answerRepo struct {
	*postgres.Postgres
}

func NewAnswerRepo(db *postgres.Postgres) answer.Repo {
	return &answerRepo{db}
}

func (a *answerRepo) Create(ctx context.Context, answer *models.Answer) (*models.Answer, error) {
	answerDB, err := answerBLToDB(answer)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	sql, args, err := a.Builder.
		Insert("answer_").
		Columns("question_id_, pool_answer_id_, value_").
		Values(answerDB.QuestionId, answerDB.PoolAnswerId, answerDB.Value).
		Suffix("RETURNING \"id_\"").
		ToSql()
	if err != nil {
		return nil, err
	}

	err = a.Pool.QueryRow(ctx, sql, args...).Scan(&answerDB.Id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == postgres.PermDenied {
			return nil, errs.ErrForbidden
		}

		return nil, err
	}

	return answerDBToBL(answerDB)
}

func (a *answerRepo) GetByPoolAnswerId(ctx context.Context, pool_answer_id string, sets types.GetSets) ([]*models.Answer, error) {
	intid, err := strconv.Atoi(pool_answer_id)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	sql, args, err := a.Builder.
		Select("id_, question_id_, value_").
		From("answer_").
		Where(squirrel.Eq{"pool_answer_id_": intid}).
		Limit(sets.Limit).
		Offset(sets.Offset).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := a.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*models.Answer, 0)

	for rows.Next() {
		modelDB := AnswerDB{PoolAnswerId: intid}

		err = rows.Scan(&modelDB.Id, &modelDB.QuestionId, &modelDB.Value)
		if err != nil {
			return nil, err
		}

		answerBL, err := answerDBToBL(&modelDB)
		if err != nil {
			return nil, err
		}

		res = append(res, answerBL)
	}

	return res, nil
}

func answerDBToBL(answerDB *AnswerDB) (*models.Answer, error) {
	return &models.Answer{
		Id:             strconv.Itoa(answerDB.Id),
		Question_id:    strconv.Itoa(answerDB.QuestionId),
		Pool_answer_id: strconv.Itoa(answerDB.PoolAnswerId),
		Value:          answerDB.Value,
	}, nil
}

func answerBLToDB(answerBL *models.Answer) (*AnswerDB, error) {
	var (
		err error
		id  int
	)

	if answerBL.Id != "" {
		id, err = strconv.Atoi(answerBL.Id)
		if err != nil {
			return nil, err
		}
	}

	var qid int
	if answerBL.Question_id != "" {
		qid, err = strconv.Atoi(answerBL.Question_id)
		if err != nil {
			return nil, err
		}
	}

	var paid int
	if answerBL.Pool_answer_id != "" {
		paid, err = strconv.Atoi(answerBL.Pool_answer_id)
		if err != nil {
			return nil, err
		}
	}

	return &AnswerDB{
		Id:           id,
		QuestionId:   qid,
		PoolAnswerId: paid,
		Value:        answerBL.Value,
	}, nil
}
