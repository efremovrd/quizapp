package usecase

import (
	"context"
	"quizapp/internal/answer"
	"quizapp/internal/form"
	"quizapp/internal/poolanswer"
	"quizapp/models"
	"quizapp/pkg/types"
)

type answerUseCase struct {
	answerRepo answer.Repo
	formRepo   form.Repo
	paRepo     poolanswer.Repo
}

func NewAnswerUseCase(answerRepo answer.Repo, formRepo form.Repo, paRepo poolanswer.Repo) answer.UseCase {
	return &answerUseCase{
		answerRepo: answerRepo,
		formRepo:   formRepo,
		paRepo:     paRepo,
	}
}

func (answerUC *answerUseCase) GetByPoolAnswerId(ctx context.Context, pool_answer_id string, sets types.GetSets) ([]*models.Answer, error) {
	foundpa, err := answerUC.paRepo.GetById(ctx, pool_answer_id)
	if err != nil {
		return nil, err
	}

	err = answerUC.formRepo.ValidateIsOwner(ctx, foundpa.Form_id)
	if err != nil {
		return nil, err
	}

	return answerUC.answerRepo.GetByPoolAnswerId(ctx, pool_answer_id, sets)
}
