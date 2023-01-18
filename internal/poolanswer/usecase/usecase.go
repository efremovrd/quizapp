package usecase

import (
	"context"
	"quiz-app/internal/answer"
	"quiz-app/internal/form"
	"quiz-app/internal/poolanswer"
	"quiz-app/models"
	"quiz-app/pkg/types"
)

type poolAnswerUseCase struct {
	poolAnswerRepo poolanswer.Repo
	answerRepo     answer.Repo
	formRepo       form.Repo
}

func NewPoolAnswerUseCase(poolAnswerRepo poolanswer.Repo, answerRepo answer.Repo, formRepo form.Repo) poolanswer.UseCase {
	return &poolAnswerUseCase{
		poolAnswerRepo: poolAnswerRepo,
		answerRepo:     answerRepo,
		formRepo:       formRepo,
	}
}

func (pauc *poolAnswerUseCase) Create(ctx context.Context, pool_answer *models.PoolAnswer, answers []*models.Answer) (*models.PoolAnswer, []*models.Answer, error) {
	createdpoolanswer, err := pauc.poolAnswerRepo.Create(ctx, pool_answer)
	if err != nil {
		return nil, nil, err
	}

	var i int
	for ; i < len(answers) && err == nil; i++ {
		answers[i].Pool_answer_id = createdpoolanswer.Id
		answers[i], err = pauc.answerRepo.Create(ctx, answers[i])
	}

	if err != nil {
		pauc.poolAnswerRepo.Delete(ctx, createdpoolanswer.Id)
		return nil, nil, err
	}

	return createdpoolanswer, answers, err
}

func (pauc *poolAnswerUseCase) GetByFormId(ctx context.Context, form_id string, sets types.GetSets) ([]*models.PoolAnswer, error) {
	err := pauc.formRepo.ValidateIsOwner(ctx, form_id)
	if err != nil {
		return nil, err
	}

	return pauc.poolAnswerRepo.GetByFormId(ctx, form_id, sets)
}

func (pauc *poolAnswerUseCase) GetById(ctx context.Context, id string) (*models.PoolAnswer, error) {
	return pauc.poolAnswerRepo.GetById(ctx, id)
}
